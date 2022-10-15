package http

import (
	"html/template"
	"strings"
)

func replace(input, from, to string) string {
	return strings.Replace(input, from, to, -1)
}

var funcMap = template.FuncMap{
	"replace": replace,
}
