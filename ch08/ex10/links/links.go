package links

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

func Extract(ctx context.Context, uri string) ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("getting %s: %s", uri, resp.Status)
	}

	if strings.Split(resp.Header.Get("Content-Type"), ";")[0] != "text/html" {
		p := toFilePath(toExplicit(resp.Request.URL))
		os.MkdirAll(filepath.Dir(p), 0777)
		f, err := os.Create(p)
		if err != nil {
			return nil, fmt.Errorf("failed to create file: %v", err)
		}
		if _, err := io.Copy(f, resp.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", uri, err)
	}

	var links []string
	attrMap := map[string]string{
		"a":      "href",
		"link":   "href",
		"img":    "src",
		"script": "src",
	}

	forEachNode(doc, func(n *html.Node) {
		if n.Type != html.ElementNode {
			return
		}
		uri, ok := getAttr(n, attrMap[n.Data])
		if !ok {
			return
		}
		u, err := resp.Request.URL.Parse(*uri)
		if err != nil {
			return
		}
		if u.Host != resp.Request.URL.Host {
			return
		}
		links = append(links, u.String())

		*uri = toExplicit(u).Path
	}, nil)

	p := toFilePath(toExplicit(resp.Request.URL))
	os.MkdirAll(filepath.Dir(p), 0777)
	if f, err := os.Create(p); err == nil {
		defer f.Close()
		html.Render(f, doc)
	}

	return links, nil
}

func getAttr(n *html.Node, key string) (*string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return &attr.Val, true
		}
	}
	return nil, false
}

func toFilePath(u *url.URL) string {
	return filepath.Join(".", filepath.FromSlash(u.Path))
}

func toExplicit(u *url.URL) *url.URL {
	x := *u
	if x.Path == "" || strings.HasSuffix(x.Path, "/") {
		x.Path = "index.html"
	}
	x.Path = path.Join("/", x.Host, x.Path)
	return &x
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil {
		post(n)
	}
}
