package xmltree

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

type Node interface {
	String() string
} // CharData or *Element

type CharData string

func (c CharData) String() string {
	return string(c)
}

type Element struct {
	Type     xml.Name
	Attr     []xml.Attr
	Children []Node
}

func nameToString(name xml.Name) string {
	if name.Space != "" {
		return name.Space + ":" + name.Local
	} else {
		return name.Local
	}
}

func (el *Element) String() string {
	if el == nil {
		return "<nil>"
	}

	var b bytes.Buffer
	b.WriteString("<" + nameToString(el.Type))
	for _, attr := range el.Attr {
		fmt.Fprintf(&b, " %s=%q", nameToString(attr.Name), attr.Value)
	}
	b.WriteByte('>')
	for _, child := range el.Children {
		b.WriteString(child.String())
	}
	b.WriteString("</" + nameToString(el.Type) + ">")

	return b.String()
}

func Build(r io.Reader) ([]Node, error) {
	dec := xml.NewDecoder(r)

	stack := []*Element{{}}
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		switch tok := tok.(type) {
		case xml.StartElement:
			el := &Element{
				Type: tok.Name,
				Attr: tok.Attr,
			}
			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, el)
			stack = append(stack, el)
		case xml.CharData:
			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, CharData(tok))
		case xml.EndElement:
			stack = stack[:len(stack)-1]
		}
	}
	return stack[0].Children, nil
}
