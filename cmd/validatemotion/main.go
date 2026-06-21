package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"animals-desktop/internal/catalog"
)

const (
	frameW      = 96
	frameH      = 64
	totalFrames = 62
	motionSets  = 10
)

type variantReport struct {
	Variant          string   `json:"variant"`
	Species          string   `json:"species"`
	SourceStatus     string   `json:"source_status"`
	MotionSource     string   `json:"motion_source"`
	MotionSetSources []string `json:"motion_set_sources"`
	FramesPerSet     int      `json:"frames_per_set"`
	RuntimeSets      int      `json:"runtime_sets"`
	UniqueSetHashes  int      `json:"unique_set_hashes"`
	ReleaseReady     bool     `json:"release_ready"`
	Warnings         []string `json:"warnings,omitempty"`
}

func main() {
	variantID := flag.String("variant", "", "optional variant ID to validate")
	runtimeOnly := flag.Bool("runtime-only", false, "validate runtime variants only")
	requireAccepted := flag.Bool("require-accepted", false, "fail if selected variants are not accepted release-quality motion sources")
	flag.Parse()

	variants := catalog.Variants
	if *runtimeOnly {
		variants = catalog.RuntimeVariants()
	}
	if *variantID != "" {
		variant, ok := catalog.VariantByID(*variantID)
		if !ok {
			fatalf("unknown variant %q", *variantID)
		}
		variants = []catalog.Variant{variant}
	}

	reports := make([]variantReport, 0)
	hasFailure := false
	for _, variant := range variants {
		if variant.MotionSourcePath == "" {
			if *requireAccepted {
				fmt.Fprintf(os.Stderr, "%s has no motion source path\n", variant.ID)
				hasFailure = true
			}
			continue
		}
		report, err := validateVariant(variant)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", variant.ID, err)
			hasFailure = true
			continue
		}
		if *requireAccepted && !report.ReleaseReady {
			fmt.Fprintf(os.Stderr, "%s is not release-ready: source status %s\n", variant.ID, variant.SourceStatus)
			hasFailure = true
		}
		reports = append(reports, report)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(reports); err != nil {
		fatalf("write report: %v", err)
	}
	if hasFailure {
		os.Exit(1)
	}
}

func validateVariant(variant catalog.Variant) (variantReport, error) {
	paths, err := motionSourceSheetPaths(variant.MotionSourcePath)
	if err != nil {
		return variantReport{}, err
	}

	hashes := map[string]bool{}
	for _, path := range paths {
		sum, err := validateSheet(path)
		if err != nil {
			return variantReport{}, err
		}
		hashes[sum] = true
	}

	report := variantReport{
		Variant:          variant.ID,
		Species:          variant.SpeciesID,
		SourceStatus:     variant.SourceStatus,
		MotionSource:     filepath.ToSlash(variant.MotionSourcePath),
		MotionSetSources: slashPaths(paths),
		FramesPerSet:     totalFrames,
		RuntimeSets:      len(paths),
		UniqueSetHashes:  len(hashes),
		ReleaseReady:     variant.SourceStatus == catalog.SourceStatusMotionAccepted && len(paths) == motionSets,
	}
	if variant.SourceStatus == catalog.SourceStatusMotionDraft {
		report.Warnings = append(report.Warnings, "motion source is draft and must not be released")
	}
	if len(paths) < motionSets {
		report.Warnings = append(report.Warnings, fmt.Sprintf("motion source has %d source set(s), accepted release requires %d", len(paths), motionSets))
	}
	if len(paths) == motionSets && len(hashes) < motionSets {
		report.Warnings = append(report.Warnings, "one or more motion source sheets are byte-identical")
	}
	return report, nil
}

func motionSourceSheetPaths(set00Path string) ([]string, error) {
	if set00Path == "" {
		return nil, fmt.Errorf("motion source path is empty")
	}
	if !strings.Contains(set00Path, "set00") {
		if _, err := os.Stat(set00Path); err != nil {
			return nil, err
		}
		return []string{set00Path}, nil
	}
	paths := make([]string, 0, motionSets)
	for set := 0; set < motionSets; set++ {
		path := strings.Replace(set00Path, "set00", fmt.Sprintf("set%02d", set), 1)
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) && set > 0 {
				return []string{set00Path}, nil
			}
			return nil, err
		}
		paths = append(paths, path)
	}
	return paths, nil
}

func validateSheet(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("decode PNG %s: %w", path, err)
	}
	bounds := img.Bounds()
	wantW := frameW * totalFrames
	if bounds.Dx() != wantW || bounds.Dy() != frameH {
		return "", fmt.Errorf("%s bounds = %dx%d, want %dx%d", path, bounds.Dx(), bounds.Dy(), wantW, frameH)
	}
	for frame := 0; frame < totalFrames; frame++ {
		frameRect := image.Rect(bounds.Min.X+frame*frameW, bounds.Min.Y, bounds.Min.X+(frame+1)*frameW, bounds.Min.Y+frameH)
		content := alphaBounds(img, frameRect)
		if content.Empty() {
			return "", fmt.Errorf("%s frame %02d has no visible alpha", path, frame)
		}
		if content == frameRect {
			return "", fmt.Errorf("%s frame %02d has no transparent background", path, frame)
		}
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func alphaBounds(img image.Image, rect image.Rectangle) image.Rectangle {
	rect = rect.Intersect(img.Bounds())
	minX, minY := rect.Max.X, rect.Max.Y
	maxX, maxY := rect.Min.X, rect.Min.Y
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a <= 0x0800 {
				continue
			}
			if x < minX {
				minX = x
			}
			if y < minY {
				minY = y
			}
			if x+1 > maxX {
				maxX = x + 1
			}
			if y+1 > maxY {
				maxY = y + 1
			}
		}
	}
	if maxX <= minX || maxY <= minY {
		return image.Rectangle{}
	}
	return image.Rect(minX, minY, maxX, maxY)
}

func slashPaths(paths []string) []string {
	out := make([]string, len(paths))
	for i, path := range paths {
		out[i] = filepath.ToSlash(path)
	}
	return out
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
