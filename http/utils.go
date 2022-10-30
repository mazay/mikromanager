package http

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"time"

	"github.com/mazay/mikromanager/utils"
)

var funcMap = template.FuncMap{
	"replace":      replace,
	"timeAgo":      timeAgo,
	"memoryUsage":  memoryUsage,
	"getExportUrl": getExportUrl,
	"containsInt":  containsInt,
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

func containsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// chunkSliceOfObjects accepts slices of Export, Credentials or Device objects and a chunk size
// and returns chunks of the input objects
func chunkSliceOfObjects[obj utils.Export | utils.Credentials | utils.Device](slice []*obj, chunkSize int) [][]*obj {
	var chunks [][]*obj
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		// necessary check to avoid slicing beyond
		// slice capacity
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}
