package dto

import (
	"math"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
)

type PaginationReq struct {
	Page  int `form:"page,default=1"   binding:"min=1"`
	Limit int `form:"limit,default=20" binding:"min=1,max=100"`
}

func (p *PaginationReq) Offset() int {
	if p.Page < 1 {
		return 0
	}
	return (p.Page - 1) * p.Limit
}

func ParsePagination(c *gin.Context) PaginationReq {
	var p PaginationReq
	if err := c.ShouldBindQuery(&p); err != nil {
		return PaginationReq{Page: DefaultPage, Limit: DefaultLimit}
	}
	if p.Page < 1 {
		p.Page = DefaultPage
	}
	if p.Limit < 1 || p.Limit > MaxLimit {
		p.Limit = DefaultLimit
	}
	return p
}

// PaginatedRes generic sayfalı yanıt — Go 1.18+ generics
type PaginatedRes[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

func NewPaginatedRes[T any](data []T, total int64, page, limit int) PaginatedRes[T] {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return PaginatedRes[T]{
		Data:       data,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}
