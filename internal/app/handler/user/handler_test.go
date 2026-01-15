package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gin_demo/internal/domain/service"
	"gin_demo/internal/repository"
	"gin_demo/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService 是 UserService 的 mock 实现
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(ctx context.Context, input service.RegisterInput) (repository.User, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(repository.User), args.Error(1)
}

func (m *MockUserService) Login(ctx context.Context, input service.LoginInput) (repository.User, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(repository.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, userID int64) (repository.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(repository.User), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(ctx context.Context, email string) (repository.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(repository.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, input service.UpdateUserInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockUserService) ChangePassword(ctx context.Context, input service.ChangePasswordInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserService) ListUsers(ctx context.Context, limit, offset int32) ([]repository.User, int64, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]repository.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserService) TransferUserData(ctx context.Context, fromUserID, toUserID int64) error {
	args := m.Called(ctx, fromUserID, toUserID)
	return args.Error(0)
}

func (m *MockUserService) BatchUpdateUsers(ctx context.Context, updates []service.BatchUpdateInput) error {
	args := m.Called(ctx, updates)
	return args.Error(0)
}

// setupTestHandler 设置测试 Handler
func setupTestHandler() (*Handler, *MockUserService, *auth.DefaultJWTManager) {
	mockService := new(MockUserService)
	jwtManager := auth.NewDefaultJWTManager("test-secret", 1*time.Hour)
	handler := NewHandler(mockService, jwtManager)
	
	gin.SetMode(gin.TestMode)
	
	return handler, mockService, jwtManager
}

// TestHandler_Register 测试用户注册
func TestHandler_Register(t *testing.T) {
	handler, mockService, _ := setupTestHandler()

	t.Run("成功注册", func(t *testing.T) {
		// 准备请求
		reqBody := RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)

		// 创建测试上下文
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		// Mock 服务层
		expectedUser := repository.User{
			ID:        1,
			Username:  "testuser",
			Email:     "test@example.com",
			Status:    1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockService.On("Register", mock.Anything, mock.AnythingOfType("service.RegisterInput")).
			Return(expectedUser, nil)

		// 执行 Handler
		handler.Register(c)

		// 断言
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), response["code"])
		
		mockService.AssertExpectations(t)
	})

	t.Run("参数验证失败", func(t *testing.T) {
		// 准备无效请求
		reqBody := RegisterRequest{
			Username: "te", // 太短
			Email:    "invalid-email",
			Password: "123", // 太短
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		// 执行 Handler
		handler.Register(c)

		// 断言
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("用户已存在", func(t *testing.T) {
		// 重新创建 handler 和 mock（避免之前的 mock 影响）
		handler, mockService, _ := setupTestHandler()
		
		reqBody := RegisterRequest{
			Username: "existinguser",
			Email:    "existing@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		// Mock 服务层返回已存在错误
		mockService.On("Register", mock.Anything, mock.AnythingOfType("service.RegisterInput")).
			Return(repository.User{}, service.ErrUserExists)

		// 执行 Handler
		handler.Register(c)

		// 断言
		assert.Equal(t, http.StatusConflict, w.Code)
		mockService.AssertExpectations(t)
	})
}

// TestHandler_Login 测试用户登录
func TestHandler_Login(t *testing.T) {
	handler, mockService, jwtManager := setupTestHandler()

	t.Run("成功登录", func(t *testing.T) {
		reqBody := LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		// Mock 服务层
		expectedUser := repository.User{
			ID:        1,
			Username:  "testuser",
			Email:     "test@example.com",
			Status:    1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockService.On("Login", mock.Anything, mock.AnythingOfType("service.LoginInput")).
			Return(expectedUser, nil)

		// 执行 Handler
		handler.Login(c)

		// 断言
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), response["code"])
		
		// 验证返回了 token
		data := response["data"].(map[string]interface{})
		assert.NotEmpty(t, data["token"])
		
		// 验证 token 有效性
		token := data["token"].(string)
		claims, err := jwtManager.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, claims.UserID)
		
		mockService.AssertExpectations(t)
	})

	t.Run("用户不存在", func(t *testing.T) {
		handler, mockService, _ := setupTestHandler()
		
		reqBody := LoginRequest{
			Email:    "notexist@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		mockService.On("Login", mock.Anything, mock.AnythingOfType("service.LoginInput")).
			Return(repository.User{}, service.ErrUserNotFound)

		handler.Login(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("密码错误", func(t *testing.T) {
		handler, mockService, _ := setupTestHandler()
		
		reqBody := LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		mockService.On("Login", mock.Anything, mock.AnythingOfType("service.LoginInput")).
			Return(repository.User{}, service.ErrInvalidPassword)

		handler.Login(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}

// TestHandler_GetProfile 测试获取当前用户信息
func TestHandler_GetProfile(t *testing.T) {
	handler, mockService, jwtManager := setupTestHandler()

	t.Run("成功获取个人信息", func(t *testing.T) {
		// 生成测试 Token
		userID := int64(1)
		token, _ := jwtManager.GenerateToken(userID)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/users/me", nil)
		c.Request.Header.Set("Authorization", "Bearer "+token)
		
		// 模拟认证中间件设置的 user_id
		c.Set("user_id", userID)

		expectedUser := repository.User{
			ID:        userID,
			Username:  "testuser",
			Email:     "test@example.com",
			Status:    1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockService.On("GetUserByID", mock.Anything, userID).Return(expectedUser, nil)

		handler.GetProfile(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("未认证", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/users/me", nil)
		// 不设置 user_id，模拟未认证

		handler.GetProfile(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

// TestHandler_UpdateProfile 测试更新当前用户信息
func TestHandler_UpdateProfile(t *testing.T) {
	handler, mockService, jwtManager := setupTestHandler()

	t.Run("成功更新个人信息", func(t *testing.T) {
		userID := int64(1)
		token, _ := jwtManager.GenerateToken(userID)

		reqBody := UpdateUserRequest{
			Username: "newusername",
			Email:    "newemail@example.com",
			Avatar:   "https://example.com/avatar.jpg",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PUT", "/users/me", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("Authorization", "Bearer "+token)
		c.Set("user_id", userID)

		mockService.On("UpdateUser", mock.Anything, mock.AnythingOfType("service.UpdateUserInput")).
			Return(nil)

		handler.UpdateProfile(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}

// TestHandler_ChangePassword 测试修改密码
func TestHandler_ChangePassword(t *testing.T) {
	handler, mockService, jwtManager := setupTestHandler()

	t.Run("成功修改密码", func(t *testing.T) {
		userID := int64(1)
		token, _ := jwtManager.GenerateToken(userID)

		reqBody := ChangePasswordRequest{
			OldPassword: "oldpassword",
			NewPassword: "newpassword123",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PUT", "/users/me/password", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("Authorization", "Bearer "+token)
		c.Set("user_id", userID)

		mockService.On("ChangePassword", mock.Anything, mock.AnythingOfType("service.ChangePasswordInput")).
			Return(nil)

		handler.ChangePassword(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("旧密码错误", func(t *testing.T) {
		handler, mockService, jwtManager := setupTestHandler()
		
		userID := int64(1)
		token, _ := jwtManager.GenerateToken(userID)

		reqBody := ChangePasswordRequest{
			OldPassword: "wrongpassword",
			NewPassword: "newpassword123",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PUT", "/users/me/password", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Header.Set("Authorization", "Bearer "+token)
		c.Set("user_id", userID)

		mockService.On("ChangePassword", mock.Anything, mock.AnythingOfType("service.ChangePasswordInput")).
			Return(service.ErrInvalidPassword)

		handler.ChangePassword(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}

// TestHandler_GetUser 测试获取指定用户（管理员）
func TestHandler_GetUser(t *testing.T) {
	handler, mockService, _ := setupTestHandler()

	t.Run("成功获取用户", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/users/1", nil)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		expectedUser := repository.User{
			ID:        1,
			Username:  "testuser",
			Email:     "test@example.com",
			Status:    1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockService.On("GetUserByID", mock.Anything, int64(1)).Return(expectedUser, nil)

		handler.GetUser(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("用户不存在", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/users/999", nil)
		c.Params = gin.Params{{Key: "id", Value: "999"}}

		mockService.On("GetUserByID", mock.Anything, int64(999)).
			Return(repository.User{}, service.ErrUserNotFound)

		handler.GetUser(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("无效的用户ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/users/invalid", nil)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		handler.GetUser(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestHandler_ListUsers 测试用户列表
func TestHandler_ListUsers(t *testing.T) {
	handler, mockService, _ := setupTestHandler()

	t.Run("成功获取用户列表", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/users?page=1&size=10", nil)

		expectedUsers := []repository.User{
			{ID: 1, Username: "user1", Email: "user1@example.com"},
			{ID: 2, Username: "user2", Email: "user2@example.com"},
			{ID: 3, Username: "user3", Email: "user3@example.com"},
		}
		mockService.On("ListUsers", mock.Anything, int32(10), int32(0)).
			Return(expectedUsers, int64(100), nil)

		handler.ListUsers(c)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), response["code"])
		
		mockService.AssertExpectations(t)
	})

	t.Run("自定义分页参数", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/users?page=2&size=20", nil)

		mockService.On("ListUsers", mock.Anything, int32(20), int32(20)).
			Return([]repository.User{}, int64(100), nil)

		handler.ListUsers(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}

// TestHandler_DeleteUser 测试删除用户
func TestHandler_DeleteUser(t *testing.T) {
	handler, mockService, _ := setupTestHandler()

	t.Run("成功删除用户", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("DELETE", "/users/1", nil)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		mockService.On("DeleteUser", mock.Anything, int64(1)).Return(nil)

		handler.DeleteUser(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("用户不存在", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("DELETE", "/users/999", nil)
		c.Params = gin.Params{{Key: "id", Value: "999"}}

		mockService.On("DeleteUser", mock.Anything, int64(999)).
			Return(service.ErrUserNotFound)

		handler.DeleteUser(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

// TestHandler_UpdateUser 测试更新指定用户（管理员）
func TestHandler_UpdateUser(t *testing.T) {
	handler, mockService, _ := setupTestHandler()

	t.Run("成功更新用户", func(t *testing.T) {
		reqBody := UpdateUserRequest{
			Username: "updateduser",
			Email:    "updated@example.com",
			Avatar:   "https://example.com/new-avatar.jpg",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		mockService.On("UpdateUser", mock.Anything, mock.AnythingOfType("service.UpdateUserInput")).
			Return(nil)

		handler.UpdateUser(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}
