package treesort

import (
	"bytes"
	"strconv"
)

type tree struct {
	value       int
	left, right *tree
}

func (t *tree) String() string {
	if t == nil {
		return ""
	}
	var b bytes.Buffer

	if s := t.left.String(); len(s) > 0 {
		b.WriteString(s)
		b.WriteByte(' ')
	}
	b.WriteString(strconv.Itoa(t.value))
	if s := t.right.String(); len(s) > 0 {
		b.WriteByte(' ')
		b.WriteString(s)
	}
	return b.String()
}

func Sort(values []int) {
	var root *tree
	for _, v := range values {
		root = add(root, v)
	}
	appendValues(values[:0], root)
}

func appendValues(values []int, t *tree) []int {
	if t != nil {
		values = appendValues(values, t.left)
		values = append(values, t.value)
		values = appendValues(values, t.right)
	}
	return values
}

func add(t *tree, value int) *tree {
	if t == nil {
		t = new(tree)
		t.value = value
		return t
	}
	if value < t.value {
		t.left = add(t.left, value)
	} else {
		t.right = add(t.right, value)
	}
	return t
}
