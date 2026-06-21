package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
)

const (
	frameW      = 96
	frameH      = 64
	totalFrames = 62
	motionSets  = 10
)

type auditReport struct {
	Root        string      `json:"root,omitempty"`
	FramesDir   string      `json:"frames_dir,omitempty"`
	Pattern     string      `json:"pattern"`
	Strict      bool        `json:"strict"`
	FrameWidth  int         `json:"frame_width"`
	FrameHeight int         `json:"frame_height"`
	SetCount    int         `json:"set_count"`
	FrameCount  int         `json:"frame_count"`
	Valid       int         `json:"valid"`
	Missing     int         `json:"missing"`
	Invalid     int         `json:"invalid"`
	Warnings    int         `json:"warnings"`
	Sets        []setReport `json:"sets"`
}

type setReport struct {
	Set       string        `json:"set"`
	Dir       string        `json:"dir"`
	Valid     int           `json:"valid"`
	Missing   int           `json:"missing"`
	Invalid   int           `json:"invalid"`
	Warnings  int           `json:"warnings"`
	Frames    []frameReport `json:"frames"`
	Completed bool          `json:"completed"`
}

type frameReport struct {
	Frame    int      `json:"frame"`
	Path     string   `json:"path"`
	Status   string   `json:"status"`
	Content  rectJSON `json:"content,omitempty"`
	Error    string   `json:"error,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

type rectJSON struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

func main() {
	root := flag.String("root", "", "optional root containing set00 through set09 frame directories")
	framesDir := flag.String("frames-dir", "", "single set directory containing frame PNGs")
	pattern := flag.String("pattern", "frame-%02d.png", "frame filename pattern with one integer verb")
	reportPath := flag.String("report", "", "optional JSON report path")
	strict := flag.Bool("strict", false, "exit non-zero unless every expected frame is valid")
	flag.Parse()

	if (*root == "") == (*framesDir == "") {
		fatalf("provide exactly one of -root or -frames-dir")
	}

	report, err := audit(*root, *framesDir, *pattern, *strict)
	if err != nil {
		fatalf("%v", err)
	}
	if *reportPath != "" {
		if err := writeReport(*reportPath, report); err != nil {
			fatalf("write report: %v", err)
		}
	}
	fmt.Printf("audited %d frame slots: valid=%d missing=%d invalid=%d warnings=%d\n", report.FrameCount, report.Valid, report.Missing, report.Invalid, report.Warnings)
	if *strict && (report.Missing > 0 || report.Invalid > 0) {
		os.Exit(1)
	}
}

func audit(root string, framesDir string, pattern string, strict bool) (auditReport, error) {
	report := auditReport{
		Root:        filepath.ToSlash(root),
		FramesDir:   filepath.ToSlash(framesDir),
		Pattern:     pattern,
		Strict:      strict,
		FrameWidth:  frameW,
		FrameHeight: frameH,
		Sets:        []setReport{},
	}
	if root != "" {
		for set := 0; set < motionSets; set++ {
			setName := fmt.Sprintf("set%02d", set)
			setDir := filepath.Join(root, setName)
			setReport := auditSet(setName, setDir, pattern)
			addSet(&report, setReport)
		}
		return report, nil
	}
	setReport := auditSet(filepath.Base(framesDir), framesDir, pattern)
	addSet(&report, setReport)
	return report, nil
}

func addSet(report *auditReport, set setReport) {
	report.SetCount++
	report.FrameCount += totalFrames
	report.Valid += set.Valid
	report.Missing += set.Missing
	report.Invalid += set.Invalid
	report.Warnings += set.Warnings
	report.Sets = append(report.Sets, set)
}

func auditSet(setName string, framesDir string, pattern string) setReport {
	report := setReport{
		Set:    setName,
		Dir:    filepath.ToSlash(framesDir),
		Frames: make([]frameReport, 0, totalFrames),
	}
	for frame := 0; frame < totalFrames; frame++ {
		framePath := filepath.Join(framesDir, fmt.Sprintf(pattern, frame))
		frameReport := auditFrame(frame, framePath)
		switch frameReport.Status {
		case "valid":
			report.Valid++
		case "missing":
			report.Missing++
		default:
			report.Invalid++
		}
		report.Warnings += len(frameReport.Warnings)
		report.Frames = append(report.Frames, frameReport)
	}
	report.Completed = report.Valid == totalFrames && report.Missing == 0 && report.Invalid == 0
	return report
}

func auditFrame(frame int, path string) frameReport {
	report := frameReport{
		Frame:  frame,
		Path:   filepath.ToSlash(path),
		Status: "valid",
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			report.Status = "missing"
			report.Error = "file does not exist"
			return report
		}
		report.Status = "invalid"
		report.Error = err.Error()
		return report
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		report.Status = "invalid"
		report.Error = fmt.Sprintf("decode PNG: %v", err)
		return report
	}
	bounds := img.Bounds()
	if bounds.Dx() != frameW || bounds.Dy() != frameH {
		report.Status = "invalid"
		report.Error = fmt.Sprintf("bounds = %dx%d, want %dx%d", bounds.Dx(), bounds.Dy(), frameW, frameH)
		return report
	}
	content := alphaBounds(img, bounds)
	if content.Empty() {
		report.Status = "invalid"
		report.Error = "no visible alpha"
		return report
	}
	if content == bounds {
		report.Status = "invalid"
		report.Error = "no transparent background"
		return report
	}
	report.Content = rectToJSON(content)
	report.Warnings = frameWarnings(content, bounds)
	return report
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

func frameWarnings(content image.Rectangle, bounds image.Rectangle) []string {
	warnings := []string{}
	if content.Min.X <= bounds.Min.X || content.Max.X >= bounds.Max.X {
		warnings = append(warnings, "alpha touches horizontal canvas edge")
	}
	if content.Min.Y <= bounds.Min.Y || content.Max.Y >= bounds.Max.Y {
		warnings = append(warnings, "alpha touches vertical canvas edge")
	}
	return warnings
}

func rectToJSON(rect image.Rectangle) rectJSON {
	return rectJSON{X: rect.Min.X, Y: rect.Min.Y, W: rect.Dx(), H: rect.Dy()}
}

func writeReport(path string, report auditReport) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
