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
	artifactWarnings := flag.Bool("artifact-warnings", false, "warn about likely visual artifacts such as long low horizontal alpha runs")
	flag.Parse()

	if (*root == "") == (*framesDir == "") {
		fatalf("provide exactly one of -root or -frames-dir")
	}

	report, err := audit(*root, *framesDir, *pattern, *strict, *artifactWarnings)
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

func audit(root string, framesDir string, pattern string, strict bool, artifactWarnings bool) (auditReport, error) {
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
			setReport := auditSet(setName, setDir, pattern, artifactWarnings)
			addSet(&report, setReport)
		}
		return report, nil
	}
	setReport := auditSet(filepath.Base(framesDir), framesDir, pattern, artifactWarnings)
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

func auditSet(setName string, framesDir string, pattern string, artifactWarnings bool) setReport {
	report := setReport{
		Set:    setName,
		Dir:    filepath.ToSlash(framesDir),
		Frames: make([]frameReport, 0, totalFrames),
	}
	for frame := 0; frame < totalFrames; frame++ {
		framePath := filepath.Join(framesDir, fmt.Sprintf(pattern, frame))
		frameReport := auditFrame(frame, framePath, artifactWarnings)
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

func auditFrame(frame int, path string, artifactWarnings bool) frameReport {
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
	report.Warnings = frameWarnings(img, content, bounds, artifactWarnings)
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

func frameWarnings(img image.Image, content image.Rectangle, bounds image.Rectangle, artifactWarnings bool) []string {
	warnings := []string{}
	if content.Min.X <= bounds.Min.X || content.Max.X >= bounds.Max.X {
		warnings = append(warnings, "alpha touches horizontal canvas edge")
	}
	if content.Min.Y <= bounds.Min.Y || content.Max.Y >= bounds.Max.Y {
		warnings = append(warnings, "alpha touches vertical canvas edge")
	}
	if artifactWarnings {
		warnings = append(warnings, artifactWarningsForFrame(img, content)...)
	}
	return warnings
}

func artifactWarningsForFrame(img image.Image, content image.Rectangle) []string {
	if content.Empty() {
		return nil
	}
	warnings := []string{}
	maxRun := 0
	maxRow := -1
	startY := maxInt(content.Min.Y, content.Max.Y-3)
	for y := startY; y < content.Max.Y; y++ {
		run := 0
		for x := content.Min.X; x < content.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0x0800 {
				run++
				if run > maxRun {
					maxRun = run
					maxRow = y
				}
				continue
			}
			run = 0
		}
	}
	threshold := maxInt(18, content.Dx()/3)
	if maxRun >= threshold {
		warnings = append(warnings, fmt.Sprintf("possible floor/shadow artifact: alpha run length %d at y=%d near lower content", maxRun, maxRow))
	}
	if row, run, drop := lowerShelfArtifact(img, content); row >= 0 {
		warnings = append(warnings, fmt.Sprintf("possible lower ledge/shelf artifact: alpha run length %d at y=%d with %dpx row drop below", run, row, drop))
	}
	if componentCount, detachedCount, detachedArea, largestDetached := alphaComponentStats(img, content); detachedCount > 0 {
		warnings = append(warnings, fmt.Sprintf("disconnected alpha components: components=%d detached=%d detached_area=%d largest_detached=%d", componentCount, detachedCount, detachedArea, largestDetached))
	}
	if holeCount, holeArea, largestHole := transparentHoleStats(img, content); holeCount > 0 {
		warnings = append(warnings, fmt.Sprintf("transparent pinholes: holes=%d area=%d largest=%d", holeCount, holeArea, largestHole))
	}
	return warnings
}

func lowerShelfArtifact(img image.Image, content image.Rectangle) (row int, run int, drop int) {
	row = -1
	if content.Dx() < 24 || content.Dy() < 24 {
		return row, 0, 0
	}
	startY := content.Min.Y + (content.Dy()*2)/3
	endY := content.Max.Y - 4
	if endY <= startY {
		return row, 0, 0
	}
	runThreshold := maxInt(24, (content.Dx()*2)/3)
	dropThreshold := maxInt(18, content.Dx()/3)
	for y := startY; y < endY; y++ {
		count, maxRun := rowAlphaStats(img, content, y)
		if maxRun < runThreshold {
			continue
		}
		minBelow := count
		for belowY := y + 1; belowY < content.Max.Y && belowY <= y+4; belowY++ {
			belowCount, _ := rowAlphaStats(img, content, belowY)
			if belowCount < minBelow {
				minBelow = belowCount
			}
		}
		rowDrop := count - minBelow
		if rowDrop >= dropThreshold {
			return y, maxRun, rowDrop
		}
	}
	return row, 0, 0
}

func rowAlphaStats(img image.Image, content image.Rectangle, y int) (count int, maxRun int) {
	run := 0
	for x := content.Min.X; x < content.Max.X; x++ {
		if alphaVisible(img, x, y) {
			count++
			run++
			if run > maxRun {
				maxRun = run
			}
			continue
		}
		run = 0
	}
	return count, maxRun
}

func alphaComponentStats(img image.Image, content image.Rectangle) (componentCount int, detachedCount int, detachedArea int, largestDetached int) {
	width := content.Dx()
	height := content.Dy()
	if width <= 0 || height <= 0 {
		return 0, 0, 0, 0
	}
	visited := make([]bool, width*height)
	componentAreas := []int{}
	for y := content.Min.Y; y < content.Max.Y; y++ {
		for x := content.Min.X; x < content.Max.X; x++ {
			idx := (y-content.Min.Y)*width + (x - content.Min.X)
			if visited[idx] || !alphaVisible(img, x, y) {
				continue
			}
			area := floodAlphaComponent(img, content, x, y, visited)
			componentAreas = append(componentAreas, area)
		}
	}
	if len(componentAreas) <= 1 {
		return len(componentAreas), 0, 0, 0
	}
	mainAreaIndex := -1
	mainArea := 0
	for i, area := range componentAreas {
		if area > mainArea {
			mainAreaIndex = i
			mainArea = area
		}
	}
	for i, area := range componentAreas {
		if i == mainAreaIndex {
			continue
		}
		if area < 4 {
			continue
		}
		detachedCount++
		detachedArea += area
		if area > largestDetached {
			largestDetached = area
		}
	}
	return len(componentAreas), detachedCount, detachedArea, largestDetached
}

func floodAlphaComponent(img image.Image, content image.Rectangle, startX int, startY int, visited []bool) int {
	width := content.Dx()
	stack := []image.Point{{X: startX, Y: startY}}
	area := 0
	for len(stack) > 0 {
		point := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if !point.In(content) {
			continue
		}
		idx := (point.Y-content.Min.Y)*width + (point.X - content.Min.X)
		if visited[idx] || !alphaVisible(img, point.X, point.Y) {
			continue
		}
		visited[idx] = true
		area++
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 {
					continue
				}
				stack = append(stack, image.Point{X: point.X + dx, Y: point.Y + dy})
			}
		}
	}
	return area
}

func transparentHoleStats(img image.Image, content image.Rectangle) (holeCount int, holeArea int, largestHole int) {
	width := content.Dx()
	height := content.Dy()
	if width <= 0 || height <= 0 {
		return 0, 0, 0
	}
	visited := make([]bool, width*height)
	for y := content.Min.Y; y < content.Max.Y; y++ {
		for x := content.Min.X; x < content.Max.X; x++ {
			idx := (y-content.Min.Y)*width + (x - content.Min.X)
			if visited[idx] || alphaVisible(img, x, y) {
				continue
			}
			area, touchesBoundary := floodTransparentComponent(img, content, x, y, visited)
			if touchesBoundary {
				continue
			}
			holeCount++
			holeArea += area
			if area > largestHole {
				largestHole = area
			}
		}
	}
	return holeCount, holeArea, largestHole
}

func floodTransparentComponent(img image.Image, content image.Rectangle, startX int, startY int, visited []bool) (area int, touchesBoundary bool) {
	width := content.Dx()
	stack := []image.Point{{X: startX, Y: startY}}
	for len(stack) > 0 {
		point := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if !point.In(content) {
			continue
		}
		idx := (point.Y-content.Min.Y)*width + (point.X - content.Min.X)
		if visited[idx] || alphaVisible(img, point.X, point.Y) {
			continue
		}
		visited[idx] = true
		area++
		if point.X == content.Min.X || point.X == content.Max.X-1 || point.Y == content.Min.Y || point.Y == content.Max.Y-1 {
			touchesBoundary = true
		}
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 {
					continue
				}
				stack = append(stack, image.Point{X: point.X + dx, Y: point.Y + dy})
			}
		}
	}
	return area, touchesBoundary
}

func alphaVisible(img image.Image, x int, y int) bool {
	_, _, _, a := img.At(x, y).RGBA()
	return a > 0x0800
}

func rectToJSON(rect image.Rectangle) rectJSON {
	return rectJSON{X: rect.Min.X, Y: rect.Min.Y, W: rect.Dx(), H: rect.Dy()}
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
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
