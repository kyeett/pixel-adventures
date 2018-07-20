package main

import (
	"image/color"
	"math"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

const (
	maxSteps     = 50
	height       = 300
	width        = 300
	sleepTime    = 200
	triangleSize = height / 4
)

type state int

const (
	shrinkSquare = iota + 1
	spinTriangles
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

	//Initialize
	offset := pixel.V(width/2, height/2)
	midR := pixel.R(0, 0, width, height)
	midR = midR.Resized(offset, midR.Size().Scaled(0.5*1/math.Sqrt2)) //.Resized(offset, pixel.V(0.5, 0.5))
	fullR := pixel.R(0, 0, width, height)

	backgroundColor := color.RGBA{228, 0, 153, 255}
	foregroundColor := color.RGBA{237, 237, 237, 255}

	currentState := shrinkSquare
	var step float64
	for !win.Closed() {
		canvas.Clear(pixel.Alpha(0))

		// clear the canvas to be totally transparent and set the xor compose method
		//canvas.Clear(pixel.Alpha(0))
		//canvas.SetComposeMethod(pixel.ComposeXor)

		imd.Clear()

		if currentState == shrinkSquare {
			if step == maxSteps {
				step, currentState = 0, spinTriangles
				time.Sleep(sleepTime * time.Millisecond)
				continue
			}

			//Draw shrinking square
			imd.Color = foregroundColor
			imd.SetMatrix(pixel.IM)
			shrinkR := fullR.Resized(offset, fullR.Size().Scaled(1-smoothSwing(step/maxSteps)*0.5))
			imd.Push(shrinkR.Min, shrinkR.Max)
			imd.Rectangle(0)

			//Draw central square
			imd.Color = backgroundColor
			imd.SetMatrix(pixel.IM.Rotated(offset, math.Pi/4))
			imd.Push(midR.Min, midR.Max)
			imd.Rectangle(0)

			canvas.SetComposeMethod(pixel.ComposeOver)

		} else {
			// Spin Triangles
			if step == maxSteps {
				step, currentState = 0, shrinkSquare
				foregroundColor, backgroundColor = backgroundColor, foregroundColor
				time.Sleep(sleepTime * time.Millisecond)
				continue
			}

			for i := float64(0); i < 4; i++ {
				imd.Color = foregroundColor
				imd.SetMatrix(pixel.IM.Rotated(t.Center(), (1-smoothSwing(step/maxSteps))*math.Pi).Moved(offset).Rotated(offset, i*math.Pi/2))
				imd.Push(t.vertices()...)
				imd.Polygon(0)
			}

			canvas.SetComposeMethod(pixel.ComposeXor)
		}

		imd.Color = foregroundColor

		imd.Draw(canvas)
		win.Clear(backgroundColor)

		//Draw everything to the screen
		canvas.Draw(win, pixel.IM.Moved(offset))
		win.Update()
		step++
	}
}

func smoothSwing(x float64) float64 {
	return math.Pow(0.5-math.Cos(x*math.Pi)/2, 2)
}

func main() {
	pixelgl.Run(run)
}
