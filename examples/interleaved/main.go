package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	height   = 500.0
	width    = 500.0
	nLines   = 4
	gridSize = 20
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

	offset := cfg.Bounds.Center()
	for !win.Closed() {
		win.Clear(colornames.Black)

		p := drawPicture()
		s := pixel.NewSprite(p, p.Bounds())
		s.Draw(win, pixel.IM.Moved(offset))

		win.Update()
	}
}

func drawPicture() *pixel.PictureData {
	gridColors := []color.Color{colornames.Whitesmoke, colornames.Darkgray}
	m := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(m, m.Bounds(), &image.Uniform{colornames.Darkblue}, image.ZP, draw.Src)

	for gx := 0; gx < nLines; gx++ {
		for gy := 0; gy < nLines; gy++ {
			rect := image.Rect(gx*gridSize, gy*gridSize, (gx+1)*gridSize, (gy+1)*gridSize).Add(image.Pt((width-gridSize*nLines)/2, (height-gridSize*nLines)/2))
			draw.Draw(m,
				rect,
				&image.Uniform{gridColors[(gx+gy)%len(gridColors)]},
				image.ZP,
				draw.Src)
		}
	}

	return pixel.PictureDataFromImage(m)
}

func main() {
	pixelgl.Run(run)
}
