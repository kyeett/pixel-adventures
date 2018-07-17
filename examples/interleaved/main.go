package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
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
	nLines         = 4
	tailLength     = 16
	lineWidth      = 10
	tickSize       = 500000
	stepSize       = 0.5
	drawGrid       = false
	drawLineBorder = false
)

var (
	backgroundColor = color.Black
)

func run() {
	cfg := pixelgl.WindowConfig{
		Bounds: pixel.R(0, 0, height, width),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(fmt.Sprint("Window setup failed:", err))
	}
	win.SetPos(pixel.V(100, 100))

	linesOffset := cfg.Bounds.Center().Add(
		pixel.V(-gridSize*(nLines-0.5), -gridSize*(nLines-0.5))) // -gridSize*nLines/2))
	centerOffset := cfg.Bounds.Center()

	p := drawPicture()
	s := pixel.NewSprite(p, p.Bounds())

	tick := time.Tick(tickSize * time.Millisecond)
	for !win.Closed() {
		imd := drawLines(linesOffset)
		win.Clear(backgroundColor)

		if drawGrid {
			s.Draw(win, pixel.IM.Moved(centerOffset))
		}

		imd.Draw(win)
		//centerPoint(centerOffset).Draw(win)

		win.Update()

		<-tick
	}
}

var lines []*line

func createLine(nLines, i float64) *line {
	colorPalete := []color.Color{
		colornames.Darkorchid,
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
		color: colorPalete[int(i)],
		vs:    vs,
		zMask: newZMask(nLines, i),
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
	imd.SetMatrix(pixel.IM.Moved(offset).ScaledXY(center, pixel.V(1, -1))) //.Rotated(center, math.Pi/4))

	/*var visibleLines []pixel.Vec
	for c := 0; c < tailLength; c++ {
		wrappedIndex := wrapIndex(l.head-c, len(l.vs))
		l.vs[l.zMask[0]]
		visibleLines = append(visibleLines, l.vs[wrappedIndex].Scaled(gridSize))
		}*/
	//imd.Push(visibleLines...)

	/*		if drawLineBorder {
			imd.Color = color.White
			//	imd.Push(visibleLines...)
			imd.Line(2 * lineWidth)
		}*/

	for _, l := range lines {
		imd.Color = l.color
		for _, v := range l.bottomLayer {
			imd.Push(l.vs[v[0]].Scaled(gridSize), l.vs[v[1]].Scaled(gridSize))
			imd.Line(lineWidth)
		}
	}
	for _, l := range lines {
		imd.Color = l.color
		for _, v := range l.topLayer {
			imd.Push(l.vs[v[0]].Scaled(gridSize), l.vs[v[1]].Scaled(gridSize))
			imd.Line(lineWidth)
		}
	}
	/*
		imd.Push(l.vs[l.head].Scaled(gridSize))
		imd.Color = colornames.Peachpuff
		imd.Circle(10, 10)

		l.head = (l.head + 1) % (len(l.vs) - 1)
	}*/

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
	head        int
	bottomLayer [][]int
	topLayer    [][]int
}

func main() {
	lines = createLines(nLines)
	pixelgl.Run(run)
}
