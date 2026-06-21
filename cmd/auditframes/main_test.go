package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestAuditReportsPartialFrameSet(t *testing.T) {
	root := t.TempDir()
	framesDir := filepath.Join(root, "set00")
	writeFrame(t, filepath.Join(framesDir, "frame-00.png"), frameW, frameH, false)
	writeFrame(t, filepath.Join(framesDir, "frame-01.png"), frameW+1, frameH, false)

	report, err := audit("", framesDir, "frame-%02d.png", false)
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

	report, err := audit("", framesDir, "frame-%02d.png", false)
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

	report, err := audit(root, "", "frame-%02d.png", false)
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

func minInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
