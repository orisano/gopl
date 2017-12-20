package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
)

const (
	cells   = 100
	xyrange = 30.0
	angle   = math.Pi / 6

	minZ = -0.22
	maxZ = 0.99
)

var sin30, cos30 = math.Sincos(angle)

type SVGOption struct {
	Height int
	Width  int
	Color  string
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		svgOption := &SVGOption{
			Width:  620,
			Height: 300,
		}
		if width := r.FormValue("w"); len(width) > 0 {
			if x, err := strconv.Atoi(width); err == nil {
				svgOption.Width = x
			} else {
				log.Print(err)
			}
		}
		if height := r.FormValue("h"); len(height) > 0 {
			if x, err := strconv.Atoi(height); err == nil {
				svgOption.Height = x
			} else {
				log.Print(err)
			}
		}
		svgOption.Color = r.FormValue("c")

		w.Header().Set("Content-Type", "image/svg+xml")
		writeSVG(w, svgOption)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func writeSVG(w io.Writer, option *SVGOption) {
	width, height := option.Width, option.Height
	fmt.Fprintf(w, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, az := corner(i+1, j, width, height)
			bx, by, bz := corner(i, j, width, height)
			cx, cy, cz := corner(i, j+1, width, height)
			dx, dy, dz := corner(i+1, j+1, width, height)
			ok := func() bool {
				for _, x := range []float64{ax, ay, az, bx, by, bz, cx, cy, cz, dx, dy, dz} {
					if math.IsInf(x, 0) || math.IsNaN(x) {
						return false
					}
				}
				return true
			}()
			if ok {
				color := option.Color
				if len(color) == 0 {
					z := (az + bz + cz + dz) / 4
					r := int((z - minZ) / (maxZ - minZ) * 255)
					b := 255 - r
					color = fmt.Sprintf("rgb(%d,0,%d)", r, b)
				}
				fmt.Fprintf(w, "<polygon points='%g,%g %g,%g %g,%g %g,%g' fill='%s' />\n", ax, ay, bx, by, cx, cy, dx, dy, color)
			}
		}
	}
	fmt.Fprintln(w, "</svg>")
}

func corner(i, j, w, h int) (float64, float64, float64) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	z := f(x, y)

	xyscale := float64(w) / 2 / xyrange
	zscale := float64(h) * 0.4

	sx := float64(w)/2 + (x-y)*cos30*xyscale
	sy := float64(h)/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, z
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y)
	return math.Sin(r) / r
}
