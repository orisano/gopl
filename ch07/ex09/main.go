package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/orisano/gopl/ch07/ex08"
)

type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return d
}

var table = struct {
	sync.RWMutex
	tracks []*Track
	sorter *ex08.TableSorter
}{
	tracks: []*Track{
		{"Go", "Dlilah", "From the Roots Up", 2012, length("3m38s")},
		{"Go", "Moby", "Moby", 1992, length("3m37s")},
		{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
		{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
	},
	sorter: &ex08.TableSorter{},
}

var tmpl = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <title>Tracks</title>
</head>
<body>
    <h1>Tracks</h1>
    <table>
        <tr>
            <th><a href="?sort=title">Title</a></th>
            <th><a href="?sort=artist">Artist</a></th>
            <th><a href="?sort=album">Album</a></th>
            <th><a href="?sort=year">Year</a></th>
            <th><a href="?sort=length">Length</a></th>
        </tr>
        {{ range . }}
            <tr>
                <td>{{ .Title }}</td>
                <td>{{ .Artist }}</td>
                <td>{{ .Album }}</td>
                <td>{{ .Year }}</td>
                <td>{{ .Length }}</td>
            </tr>
        {{ end }}
    </table>
</body>
</html>
`))

func printTracks(w io.Writer, tracks []*Track) {
	tmpl.Execute(w, tracks)
}

type customSorter struct {
	t    []*Track
	less func(x, y *Track) bool
}

func (x customSorter) Len() int           { return len(x.t) }
func (x customSorter) Less(i, j int) bool { return x.less(x.t[i], x.t[j]) }
func (x customSorter) Swap(i, j int)      { x.t[i], x.t[j] = x.t[j], x.t[i] }

func main() {
	lessMap := map[string]func(x, y *Track) bool{
		"title":  func(x, y *Track) bool { return x.Title < y.Title },
		"artist": func(x, y *Track) bool { return x.Artist < y.Artist },
		"album":  func(x, y *Track) bool { return x.Album < y.Album },
		"year":   func(x, y *Track) bool { return x.Year < y.Year },
		"length": func(x, y *Track) bool { return x.Length < y.Length },
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		column := r.FormValue("sort")
		if len(column) > 0 {
			if less, ok := lessMap[column]; ok {
				table.Lock()
				table.sorter.SetSorter(customSorter{table.tracks, less})
				sort.Sort(table.sorter)
				table.Unlock()
			}
		}

		table.RLock()
		printTracks(w, table.tracks)
		table.RUnlock()
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
