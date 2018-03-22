package popcount

var pc [256]byte
var pc15 [32 * 1024]byte
var pc16 [64 * 1024]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
	for i := range pc15 {
		pc15[i] = pc15[i/2] + byte(i&1)
	}
	for i := range pc16 {
		pc16[i] = pc16[i/2] + byte(i&1)
	}
}

func PopCount(x uint64) int {
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}

func PopCountNaive(x uint64) int {
	r := uint64(0)
	for i := uint(0); i < 64; i++ {
		r += (x >> i) & 1
	}
	return int(r)
}

func PopCountBitMagic(x uint64) int {
	cnt := 0
	for x > 0 {
		x &= x - 1
		cnt++
	}
	return cnt
}

func PopCountParallel(x uint64) int {
	x = (x&0xAAAAAAAAAAAAAAAA)>>1 + (x & 0x5555555555555555)
	x = (x&0xCCCCCCCCCCCCCCCC)>>2 + (x & 0x3333333333333333)
	x = (x&0xF0F0F0F0F0F0F0F0)>>4 + (x & 0x0F0F0F0F0F0F0F0F)
	x = (x&0xFF00FF00FF00FF00)>>8 + (x & 0x00FF00FF00FF00FF)
	x = (x&0xFFFF0000FFFF0000)>>16 + (x & 0x0000FFFF0000FFFF)
	return int((x&0xFFFFFFFF00000000)>>32 + (x & 0x00000000FFFFFFFF))
}

func PopCount16(x uint64) int {
	return int(pc16[(x>>(0*16))&0xffff] +
		pc16[(x>>(1*16))&0xffff] +
		pc16[(x>>(2*16))&0xffff] +
		pc16[(x>>(3*16))&0xffff])
}

func PopCount15(x uint64) int {
	return int(pc15[(x>>(0*16+1))&0x7fff] + byte((x>>(0*16))&1) +
		pc15[(x>>(1*16+1))&0x7fff] + byte((x>>(1*16))&1) +
		pc15[(x>>(2*16+1))&0x7fff] + byte((x>>(2*16))&1) +
		pc15[(x>>(3*16+1))&0x7fff] + byte((x>>(3*16))&1))
}

func PopCountCPU(x uint64) int
