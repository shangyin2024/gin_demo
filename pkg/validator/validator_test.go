package validator

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestValidateID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name    string
		id      string
		wantID  int64
		wantOK  bool
	}{
		{"Valid ID", "123", 123, true},
		{"Valid ID 1", "1", 1, true},
		{"Invalid - Zero", "0", 0, false},
		{"Invalid - Negative", "-1", 0, false},
		{"Invalid - String", "abc", 0, false},
		{"Invalid - Float", "12.3", 0, false},
		{"Invalid - Empty", "", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试 context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: tt.id}}

			// 执行测试
			gotID, gotOK := ValidateID(c, "id")

			// 验证结果
			assert.Equal(t, tt.wantID, gotID, "ID should match")
			assert.Equal(t, tt.wantOK, gotOK, "OK should match")
		})
	}
}

func TestValidatePagination(t *testing.T) {
	tests := []struct {
		name         string
		page         int
		pageSize     int
		wantPage     int
		wantPageSize int
	}{
		{"Valid params", 2, 30, 2, 30},
		{"Default page", 0, 20, 1, 20},
		{"Default page negative", -1, 20, 1, 20},
		{"Default pageSize", 1, 0, 1, 20},
		{"Default pageSize negative", 1, -10, 1, 20},
		{"Max pageSize exceeded", 1, 200, 1, 100},
		{"All defaults", 0, 0, 1, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPage, gotPageSize := ValidatePagination(tt.page, tt.pageSize)

			assert.Equal(t, tt.wantPage, gotPage, "Page should match")
			assert.Equal(t, tt.wantPageSize, gotPageSize, "PageSize should match")
		})
	}
}

func TestCalculateOffset(t *testing.T) {
	tests := []struct {
		name       string
		page       int
		pageSize   int
		wantOffset int
	}{
		{"Page 1", 1, 20, 0},
		{"Page 2", 2, 20, 20},
		{"Page 3", 3, 10, 20},
		{"Page 10", 10, 5, 45},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOffset := CalculateOffset(tt.page, tt.pageSize)
			assert.Equal(t, tt.wantOffset, gotOffset, "Offset should match")
		})
	}
}
