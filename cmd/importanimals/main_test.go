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
	got := normalizeSource(src, content, profileFor("hamster"))
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
		got := seedFrame(src, frame, 3, "gecko")
		if got.Bounds().Dx() != frameW || got.Bounds().Dy() != frameH {
			t.Fatalf("frame %d bounds = %v", frame, got.Bounds())
		}
		if alphaBounds(got).Empty() {
			t.Fatalf("frame %d has no visible pixels", frame)
		}
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
		}
	}
}

func twoDigits(v int) string {
	return string([]byte{'0' + byte(v/10), '0' + byte(v%10)})
}
