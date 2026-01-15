package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"gin_demo/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository 是 UserRepository 的 mock 实现
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, userID int64) (repository.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(repository.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (repository.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(repository.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (repository.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(repository.User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, params repository.CreateUserParams) (repository.User, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(repository.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, params repository.UpdateUserParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUserPassword(ctx context.Context, userID int64, password string) error {
	args := m.Called(ctx, userID, password)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepository) WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

func (m *MockUserRepository) BatchExecInTx(ctx context.Context, ops []func(ctx context.Context, tx *sql.Tx) error) error {
	args := m.Called(ctx, ops)
	return args.Error(0)
}

func (m *MockUserRepository) WithTxRepo(tx *sql.Tx) repository.UserRepositoryInterface {
	args := m.Called(tx)
	return args.Get(0).(repository.UserRepositoryInterface)
}

func (m *MockUserRepository) ListUsers(ctx context.Context, limit, offset int32) ([]repository.User, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]repository.User), args.Error(1)
}

func (m *MockUserRepository) CountUsers(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// TestUserService_Register 测试用户注册
func TestUserService_Register(t *testing.T) {
	ctx := context.Background()

	t.Run("成功注册", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		// Mock 数据
		input := RegisterInput{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}

		expectedUser := repository.User{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
			Status:   1,
		}

		// 设置 mock 期望
		mockRepo.On("GetUserByEmail", ctx, input.Email).Return(repository.User{}, sql.ErrNoRows)
		mockRepo.On("GetUserByUsername", ctx, input.Username).Return(repository.User{}, sql.ErrNoRows)
		mockRepo.On("CreateUser", ctx, mock.AnythingOfType("repository.CreateUserParams")).Return(expectedUser, nil)

		// 执行测试
		user, err := service.Register(ctx, input)

		// 断言
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.Username, user.Username)
		assert.Equal(t, expectedUser.Email, user.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("邮箱已存在", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		input := RegisterInput{
			Username: "testuser",
			Email:    "existing@example.com",
			Password: "password123",
		}

		existingUser := repository.User{
			ID:    1,
			Email: "existing@example.com",
		}

		mockRepo.On("GetUserByEmail", ctx, input.Email).Return(existingUser, nil)

		// 执行测试
		_, err := service.Register(ctx, input)

		// 断言
		assert.Error(t, err)
		assert.Equal(t, ErrUserExists, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("用户名已存在", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		input := RegisterInput{
			Username: "existinguser",
			Email:    "new@example.com",
			Password: "password123",
		}

		existingUser := repository.User{
			ID:       1,
			Username: "existinguser",
		}

		mockRepo.On("GetUserByEmail", ctx, input.Email).Return(repository.User{}, sql.ErrNoRows)
		mockRepo.On("GetUserByUsername", ctx, input.Username).Return(existingUser, nil)

		// 执行测试
		_, err := service.Register(ctx, input)

		// 断言
		assert.Error(t, err)
		assert.Equal(t, ErrUserExists, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("参数验证失败", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		testCases := []struct {
			name  string
			input RegisterInput
		}{
			{"空用户名", RegisterInput{Username: "", Email: "test@example.com", Password: "password123"}},
			{"空邮箱", RegisterInput{Username: "testuser", Email: "", Password: "password123"}},
			{"空密码", RegisterInput{Username: "testuser", Email: "test@example.com", Password: ""}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := service.Register(ctx, tc.input)
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidInput, err)
			})
		}
	})
}

// TestUserService_Login 测试用户登录
func TestUserService_Login(t *testing.T) {
	ctx := context.Background()

	t.Run("成功登录", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		input := LoginInput{
			Email:    "test@example.com",
			Password: "password123",
		}

		// 生成正确的 bcrypt 哈希
		hashedPasswordBytes, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		hashedPassword := string(hashedPasswordBytes)

		user := repository.User{
			ID:       1,
			Email:    "test@example.com",
			Password: hashedPassword,
			Status:   1,
		}

		mockRepo.On("GetUserByEmail", ctx, input.Email).Return(user, nil)

		// 执行测试
		resultUser, err := service.Login(ctx, input)

		// 断言
		assert.NoError(t, err)
		assert.Equal(t, user.ID, resultUser.ID)
		assert.Equal(t, user.Email, resultUser.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("用户不存在", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		input := LoginInput{
			Email:    "notexist@example.com",
			Password: "password123",
		}

		mockRepo.On("GetUserByEmail", ctx, input.Email).Return(repository.User{}, sql.ErrNoRows)

		// 执行测试
		_, err := service.Login(ctx, input)

		// 断言
		assert.Error(t, err)
		assert.Equal(t, ErrUserNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("密码错误", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		input := LoginInput{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		// 生成正确的 bcrypt 哈希（正确密码是 password123）
		hashedPasswordBytes, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		hashedPassword := string(hashedPasswordBytes)

		user := repository.User{
			ID:       1,
			Email:    "test@example.com",
			Password: hashedPassword,
			Status:   1,
		}

		mockRepo.On("GetUserByEmail", ctx, input.Email).Return(user, nil)

		// 执行测试
		_, err := service.Login(ctx, input)

		// 断言
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidPassword, err)
		mockRepo.AssertExpectations(t)
	})
}

// TestUserService_GetUserByID 测试通过ID获取用户
func TestUserService_GetUserByID(t *testing.T) {
	ctx := context.Background()

	t.Run("成功获取用户", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		userID := int64(1)
		expectedUser := repository.User{
			ID:       userID,
			Username: "testuser",
			Email:    "test@example.com",
			Status:   1,
		}

		mockRepo.On("GetUserByID", ctx, userID).Return(expectedUser, nil)

		// 执行测试
		user, err := service.GetUserByID(ctx, userID)

		// 断言
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("用户不存在", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		userID := int64(999)

		mockRepo.On("GetUserByID", ctx, userID).Return(repository.User{}, sql.ErrNoRows)

		// 执行测试
		_, err := service.GetUserByID(ctx, userID)

		// 断言
		assert.Error(t, err)
		assert.Equal(t, ErrUserNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

// TestUserService_UpdateUser 测试更新用户信息
func TestUserService_UpdateUser(t *testing.T) {
	ctx := context.Background()

	t.Run("成功更新用户", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		userID := int64(1)
		newUsername := "newusername"
		newEmail := "newemail@example.com"

		currentUser := repository.User{
			ID:       userID,
			Username: "oldusername",
			Email:    "oldemail@example.com",
			Status:   1,
		}

		input := UpdateUserInput{
			UserID:   userID,
			Username: &newUsername,
			Email:    &newEmail,
		}

		mockRepo.On("GetUserByID", ctx, userID).Return(currentUser, nil)
		mockRepo.On("GetUserByEmail", ctx, newEmail).Return(repository.User{}, sql.ErrNoRows)
		mockRepo.On("UpdateUser", ctx, mock.AnythingOfType("repository.UpdateUserParams")).Return(nil)

		// 执行测试
		err := service.UpdateUser(ctx, input)

		// 断言
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("用户不存在", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		userID := int64(999)
		newUsername := "newusername"

		input := UpdateUserInput{
			UserID:   userID,
			Username: &newUsername,
		}

		mockRepo.On("GetUserByID", ctx, userID).Return(repository.User{}, sql.ErrNoRows)

		// 执行测试
		err := service.UpdateUser(ctx, input)

		// 断言
		assert.Error(t, err)
		assert.Equal(t, ErrUserNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("邮箱已被占用", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		userID := int64(1)
		newEmail := "existing@example.com"

		currentUser := repository.User{
			ID:       userID,
			Username: "testuser",
			Email:    "test@example.com",
			Status:   1,
		}

		existingUser := repository.User{
			ID:    2, // 不同的用户ID
			Email: "existing@example.com",
		}

		input := UpdateUserInput{
			UserID: userID,
			Email:  &newEmail,
		}

		mockRepo.On("GetUserByID", ctx, userID).Return(currentUser, nil)
		mockRepo.On("GetUserByEmail", ctx, newEmail).Return(existingUser, nil)

		// 执行测试
		err := service.UpdateUser(ctx, input)

		// 断言
		assert.Error(t, err)
		assert.Equal(t, ErrUserExists, err)
		mockRepo.AssertExpectations(t)
	})
}

// TestUserService_ChangePassword 测试修改密码
func TestUserService_ChangePassword(t *testing.T) {
	ctx := context.Background()

	t.Run("成功修改密码", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		userID := int64(1)
		oldPassword := "password123"
		newPassword := "newpassword456"

		// 生成正确的 bcrypt 哈希
		hashedOldPasswordBytes, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		hashedOldPassword := string(hashedOldPasswordBytes)

		user := repository.User{
			ID:       userID,
			Password: hashedOldPassword,
			Status:   1,
		}

		input := ChangePasswordInput{
			UserID:      userID,
			OldPassword: oldPassword,
			NewPassword: newPassword,
		}

		mockRepo.On("GetUserByID", ctx, userID).Return(user, nil)
		mockRepo.On("UpdateUserPassword", ctx, userID, mock.AnythingOfType("string")).Return(nil)

		// 执行测试
		err := service.ChangePassword(ctx, input)

		// 断言
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("旧密码错误", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		userID := int64(1)
		wrongOldPassword := "wrongpassword"
		newPassword := "newpassword456"

		// 生成正确的 bcrypt 哈希（正确密码是 password123）
		hashedOldPasswordBytes, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		hashedOldPassword := string(hashedOldPasswordBytes)

		user := repository.User{
			ID:       userID,
			Password: hashedOldPassword,
			Status:   1,
		}

		input := ChangePasswordInput{
			UserID:      userID,
			OldPassword: wrongOldPassword,
			NewPassword: newPassword,
		}

		mockRepo.On("GetUserByID", ctx, userID).Return(user, nil)

		// 执行测试
		err := service.ChangePassword(ctx, input)

		// 断言
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidPassword, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("参数验证失败", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		testCases := []struct {
			name  string
			input ChangePasswordInput
		}{
			{"空旧密码", ChangePasswordInput{UserID: 1, OldPassword: "", NewPassword: "newpass"}},
			{"空新密码", ChangePasswordInput{UserID: 1, OldPassword: "oldpass", NewPassword: ""}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := service.ChangePassword(ctx, tc.input)
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidInput, err)
			})
		}
	})
}

// TestUserService_DeleteUser 测试删除用户
func TestUserService_DeleteUser(t *testing.T) {
	ctx := context.Background()

	t.Run("成功删除用户", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		userID := int64(1)
		user := repository.User{
			ID:     userID,
			Status: 1,
		}

		mockRepo.On("GetUserByID", ctx, userID).Return(user, nil)
		mockRepo.On("DeleteUser", ctx, userID).Return(nil)

		// 执行测试
		err := service.DeleteUser(ctx, userID)

		// 断言
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("用户不存在", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		userID := int64(999)

		mockRepo.On("GetUserByID", ctx, userID).Return(repository.User{}, sql.ErrNoRows)

		// 执行测试
		err := service.DeleteUser(ctx, userID)

		// 断言
		assert.Error(t, err)
		assert.Equal(t, ErrUserNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("删除失败", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		userID := int64(1)
		user := repository.User{
			ID:     userID,
			Status: 1,
		}

		dbError := errors.New("database error")

		mockRepo.On("GetUserByID", ctx, userID).Return(user, nil)
		mockRepo.On("DeleteUser", ctx, userID).Return(dbError)

		// 执行测试
		err := service.DeleteUser(ctx, userID)

		// 断言
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// TestUserService_ListUsers 测试用户列表
func TestUserService_ListUsers(t *testing.T) {
	ctx := context.Background()

	t.Run("成功获取用户列表", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		limit := int32(10)
		offset := int32(0)
		total := int64(50)

		expectedUsers := []repository.User{
			{ID: 1, Username: "user1", Email: "user1@example.com"},
			{ID: 2, Username: "user2", Email: "user2@example.com"},
			{ID: 3, Username: "user3", Email: "user3@example.com"},
		}

		mockRepo.On("ListUsers", ctx, limit, offset).Return(expectedUsers, nil)
		mockRepo.On("CountUsers", ctx).Return(total, nil)

		// 执行测试
		users, count, err := service.ListUsers(ctx, limit, offset)

		// 断言
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.Equal(t, total, count)
		assert.Len(t, users, 3)
		mockRepo.AssertExpectations(t)
	})

	t.Run("空列表", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		limit := int32(10)
		offset := int32(100)
		total := int64(50)

		mockRepo.On("ListUsers", ctx, limit, offset).Return([]repository.User{}, nil)
		mockRepo.On("CountUsers", ctx).Return(total, nil)

		// 执行测试
		users, count, err := service.ListUsers(ctx, limit, offset)

		// 断言
		assert.NoError(t, err)
		assert.Empty(t, users)
		assert.Equal(t, total, count)
		mockRepo.AssertExpectations(t)
	})
}
