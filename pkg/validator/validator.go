package validator

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// ValidateID 验证 ID 参数（必须 > 0）
func ValidateID(c *gin.Context, paramName string) (int64, bool) {
	idStr := c.Param(paramName)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

// ValidatePagination 验证分页参数
func ValidatePagination(page, pageSize int) (int, int) {
	// 默认值
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	// 最大值限制
	if pageSize > 100 {
		pageSize = 100
	}

	return page, pageSize
}

// CalculateOffset 计算分页偏移量
func CalculateOffset(page, pageSize int) int {
	return (page - 1) * pageSize
}
