package main

import (
	"image"
	"image/color"
	"testing"
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
