package main

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
)

func gridPicture() *pixel.PictureData {
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

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
