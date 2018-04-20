package sexpr

import (
	"reflect"
	"testing"
)

func TestMarshal(t *testing.T) {
	type Movie struct {
		Title, Subtitle string
		Year            int
		Oscars          []string
	}

	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)",
			"Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
	}
	b, err := Marshal(strangelove)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	expected := `((Title "Dr. Strangelove") (Subtitle "How I Learned to Stop Worrying and Love the Bomb") (Year 1964) (Oscars ("Best Actor (Nomin.)" "Best Adapted Screenplay (Nomin.)" "Best Director (Nomin.)" "Best Picture (Nomin.)")))`
	if got != expected {
		t.Errorf("unexpected s expr. expected: %q, but got: %q", expected, got)
	}
}

func TestUnmarshal(t *testing.T) {
	t.Run("Movie", func(t *testing.T) {
		type Movie struct {
			Title, Subtitle string
			Year            int
			Oscars          []string
		}

		expected := Movie{
			Title:    "Dr. Strangelove",
			Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
			Year:     1964,
			Oscars: []string{
				"Best Actor (Nomin.)",
				"Best Adapted Screenplay (Nomin.)",
				"Best Director (Nomin.)",
				"Best Picture (Nomin.)",
			},
		}

		b := []byte(`((Title "Dr. Strangelove") (Subtitle "How I Learned to Stop Worrying and Love the Bomb") (Year 1964) (Oscars ("Best Actor (Nomin.)" "Best Adapted Screenplay (Nomin.)" "Best Director (Nomin.)" "Best Picture (Nomin.)")))`)
		var got Movie
		if err := Unmarshal(b, &got); err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("unexpected Movie. expected: %#v, but got: %#v", expected, got)
		}
	})
	t.Run("Bool", func(t *testing.T) {
		var b bool
		if err := Unmarshal([]byte(`t`), &b); err != nil {
			t.Fatal(err)
		}
		if b != true {
			t.Error("unexpected boolean. expected: true, but got: false")
		}
	})

	t.Run("Complex", func(t *testing.T) {
		var c complex128
		if err := Unmarshal([]byte(`#C(1.0 2.0)`), &c); err != nil {
			t.Fatal(err)
		}
		if c != 1+2i {
			t.Errorf("unexpected complex. expected: 1+2i, but got: %v", c)
		}
	})

	t.Run("Float", func(t *testing.T) {
		var f float64
		if err := Unmarshal([]byte(`2.5`), &f); err != nil {
			t.Fatal(err)
		}
		if f != 2.5 {
			t.Errorf("unexpected float. expected: 2.5, but got: %v", f)
		}
	})

	t.Run("InterfaceSlice", func(t *testing.T) {
		var i interface{}
		if err := Unmarshal([]byte(`("[]int" (1 2 3 4 5))`), &i); err != nil {
			t.Fatal(err)
		}
		expected := []int{1, 2, 3, 4, 5}
		if !reflect.DeepEqual(i, expected) {
			t.Errorf("unexpected value. expected: %#v, but got: %#v", expected, i)
		}
	})
	t.Run("InterfaceArray", func(t *testing.T) {
		var i interface{}
		if err := Unmarshal([]byte(`("[6]byte" (5 4 3 2 1 0))`), &i); err != nil {
			t.Fatal(err)
		}
		expected := [6]byte{5, 4, 3, 2, 1, 0}
		if !reflect.DeepEqual(i, expected) {
			t.Errorf("unexpected value. expected: %#v, but got: %#v", expected, i)
		}
	})
	t.Run("InterfaceMap", func(t *testing.T) {
		var i interface{}
		if err := Unmarshal([]byte(`("map[string]int" (("a" 1) ("aa" 2) ("日本語" 3)))`), &i); err != nil {
			t.Fatal(err)
		}
		expected := map[string]int{
			"a":   1,
			"aa":  2,
			"日本語": 3,
		}
		if !reflect.DeepEqual(i, expected) {
			t.Errorf("unexpected value. expected: %#v, but got: %#v", expected, i)
		}
	})
	t.Run("InterfaceMapArrayKey", func(t *testing.T) {
		var i interface{}
		if err := Unmarshal([]byte(`("map[[1]byte]int" (((1) 1) ((2) 2) ((3) 3)))`), &i); err != nil {
			t.Fatal(err)
		}
		expected := map[[1]byte]int{
			[1]byte{1}: 1,
			[1]byte{2}: 2,
			[1]byte{3}: 3,
		}
		if !reflect.DeepEqual(i, expected) {
			t.Errorf("unexpected value. expected: %#v, but got: %#v", expected, i)
		}
	})
	t.Run("InterfaceMapSliceValue", func(t *testing.T) {
		var i interface{}
		if err := Unmarshal([]byte(`("map[int][]int" ((-10000 (1 2 3)) (200000 (4 5 6)) (3000000000 (7 8 9)))`), &i); err != nil {
			t.Fatal(err)
		}
		expected := map[int][]int{
			-10000:     {1, 2, 3},
			200000:     {4, 5, 6},
			3000000000: {7, 8, 9},
		}
		if !reflect.DeepEqual(i, expected) {
			t.Errorf("unexpected value. expected: %#v, but got: %#v", expected, i)
		}
	})
	t.Run("InterfaceMapInterface", func(t *testing.T) {
		var i interface{}
		if err := Unmarshal([]byte(`("map[string]interface{}" (("a" ("[]int" (1 2 3))) ("b" ("string" "hello"))))`), &i); err != nil {
			t.Fatal(err)
		}
		expected := map[string]interface{}{
			"a":     []int{1, 2, 3},
			"b":     "hello",
		}
		if !reflect.DeepEqual(i, expected) {
			t.Errorf("unexpected value. expected: %#v, but got: %#v", expected, i)
		}
	})
}
