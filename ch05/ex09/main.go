package main

import (
	"regexp"
)

var re = regexp.MustCompile(`\$\w*`)

func expand(s string, f func(string) string) string {
	return re.ReplaceAllStringFunc(s, func(s string) string {
		return f(s[1:])
	})
}
