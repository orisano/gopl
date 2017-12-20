package main

import (
	"fmt"
	"math"
)

const (
	width, height = 600, 320
	cells         = 100
	xyrange       = 30.0
	xyscale       = width / 2 / xyrange
	zscale        = height * 0.4
	angle         = math.Pi / 6

	minZ = -0.22
	maxZ = 0.99
)

var sin30, cos30 = math.Sincos(angle)

func main() {
	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, az := corner(i+1, j)
			bx, by, bz := corner(i, j)
			cx, cy, cz := corner(i, j+1)
			dx, dy, dz := corner(i+1, j+1)
			ok := func() bool {
				for _, x := range []float64{ax, ay, az, bx, by, bz, cx, cy, cz, dx, dy, dz} {
					if math.IsInf(x, 0) || math.IsNaN(x) {
						return false
					}
				}
				return true
			}()
			if ok {
				z := (az + bz + cz + dz) / 4
				r := int((z - minZ) / (maxZ - minZ) * 255)
				b := 255 - r
				fmt.Printf("<polygon points='%g,%g %g,%g %g,%g %g,%g' fill='rgb(%d,0,%d)' />\n", ax, ay, bx, by, cx, cy, dx, dy, r, b)
			}
		}
	}
	fmt.Printf("</svg>")
}

func corner(i, j int) (float64, float64, float64) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	z := f(x, y)

	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, z
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y)
	return math.Sin(r) / r
}
