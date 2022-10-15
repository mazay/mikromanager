package http

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"
)

var funcMap = template.FuncMap{
	"replace":     replace,
	"timeAgo":     timeAgo,
	"memoryUsage": memoryUsage,
}

func replace(input, from, to string) string {
	return strings.Replace(input, from, to, -1)
}

func timeAgo(t time.Time) string {
	diff := time.Now().Sub(t)
	out := time.Time{}.Add(diff)
	return out.Format("15:04:05")
}

func memoryUsage(total string, free string) string {
	totalInt, err := strconv.Atoi(total)
	if err != nil {
		return err.Error()
	}
	freeInt, err := strconv.Atoi(free)
	if err != nil {
		return err.Error()
	}
	usage := (float64(totalInt) - float64(freeInt)) / float64(totalInt) * 100
	return fmt.Sprintf("%.2f", usage)
}
