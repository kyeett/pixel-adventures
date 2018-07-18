package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/kyeett/colorschemes"
)

var (
	backgroundColor = color.Black
	slownessSlice   = []float64{4, 3, 2, 3, 4, 8, 8}
)

const (
	height    = 500.0
	width     = 500.0
	frameTick = 10
	gridSize  = 40
	nLines    = 5
)

var cfg = Config{
	window: pixelgl.WindowConfig{
		Bounds: pixel.R(0, 0, height, width),
	},
	toggle: ToggleOptions{
		grid:        false,
		lineBorder:  true,
		rotatedView: false,
	},
	line: LineOptions{
		width:              12,
		borderWidth:        2,
		tailLength:         7 * nLines,
		tailShrinkStart:    0.5 * 8 * nLines,
		tailShrinkStepSize: 0.5,
		coordinateStepSize: 0.5,
	},
}

func run() {
	if nLines > 7 {
		panic("Currently max supported lines is 7")
	}

	lines := newLines(nLines)

	win, err := pixelgl.NewWindow(cfg.window)
	if err != nil {
		panic(fmt.Sprint("Window setup failed:", err))
	}
	win.SetPos(pixel.V(100, 100))

	linesOffset := cfg.window.Bounds.Center().Add(pixel.V(-gridSize*(nLines-0.5), -gridSize*(nLines-0.5)))
	gridOffset := cfg.window.Bounds.Center()

	// Initial background grid
	p := gridPicture()
	s := pixel.NewSprite(p, p.Bounds())

	tick := time.Tick(frameTick * time.Millisecond)
	var imdLines *imdraw.IMDraw
	paletteIndex := 0
	for !win.Closed() {
		imdLines = lines.prepareDraw(linesOffset)

		win.Clear(backgroundColor)
		if cfg.toggle.grid {
			s.Draw(win, pixel.IM.Moved(gridOffset))
		}
		imdLines.Draw(win)

		//Handle key presses
		if win.JustPressed(pixelgl.KeyUp) {
			paletteIndex = (paletteIndex - 1 + len(colorschemes.Palettes)) % len(colorschemes.Palettes)
			lines.updateColors(colorschemes.Palettes[paletteIndex])
		}
		if win.JustPressed(pixelgl.KeyDown) {
			paletteIndex = (paletteIndex + 1) % len(colorschemes.Palettes)
			lines.updateColors(colorschemes.Palettes[paletteIndex])
		}
		if win.JustPressed(pixelgl.KeyV) {
			cfg.toggle.rotatedView = !cfg.toggle.rotatedView
		}
		if win.JustPressed(pixelgl.KeyB) {
			cfg.toggle.lineBorder = !cfg.toggle.lineBorder
		}
		if win.JustPressed(pixelgl.KeyG) {
			cfg.toggle.grid = !cfg.toggle.grid
		}
		if win.JustPressed(pixelgl.KeyEscape) {
			return
		}

		win.Update()
		<-tick
	}
}

func main() {
	pixelgl.Run(run)
}

type Config struct {
	window pixelgl.WindowConfig
	toggle ToggleOptions
	line   LineOptions
}

type LineOptions struct {
	width              int
	borderWidth        float64
	tailLength         int
	tailShrinkStart    int
	tailShrinkStepSize float64
	coordinateStepSize float64
}

type ToggleOptions struct {
	lineBorder  bool
	grid        bool
	rotatedView bool
}
