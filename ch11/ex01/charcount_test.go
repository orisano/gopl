package charcount

import (
	"bytes"
	"reflect"
	"testing"
)

func TestCharCount(t *testing.T) {
	tests := []struct {
		input []byte

		counts  map[rune]int
		utflen  []int
		invalid int
	}{
		{
			input: nil,

			counts:  map[rune]int{},
			utflen:  []int{0, 0, 0, 0, 0},
			invalid: 0,
		},
		{
			input: []byte("a"),

			counts:  map[rune]int{'a': 1},
			utflen:  []int{0, 1, 0, 0, 0},
			invalid: 0,
		},
		{
			input: []byte("é"),

			counts:  map[rune]int{'é': 1},
			utflen:  []int{0, 0, 1, 0, 0},
			invalid: 0,
		},
		{
			input: []byte("あ"),

			counts:  map[rune]int{'あ': 1},
			utflen:  []int{0, 0, 0, 1, 0},
			invalid: 0,
		},
		{
			input: []byte("𩸽"),

			counts:  map[rune]int{'𩸽': 1},
			utflen:  []int{0, 0, 0, 0, 1},
			invalid: 0,
		},
		{
			input: []byte("あa"),

			counts:  map[rune]int{'あ': 1, 'a': 1},
			utflen:  []int{0, 1, 0, 1, 0},
			invalid: 0,
		},
		{
			input: []byte{0x80, 0x80},

			counts:  map[rune]int{},
			utflen:  []int{0, 0, 0, 0, 0},
			invalid: 2,
		},
	}

	for _, test := range tests {
		counts, utflen, invalid, err := CharCount(bytes.NewReader(test.input))
		if err != nil {
			t.Error(err)
			continue
		}
		if got, expected := counts, test.counts; !reflect.DeepEqual(got, expected) {
			t.Errorf("unexpected counts. expected: %#v, but got: %#v", expected, got)
		}
		if got, expected := utflen, test.utflen; !reflect.DeepEqual(got, expected) {
			t.Errorf("unexpected utflen. expected: %#v, but got: %#v", expected, got)
		}
		if got, expected := invalid, test.invalid; got != expected {
			t.Errorf("unexpected invalid. expected: %#v, but got: %#v", expected, got)
		}
	}
}
