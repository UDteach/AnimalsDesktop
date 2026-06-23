package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrepareFrameKeepsTransparentSource(t *testing.T) {
	root := t.TempDir()
	srcPath := filepath.Join(root, "transparent.png")
	img := image.NewRGBA(image.Rect(0, 0, 160, 100))
	for y := 30; y < 82; y++ {
		for x := 28; x < 140; x++ {
			img.SetRGBA(x, y, color.RGBA{R: 140, G: 140, B: 136, A: 255})
		}
	}
	writeTestPNG(t, srcPath, img)

	outPath := filepath.Join(root, "frame.png")
	report, err := prepareFrame(srcPath, outPath, 18)
	if err != nil {
		t.Fatalf("prepareFrame() error = %v", err)
	}
	if report.BackgroundMode != "source-alpha" || report.BackgroundRemoved {
		t.Fatalf("background report = %s/%v", report.BackgroundMode, report.BackgroundRemoved)
	}
	out := openTestPNG(t, outPath)
	if got := out.Bounds(); got.Dx() != frameW || got.Dy() != frameH {
		t.Fatalf("output bounds = %v", got)
	}
	if alphaBounds(out, out.Bounds()).Empty() {
		t.Fatalf("output has no visible alpha")
	}
	if out.RGBAAt(0, 0).A != 0 {
		t.Fatalf("output corner alpha = %d, want transparent", out.RGBAAt(0, 0).A)
	}
}

func TestPrepareFrameRemovesUniformOpaqueBackground(t *testing.T) {
	root := t.TempDir()
	srcPath := filepath.Join(root, "opaque.png")
	img := image.NewRGBA(image.Rect(0, 0, 180, 120))
	for y := 0; y < 120; y++ {
		for x := 0; x < 180; x++ {
			img.SetRGBA(x, y, color.RGBA{R: 12, G: 240, B: 12, A: 255})
		}
	}
	for y := 35; y < 95; y++ {
		for x := 42; x < 150; x++ {
			img.SetRGBA(x, y, color.RGBA{R: 130, G: 130, B: 128, A: 255})
		}
	}
	writeTestPNG(t, srcPath, img)

	outPath := filepath.Join(root, "frame.png")
	report, err := prepareFrame(srcPath, outPath, 18)
	if err != nil {
		t.Fatalf("prepareFrame() error = %v", err)
	}
	if report.BackgroundMode != "uniform-edge-rgb" || !report.BackgroundRemoved {
		t.Fatalf("background report = %s/%v", report.BackgroundMode, report.BackgroundRemoved)
	}
	out := openTestPNG(t, outPath)
	if out.RGBAAt(0, 0).A != 0 {
		t.Fatalf("output corner alpha = %d, want transparent", out.RGBAAt(0, 0).A)
	}
	if alphaBounds(out, out.Bounds()).Empty() {
		t.Fatalf("output has no visible alpha")
	}
}

func TestPrepareFrameRemovesChromaGreenBackground(t *testing.T) {
	root := t.TempDir()
	srcPath := filepath.Join(root, "chroma.png")
	img := image.NewRGBA(image.Rect(0, 0, 180, 120))
	for y := 0; y < 120; y++ {
		for x := 0; x < 180; x++ {
			g := byte(210 + (x+y)%35)
			img.SetRGBA(x, y, color.RGBA{R: 8, G: g, B: 12, A: 255})
		}
	}
	for y := 35; y < 95; y++ {
		for x := 42; x < 150; x++ {
			img.SetRGBA(x, y, color.RGBA{R: 130, G: 130, B: 128, A: 255})
		}
	}
	writeTestPNG(t, srcPath, img)

	outPath := filepath.Join(root, "frame.png")
	report, err := prepareFrameWithMode(srcPath, outPath, "chroma-green", 18)
	if err != nil {
		t.Fatalf("prepareFrameWithMode() error = %v", err)
	}
	if report.BackgroundMode != "chroma-green" || !report.BackgroundRemoved {
		t.Fatalf("background report = %s/%v", report.BackgroundMode, report.BackgroundRemoved)
	}
	out := openTestPNG(t, outPath)
	if out.RGBAAt(0, 0).A != 0 {
		t.Fatalf("output corner alpha = %d, want transparent", out.RGBAAt(0, 0).A)
	}
	if alphaBounds(out, out.Bounds()).Empty() {
		t.Fatalf("output has no visible alpha")
	}
}

func TestPrepareFrameRejectsChromaGreenPinholes(t *testing.T) {
	root := t.TempDir()
	srcPath := filepath.Join(root, "chroma-pinholes.png")
	img := image.NewRGBA(image.Rect(0, 0, 180, 120))
	for y := 0; y < 120; y++ {
		for x := 0; x < 180; x++ {
			img.SetRGBA(x, y, color.RGBA{G: 255, A: 255})
		}
	}
	for y := 35; y < 95; y++ {
		for x := 42; x < 150; x++ {
			img.SetRGBA(x, y, color.RGBA{R: 130, G: 130, B: 128, A: 255})
		}
	}
	for y := 58; y < 60; y++ {
		for x := 86; x < 88; x++ {
			img.SetRGBA(x, y, color.RGBA{G: 255, A: 255})
		}
	}
	writeTestPNG(t, srcPath, img)

	_, err := prepareFrameWithMode(srcPath, filepath.Join(root, "frame.png"), "chroma-green", 18)
	if err == nil {
		t.Fatalf("prepareFrameWithMode() succeeded for chroma pinholes")
	}
	if !strings.Contains(err.Error(), "transparent pinholes") {
		t.Fatalf("prepareFrameWithMode() error = %v, want transparent pinholes failure", err)
	}
}

func TestPrepareFrameRejectsCheckerBackground(t *testing.T) {
	root := t.TempDir()
	srcPath := filepath.Join(root, "checker.png")
	img := image.NewRGBA(image.Rect(0, 0, 120, 80))
	for y := 0; y < 80; y++ {
		for x := 0; x < 120; x++ {
			if (x/8+y/8)%2 == 0 {
				img.SetRGBA(x, y, color.RGBA{R: 238, G: 238, B: 238, A: 255})
			} else {
				img.SetRGBA(x, y, color.RGBA{R: 190, G: 190, B: 190, A: 255})
			}
		}
	}
	for y := 24; y < 60; y++ {
		for x := 34; x < 92; x++ {
			img.SetRGBA(x, y, color.RGBA{R: 130, G: 130, B: 128, A: 255})
		}
	}
	writeTestPNG(t, srcPath, img)

	_, err := prepareFrame(srcPath, filepath.Join(root, "frame.png"), 18)
	if err == nil {
		t.Fatalf("prepareFrame() succeeded for checker background")
	}
	if !strings.Contains(err.Error(), "checker") {
		t.Fatalf("prepareFrame() error = %v, want checker/noisy failure", err)
	}
}

func TestPrepareFrameRejectsUnknownBackgroundMode(t *testing.T) {
	root := t.TempDir()
	srcPath := filepath.Join(root, "transparent.png")
	img := image.NewRGBA(image.Rect(0, 0, 120, 80))
	img.SetRGBA(40, 40, color.RGBA{R: 130, G: 130, B: 128, A: 255})
	writeTestPNG(t, srcPath, img)

	_, err := prepareFrameWithMode(srcPath, filepath.Join(root, "frame.png"), "magic", 18)
	if err == nil {
		t.Fatalf("prepareFrameWithMode() succeeded for unknown mode")
	}
	if !strings.Contains(err.Error(), "unknown") {
		t.Fatalf("prepareFrameWithMode() error = %v, want unknown-mode failure", err)
	}
}

func TestDespillGreenNeutralizesRemainingGreenEdges(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 1))
	img.SetRGBA(0, 0, color.RGBA{R: 50, G: 120, B: 55, A: 255})
	img.SetRGBA(1, 0, color.RGBA{R: 80, G: 82, B: 78, A: 255})

	despillGreen(img)

	if got := img.RGBAAt(0, 0); got.G != 55 {
		t.Fatalf("green spill pixel = %#v, want green channel clamped to max red/blue", got)
	}
	if got := img.RGBAAt(1, 0); got.G != 82 {
		t.Fatalf("neutral pixel changed to %#v", got)
	}
}

func TestClearTransparentRGBRemovesHiddenChroma(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.SetRGBA(0, 0, color.RGBA{R: 0, G: 255, B: 0, A: 0})

	clearTransparentRGB(img)

	if got := img.RGBAAt(0, 0); got != (color.RGBA{}) {
		t.Fatalf("transparent pixel = %#v, want zero RGB", got)
	}
}

func TestPrepareFrameRejectsIncompleteBackgroundRemoval(t *testing.T) {
	root := t.TempDir()
	srcPath := filepath.Join(root, "green-gradient.png")
	img := image.NewRGBA(image.Rect(0, 0, 140, 90))
	for y := 0; y < 90; y++ {
		for x := 0; x < 140; x++ {
			img.SetRGBA(x, y, color.RGBA{R: 0, G: 240, B: 0, A: 255})
		}
	}
	for y := 8; y < 82; y++ {
		for x := 0; x < 9; x++ {
			img.SetRGBA(x, y, color.RGBA{R: 0, G: 160, B: 0, A: 255})
		}
	}
	for y := 30; y < 62; y++ {
		for x := 42; x < 98; x++ {
			img.SetRGBA(x, y, color.RGBA{R: 130, G: 130, B: 128, A: 255})
		}
	}
	writeTestPNG(t, srcPath, img)

	_, err := prepareFrame(srcPath, filepath.Join(root, "frame.png"), 18)
	if err == nil {
		t.Fatalf("prepareFrame() succeeded for incomplete background removal")
	}
	if !strings.Contains(err.Error(), "content touches source canvas edge") {
		t.Fatalf("prepareFrame() error = %v, want source-edge failure", err)
	}
}

func writeTestPNG(t *testing.T, path string, img image.Image) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("encode: %v", err)
	}
}

func openTestPNG(t *testing.T, path string) *image.RGBA {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	rgba, ok := img.(*image.RGBA)
	if ok {
		return rgba
	}
	out := image.NewRGBA(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			out.Set(x, y, img.At(x, y))
		}
	}
	return out
}
