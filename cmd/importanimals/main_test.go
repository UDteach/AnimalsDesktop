package main

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"animals-desktop/internal/catalog"
)

func TestNormalizeSourceUsesFixedCanvas(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 200, 120))
	for y := 35; y < 95; y++ {
		for x := 20; x < 180; x++ {
			src.SetRGBA(x, y, color.RGBA{R: 150, G: 110, B: 70, A: 255})
		}
	}
	content := alphaBounds(src)
	if content.Empty() {
		t.Fatalf("alphaBounds returned empty content")
	}
	got := normalizeSource(src, content, profileFor(catalog.MotionProfileSmallRodentScurry))
	if got.Bounds().Dx() != frameW || got.Bounds().Dy() != frameH {
		t.Fatalf("normalized bounds = %v, want %dx%d", got.Bounds(), frameW, frameH)
	}
	if alphaBounds(got).Empty() {
		t.Fatalf("normalized source has no visible pixels")
	}
}

func TestSeedFrameKeepsSpriteInCanvas(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, frameW, frameH))
	for y := 30; y < 52; y++ {
		for x := 20; x < 76; x++ {
			src.SetRGBA(x, y, color.RGBA{R: 110, G: 110, B: 110, A: 255})
		}
	}
	for frame := 0; frame < totalFrames; frame++ {
		got := seedFrame(src, frame, 3, catalog.MotionProfileGeckoCrawl)
		if got.Bounds().Dx() != frameW || got.Bounds().Dy() != frameH {
			t.Fatalf("frame %d bounds = %v", frame, got.Bounds())
		}
		if alphaBounds(got).Empty() {
			t.Fatalf("frame %d has no visible pixels", frame)
		}
	}
}

func TestImportVariantUsesMotionSourceSheet(t *testing.T) {
	root := t.TempDir()
	sourcePath := filepath.Join(root, "source.png")
	src := image.NewRGBA(image.Rect(0, 0, 120, 80))
	for y := 20; y < 60; y++ {
		for x := 20; x < 100; x++ {
			src.SetRGBA(x, y, color.RGBA{R: 140, G: 140, B: 135, A: 255})
		}
	}
	if err := writePNG(sourcePath, src); err != nil {
		t.Fatalf("write source: %v", err)
	}

	motionPath := filepath.Join(root, "motion-source.png")
	motion := image.NewRGBA(image.Rect(0, 0, frameW*totalFrames, frameH))
	for frame := 0; frame < totalFrames; frame++ {
		xOffset := frame * frameW
		motion.SetRGBA(xOffset+8, 8, color.RGBA{R: byte(frame), G: 40, B: 90, A: 255})
		for y := 34; y < 54; y++ {
			for x := 24; x < 72; x++ {
				motion.SetRGBA(xOffset+x, y, color.RGBA{R: 160, G: 160, B: 150, A: 255})
			}
		}
	}
	if err := writePNG(motionPath, motion); err != nil {
		t.Fatalf("write motion source: %v", err)
	}

	variant := catalog.Variant{
		ID:               "chinchilla_standard_gray",
		SpeciesID:        "chinchilla",
		BreedOrMorph:     "Chinchilla",
		Color:            "standard gray",
		PopularityTier:   1,
		MotionProfile:    catalog.MotionProfileSmallRodentScurry,
		SourceStatus:     catalog.SourceStatusMotionDraft,
		SpriteBase:       "test_chinchilla",
		SeedStage:        true,
		SourcePath:       sourcePath,
		MotionSourcePath: motionPath,
	}
	outDir := filepath.Join(root, "sprites")
	report, err := importVariant(variant, outDir, filepath.Join(root, "generated"))
	if err != nil {
		t.Fatalf("importVariant() error = %v", err)
	}
	if report.MotionFrames != totalFrames || report.MotionSets != 1 || report.MotionSource == "" {
		t.Fatalf("motion report = frames:%d sets:%d source:%q", report.MotionFrames, report.MotionSets, report.MotionSource)
	}
	if len(report.Outputs) != motionSets {
		t.Fatalf("outputs = %d, want %d", len(report.Outputs), motionSets)
	}
	sheet, err := openPNG(filepath.Join(outDir, "test_chinchilla_set00.png"))
	if err != nil {
		t.Fatalf("open sheet: %v", err)
	}
	got := sheet.RGBAAt(12*frameW+8, 8)
	if got.R != 12 || got.G != 40 || got.B != 90 || got.A != 255 {
		t.Fatalf("frame 12 marker = %#v, want imported motion source marker", got)
	}
}

func TestProceduralSourceHasVisibleContent(t *testing.T) {
	src := proceduralSource(catalog.Variant{
		ID:        "ferret_sable",
		SpeciesID: "ferret",
		Shape:     "ferret",
		TintHex:   "8b6746",
		AccentHex: "ece0c8",
	})
	content := alphaBounds(src)
	if content.Empty() {
		t.Fatalf("procedural source content is empty")
	}
	if content.Dx() < 180 || content.Dy() < 50 {
		t.Fatalf("procedural source content too small: %v", content)
	}
}

func TestTintSourcePreservesAlpha(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 8, 8))
	src.SetRGBA(2, 2, color.RGBA{R: 200, G: 180, B: 120, A: 190})
	tint, err := parseHexColor("336699")
	if err != nil {
		t.Fatalf("parseHexColor() error = %v", err)
	}
	got := tintSource(src, tint, nil)
	if got.RGBAAt(2, 2).A != 190 {
		t.Fatalf("tinted alpha = %d, want 190", got.RGBAAt(2, 2).A)
	}
	if got.RGBAAt(0, 0).A != 0 {
		t.Fatalf("transparent pixel alpha = %d, want 0", got.RGBAAt(0, 0).A)
	}
}

func TestSeedVariantGeneratedAssetsExist(t *testing.T) {
	for _, variant := range catalog.SeedVariants() {
		source := filepath.Join("..", "..", "assets", "source", "animals", "generated", variant.SpriteBase+"-source.png")
		if _, err := os.Stat(source); err != nil {
			t.Fatalf("missing generated source for %s: %v", variant.ID, err)
		}
		for set := 0; set < motionSets; set++ {
			path := filepath.Join("..", "..", "assets", "sprites", variant.SpriteBase+"_set"+twoDigits(set)+".png")
			if _, err := os.Stat(path); err != nil {
				t.Fatalf("missing sprite sheet for %s set %02d: %v", variant.ID, set, err)
			}
			sheet, err := openPNG(path)
			if err != nil {
				t.Fatalf("open sprite sheet for %s set %02d: %v", variant.ID, set, err)
			}
			if got := sheet.Bounds(); got.Dx() != frameW*totalFrames || got.Dy() != frameH {
				t.Fatalf("sheet bounds for %s set %02d = %v", variant.ID, set, got)
			}
			for frame := 0; frame < totalFrames; frame++ {
				frameRect := image.Rect(frame*frameW, 0, (frame+1)*frameW, frameH)
				if alphaBounds(sheet.SubImage(frameRect)).Empty() {
					t.Fatalf("empty frame for %s set %02d frame %02d", variant.ID, set, frame)
				}
			}
		}
	}
}

func TestMotionProfilesHaveDistinctEcologyOffsets(t *testing.T) {
	rabbitDx, rabbitDy := motionOffset(27, 0, catalog.MotionProfileRabbitHop)
	_, snakeDy := motionOffset(27, 0, catalog.MotionProfileSnakeSlither)
	_, tortoiseDy := motionOffset(14, 0, catalog.MotionProfileTortoisePlod)
	dogDx, dogDy := motionOffset(13, 0, catalog.MotionProfileDogTrot)

	if rabbitDy >= 0 || rabbitDx != 0 {
		t.Fatalf("rabbit action offset = (%d,%d), want upward hop without horizontal drift", rabbitDx, rabbitDy)
	}
	if snakeDy != 0 {
		t.Fatalf("snake action dy = %d, want low slither", snakeDy)
	}
	if tortoiseDy != 0 {
		t.Fatalf("tortoise fast dy = %d, want no vertical bob", tortoiseDy)
	}
	if dogDy >= 0 || dogDx == 0 {
		t.Fatalf("dog trot offset = (%d,%d), want horizontal trot with lift", dogDx, dogDy)
	}
}

func twoDigits(v int) string {
	return string([]byte{'0' + byte(v/10), '0' + byte(v%10)})
}
