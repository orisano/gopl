package popcount

import "sync"

var once sync.Once
var pc [256]byte

func precalc() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

func PopCount(x uint64) int {
	once.Do(precalc)
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}

func PopCountLoop(x uint64) int {
	once.Do(precalc)
	var r int
	for i := uint(0); i < 8; i++ {
		r += int(pc[byte(x>>(i*8))])
	}
	return r
}
