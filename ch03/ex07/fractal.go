package main

import (
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
)

const (
	EPS = 1e-5
)

func main() {
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 1024, 1024
	)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			z := complex(x, y)
			img.Set(px, py, newton(z))
		}
	}
	png.Encode(os.Stdout, img)
}

func newton(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	for n := uint8(0); n < iterations; n++ {
		z -= (z*z*z*z - 1) / (4 * z * z * z)
		switch {
		case cmplx.Abs(1-z) < EPS:
			return color.RGBA{255 - contrast*n, 0, 0, 0xff}
		case cmplx.Abs(-1-z) < EPS:
			return color.RGBA{0, 255 - contrast*n, 0, 0xff}
		case cmplx.Abs(1i-z) < EPS:
			return color.RGBA{0, 0, 255 - contrast*n, 0xff}
		case cmplx.Abs(-1i-z) < EPS:
			return color.RGBA{128, 0, 255 - contrast*n, 0xff}
		}
	}
	return color.Black
}
