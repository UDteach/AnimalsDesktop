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

func TestAuditReportsPartialFrameSet(t *testing.T) {
	root := t.TempDir()
	framesDir := filepath.Join(root, "set00")
	writeFrame(t, filepath.Join(framesDir, "frame-00.png"), frameW, frameH, false)
	writeFrame(t, filepath.Join(framesDir, "frame-01.png"), frameW+1, frameH, false)

	report, err := audit("", framesDir, "frame-%02d.png", false, false)
	if err != nil {
		t.Fatalf("audit() error = %v", err)
	}
	if report.Valid != 1 || report.Invalid != 1 || report.Missing != totalFrames-2 {
		t.Fatalf("report counts valid/invalid/missing = %d/%d/%d", report.Valid, report.Invalid, report.Missing)
	}
	if report.Sets[0].Completed {
		t.Fatalf("partial set reported complete")
	}
}

func TestAuditRejectsOpaqueFrame(t *testing.T) {
	root := t.TempDir()
	framesDir := filepath.Join(root, "set00")
	writeFrame(t, filepath.Join(framesDir, "frame-00.png"), frameW, frameH, true)

	report, err := audit("", framesDir, "frame-%02d.png", false, false)
	if err != nil {
		t.Fatalf("audit() error = %v", err)
	}
	frame := report.Sets[0].Frames[0]
	if frame.Status != "invalid" || frame.Error != "no transparent background" {
		t.Fatalf("frame status/error = %q/%q", frame.Status, frame.Error)
	}
}

func TestAuditRootScansTenSets(t *testing.T) {
	root := t.TempDir()
	writeFrame(t, filepath.Join(root, "set03", "frame-12.png"), frameW, frameH, false)

	report, err := audit(root, "", "frame-%02d.png", false, false)
	if err != nil {
		t.Fatalf("audit() error = %v", err)
	}
	if report.SetCount != motionSets || report.FrameCount != motionSets*totalFrames {
		t.Fatalf("report set/frame count = %d/%d", report.SetCount, report.FrameCount)
	}
	if report.Valid != 1 {
		t.Fatalf("valid frames = %d, want 1", report.Valid)
	}
}

func TestAuditArtifactWarningsDetectsLowHorizontalRun(t *testing.T) {
	root := t.TempDir()
	framesDir := filepath.Join(root, "set00")
	path := filepath.Join(framesDir, "frame-00.png")
	writeFrame(t, path, frameW, frameH, false)
	addHorizontalRun(t, path, 16, 55, 64)

	report, err := audit("", framesDir, "frame-%02d.png", false, true)
	if err != nil {
		t.Fatalf("audit() error = %v", err)
	}
	frame := report.Sets[0].Frames[0]
	if frame.Status != "valid" {
		t.Fatalf("frame status = %q", frame.Status)
	}
	if len(frame.Warnings) == 0 {
		t.Fatalf("expected artifact warning")
	}
}

func TestAuditArtifactWarningsDetectsDisconnectedComponents(t *testing.T) {
	root := t.TempDir()
	framesDir := filepath.Join(root, "set00")
	path := filepath.Join(framesDir, "frame-00.png")
	writeFrame(t, path, frameW, frameH, false)
	addDetachedBlock(t, path, 8, 57, 4, 4)

	report, err := audit("", framesDir, "frame-%02d.png", false, true)
	if err != nil {
		t.Fatalf("audit() error = %v", err)
	}
	frame := report.Sets[0].Frames[0]
	warnings := strings.Join(frame.Warnings, "\n")
	if !strings.Contains(warnings, "disconnected alpha components") {
		t.Fatalf("warnings = %q, want disconnected alpha components warning", warnings)
	}
}

func TestAuditArtifactWarningsDetectsTransparentPinholes(t *testing.T) {
	root := t.TempDir()
	framesDir := filepath.Join(root, "set00")
	path := filepath.Join(framesDir, "frame-00.png")
	writeFrame(t, path, frameW, frameH, false)
	addTransparentHole(t, path, 42, 36, 2, 2)

	report, err := audit("", framesDir, "frame-%02d.png", false, true)
	if err != nil {
		t.Fatalf("audit() error = %v", err)
	}
	frame := report.Sets[0].Frames[0]
	warnings := strings.Join(frame.Warnings, "\n")
	if !strings.Contains(warnings, "transparent pinholes") {
		t.Fatalf("warnings = %q, want transparent pinholes warning", warnings)
	}
}

func TestAuditArtifactWarningsDetectsLowerShelf(t *testing.T) {
	root := t.TempDir()
	framesDir := filepath.Join(root, "set00")
	path := filepath.Join(framesDir, "frame-00.png")
	writeFrame(t, path, frameW, frameH, false)
	addLowerShelf(t, path)

	report, err := audit("", framesDir, "frame-%02d.png", false, true)
	if err != nil {
		t.Fatalf("audit() error = %v", err)
	}
	frame := report.Sets[0].Frames[0]
	warnings := strings.Join(frame.Warnings, "\n")
	if !strings.Contains(warnings, "possible lower ledge/shelf artifact") {
		t.Fatalf("warnings = %q, want lower shelf artifact warning", warnings)
	}
}

func writeFrame(t *testing.T, path string, w int, h int, opaque bool) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	if opaque {
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				img.SetRGBA(x, y, color.RGBA{R: 210, G: 210, B: 210, A: 255})
			}
		}
	} else {
		for y := 24; y < minInt(54, h); y++ {
			for x := 20; x < minInt(76, w); x++ {
				img.SetRGBA(x, y, color.RGBA{R: 130, G: 130, B: 130, A: 255})
			}
		}
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

func addDetachedBlock(t *testing.T, path string, x int, y int, w int, h int) {
	t.Helper()
	updateFrame(t, path, func(rgba *image.RGBA) {
		for py := y; py < y+h; py++ {
			for px := x; px < x+w; px++ {
				rgba.SetRGBA(px, py, color.RGBA{R: 70, G: 60, B: 50, A: 230})
			}
		}
	})
}

func addTransparentHole(t *testing.T, path string, x int, y int, w int, h int) {
	t.Helper()
	updateFrame(t, path, func(rgba *image.RGBA) {
		for py := y; py < y+h; py++ {
			for px := x; px < x+w; px++ {
				rgba.SetRGBA(px, py, color.RGBA{})
			}
		}
	})
}

func addLowerShelf(t *testing.T, path string) {
	t.Helper()
	updateFrame(t, path, func(rgba *image.RGBA) {
		for px := 8; px < 88; px++ {
			rgba.SetRGBA(px, 50, color.RGBA{R: 90, G: 84, B: 76, A: 190})
		}
		for py := 54; py < 58; py++ {
			for px := 30; px < 34; px++ {
				rgba.SetRGBA(px, py, color.RGBA{R: 80, G: 76, B: 70, A: 220})
			}
		}
	})
}

func addHorizontalRun(t *testing.T, path string, x int, y int, length int) {
	t.Helper()
	updateFrame(t, path, func(rgba *image.RGBA) {
		for px := x; px < x+length; px++ {
			rgba.SetRGBA(px, y, color.RGBA{R: 80, G: 70, B: 60, A: 180})
		}
	})
}

func updateFrame(t *testing.T, path string, draw func(*image.RGBA)) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	img, err := png.Decode(f)
	if closeErr := f.Close(); closeErr != nil {
		t.Fatalf("close: %v", closeErr)
	}
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	rgba := image.NewRGBA(img.Bounds())
	for py := img.Bounds().Min.Y; py < img.Bounds().Max.Y; py++ {
		for px := img.Bounds().Min.X; px < img.Bounds().Max.X; px++ {
			rgba.Set(px, py, img.At(px, py))
		}
	}
	draw(rgba)
	out, err := os.Create(path)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer out.Close()
	if err := png.Encode(out, rgba); err != nil {
		t.Fatalf("encode: %v", err)
	}
}

func minInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
