package http

import (
	"net/url"
	"strconv"
)

type Page struct {
	Url    string
	Number int
}

type Pagination struct {
	Pages    []*Page
	Previous *Page
	Next     *Page
	Breaks   []int
}

func getPaginatedUrl(u url.URL, pageNum string) string {
	values := u.Query()
	values.Set("page_id", pageNum)
	u.RawQuery = values.Encode()
	return u.String()
}

func isHidden(num, current, total int) bool {
	var offset = 2

	if num == current || num == total || num == 1 {
		return false
	}

	if num >= current-offset && num <= current+offset {
		return false
	}

	return true
}

func (p *Pagination) paginate(u url.URL, current int, total int) {
	var pages []*Page

	for i := 0; i < total; i++ {
		pageNum := i + 1
		if !isHidden(pageNum, current, total) {
			page := &Page{Url: getPaginatedUrl(u, strconv.Itoa(pageNum)), Number: pageNum}
			pages = append(pages, page)
		} else {
			p.Breaks = append(p.Breaks, pageNum-1)
		}
	}
	p.Pages = pages

	if current > 1 {
		pageNum := current - 1
		p.Previous = &Page{Url: getPaginatedUrl(u, strconv.Itoa(pageNum)), Number: pageNum}
	}

	if current < total {
		pageNum := current + 1
		p.Next = &Page{Url: getPaginatedUrl(u, strconv.Itoa(pageNum)), Number: pageNum}
	}
}
