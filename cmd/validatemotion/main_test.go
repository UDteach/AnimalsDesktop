package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"animals-desktop/internal/catalog"
)

func TestValidateVariantReportsDraftNotReleaseReady(t *testing.T) {
	root := t.TempDir()
	set00Path := writeMotionSourceFamily(t, root)
	variant := catalog.Variant{
		ID:               "chinchilla_standard_gray",
		SpeciesID:        "chinchilla",
		SourceStatus:     catalog.SourceStatusMotionDraft,
		MotionSourcePath: set00Path,
	}

	report, err := validateVariant(variant)
	if err != nil {
		t.Fatalf("validateVariant() error = %v", err)
	}
	if report.RuntimeSets != motionSets || report.FramesPerSet != totalFrames {
		t.Fatalf("report sets/frames = %d/%d, want %d/%d", report.RuntimeSets, report.FramesPerSet, motionSets, totalFrames)
	}
	if report.ReleaseReady {
		t.Fatalf("draft motion source reported release-ready")
	}
	if len(report.Warnings) == 0 {
		t.Fatalf("draft motion source should report a warning")
	}
}

func TestValidateVariantReportsAcceptedReleaseReady(t *testing.T) {
	root := t.TempDir()
	set00Path := writeMotionSourceFamily(t, root)
	variant := catalog.Variant{
		ID:               "chinchilla_standard_gray",
		SpeciesID:        "chinchilla",
		SourceStatus:     catalog.SourceStatusMotionAccepted,
		MotionSourcePath: set00Path,
	}

	report, err := validateVariant(variant)
	if err != nil {
		t.Fatalf("validateVariant() error = %v", err)
	}
	if !report.ReleaseReady {
		t.Fatalf("accepted motion source reported not release-ready")
	}
	if report.UniqueSetHashes != motionSets {
		t.Fatalf("unique set hashes = %d, want %d", report.UniqueSetHashes, motionSets)
	}
}

func TestValidateVariantAllowsSingleDraftSourceSet(t *testing.T) {
	root := t.TempDir()
	set00Path := filepath.Join(root, "animal-set00-source.png")
	sheet := image.NewRGBA(image.Rect(0, 0, frameW*totalFrames, frameH))
	for frame := 0; frame < totalFrames; frame++ {
		x0 := frame * frameW
		sheet.SetRGBA(x0+6, 6, color.RGBA{R: 80, G: byte(frame), B: 80, A: 255})
		for y := 28; y < 52; y++ {
			for x := 20; x < 76; x++ {
				sheet.SetRGBA(x0+x, y, color.RGBA{R: 150, G: 150, B: 145, A: 255})
			}
		}
	}
	writePNG(t, set00Path, sheet)

	report, err := validateVariant(catalog.Variant{
		ID:               "chinchilla_standard_gray",
		SpeciesID:        "chinchilla",
		SourceStatus:     catalog.SourceStatusMotionDraft,
		MotionSourcePath: set00Path,
	})
	if err != nil {
		t.Fatalf("validateVariant() error = %v", err)
	}
	if report.RuntimeSets != 1 {
		t.Fatalf("runtime sets = %d, want single draft source set", report.RuntimeSets)
	}
	if report.ReleaseReady {
		t.Fatalf("single draft source set reported release-ready")
	}
	if len(report.Warnings) < 2 {
		t.Fatalf("warnings = %v, want draft and source-set count warnings", report.Warnings)
	}
}

func TestValidateVariantRejectsOpaqueMotionFrameBackground(t *testing.T) {
	root := t.TempDir()
	set00Path := filepath.Join(root, "animal-set00-source.png")
	sheet := image.NewRGBA(image.Rect(0, 0, frameW*totalFrames, frameH))
	for y := 0; y < frameH; y++ {
		for x := 0; x < frameW*totalFrames; x++ {
			sheet.SetRGBA(x, y, color.RGBA{R: 230, G: 230, B: 230, A: 255})
		}
	}
	writePNG(t, set00Path, sheet)

	_, err := validateVariant(catalog.Variant{
		ID:               "chinchilla_standard_gray",
		SpeciesID:        "chinchilla",
		SourceStatus:     catalog.SourceStatusMotionDraft,
		MotionSourcePath: set00Path,
	})
	if err == nil {
		t.Fatalf("validateVariant() succeeded for opaque motion source")
	}
	if !strings.Contains(err.Error(), "transparent background") {
		t.Fatalf("validateVariant() error = %v, want transparent-background failure", err)
	}
}

func writeMotionSourceFamily(t *testing.T, root string) string {
	t.Helper()
	set00Path := filepath.Join(root, "animal-set00-source.png")
	for set := 0; set < motionSets; set++ {
		sheet := image.NewRGBA(image.Rect(0, 0, frameW*totalFrames, frameH))
		for frame := 0; frame < totalFrames; frame++ {
			x0 := frame * frameW
			sheet.SetRGBA(x0+6, 6, color.RGBA{R: byte(set), G: byte(frame), B: 80, A: 255})
			for y := 28; y < 52; y++ {
				for x := 20; x < 76; x++ {
					sheet.SetRGBA(x0+x, y, color.RGBA{R: 150, G: 150, B: 145, A: 255})
				}
			}
		}
		path := filepath.Join(root, "animal-set"+twoDigits(set)+"-source.png")
		writePNG(t, path, sheet)
	}
	return set00Path
}

func writePNG(t *testing.T, path string, img image.Image) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create png: %v", err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}
}

func twoDigits(v int) string {
	return string([]byte{'0' + byte(v/10), '0' + byte(v%10)})
}
