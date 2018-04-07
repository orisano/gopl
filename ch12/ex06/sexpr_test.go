package sexpr_test

import (
	"testing"

	"github.com/orisano/gopl/ch12/ex06"
)

func TestMarshal(t *testing.T) {
	var x interface{} = []int{1, 2, 3}
	b, err := sexpr.Marshal(&x)
	if err != nil {
		t.Fatal(err)
	}
	if got, expected := string(b), `("[]int" (1
          2
          3))`; got != expected {
		t.Errorf("unexpected bytes. \nexpected: \n%v\nbut got: \n%v", expected, got)
	}
}

func TestMarshalPretty(t *testing.T) {
	type Movie struct {
		Title, Subtitle string
		Year            int
		Oscars          []string
		Sequel          *string
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
	b, err := sexpr.Marshal(strangelove)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	expected := `((Title "Dr. Strangelove")
 (Subtitle "How I Learned to Stop Worrying and Love the Bomb")
 (Year 1964)
 (Oscars ("Best Actor (Nomin.)"
          "Best Adapted Screenplay (Nomin.)"
          "Best Director (Nomin.)"
          "Best Picture (Nomin.)"))
 (Sequel nil))`

	if got != expected {
		t.Errorf("unexpected bytes. \nexpected: \n%v\nbut got: \n%v", expected, got)
	}
}

func TestMarshalOmitEmpty(t *testing.T) {
	type Movie struct {
		Title, Subtitle string
		Year            int
		Oscars          []string
		Sequel          *string
	}
	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     0,
		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)",
			"Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
	}
	b, err := sexpr.Marshal(strangelove)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	expected := `((Title "Dr. Strangelove")
 (Subtitle "How I Learned to Stop Worrying and Love the Bomb")
 (Oscars ("Best Actor (Nomin.)"
          "Best Adapted Screenplay (Nomin.)"
          "Best Director (Nomin.)"
          "Best Picture (Nomin.)")))`

	if got != expected {
		t.Errorf("unexpected bytes. \nexpected: \n%v\nbut got: \n%v", expected, got)
	}
}
