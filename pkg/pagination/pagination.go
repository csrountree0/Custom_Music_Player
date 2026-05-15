package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Params struct {
	Page  int
	Limit int
}

func FromQuery(c *gin.Context) Params {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 1
	}
	if limit > 100 {
		limit = 100
	}

	return Params{Page: page, Limit: limit}
}

func (p Params) Offset() int {
	return (p.Page - 1) * p.Limit
}

type Page[T any] struct {
	Items []T   `json:"items"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}
