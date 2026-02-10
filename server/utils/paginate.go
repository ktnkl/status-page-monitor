package utils

import (
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type Pagination struct {
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

func Paginate(r *http.Request, pagination *Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		if page <= 0 {
			page = 1
		}

		limit, _ := strconv.Atoi(q.Get("limit"))
		switch {
		case limit > 100:
			limit = 100
		case limit <= 0:
			limit = 10
		}

		pagination.Page = page
		pagination.Limit = limit

		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}
