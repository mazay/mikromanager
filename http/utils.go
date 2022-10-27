package http

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"time"
)

var funcMap = template.FuncMap{
	"replace":      replace,
	"timeAgo":      timeAgo,
	"memoryUsage":  memoryUsage,
	"getExportUrl": getExportUrl,
}

func replace(input, from, to string) string {
	return strings.Replace(input, from, to, -1)
}

func timeAgo(t time.Time) string {
	diff := time.Since(t)
	out := time.Time{}.Add(diff)
	return out.Format("15:04:05")
}

func memoryUsage(total int64, free int64) string {
	usage := (float64(total) - float64(free)) / float64(total) * 100
	return fmt.Sprintf("%.2f", usage)
}

func getExportUrl(root string, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return filepath.Join("backups", rel)
}
