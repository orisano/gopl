package main

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"math/cmplx"
	"os"
)

const (
	xmin, ymin, xmax, ymax = -2, -2, +2, +2
	width, height          = 1024, 1024
)

func main() {
	run(os.Stdout, 4)
}

func naive(w io.Writer) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		for px := 0; px < width; px++ {
			y := float64(py)/height*(ymax-ymin) + ymin
			x := float64(px)/width*(xmax-xmin) + xmin
			z := complex(x, y)
			img.Set(px, py, mandelbrot(z))
		}
	}
	png.Encode(w, img)
}

func run(w io.Writer, parallels int) {
	type Pos struct {
		i, j int
		c    color.Color
	}

	pch := make(chan *Pos, parallels)
	mch := make(chan *Pos, parallels)
	for i := 0; i < parallels; i++ {
		go func() {
			for p := range pch {
				y := float64(p.i)/height*(ymax-ymin) + ymin
				x := float64(p.j)/width*(xmax-xmin) + xmin
				z := complex(x, y)
				mch <- &Pos{p.i, p.j, mandelbrot(z)}
			}
		}()
	}

	go func() {
		for py := 0; py < height; py++ {
			for px := 0; px < width; px++ {
				pch <- &Pos{i: py, j: px}
			}
		}
		close(pch)
	}()

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for i := 0; i < height*width; i++ {
		p := <-mch
		img.Set(p.j, p.i, p.c)
	}
	png.Encode(w, img)
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.YCbCr{255 - contrast*n, uint8(real(v)), uint8(imag(v))}
		}
	}
	return color.Black
}
