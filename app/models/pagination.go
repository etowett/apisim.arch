package models

import (
	null "gopkg.in/guregu/null.v4"
)

type Filter struct {
	Page     int
	Per      int
	Term     string
	FromTime null.Time
	ToTime   null.Time
	UserID   null.Int
	Deleted  null.Bool
}

type Pagination struct {
	Count       int      `json:"count"`
	HasNextPage bool     `json:"has_next_page"`
	HasPrevPage bool     `json:"has_prev_page"`
	NextPage    null.Int `json:"next_page"`
	NumPages    int      `json:"num_pages"`
	Page        int      `json:"page"`
	Per         int      `json:"per"`
	PrevPage    null.Int `json:"prev_page"`
}

func NewPagination(count, page, per int) Pagination {

	var prevPage, nextPage null.Int

	hasPrevPage := page > 1
	if hasPrevPage {
		prevPage = null.IntFrom(int64(page - 1))
	}

	numPages := count / per
	if count == 0 {
		numPages = 1
	} else if count%per != 0 {
		numPages++
	}

	hasNextPage := page < numPages
	if hasNextPage {
		nextPage = null.IntFrom(int64(page + 1))
	}

	return Pagination{
		Count:       count,
		HasNextPage: hasNextPage,
		HasPrevPage: hasPrevPage,
		NextPage:    nextPage,
		NumPages:    numPages,
		Page:        page,
		Per:         per,
		PrevPage:    prevPage,
	}
}
