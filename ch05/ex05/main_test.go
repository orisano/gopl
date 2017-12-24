package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func buildHTML(body string) string {
	return "<html><head></head><body>" + body + "</body></html>"
}

func TestCountWordsAndImages(t *testing.T) {
	type result struct {
		words, images int
	}

	ts := []struct {
		HTML     string
		expected result
	}{
		{
			HTML:     buildHTML(""),
			expected: result{words: 0, images: 0},
		},
		{
			HTML:     buildHTML("<p>foo bar</p>"),
			expected: result{words: 2, images: 0},
		},
		{
			HTML:     buildHTML("<h1>It works!</h1><p>foo bar</p>"),
			expected: result{words: 4, images: 0},
		},
		{
			HTML:     buildHTML(`<img src="https://gopl.example/count1.png" />`),
			expected: result{words: 0, images: 1},
		},
		{
			HTML: buildHTML(`
<h1>It works!</h1>
<img src="https://gopl.example/count1.png" />
<div>
	<p>foo</p>
	<span><p>bar</p></span>
	<img src="https://gopl.example/count2.png" />
	<a href="#"><img src="https://gopl.example/count3.png" /></a>
</div>`),
			expected: result{words: 4, images: 3},
		},
	}
	for _, tc := range ts {
		t.Run("Case", func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
				io.WriteString(w, tc.HTML)
			})
			s := httptest.NewServer(mux)
			defer s.Close()

			words, images, err := CountWordsAndImages(s.URL)
			if err != nil {
				t.Fatal(err)
			}
			if words != tc.expected.words || images != tc.expected.images {
				t.Errorf("unexpected result. expected: (%v,%v), but got: (%v,%v)", tc.expected.words, tc.expected.images, words, images)
			}
		})

	}
}
