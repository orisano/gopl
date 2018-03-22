package charcount

import (
	"bufio"
	"io"
	"unicode"
	"unicode/utf8"
)

func CharCount(r io.Reader) (map[rune]int, []int, int, error) {
	counts := make(map[rune]int)
	var utflen [utf8.UTFMax + 1]int
	invalid := 0

	in := bufio.NewReader(r)
	for {
		r, n, err := in.ReadRune()
		if err == io.EOF {
			return counts, utflen[:], invalid, nil
		}
		if err != nil {
			return nil, nil, 0, err
		}
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}
		counts[r]++
		utflen[n]++
	}
}
