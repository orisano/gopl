package main

import "unicode/utf8"

func reverseBytes(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

func rotateLeftBytes(b []byte, p int) {
	reverseBytes(b[:p])
	reverseBytes(b[p:])
	reverseBytes(b)
}

func ReverseUTF8(b []byte) {
	if len(b) == 0 {
		return
	}

	_, size := utf8.DecodeRune(b)
	rotateLeftBytes(b, size)
	ReverseUTF8(b[:len(b)-size])
}
