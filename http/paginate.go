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
}

func getPaginatedUrl(u url.URL, pageNum string) string {
	values := u.Query()
	values.Set("page_id", pageNum)
	u.RawQuery = values.Encode()
	return u.String()
}

func (p *Pagination) paginate(u url.URL, current int, total int) {
	var pages []*Page

	for i := 0; i < total; i++ {
		pageNum := i + 1
		page := &Page{Url: getPaginatedUrl(u, strconv.Itoa(pageNum)), Number: pageNum}
		pages = append(pages, page)
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
