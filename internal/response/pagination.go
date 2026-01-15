package response

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaginationRequest 分页请求参数
// 支持两种方式：
// 1. page + size (推荐，对前端友好)
// 2. offset + limit (直接使用数据库参数)
type PaginationRequest struct {
	Page     int `form:"page"`   // 页码（从 1 开始）
	PageSize int `form:"size"`   // 每页数量
	Offset   int `form:"offset"` // 偏移量（直接方式）
	Limit    int `form:"limit"`  // 限制数量（直接方式）
}

// PaginationResponse 分页响应
type PaginationResponse struct {
	Page       int   `json:"page"`        // 当前页码
	PageSize   int   `json:"page_size"`   // 每页数量
	Total      int64 `json:"total"`       // 总记录数
	TotalPages int   `json:"total_pages"` // 总页数
}

// GetPagination 从请求中获取分页参数
// 优先使用 page/size，如果没有则使用 offset/limit
func GetPagination(c *gin.Context) PaginationRequest {
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	// 方式1：使用 page/size
	if page > 0 || size > 0 {
		if page < 1 {
			page = 1
		}
		if size < 1 {
			size = 10
		} else if size > 100 {
			size = 100
		}
		return PaginationRequest{
			Page:     page,
			PageSize: size,
		}
	}

	// 方式2：使用 offset/limit
	if offset >= 0 || limit > 0 {
		if offset < 0 {
			offset = 0
		}
		if limit < 1 {
			limit = 10
		} else if limit > 100 {
			limit = 100
		}
		// 从 offset/limit 反推 page/size
		page = (offset / limit) + 1
		size = limit
		return PaginationRequest{
			Page:     page,
			PageSize: size,
			Offset:   offset,
			Limit:    limit,
		}
	}

	// 默认值
	return PaginationRequest{
		Page:     1,
		PageSize: 10,
	}
}

// GetOffset 计算数据库查询的偏移量
func (p PaginationRequest) GetOffset() int32 {
	return int32((p.Page - 1) * p.PageSize)
}

// GetLimit 获取查询限制数量
func (p PaginationRequest) GetLimit() int32 {
	return int32(p.PageSize)
}

// NewPaginationResponse 创建分页响应
func NewPaginationResponse(page, pageSize int, total int64) PaginationResponse {
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	if totalPages == 0 {
		totalPages = 1
	}

	return PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}

// ListResponse 通用列表响应（带分页）
type ListResponse[T any] struct {
	Items      []T                `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
}

// NewListResponse 创建列表响应
func NewListResponse[T any](items []T, pagination PaginationResponse) ListResponse[T] {
	if items == nil {
		items = []T{} // 避免返回 null
	}
	return ListResponse[T]{
		Items:      items,
		Pagination: pagination,
	}
}
