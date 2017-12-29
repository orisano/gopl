package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	urlStr := flag.String("url", "https://golang.org", "crawl start url")
	flag.Parse()

	u, err := url.Parse(*urlStr)
	if err != nil {
		log.Fatal(err)
	}

	breadthFirst(func(item string) []string {
		log.Print(item)
		return crawl(u.Host, item)
	}, []string{*urlStr})
}

func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

func crawl(origin, url string) []string {
	resp, err := http.Get(url)
	if err != nil {
		log.Print(err)
		return nil
	}
	defer resp.Body.Close()

	var body io.Reader = resp.Body
	reqURL := resp.Request.URL

	if reqURL.Host == origin {
		fpath := getFilePath(reqURL)

		os.MkdirAll(filepath.Dir(fpath), os.ModePerm)

		f, err := os.Create(fpath)
		if err == nil {
			defer f.Close()
			body = io.TeeReader(resp.Body, f)
		} else {
			log.Print(err)
		}
	}

	var links []string
	if ct := resp.Header.Get("Content-Type"); ct == "text/html" || strings.HasPrefix(ct, "text/html;") {
		x, err := extractLinks(body, reqURL)
		if err != nil {
			log.Print(err)
		} else {
			links = x
		}
	}
	io.Copy(ioutil.Discard, body)
	return links
}

func extractLinks(r io.Reader, u *url.URL) ([]string, error) {
	var links []string
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	forEachNode(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key != "href" {
					continue
				}
				link, err := u.Parse(attr.Val)
				if err != nil {
					log.Print(err)
					continue
				}
				links = append(links, link.String())
			}
		}
	}, nil)

	return links, nil
}

func getFilePath(u *url.URL) string {
	p := u.Path
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	if strings.HasSuffix(p, "/") {
		p = p + "index.html"
	}
	return filepath.FromSlash("./" + u.Host + p)
}

func forEachNode(doc *html.Node, pre, post func(*html.Node)) {
	if pre != nil {
		pre(doc)
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(doc)
	}
}
