package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	stddraw "image/draw"
	"image/png"
	"math"
	"os"
	"path/filepath"

	xdraw "golang.org/x/image/draw"
)

const (
	frameW    = 96
	frameH    = 64
	targetW   = 88
	targetH   = 52
	baselineY = 58
)

type prepareReport struct {
	Source            string   `json:"source"`
	Output            string   `json:"output"`
	SourceWidth       int      `json:"source_width"`
	SourceHeight      int      `json:"source_height"`
	OutputWidth       int      `json:"output_width"`
	OutputHeight      int      `json:"output_height"`
	BackgroundMode    string   `json:"background_mode"`
	BackgroundRemoved bool     `json:"background_removed"`
	Content           rectJSON `json:"content"`
	OutputContent     rectJSON `json:"output_content"`
	Warnings          []string `json:"warnings,omitempty"`
}

type rectJSON struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

func main() {
	srcPath := flag.String("src", "", "source candidate PNG")
	outPath := flag.String("out", "", "output 96x64 transparent PNG")
	reportPath := flag.String("report", "", "optional JSON report path")
	tolerance := flag.Float64("tolerance", 18, "RGB distance tolerance for uniform edge background removal")
	flag.Parse()

	if *srcPath == "" {
		fatalf("-src is required")
	}
	if *outPath == "" {
		fatalf("-out is required")
	}

	report, err := prepareFrame(*srcPath, *outPath, *tolerance)
	if err != nil {
		fatalf("%v", err)
	}
	if *reportPath != "" {
		if err := writeReport(*reportPath, report); err != nil {
			fatalf("write report: %v", err)
		}
	}
	fmt.Printf("prepared %s -> %s\n", filepath.ToSlash(*srcPath), filepath.ToSlash(*outPath))
}

func prepareFrame(srcPath string, outPath string, tolerance float64) (prepareReport, error) {
	src, err := openPNG(srcPath)
	if err != nil {
		return prepareReport{}, err
	}
	bounds := src.Bounds()
	report := prepareReport{
		Source:       filepath.ToSlash(srcPath),
		Output:       filepath.ToSlash(outPath),
		SourceWidth:  bounds.Dx(),
		SourceHeight: bounds.Dy(),
		OutputWidth:  frameW,
		OutputHeight: frameH,
	}

	cleaned := cloneRGBA(src)
	if hasTransparentAlpha(cleaned) {
		report.BackgroundMode = "source-alpha"
	} else {
		if err := removeUniformEdgeBackground(cleaned, tolerance); err != nil {
			return prepareReport{}, err
		}
		report.BackgroundMode = "uniform-edge-rgb"
		report.BackgroundRemoved = true
	}

	content := alphaBounds(cleaned, cleaned.Bounds())
	if content.Empty() {
		return prepareReport{}, fmt.Errorf("%s has no visible alpha after background preparation", srcPath)
	}
	if content == cleaned.Bounds() {
		return prepareReport{}, fmt.Errorf("%s still has no transparent background after preparation", srcPath)
	}
	report.Content = rectToJSON(content)
	report.Warnings = frameWarnings(content, cleaned.Bounds())
	if len(report.Warnings) > 0 {
		return prepareReport{}, fmt.Errorf("%s content touches source canvas edge after preparation; background removal or source crop is not clean", srcPath)
	}

	out := fitContent(cleaned, content)
	outContent := alphaBounds(out, out.Bounds())
	if outContent.Empty() {
		return prepareReport{}, fmt.Errorf("%s produced empty output", srcPath)
	}
	if outContent == out.Bounds() {
		return prepareReport{}, fmt.Errorf("%s produced output with no transparent margin", srcPath)
	}
	report.OutputContent = rectToJSON(outContent)
	report.Warnings = append(report.Warnings, frameWarnings(outContent, out.Bounds())...)

	if err := writePNG(outPath, out); err != nil {
		return prepareReport{}, err
	}
	return report, nil
}

func openPNG(path string) (*image.RGBA, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	src, err := png.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decode PNG %s: %w", path, err)
	}
	bounds := src.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	stddraw.Draw(out, out.Bounds(), src, bounds.Min, stddraw.Src)
	return out, nil
}

func cloneRGBA(src *image.RGBA) *image.RGBA {
	out := image.NewRGBA(src.Bounds())
	copy(out.Pix, src.Pix)
	return out
}

func hasTransparentAlpha(img *image.RGBA) bool {
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			if img.RGBAAt(x, y).A < 250 {
				return true
			}
		}
	}
	return false
}

func removeUniformEdgeBackground(img *image.RGBA, tolerance float64) error {
	bounds := img.Bounds()
	corners := []color.RGBA{
		img.RGBAAt(bounds.Min.X, bounds.Min.Y),
		img.RGBAAt(bounds.Max.X-1, bounds.Min.Y),
		img.RGBAAt(bounds.Min.X, bounds.Max.Y-1),
		img.RGBAAt(bounds.Max.X-1, bounds.Max.Y-1),
	}
	bg := averageColor(corners)
	for _, corner := range corners {
		if colorDistance(corner, bg) > tolerance {
			return fmt.Errorf("background corners differ; likely checker or noisy background")
		}
	}

	w, h := bounds.Dx(), bounds.Dy()
	seen := make([]bool, w*h)
	queue := make([]image.Point, 0, w*2+h*2)
	add := func(x, y int) {
		if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
			return
		}
		idx := (y-bounds.Min.Y)*w + (x - bounds.Min.X)
		if seen[idx] {
			return
		}
		seen[idx] = true
		queue = append(queue, image.Point{X: x, Y: y})
	}
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		add(x, bounds.Min.Y)
		add(x, bounds.Max.Y-1)
	}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		add(bounds.Min.X, y)
		add(bounds.Max.X-1, y)
	}

	removed := 0
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		c := img.RGBAAt(p.X, p.Y)
		if colorDistance(c, bg) > tolerance {
			continue
		}
		c.A = 0
		img.SetRGBA(p.X, p.Y, c)
		removed++
		add(p.X+1, p.Y)
		add(p.X-1, p.Y)
		add(p.X, p.Y+1)
		add(p.X, p.Y-1)
	}
	if removed == 0 {
		return fmt.Errorf("no edge background pixels matched")
	}
	return nil
}

func averageColor(colors []color.RGBA) color.RGBA {
	var r, g, b int
	for _, c := range colors {
		r += int(c.R)
		g += int(c.G)
		b += int(c.B)
	}
	n := len(colors)
	return color.RGBA{R: byte(r / n), G: byte(g / n), B: byte(b / n), A: 255}
}

func colorDistance(a color.RGBA, b color.RGBA) float64 {
	dr := float64(int(a.R) - int(b.R))
	dg := float64(int(a.G) - int(b.G))
	db := float64(int(a.B) - int(b.B))
	return math.Sqrt(dr*dr + dg*dg + db*db)
}

func fitContent(src *image.RGBA, content image.Rectangle) *image.RGBA {
	cropped := image.NewRGBA(image.Rect(0, 0, content.Dx(), content.Dy()))
	stddraw.Draw(cropped, cropped.Bounds(), src, content.Min, stddraw.Src)

	scale := math.Min(float64(targetW)/float64(content.Dx()), float64(targetH)/float64(content.Dy()))
	if scale > 1 {
		scale = 1
	}
	w := maxInt(1, int(math.Round(float64(content.Dx())*scale)))
	h := maxInt(1, int(math.Round(float64(content.Dy())*scale)))
	resized := image.NewRGBA(image.Rect(0, 0, w, h))
	xdraw.CatmullRom.Scale(resized, resized.Bounds(), cropped, cropped.Bounds(), xdraw.Over, nil)

	out := image.NewRGBA(image.Rect(0, 0, frameW, frameH))
	x := (frameW - w) / 2
	y := baselineY - h
	if y < 0 {
		y = 0
	}
	stddraw.Draw(out, image.Rect(x, y, x+w, y+h), resized, image.Point{}, stddraw.Over)
	return out
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

func writeReport(path string, report prepareReport) error {
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

func maxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
