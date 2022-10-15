package http

import (
	"html/template"
	"strings"
	"time"
)

var funcMap = template.FuncMap{
	"replace": replace,
	"timeAgo": timeAgo,
}

func replace(input, from, to string) string {
	return strings.Replace(input, from, to, -1)
}

func timeAgo(t time.Time) string {
	diff := time.Now().Sub(t)
	out := time.Time{}.Add(diff)
	return out.Format("15:04:05")
}
