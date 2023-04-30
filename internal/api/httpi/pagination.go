package httpi

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaginationRequest struct {
	Sort     string `json:"sort"`
	Page     int64  `json:"page"`
	PageSize int64  `json:"page_size"`
	All      bool   `json:"all"`
}

type PaginationResponse struct {
	TotalCount   int64 `json:"total_count"`
	CurrentCount int64 `json:"current_count"`
	TotalPages   int64 `json:"total_pages"`
	CurrentPage  int64 `json:"current_page"`
}

func BuildPagination(ctx *gin.Context) PaginationRequest {
	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("page_size")
	sort := ctx.Query("sort")
	if sort == "" {
		sort = "created_at DESC"
	}

	var all bool
	if pageSizeStr == "-1" {
		all = true
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	return PaginationRequest{
		Page:     int64(page),
		PageSize: int64(pageSize),
		Sort:     sort,
		All:      all,
	}
}
