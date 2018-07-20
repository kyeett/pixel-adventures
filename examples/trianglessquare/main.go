package main

import (
	"image/color"
	"math"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

const (
	height       = 200
	width        = 200
	triangleSize = 50
)

type triangle struct {
	pixel.Rect
	color color.Color
}

func (t triangle) vertices() []pixel.Vec {
	return []pixel.Vec{
		pixel.V(0, 0),
		pixel.V(t.W(), 0),
		pixel.V(0, t.H()),
	}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Bounds: pixel.R(0, 0, height, width),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetPos(pixel.V(0, 0))

	imd := imdraw.New(nil)

	canvas := pixelgl.NewCanvas(win.Bounds())

	//	ts := NewTriangles()
	t := triangle{
		pixel.R(0, 0, triangleSize, triangleSize),
		colornames.Red,
	}

	var step float64
	//start := time.Now()
	for !win.Closed() {
		step++
		// in case window got resized, we also need to resize our canvas
		canvas.SetBounds(win.Bounds())
		canvas.Clear(pixel.Alpha(0))

		//offset := math.Sin(time.Since(start).Seconds()) * 300

		// clear the canvas to be totally transparent and set the xor compose method
		//canvas.Clear(pixel.Alpha(0))
		//canvas.SetComposeMethod(pixel.ComposeXor)

		offset := pixel.V(width/2, height/2)
		/* 		imd.SetMatrix(pixel.IM.Rotated(t.Center(), step/30*math.Pi).Moved(offset)) */
		imd.Clear()
		imd.Color = t.color
		for i := float64(0); i < 4; i++ {

			imd.SetMatrix(pixel.IM.Moved(offset).Rotated(offset, i*math.Pi/2))
			imd.Push(t.vertices()...)
			imd.Polygon(0)
		}
		imd.Draw(canvas)

		//Draw everything to the screen
		win.Clear(color.Black)
		canvas.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		win.Update()

		/*
			// red circle
			imd.Clear()
			imd.Color = pixel.RGB(1, 0, 0)
			imd.Push(win.Bounds().Center().Add(pixel.V(-offset, 0)))
			imd.Circle(200, 0)
			imd.Draw(canvas)

			// blue circle
			imd.Clear()
			imd.Color = pixel.RGB(0, 0, 1)
			imd.Push(win.Bounds().Center().Add(pixel.V(offset, 0)))
			imd.Circle(150, 0)
			imd.Draw(canvas)

			// yellow circle
			imd.Clear()
			imd.Color = pixel.RGB(1, 1, 0)
			imd.Push(win.Bounds().Center().Add(pixel.V(0, -offset)))
			imd.Circle(100, 0)
			imd.Draw(canvas)

			// magenta circle
			imd.Clear()
			imd.Color = pixel.RGB(1, 0, 1)
			imd.Push(win.Bounds().Center().Add(pixel.V(0, offset)))
			imd.Circle(50, 0)
			imd.Draw(canvas)

			win.Clear(colornames.Green)
			canvas.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
			win.Update()
		*/
	}
}

func main() {
	pixelgl.Run(run)
}
