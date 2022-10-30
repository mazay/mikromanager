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

func (p *Pagination) paginate(u url.URL, current int, total int) {
	var pages []*Page

	for i := 0; i < total; i++ {
		pageNum := i + 1
		values := u.Query()
		values.Set("page_id", strconv.Itoa(pageNum))
		u.RawQuery = values.Encode()

		page := &Page{Url: u.String(), Number: pageNum}
		pages = append(pages, page)
	}

	p.Pages = pages

	if current > 1 {
		pageNum := current - 1
		values := u.Query()
		values.Set("page_id", strconv.Itoa(pageNum))
		u.RawQuery = values.Encode()

		p.Previous = &Page{Url: u.String(), Number: pageNum}
	}

	if current < total {
		pageNum := current + 1
		values := u.Query()
		values.Set("page_id", strconv.Itoa(pageNum))
		u.RawQuery = values.Encode()

		p.Next = &Page{Url: u.String(), Number: pageNum}
	}
}
