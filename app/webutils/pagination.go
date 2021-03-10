package webutils

import (
	"apisim/app/models"
	"strconv"
	"strings"

	"github.com/revel/revel"
)

func FilterFromQuery(
	c *revel.Params,
) (*models.Filter, error) {

	var err error

	page := 1
	pageQuery := strings.TrimSpace(c.Query.Get("page"))
	if pageQuery != "" {
		page, err = strconv.Atoi(pageQuery)
		if err != nil {
			page = 1
		}
	}

	per := 20
	perQuery := strings.TrimSpace(c.Query.Get("per"))
	if perQuery != "" {
		per, err = strconv.Atoi(perQuery)
		if err != nil {
			per = 20
		}
	}

	var search string
	searchQuery := strings.TrimSpace(c.Query.Get("term"))
	if searchQuery != "" {
		search = searchQuery
	}

	return &models.Filter{
		Page: page,
		Per:  per,
		Term: search,
	}, nil
}
