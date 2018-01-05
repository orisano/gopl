package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var dbMu sync.RWMutex

var tmpl = template.Must(template.New("list").Parse(`
<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <title>Items</title>
</head>
<body>
    <h1>Items</h1>
    <table>
        <tr>
            <th>Item</th>
            <th>Price</th>
        </tr>
        {{ range $item, $price := . }}
        <tr>
            <td>{{ $item }}</td>
            <td>{{ $price }}</td>
        </tr>
        {{ end }}
    </table>
</body>
</html>
`))

func main() {
	db := database{"shoes": 50, "socks": 5}
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/create", db.create)
	http.HandleFunc("/delete", db.delete)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

func (db database) print(w io.Writer) {
	tmpl.Execute(w, db)
}

func (db database) list(w http.ResponseWriter, req *http.Request) {
	dbMu.RLock()
	defer dbMu.RUnlock()

	db.print(w)
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	dbMu.RLock()
	defer dbMu.RUnlock()

	item := req.URL.Query().Get("item")
	price, ok := db[item]
	if !ok {
		msg := fmt.Sprintf("no such item: %q", item)
		http.Error(w, msg, http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "%s\n", price)
}

func (db database) update(w http.ResponseWriter, req *http.Request) {
	dbMu.Lock()
	defer dbMu.Unlock()

	priceStr := req.URL.Query().Get("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		msg := fmt.Sprintf("price invalid format %q: %v", priceStr, err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	item := req.URL.Query().Get("item")
	if _, ok := db[item]; !ok {
		msg := fmt.Sprintf("not such item: %q", item)
		http.Error(w, msg, http.StatusNotFound)
		return
	}

	db[item] = dollars(price)
	db.print(w)
}

func (db database) create(w http.ResponseWriter, req *http.Request) {
	dbMu.Lock()
	defer dbMu.Unlock()

	priceStr := req.URL.Query().Get("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		msg := fmt.Sprintf("price invalid format %q: %v", priceStr, err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	item := req.URL.Query().Get("item")
	if _, ok := db[item]; ok {
		msg := fmt.Sprintf("already exists item: %q", item)
		http.Error(w, msg, http.StatusConflict)
		return
	}

	db[item] = dollars(price)
	db.print(w)
}

func (db database) delete(w http.ResponseWriter, req *http.Request) {
	dbMu.Lock()
	defer dbMu.Unlock()

	item := req.URL.Query().Get("item")
	if _, ok := db[item]; !ok {
		msg := fmt.Sprintf("not such item: %q", item)
		http.Error(w, msg, http.StatusNotFound)
		return
	}

	delete(db, item)
	db.print(w)
}
