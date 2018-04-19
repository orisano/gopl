package params

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestUnpack(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		type testT struct {
			A string
			B string `http:"C,^[0-9]+$"`
		}
		var s testT
		req := &http.Request{}
		req.Form = url.Values{}
		req.Form.Add("a", "abc")
		req.Form.Add("C", "100000")

		if err := Unpack(req, &s); err != nil {
			t.Fatal(err)
		}

		expected := testT{"abc", "100000"}
		if !reflect.DeepEqual(s, expected) {
			t.Errorf("unexpecetd result. expected: %v, but got: %v", expected, s)
		}
	})

	t.Run("Failed", func(t *testing.T) {
		type testT struct {
			A string
			B string `http:"C,^[0-9]+$"`
		}
		var s testT
		req := &http.Request{}
		req.Form = url.Values{}
		req.Form.Add("a", "abc")
		req.Form.Add("C", "10a000")

		if err := Unpack(req, &s); err == nil {
			t.Fatal("must be error")
		}
	})
}

func TestPack(t *testing.T) {
	tests := []struct {
		in       interface{}
		expected string
	}{
		{
			in: &struct {
				A string
			}{
				"lower",
			},
			expected: "a=lower",
		},
		{
			in: &struct {
				one int
				two int
			}{
				1, 2,
			},
			expected: "one=1&two=2",
		},
		{
			in: &struct {
				b bool
			}{
				true,
			},
			expected: "b=true",
		},
		{
			in: &struct {
				q []string
			}{
				[]string{"hello", "world"},
			},
			expected: "q=hello&q=world",
		},
		{
			in: &struct {
				tag int `http:"TAG"`
			}{
				42,
			},
			expected: "TAG=42",
		},
	}

	for _, test := range tests {
		if got := Pack(test.in); got != test.expected {
			t.Errorf("unexpected query string. expected: %v, but got: %v", test.expected, got)
		}
	}
}
