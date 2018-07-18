package main

import (
	"image/color"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

func newLine(nLines, i float64) *line {
	colorPalete := defaultColorPalette
	maxX := 2*i + 1
	maxY := 2*nLines - (2*i + 1)
	start := pixel.V(float64(nLines-1-i), i)

	points := []pixel.Vec{}

	for dx := float64(0); dx < maxX; dx += cfg.line.coordinateStepSize {
		points = append(points, start.Add(pixel.V(dx, 0)))
	}

	for dy := float64(0); dy < maxY; dy += cfg.line.coordinateStepSize {
		points = append(points, start.Add(pixel.V(maxX, dy)))
	}

	for dx := float64(maxX); dx > 0; dx -= cfg.line.coordinateStepSize {
		points = append(points, start.Add(pixel.V(dx, maxY)))
	}

	for dy := float64(maxY); dy > 0; dy -= cfg.line.coordinateStepSize {
		points = append(points, start.Add(pixel.V(0, dy)))
	}

	l := line{
		color:    colorPalete[int(i)],
		points:   points,
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
			bl = append(bl, []int{i, wrapIndex(i+1, len(l.points))})
		} else {
			tl = append(tl, []int{i, wrapIndex(i+1, len(l.points))})
		}
	}

	l.bottomLayer = bl
	l.topLayer = tl
}

// Calculate this on line creation
func (l *line) indexIsVisible(i int) bool {
	for c := 0; c < cfg.line.tailLength; c++ {
		if i == wrapIndex(l.head-c, len(l.points)) {
			return true
		}
	}

	return false
}

func (l *line) indexLineWidth(i int) float64 {

	width := float64(cfg.line.width)

	for c := 0; c < cfg.line.tailLength; c++ {
		if i == wrapIndex(l.head-c, len(l.points)) {
			return width
		}

		if c > cfg.line.tailShrinkStart && width > 0 {
			width -= cfg.line.tailShrinkStepSize
		}

		if width < 0 {
			width = 0
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

func newLines(nLines int) lineSlice {
	lines := []*line{}
	for i := float64(0); i < float64(nLines); i++ {
		lines = append(lines, newLine(float64(nLines), i))
	}
	return lines
}

func wrapIndex(i int, length int) int {
	return (i + length) % (length)
}

type lineSlice []*line

func (lines lineSlice) prepareDraw(offset pixel.Vec) *imdraw.IMDraw {
	center := pixel.V(width/2, height/2)
	imd := imdraw.New(nil)
	imd.EndShape = imdraw.RoundEndShape

	m := pixel.IM.Moved(offset).ScaledXY(center, pixel.V(1, -1))
	if cfg.toggle.rotatedView {
		m = m.Rotated(center, math.Pi/4)
	}
	imd.SetMatrix(m)

	for _, l := range lines {
		imd.Color = l.color
		for _, v := range l.bottomLayer {

			if l.indexIsVisible(v[0]) {
				imd.Push(l.points[v[0]].Scaled(gridSize), l.points[v[1]].Scaled(gridSize))
				imd.Line(l.indexLineWidth(v[0]))
			}
		}
	}

	// Top layer and head
	for _, l := range lines {

		// Top lines
		for _, v := range l.topLayer {

			if l.indexIsVisible(v[0]) {
				lineWidth := l.indexLineWidth(v[0])

				// Don't draw border on corners
				if cfg.toggle.lineBorder {
					if l.zMask[v[0]] != 2 && lineWidth > 0 {

						imd.EndShape = imdraw.NoEndShape
						if v[0] == l.head {
							imd.EndShape = imdraw.RoundEndShape
						}
						imd.Color = backgroundColor
						imd.Push(l.points[v[0]].Scaled(gridSize), l.points[v[1]].Scaled(gridSize))
						imd.Line(lineWidth + cfg.line.borderWidth*2 + 1)

					}
				}

				imd.EndShape = imdraw.RoundEndShape

				imd.Color = l.color
				imd.Push(l.points[v[0]].Scaled(gridSize), l.points[v[1]].Scaled(gridSize))
				imd.Line(lineWidth)
				if v[0] == l.head && l.zMask[v[0]] != 2 {
					imd.Color = l.color
					imd.EndShape = imdraw.NoEndShape
					new := l.points[v[0]].Sub(l.points[v[1]]).Scaled(0.5).Add(l.points[v[0]]).Scaled(gridSize)
					imd.Push(new, l.points[v[1]].Scaled(gridSize))
					imd.Line(lineWidth)
				}

			}
		}

		if l.count <= 0 {
			l.head = (l.head + 1) % (len(l.points))
			l.count = l.slowness
		}
		l.count--
	}

	return imd
}

type line struct {
	points      []pixel.Vec
	zMask       []byte
	color       color.Color
	slowness    float64
	count       float64
	head        int
	bottomLayer [][]int
	topLayer    [][]int
}

var defaultColorPalette = []color.Color{
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
