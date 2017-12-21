package main

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math/cmplx"
	"net/http"
	"strconv"
)

const (
	EPS = 1e-5
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		cx, cy := 0.0, 0.0
		scale := 1.0
		if x := request.FormValue("x"); len(x) > 0 {
			if f, err := strconv.ParseFloat(x, 64); err == nil {
				cx = f
			} else {
				log.Println(err)
			}
		}
		if y := request.FormValue("y"); len(y) > 0 {
			if f, err := strconv.ParseFloat(y, 64); err == nil {
				cy = f
			} else {
				log.Println(err)
			}
		}
		if s := request.FormValue("scale"); len(s) > 0 {
			if f, err := strconv.ParseFloat(s, 64); err == nil {
				scale = f
			} else {
				log.Println(err)
			}
		}
		renderFractal(writer, cx, cy, scale)
	})
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func renderFractal(w io.Writer, cx, cy, scale float64) {
	L := 2 / scale
	var (
		xmin, ymin, xmax, ymax = -L, -L, +L, +L
	)
	const (
		width, height = 1024, 1024
	)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin + cy
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin + cx
			z := complex(x, y)
			img.Set(px, py, newton(z))
		}
	}
	png.Encode(w, img)
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
