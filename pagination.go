package fortytwo

import (
	"math"
	"net/http"
	"strconv"
)

type Cursor int

type Pagination struct {
	Cursor   int
	PageSize int
}

func (p *Pagination) ToQuery() map[string]string {
	if p == nil {
		return nil
	}
	r := map[string]string{}
	if p.Cursor != 0 {
		r["page[number]"] = strconv.Itoa(p.Cursor)
	}

	if p.PageSize != 0 {
		r["page[size]"] = strconv.Itoa(p.PageSize)
	}

	return r
}

type PaginationResponse struct {
	Offset           int  `json:"-"`
	NumPages         int  `json:"pages"`
	NumItems         int  `json:"total"`
	ItemsPerPage     int  `json:"items"`
	CurrentPage      int  `json:"current"`
	NextPage         int  `json:"next,omitempty"`
	PrevPage         int  `json:"prev,omitempty"`
	HasPrev, HasNext bool `json:"-"`
}

// pages start at 1 - not 0
func (p *PaginationResponse) Calculate(numItems int) {
	p.NumItems = numItems

	// calculate number of pages
	d := float64(p.NumItems) / float64(p.ItemsPerPage)
	p.NumPages = int(math.Ceil(d))

	// Return the right offset
	p.Offset = (p.CurrentPage - 1) * p.ItemsPerPage

	// HasPrev, HasNext?
	p.HasPrev = p.CurrentPage > 1
	p.HasNext = p.CurrentPage < p.NumPages

	// calculate them
	if p.HasPrev {
		p.PrevPage = p.CurrentPage - 1
	}

	if p.HasNext {
		p.NextPage = p.CurrentPage + 1
	}
}

// GetPaginationInfo retrieves the pagination information from the response header.
func GetPaginationInfo(h http.Header) *PaginationResponse {
	totalItem, found := h["X-Total"]
	if !found {
		return nil
	}
	numItem, err := strconv.Atoi(totalItem[0])
	if err != nil {
		return nil
	}

	currentPage, found := h["X-Page"]
	if !found {
		return nil
	}
	currPage, err := strconv.Atoi(currentPage[0])
	if err != nil {
		return nil
	}

	itemPerPage, found := h["X-Per-Page"]
	if !found {
		return nil
	}
	itemPage, err := strconv.Atoi(itemPerPage[0])
	if err != nil {
		return nil
	}

	pag := &PaginationResponse{
		ItemsPerPage: itemPage,
		CurrentPage:  currPage,
	}

	pag.Calculate(numItem)

	return pag
}
