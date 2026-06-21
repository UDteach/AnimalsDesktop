package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"animals-desktop/internal/catalog"
	xdraw "golang.org/x/image/draw"
)

const (
	totalFrames = 62
	motionSets  = 10
	frameW      = 96
	frameH      = 64
)

type seedReport struct {
	Variant      string     `json:"variant"`
	Species      string     `json:"species"`
	Source       string     `json:"source"`
	SpriteBase   string     `json:"sprite_base"`
	SourceWidth  int        `json:"source_width"`
	SourceHeight int        `json:"source_height"`
	Content      rectReport `json:"content"`
	Outputs      []string   `json:"outputs"`
	Warnings     []string   `json:"warnings,omitempty"`
}

type rectReport struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type renderProfile struct {
	targetW  int
	targetH  int
	baseline int
	low      bool
}

func main() {
	outDir := flag.String("out", filepath.FromSlash("assets/sprites"), "output sprite directory")
	reportPath := flag.String("report", filepath.FromSlash("assets/source/animals/seed-import-report.json"), "JSON report path")
	previewPath := flag.String("preview", filepath.FromSlash("docs/assets/animalsdesktop-seed-preview.png"), "seed preview PNG path")
	flag.Parse()

	seedVariants := catalog.SeedVariants()
	reports := make([]seedReport, 0)
	for _, variant := range seedVariants {
		report, err := importVariant(variant, *outDir)
		if err != nil {
			log.Fatalf("import %s: %v", variant.ID, err)
		}
		reports = append(reports, report)
	}
	writeJSON(*reportPath, reports)
	if err := writeSeedPreview(seedVariants, *outDir, *previewPath); err != nil {
		log.Fatalf("write preview: %v", err)
	}
	fmt.Printf("imported %d seed animal variants into %d-frame sheets\n", len(reports), totalFrames)
}

func importVariant(variant catalog.Variant, outDir string) (seedReport, error) {
	src, err := openPNG(variant.SourcePath)
	if err != nil {
		return seedReport{}, err
	}
	bounds := src.Bounds()
	content := alphaBounds(src)
	report := seedReport{
		Variant:      variant.ID,
		Species:      variant.SpeciesID,
		Source:       filepath.ToSlash(variant.SourcePath),
		SpriteBase:   variant.SpriteBase,
		SourceWidth:  bounds.Dx(),
		SourceHeight: bounds.Dy(),
		Content: rectReport{
			X: content.Min.X,
			Y: content.Min.Y,
			W: content.Dx(),
			H: content.Dy(),
		},
		Outputs: make([]string, 0, motionSets),
	}
	if content.Empty() {
		return seedReport{}, fmt.Errorf("source has no visible alpha: %s", variant.SourcePath)
	}
	if content == bounds {
		report.Warnings = append(report.Warnings, "source is fully opaque; verify transparent background before final animation import")
	}

	normalized := normalizeSource(src, content, profileFor(variant.SpeciesID))
	for set := 0; set < motionSets; set++ {
		sheet := image.NewRGBA(image.Rect(0, 0, frameW*totalFrames, frameH))
		for frame := 0; frame < totalFrames; frame++ {
			sprite := seedFrame(normalized, frame, set, variant.SpeciesID)
			dst := image.Rect(frame*frameW, 0, (frame+1)*frameW, frameH)
			draw.Draw(sheet, dst, sprite, image.Point{}, draw.Over)
		}
		outPath := filepath.Join(outDir, fmt.Sprintf("%s_set%02d.png", variant.SpriteBase, set))
		if err := writePNG(outPath, sheet); err != nil {
			return seedReport{}, err
		}
		report.Outputs = append(report.Outputs, filepath.ToSlash(outPath))
	}
	return report, nil
}

func profileFor(speciesID string) renderProfile {
	switch speciesID {
	case "gecko":
		return renderProfile{targetW: 90, targetH: 30, baseline: 59, low: true}
	case "rabbit":
		return renderProfile{targetW: 82, targetH: 56, baseline: 60}
	case "dog", "cat":
		return renderProfile{targetW: 84, targetH: 54, baseline: 59}
	case "hamster", "macaroni_mouse":
		return renderProfile{targetW: 76, targetH: 48, baseline: 59}
	case "chinchilla":
		return renderProfile{targetW: 82, targetH: 50, baseline: 59}
	default:
		return renderProfile{targetW: 84, targetH: 52, baseline: 59}
	}
}

func normalizeSource(src image.Image, content image.Rectangle, profile renderProfile) *image.RGBA {
	content = content.Intersect(src.Bounds())
	if content.Empty() {
		content = src.Bounds()
	}
	scale := minFloat(
		float64(profile.targetW)/float64(maxInt(1, content.Dx())),
		float64(profile.targetH)/float64(maxInt(1, content.Dy())),
	)
	w := maxInt(1, int(float64(content.Dx())*scale+0.5))
	h := maxInt(1, int(float64(content.Dy())*scale+0.5))
	x := (frameW - w) / 2
	y := profile.baseline - h
	if profile.low {
		y = minInt(frameH-h-1, maxInt(0, y+4))
	} else {
		y = maxInt(0, minInt(frameH-h-1, y))
	}

	scaled := image.NewRGBA(image.Rect(0, 0, w, h))
	xdraw.CatmullRom.Scale(scaled, scaled.Bounds(), src, content, draw.Over, nil)
	out := image.NewRGBA(image.Rect(0, 0, frameW, frameH))
	draw.Draw(out, image.Rect(x, y, x+w, y+h), withOutline(scaled), image.Point{}, draw.Over)
	return out
}

func withOutline(src *image.RGBA) *image.RGBA {
	out := image.NewRGBA(src.Bounds())
	shadow := color.RGBA{R: 34, G: 30, B: 25, A: 78}
	for y := src.Bounds().Min.Y; y < src.Bounds().Max.Y; y++ {
		for x := src.Bounds().Min.X; x < src.Bounds().Max.X; x++ {
			if src.RGBAAt(x, y).A < 32 {
				continue
			}
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}
					px, py := x+dx, y+dy
					if !image.Pt(px, py).In(out.Bounds()) || out.RGBAAt(px, py).A >= shadow.A {
						continue
					}
					out.SetRGBA(px, py, shadow)
				}
			}
		}
	}
	draw.Draw(out, out.Bounds(), src, src.Bounds().Min, draw.Over)
	return out
}

func seedFrame(src *image.RGBA, frame int, set int, speciesID string) *image.RGBA {
	dx, dy := motionOffset(frame, set, speciesID)
	out := image.NewRGBA(image.Rect(0, 0, frameW, frameH))
	dst := src.Bounds().Add(image.Pt(dx, dy)).Intersect(out.Bounds())
	if dst.Empty() {
		return out
	}
	draw.Draw(out, dst, src, image.Pt(dst.Min.X-dx, dst.Min.Y-dy), draw.Over)
	return out
}

func motionOffset(frame int, set int, speciesID string) (int, int) {
	local := frame
	phase := set % 4
	switch {
	case frame < 4:
		return 0, []int{0, -1, 0, 0}[(local+phase)%4]
	case frame < 12:
		local = frame - 4
		return []int{-2, -1, 0, 1, 2, 1, 0, -1}[(local+phase)%8], []int{0, -1, 0, 0}[(local+phase)%4]
	case frame < 20:
		local = frame - 12
		step := []int{-3, -1, 1, 3, 2, 0, -2, -3}[(local+phase)%8]
		if speciesID == "gecko" {
			return step, 0
		}
		return step, []int{0, -1, 0, -1}[(local+phase)%4]
	case frame < 26:
		local = frame - 20
		return 0, []int{0, 1, 0, 1, 0, 0}[(local+phase)%6]
	case frame < 32:
		local = frame - 26
		if speciesID == "gecko" {
			return []int{-2, -1, 0, 1, 2, 1}[(local+phase)%6], 0
		}
		return 0, []int{0, -2, -5, -4, -1, 0}[(local+phase)%6]
	case frame < 40:
		local = frame - 32
		return []int{0, 1, 2, 1, 0, -1, -2, -1}[(local+phase)%8], 0
	case frame < 48:
		local = frame - 40
		return []int{0, 1, 0, 1, 0, -1, 0, -1}[(local+phase)%8], []int{0, 1, 0, 1, 0, 1, 0, 1}[(local+phase)%8]
	case frame < 56:
		local = frame - 48
		return 0, []int{0, -1, -1, 0, 0, 1, 1, 0}[(local+phase)%8]
	default:
		local = frame - 56
		if speciesID == "gecko" {
			return []int{-2, 0, 2, 1, -1, -2}[(local+phase)%6], 0
		}
		return []int{-2, 0, 2, 1, -1, -2}[(local+phase)%6], []int{0, -1, 0, 1, 0, -1}[(local+phase)%6]
	}
}

func openPNG(path string) (*image.RGBA, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}
	dst := image.NewRGBA(img.Bounds())
	draw.Draw(dst, dst.Bounds(), img, img.Bounds().Min, draw.Src)
	return dst, nil
}

func writePNG(path string, img image.Image) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func writeJSON(path string, report []seedReport) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		log.Fatalf("create report dir: %v", err)
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		log.Fatalf("marshal report: %v", err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		log.Fatalf("write report: %v", err)
	}
}

func writeSeedPreview(variants []catalog.Variant, outDir string, path string) error {
	if len(variants) == 0 {
		return nil
	}
	const (
		cols = 4
		pad  = 8
	)
	rows := (len(variants) + cols - 1) / cols
	w := cols*frameW + (cols+1)*pad
	h := rows*frameH + (rows+1)*pad
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	drawChecker(dst)
	for i, variant := range variants {
		sheet, err := openPNG(filepath.Join(outDir, fmt.Sprintf("%s_set00.png", variant.SpriteBase)))
		if err != nil {
			return err
		}
		cellX := pad + (i%cols)*(frameW+pad)
		cellY := pad + (i/cols)*(frameH+pad)
		draw.Draw(dst, image.Rect(cellX, cellY, cellX+frameW, cellY+frameH), sheet, image.Point{}, draw.Over)
	}
	return writePNG(path, dst)
}

func drawChecker(dst *image.RGBA) {
	a := color.RGBA{R: 235, G: 237, B: 233, A: 255}
	b := color.RGBA{R: 213, G: 217, B: 211, A: 255}
	for y := dst.Bounds().Min.Y; y < dst.Bounds().Max.Y; y++ {
		for x := dst.Bounds().Min.X; x < dst.Bounds().Max.X; x++ {
			if (x/8+y/8)%2 == 0 {
				dst.SetRGBA(x, y, a)
			} else {
				dst.SetRGBA(x, y, b)
			}
		}
	}
}

func alphaBounds(img image.Image) image.Rectangle {
	b := img.Bounds()
	minX, minY := b.Max.X, b.Max.Y
	maxX, maxY := b.Min.X, b.Min.Y
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
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

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
