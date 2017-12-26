package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

func main() {
	var values []int
	for _, x := range os.Args[1:] {
		if v, err := strconv.Atoi(x); err == nil {
			values = append(values, v)
		}
	}
	fmt.Println("max1:", max1(values...))
	if len(values) > 0 {
		fmt.Println("max2:", max2(values[0], values[1:]...))
	}

	fmt.Println("min1:", min1(values...))
	if len(values) > 0 {
		fmt.Println("min2:", min2(values[0], values[1:]...))
	}
}

func max1(values ...int) int {
	m := math.MinInt32
	for _, val := range values {
		if m < val {
			m = val
		}
	}
	return m
}

func max2(v int, values ...int) int {
	m := v
	for _, val := range values {
		if m < val {
			m = val
		}
	}
	return m
}

func min1(values ...int) int {
	m := math.MaxInt32
	for _, val := range values {
		if m > val {
			m = val
		}
	}
	return m
}

func min2(v int, values ...int) int {
	m := v
	for _, val := range values {
		if m > val {
			m = val
		}
	}
	return m
}
