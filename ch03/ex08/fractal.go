package main

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"math/big"
	"math/cmplx"
	"os"
)

const (
	EPS        = 1e-5
	iterations = 4
	contrast   = 15
)

type Fractal interface {
	Get(x, y float64) color.Color
}

type FractalFunc func(x, y float64) color.Color

func (ff FractalFunc) Get(x, y float64) color.Color {
	return ff(x, y)
}

func complex64Fractal(x, y float64) color.Color {
	z := complex64(complex(x, y))

	for n := uint8(0); n < iterations; n++ {
		z -= (z - 1/(z*z*z)) / 4
		switch {
		case cmplx.Abs(complex128(1-z)) < EPS:
			return color.RGBA{255 - contrast*n, 0, 0, 0xff}
		case cmplx.Abs(complex128(-1-z)) < EPS:
			return color.RGBA{0, 255 - contrast*n, 0, 0xff}
		case cmplx.Abs(complex128(1i-z)) < EPS:
			return color.RGBA{0, 0, 255 - contrast*n, 0xff}
		case cmplx.Abs(complex128(-1i-z)) < EPS:
			return color.RGBA{128, 0, 255 - contrast*n, 0xff}
		}
	}
	return color.Black
}

func complex128Fractal(x, y float64) color.Color {

	z := complex(x, y)

	for n := uint8(0); n < iterations; n++ {
		z -= (z - 1/(z*z*z)) / 4
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

type bigFloatComplex struct {
	Real, Imag *big.Float
}

func (bf *bigFloatComplex) Mul(rhs *bigFloatComplex) *bigFloatComplex {
	return &bigFloatComplex{
		new(big.Float).Sub(new(big.Float).Mul(bf.Real, rhs.Real), new(big.Float).Mul(bf.Imag, rhs.Imag)),
		new(big.Float).Add(new(big.Float).Mul(bf.Imag, rhs.Real), new(big.Float).Mul(bf.Real, rhs.Imag)),
	}
}

func (bf *bigFloatComplex) Sub(rhs *bigFloatComplex) *bigFloatComplex {
	return &bigFloatComplex{
		new(big.Float).Sub(bf.Real, rhs.Real),
		new(big.Float).Sub(bf.Imag, rhs.Imag),
	}
}

func (bf *bigFloatComplex) Div(rhs *bigFloatComplex) *bigFloatComplex {
	d := new(big.Float).Add(
		new(big.Float).Mul(rhs.Real, rhs.Real),
		new(big.Float).Mul(rhs.Imag, rhs.Imag),
	)
	return &bigFloatComplex{
		new(big.Float).Quo(
			new(big.Float).Add(
				new(big.Float).Mul(bf.Real, rhs.Real),
				new(big.Float).Mul(bf.Imag, rhs.Imag),
			),
			d,
		),
		new(big.Float).Quo(
			new(big.Float).Sub(
				new(big.Float).Mul(bf.Imag, rhs.Real),
				new(big.Float).Mul(bf.Real, rhs.Imag),
			),
			d,
		),
	}
}

func (bf *bigFloatComplex) IsZero() bool {
	return new(big.Float).Add(
		new(big.Float).Mul(bf.Real, bf.Real),
		new(big.Float).Mul(bf.Imag, bf.Imag),
	).Cmp(big.NewFloat(EPS)) < 0
}

func bigFloatFractal(x, y float64) color.Color {
	if x == 0 && y == 0 {
		return color.Black
	}
	z := &bigFloatComplex{
		big.NewFloat(x),
		big.NewFloat(y),
	}

	one1 := &bigFloatComplex{
		big.NewFloat(1),
		big.NewFloat(0),
	}
	one2 := &bigFloatComplex{
		big.NewFloat(-1),
		big.NewFloat(0),
	}
	one3 := &bigFloatComplex{
		big.NewFloat(0),
		big.NewFloat(1),
	}
	one4 := &bigFloatComplex{
		big.NewFloat(0),
		big.NewFloat(-1),
	}
	four := &bigFloatComplex{
		big.NewFloat(4),
		big.NewFloat(0),
	}

	for n := uint8(0); n < iterations; n++ {
		z = z.Sub(z.Sub(one1.Div(z.Mul(z).Mul(z))).Div(four))
		switch {
		case one1.Sub(z).IsZero():
			return color.RGBA{255 - contrast*n, 0, 0, 0xff}
		case one2.Sub(z).IsZero():
			return color.RGBA{0, 255 - contrast*n, 0, 0xff}
		case one3.Sub(z).IsZero():
			return color.RGBA{0, 0, 255 - contrast*n, 0xff}
		case one4.Sub(z).IsZero():
			return color.RGBA{128, 0, 255 - contrast*n, 0xff}
		}
	}
	return color.Black
}

type bigRatComplex struct {
	Real, Imag *big.Rat
}

func (bf *bigRatComplex) Mul(rhs *bigRatComplex) *bigRatComplex {
	return &bigRatComplex{
		new(big.Rat).Sub(new(big.Rat).Mul(bf.Real, rhs.Real), new(big.Rat).Mul(bf.Imag, rhs.Imag)),
		new(big.Rat).Add(new(big.Rat).Mul(bf.Imag, rhs.Real), new(big.Rat).Mul(bf.Real, rhs.Imag)),
	}
}

func (bf *bigRatComplex) Sub(rhs *bigRatComplex) *bigRatComplex {
	return &bigRatComplex{
		new(big.Rat).Sub(bf.Real, rhs.Real),
		new(big.Rat).Sub(bf.Imag, rhs.Imag),
	}
}

func (bf *bigRatComplex) Div(rhs *bigRatComplex) *bigRatComplex {
	d := new(big.Rat).Add(
		new(big.Rat).Mul(rhs.Real, rhs.Real),
		new(big.Rat).Mul(rhs.Imag, rhs.Imag),
	)
	return &bigRatComplex{
		new(big.Rat).Quo(
			new(big.Rat).Add(
				new(big.Rat).Mul(bf.Real, rhs.Real),
				new(big.Rat).Mul(bf.Imag, rhs.Imag),
			),
			d,
		),
		new(big.Rat).Quo(
			new(big.Rat).Sub(
				new(big.Rat).Mul(bf.Imag, rhs.Real),
				new(big.Rat).Mul(bf.Real, rhs.Imag),
			),
			d,
		),
	}
}

func (bf *bigRatComplex) IsZero() bool {
	return new(big.Rat).Add(
		new(big.Rat).Mul(bf.Real, bf.Real),
		new(big.Rat).Mul(bf.Imag, bf.Imag),
	).Cmp(big.NewRat(1, 1/EPS*EPS)) < 0
}

func bigRatFractal(x, y float64) color.Color {
	if x == 0 && y == 0 {
		return color.Black
	}
	z := &bigRatComplex{
		new(big.Rat).SetFloat64(x),
		new(big.Rat).SetFloat64(y),
	}

	one1 := &bigRatComplex{
		big.NewRat(1, 1),
		big.NewRat(0, 1),
	}
	one2 := &bigRatComplex{
		big.NewRat(-1, 1),
		big.NewRat(0, 1),
	}
	one3 := &bigRatComplex{
		big.NewRat(0, 1),
		big.NewRat(1, 1),
	}
	one4 := &bigRatComplex{
		big.NewRat(0, 1),
		big.NewRat(-1, 1),
	}
	four := &bigRatComplex{
		big.NewRat(4, 1),
		big.NewRat(0, 1),
	}

	for n := uint8(0); n < iterations; n++ {
		z = z.Sub(z.Sub(one1.Div(z.Mul(z).Mul(z))).Div(four))
		switch {
		case one1.Sub(z).IsZero():
			return color.RGBA{255 - contrast*n, 0, 0, 0xff}
		case one2.Sub(z).IsZero():
			return color.RGBA{0, 255 - contrast*n, 0, 0xff}
		case one3.Sub(z).IsZero():
			return color.RGBA{0, 0, 255 - contrast*n, 0xff}
		case one4.Sub(z).IsZero():
			return color.RGBA{128, 0, 255 - contrast*n, 0xff}
		}
	}
	return color.Black
}

func run(w io.Writer, fractal Fractal) {
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 256, 256
	)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			img.Set(px, py, fractal.Get(x, y))
		}
	}
	png.Encode(w, img)
}

func main() {
	run(os.Stdout, FractalFunc(complex128Fractal))
}
