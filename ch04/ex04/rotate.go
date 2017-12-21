package main

func gcd(a, b int) int {
	for b > 0 {
		a, b = b, a%b
	}
	return a
}

func Rotate(a []int, r int) {
	L := len(a)
	g := gcd(L, r)
	for i := 0; i < g; i++ {
		x := a[i]
		for j := (i + r) % L; j != i; j = (j + r) % L {
			x, a[j] = a[j], x
		}
		a[i] = x
	}
}
