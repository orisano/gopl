package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/orisano/gopl/ch02/ex02/conv"
)

var DefaultSource io.Reader = os.Stdin

func main() {
	values, err := getValues()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cv: %v\n", err.Error())
		return
	}
	for _, v := range values {
		c := conv.Celsius(v)
		f := conv.Fahrenheit(v)
		fmt.Printf("%s = %s, %s = %s\n", c, c.ToFahrenheit(), f, f.ToCelsius())

		ft := conv.Feet(v)
		m := conv.Meter(v)
		fmt.Printf("%s = %s, %s = %s\n", ft, ft.ToMeter(), m, m.ToFeet())

		p := conv.Pound(v)
		kg := conv.KiloGram(v)
		fmt.Printf("%s = %s, %s = %s\n", p, p.ToKiloGram(), kg, kg.ToPound())
	}
}

func getValues() ([]float64, error) {
	if args := os.Args[1:]; len(args) > 0 {
		return parseArgs(args)
	}

	var lines []string
	s := bufio.NewScanner(DefaultSource)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	return parseArgs(lines)
}

func parseArgs(args []string) ([]float64, error) {
	values := make([]float64, 0, len(args))
	for _, arg := range args {
		f, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			return nil, err
		}
		values = append(values, f)
	}
	return values, nil

}
