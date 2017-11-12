package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"os"
	"time"
)

var palette = []color.Color{
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0xff, 0x00, 0xff},
	color.RGBA{0x00, 0x88, 0x88, 0xff},
	color.RGBA{0x00, 0x44, 0xbb, 0xff},
	color.RGBA{0x00, 0x00, 0xff, 0xff},
	color.RGBA{0x44, 0x00, 0xbb, 0xff},
	color.RGBA{0x88, 0x00, 0x88, 0xff},
	color.RGBA{0xbb, 0x00, 0x44, 0xff},
	color.RGBA{0xff, 0x00, 0x00, 0xff},
	color.RGBA{0xbb, 0x44, 0x00, 0xff},
	color.RGBA{0x88, 0x88, 0x00, 0xff},
	color.RGBA{0x44, 0xbb, 0x00, 0xff},
}

const (
	backgroundColorIndex = 0
	lineColorIndex       = 1
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	lissajous(os.Stdout)
}

func lissajous(w io.Writer) {
	const (
		cycles  = 5
		res     = 0.001
		size    = 100
		nframes = 64
		delay   = 8
	)
	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0
	mod := len(palette) - 1
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Cos(t*freq + phase)
			c := (int(math.Abs(x+y)*10)+i)%mod + 1
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), uint8(c))
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(w, &anim)
}