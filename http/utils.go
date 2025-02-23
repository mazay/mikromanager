package http

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/mazay/mikromanager/db"
	"github.com/mazay/mikromanager/internal"
)

var funcMap = template.FuncMap{
	"replace":       replace,
	"timeAgo":       timeAgo,
	"memoryUsage":   memoryUsage,
	"containsInt":   containsInt,
	"humahizeBytes": humahizeBytes,
}

func replace(input, from, to string) string {
	return strings.Replace(input, from, to, -1)
}

func humahizeBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
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
func chunkSliceOfObjects[obj internal.Export | db.Credentials | db.Device | db.User](slice []*obj, chunkSize int) [][]*obj {
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
