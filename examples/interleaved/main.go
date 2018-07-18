package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	height         = 500.0
	width          = 500.0
	gridSize       = 40
	nLines         = 6
	tailLength     = 5 * nLines
	lineWidth      = 12
	lineWidthDelay = 0.5 * 8 * nLines
	tickSize       = 1
	stepSize       = 0.5
	drawGrid       = false
	drawLineBorder = false
)

var (
	backgroundColor = color.Black
	slownessSlice   = []float64{4, 3, 2, 3, 4, 8}
	rotateView      = true
)

func run() {
	cfg := pixelgl.WindowConfig{
		Bounds: pixel.R(0, 0, height, width),
		//VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(fmt.Sprint("Window setup failed:", err))
	}
	win.SetPos(pixel.V(100, 100))

	linesOffset := cfg.Bounds.Center().Add(
		pixel.V(-gridSize*(nLines-0.5), -gridSize*(nLines-0.5)))
	centerOffset := cfg.Bounds.Center()

	p := drawPicture()
	s := pixel.NewSprite(p, p.Bounds())

	tick := time.Tick(tickSize * time.Millisecond)
	var imd *imdraw.IMDraw
	for !win.Closed() {
		imd = drawLines(linesOffset)
		win.Clear(backgroundColor)

		if drawGrid {
			s.Draw(win, pixel.IM.Moved(centerOffset))
		}

		if win.JustPressed(pixelgl.KeyV) {
			rotateView = !rotateView
		}
		if win.JustPressed(pixelgl.KeyEscape) {
			return
		}

		imd.Draw(win)

		win.Update()

		<-tick
	}
}

var lines []*line

func createLine(nLines, i float64) *line {
	colorPalete := []color.Color{
		colornames.Darkslateblue,
		colornames.Darkgoldenrod,
		colornames.Darkgoldenrod,
		colornames.Darkgoldenrod,
		colornames.Darkgoldenrod,
		colornames.Darkorchid,
		colornames.Darkgoldenrod,
		colornames.Darkgoldenrod,
		colornames.Darkgoldenrod,
		colornames.Darkkhaki,
		colornames.Darkcyan,
		colornames.Darkgoldenrod,
		colornames.Darkgray,
		colornames.Darkblue,
		colornames.Darkgreen,
		colornames.Darkmagenta,
		colornames.Darkred,
		colornames.Dodgerblue,
		colornames.Gainsboro,
		colornames.Honeydew,
	}
	maxX := 2*i + 1
	maxY := 2*nLines - (2*i + 1)
	start := pixel.V(float64(nLines-1-i), i)

	vs := []pixel.Vec{}

	for dx := float64(0); dx < maxX; dx += stepSize {
		vs = append(vs, start.Add(pixel.V(dx, 0)))
	}

	for dy := float64(0); dy < maxY; dy += stepSize {
		vs = append(vs, start.Add(pixel.V(maxX, dy)))
	}

	for dx := float64(maxX); dx > 0; dx -= stepSize {
		vs = append(vs, start.Add(pixel.V(dx, maxY)))
	}

	for dy := float64(maxY); dy > 0; dy -= stepSize {
		vs = append(vs, start.Add(pixel.V(0, dy)))
	}

	l := line{
		color:    colorPalete[int(i)],
		vs:       vs,
		zMask:    newZMask(nLines, i),
		slowness: slownessSlice[int(i)],
	}
	l.calcLayers()

	return &l
}

// Calculate this on line creation
func (l *line) calcLayers() {
	var bl [][]int
	var tl [][]int
	for i, z := range l.zMask {
		if z == 0 {
			bl = append(bl, []int{i, wrapIndex(i+1, len(l.vs))})
		} else {
			tl = append(tl, []int{i, wrapIndex(i+1, len(l.vs))})
		}
	}

	l.bottomLayer = bl
	l.topLayer = tl
}

// Calculate this on line creation
func (l *line) indexIsVisible(i int) bool {
	for c := 0; c < tailLength; c++ {
		if i == wrapIndex(l.head-c, len(l.vs)) {
			return true
		}
	}

	return false
}

func (l *line) indexLineWidth(i int) float64 {

	lw := float64(lineWidth)

	for c := 0; c < tailLength; c++ {
		if i == wrapIndex(l.head-c, len(l.vs)) {
			return lw
		}

		if c > lineWidthDelay && lw > 0 {
			lw -= 2
		}

		if lw < 0 {
			lw = 0
		}
	}

	return 0
}

// TODO: CLEAN UP CODE!
func newZMask(nLines, i float64) []byte {
	nSquares := 4 * int(nLines-1)
	nCorners := 4

	cornersFound := 0
	cornerPositions := []int64{
		int64(2 * i),
		int64(2 * (nLines - 1 - i)),
		int64(2 * i),
		int64(2 * (nLines - 1 - i)),
	}
	zMask := []byte{2}

	for s := 0; s < nSquares+nCorners-1; s += 2 {
		// Keep track of next corner
		if cornerPositions[cornersFound] == 0 {
			zMask = append(zMask, 2, 2)
			cornersFound++
			s--
			continue
		}

		// Add a pair of squares
		zMask = append(zMask, 1, 1)
		zMask = append(zMask, 0, 0)
		cornerPositions[cornersFound] -= 2
	}
	zMask = append(zMask, 2)
	return zMask
}

func createLines(nLines int) []*line {
	for i := float64(0); i < float64(nLines); i++ {
		lines = append(lines, createLine(float64(nLines), i))
	}
	return lines
}

func centerPoint(offset pixel.Vec) *imdraw.IMDraw {
	imd := imdraw.New(nil)
	imd.Color = colornames.Red
	imd.EndShape = imdraw.RoundEndShape
	imd.SetMatrix(pixel.IM.Moved(offset))
	imd.Push(pixel.V(0, 0))
	imd.Circle(2, 0)
	return imd
}

func wrapIndex(i int, length int) int {
	return (i + length) % (length)
}

func drawLines(offset pixel.Vec) *imdraw.IMDraw {
	center := pixel.V(width/2, height/2)
	imd := imdraw.New(nil)
	imd.EndShape = imdraw.RoundEndShape

	if rotateView {
		imd.SetMatrix(pixel.IM.Moved(offset).ScaledXY(center, pixel.V(1, -1)).Rotated(center, math.Pi/4))
	} else {
		imd.SetMatrix(pixel.IM.Moved(offset).ScaledXY(center, pixel.V(1, -1)))
	}

	for _, l := range lines {
		imd.Color = l.color
		for _, v := range l.bottomLayer {

			if l.indexIsVisible(v[0]) {
				imd.Push(l.vs[v[0]].Scaled(gridSize), l.vs[v[1]].Scaled(gridSize))
				imd.Line(l.indexLineWidth(v[0]))
			}
		}
	}

	// Top layer and head
	for _, l := range lines {

		// Top lines
		if drawLineBorder {
			imd.Color = backgroundColor
			for _, v := range l.topLayer {
				imd.Push(l.vs[v[0]].Scaled(gridSize), l.vs[v[1]].Scaled(gridSize))
				imd.Line(2.2 * lineWidth)
			}
		}

		imd.Color = l.color
		// Top lines
		for _, v := range l.topLayer {

			if l.indexIsVisible(v[0]) {
				imd.Push(l.vs[v[0]].Scaled(gridSize), l.vs[v[1]].Scaled(gridSize))
				imd.Line(l.indexLineWidth(v[0]))
			}
		}

		if l.count <= 0 {
			l.head = (l.head + 1) % (len(l.vs))
			l.count = l.slowness
		}
		l.count--
	}

	return imd
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func drawPicture() *pixel.PictureData {
	gridColors := []color.Color{colornames.Whitesmoke, colornames.Darkgray}
	offset := image.Pt(width/2-gridSize*(nLines-1), height/2-gridSize*(nLines-1))
	m := image.NewRGBA(image.Rect(0, 0, width, height))

	for gy := 0; gy < 2*(nLines-1); gy++ {

		//Used to cut corners in grid
		diff := (abs(2*(gy-nLines)+3) - 1) / 2

		for gx := diff; gx < 2*(nLines-1)-diff; gx++ {

			rect := image.Rect(gx*gridSize, gy*gridSize, (gx+1)*gridSize, (gy+1)*gridSize).Add(offset)

			// nLines used to make the pattern consisten betwen odd and even number of lines
			gridColor := gridColors[(gx+gy+nLines)%len(gridColors)]

			draw.Draw(m,
				rect,
				&image.Uniform{gridColor},
				image.ZP,
				draw.Src)
		}
	}

	return pixel.PictureDataFromImage(m)
}

type line struct {
	vs          []pixel.Vec
	zMask       []byte
	color       color.Color
	slowness    float64
	count       float64
	head        int
	bottomLayer [][]int
	topLayer    [][]int
}

func main() {
	lines = createLines(nLines)
	pixelgl.Run(run)
}
