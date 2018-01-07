package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Selector interface {
	Match(xml.StartElement) bool
}

type TypeSelector struct {
	Type string
}

func (s *TypeSelector) Match(el xml.StartElement) bool {
	return s.Type == el.Name.Local
}

type AttributeSelector struct {
	Name  string
	Value string
}

func (s *AttributeSelector) Match(el xml.StartElement) bool {
	val, ok := getAttr(el, s.Name)
	return ok && s.Value == val
}

func getAttr(el xml.StartElement, name string) (string, bool) {
	for _, attr := range el.Attr {
		if attr.Name.Local != name {
			continue
		}
		return attr.Value, true
	}
	return "", false
}

func ParseSelector(s string) (Selector, error) {
	if s == "" {
		return nil, fmt.Errorf("empty selector")
	}
	switch {
	case strings.HasPrefix(s, "."):
		return &AttributeSelector{
			Name:  "class",
			Value: s[1:],
		}, nil
	case strings.HasPrefix(s, "#"):
		return &AttributeSelector{
			Name:  "id",
			Value: s[1:],
		}, nil
	default:
		return &TypeSelector{Type: s}, nil
	}
}

func main() {
	var selectors []Selector
	for _, arg := range os.Args[1:] {
		selector, err := ParseSelector(arg)
		if err != nil {
			log.Fatalf("xmlselect: invalid selector %q: %v", arg, err)
		}
		selectors = append(selectors, selector)
	}

	dec := xml.NewDecoder(os.Stdin)
	var stack []xml.StartElement
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "xmlselect: %v\n", err)
			os.Exit(1)
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			stack = append(stack, tok)
		case xml.EndElement:
			stack = stack[:len(stack)-1]
		case xml.CharData:
			if querySelector(stack, selectors) {
				var tags []string
				for _, el := range stack {
					tags = append(tags, el.Name.Local)
				}
				fmt.Printf("%s: %s\n", strings.Join(tags, " "), tok)
			}
		}
	}
}

func querySelector(elements []xml.StartElement, selectors []Selector) bool {
	for len(selectors) <= len(elements) {
		if len(selectors) == 0 {
			return true
		}
		if selectors[0].Match(elements[0]) {
			selectors = selectors[1:]
		}
		elements = elements[1:]
	}
	return false
}
