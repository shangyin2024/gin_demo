package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"gin_demo/internal/repository"
	"gin_demo/internal/response"
	"gin_demo/pkg/metrics"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = response.ErrNotFound
	// ErrUserExists 用户已存在
	ErrUserExists = response.ErrAlreadyExists
	// ErrInvalidPassword 密码错误
	ErrInvalidPassword = response.ErrInvalidPassword
	// ErrInvalidInput 输入参数错误
	ErrInvalidInput = response.ErrInvalidParams
)

// RegisterInput 注册输入参数
type RegisterInput struct {
	Username string
	Email    string
	Password string
}

// LoginInput 登录输入参数
type LoginInput struct {
	Email    string
	Password string
}

// UpdateUserInput 更新用户输入参数
type UpdateUserInput struct {
	UserID   int64
	Username *string // 使用指针表示可选字段
	Email    *string
	Avatar   *string
}

// ChangePasswordInput 修改密码输入参数
type ChangePasswordInput struct {
	UserID      int64
	OldPassword string
	NewPassword string
}

// BatchUpdateInput 批量更新输入参数
type BatchUpdateInput struct {
	UserID   int64
	Username *string
	Email    *string
	Avatar   *string
}

// UserService 用户业务逻辑接口
type UserService interface {
	// Register 用户注册
	Register(ctx context.Context, input RegisterInput) (repository.User, error)

	// Login 用户登录
	Login(ctx context.Context, input LoginInput) (repository.User, error)

	// GetUserByID 通过 ID 获取用户
	GetUserByID(ctx context.Context, userID int64) (repository.User, error)

	// GetUserByEmail 通过 Email 获取用户
	GetUserByEmail(ctx context.Context, email string) (repository.User, error)

	// UpdateUser 更新用户信息
	UpdateUser(ctx context.Context, input UpdateUserInput) error

	// ChangePassword 修改密码
	ChangePassword(ctx context.Context, input ChangePasswordInput) error

	// DeleteUser 删除用户
	DeleteUser(ctx context.Context, userID int64) error

	// ListUsers 用户列表（分页）
	ListUsers(ctx context.Context, limit, offset int32) ([]repository.User, int64, error)

	// TransferUserData 转移用户数据（示例：需要事务的复杂操作）
	TransferUserData(ctx context.Context, fromUserID, toUserID int64) error

	// BatchUpdateUsers 批量更新用户（示例：批量操作事务）
	BatchUpdateUsers(ctx context.Context, updates []BatchUpdateInput) error
}

// userService 用户业务逻辑实现
type userService struct {
	userRepo repository.UserRepositoryInterface
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repository.UserRepositoryInterface) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// Register 用户注册
func (s *userService) Register(ctx context.Context, input RegisterInput) (repository.User, error) {
	var user repository.User

	// 1. 验证输入
	if input.Username == "" || input.Email == "" || input.Password == "" {
		slog.WarnContext(ctx, "Invalid registration input",
			"has_username", input.Username != "",
			"has_email", input.Email != "",
			"has_password", input.Password != "",
		)
		return user, ErrInvalidInput
	}

	// 2. 检查 Email 是否已存在
	existingUser, err := s.userRepo.GetUserByEmail(ctx, input.Email)
	if err == nil && existingUser.ID > 0 {
		slog.WarnContext(ctx, "Registration failed: email already exists",
			"email", input.Email,
			"existing_user_id", existingUser.ID,
		)
		return user, ErrUserExists
	}

	// 3. 检查 Username 是否已存在
	existingUser, err = s.userRepo.GetUserByUsername(ctx, input.Username)
	if err == nil && existingUser.ID > 0 {
		slog.WarnContext(ctx, "Registration failed: username already exists",
			"username", input.Username,
			"existing_user_id", existingUser.ID,
		)
		return user, ErrUserExists
	}

	// 4. 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to hash password",
			"error", err,
			"username", input.Username,
		)
		return user, fmt.Errorf("service: hash password: %w", err)
	}

	// 5. 创建用户
	user, err = s.userRepo.CreateUser(ctx, repository.CreateUserParams{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		Avatar:   sql.NullString{Valid: false},
	})
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create user",
			"error", err,
			"username", input.Username,
			"email", input.Email,
		)
		metrics.RecordUserOperation("register", false)
		return user, fmt.Errorf("service: create user: %w", err)
	}

	// 记录成功日志和指标
	slog.InfoContext(ctx, "User registered successfully",
		"user_id", user.ID,
		"username", user.Username,
		"email", user.Email,
	)
	metrics.RecordUserRegistration()
	metrics.RecordUserOperation("register", true)

	return user, nil
}

// Login 用户登录
func (s *userService) Login(ctx context.Context, input LoginInput) (repository.User, error) {
	var user repository.User

	// 1. 验证输入
	if input.Email == "" || input.Password == "" {
		slog.WarnContext(ctx, "Invalid login input",
			"has_email", input.Email != "",
			"has_password", input.Password != "",
		)
		return user, ErrInvalidInput
	}

	// 2. 查询用户
	user, err := s.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.WarnContext(ctx, "Login failed: user not found",
				"email", input.Email,
			)
			return user, ErrUserNotFound
		}
		slog.ErrorContext(ctx, "Failed to get user",
			"error", err,
			"email", input.Email,
		)
		return user, fmt.Errorf("service: get user: %w", err)
	}

	// 3. 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		slog.WarnContext(ctx, "Login failed: invalid password",
			"user_id", user.ID,
			"email", input.Email,
		)
		metrics.RecordUserLogin(false)
		return user, ErrInvalidPassword
	}

	// 记录登录成功
	slog.InfoContext(ctx, "User logged in successfully",
		"user_id", user.ID,
		"email", user.Email,
		"username", user.Username,
	)
	metrics.RecordUserLogin(true)

	return user, nil
}

// GetUserByID 通过 ID 获取用户
func (s *userService) GetUserByID(ctx context.Context, userID int64) (repository.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrUserNotFound
		}
		return user, fmt.Errorf("service: get user: %w", err)
	}
	return user, nil
}

// GetUserByEmail 通过 Email 获取用户
func (s *userService) GetUserByEmail(ctx context.Context, email string) (repository.User, error) {
	if email == "" {
		return repository.User{}, ErrInvalidInput
	}

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrUserNotFound
		}
		return user, fmt.Errorf("service: get user: %w", err)
	}
	return user, nil
}

// UpdateUser 更新用户信息
func (s *userService) UpdateUser(ctx context.Context, input UpdateUserInput) error {
	// 1. 检查用户是否存在
	currentUser, err := s.userRepo.GetUserByID(ctx, input.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("service: get user: %w", err)
	}

	// 2. 准备更新参数（使用当前值作为默认值）
	username := currentUser.Username
	if input.Username != nil && *input.Username != "" {
		username = *input.Username
	}

	email := currentUser.Email
	if input.Email != nil && *input.Email != "" {
		email = *input.Email
		// 检查新 Email 是否已被占用（排除当前用户）
		existingUser, err := s.userRepo.GetUserByEmail(ctx, email)
		if err == nil && existingUser.ID != input.UserID {
			return ErrUserExists
		}
	}

	avatar := currentUser.Avatar
	if input.Avatar != nil {
		if *input.Avatar != "" {
			avatar = sql.NullString{String: *input.Avatar, Valid: true}
		} else {
			avatar = sql.NullString{Valid: false}
		}
	}

	// 3. 更新用户
	err = s.userRepo.UpdateUser(ctx, repository.UpdateUserParams{
		ID:       input.UserID,
		Username: username,
		Email:    email,
		Avatar:   avatar,
	})
	if err != nil {
		metrics.RecordUserOperation("update", false)
		return fmt.Errorf("service: update user: %w", err)
	}

	// 记录更新成功指标
	metrics.RecordUserOperation("update", true)

	return nil
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(ctx context.Context, input ChangePasswordInput) error {
	// 1. 验证输入
	if input.OldPassword == "" || input.NewPassword == "" {
		return ErrInvalidInput
	}

	// 2. 获取用户
	user, err := s.userRepo.GetUserByID(ctx, input.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("service: get user: %w", err)
	}

	// 3. 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		return ErrInvalidPassword
	}

	// 4. 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("service: hash password: %w", err)
	}

	// 5. 更新密码
	if err := s.userRepo.UpdateUserPassword(ctx, input.UserID, string(hashedPassword)); err != nil {
		return fmt.Errorf("service: update password: %w", err)
	}

	// 记录密码修改指标
	metrics.PasswordChanges.Inc()

	return nil
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(ctx context.Context, userID int64) error {
	// 1. 检查用户是否存在
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.WarnContext(ctx, "Delete failed: user not found",
				"user_id", userID,
			)
			return ErrUserNotFound
		}
		slog.ErrorContext(ctx, "Failed to get user for deletion",
			"error", err,
			"user_id", userID,
		)
		return fmt.Errorf("service: get user: %w", err)
	}

	// 2. 删除用户
	if err := s.userRepo.DeleteUser(ctx, userID); err != nil {
		slog.ErrorContext(ctx, "Failed to delete user",
			"error", err,
			"user_id", userID,
			"username", user.Username,
			"email", user.Email,
		)
		metrics.RecordUserOperation("delete", false)
		return fmt.Errorf("service: delete user: %w", err)
	}

	// 记录删除成功
	slog.InfoContext(ctx, "User deleted successfully",
		"user_id", userID,
		"username", user.Username,
		"email", user.Email,
	)
	metrics.UserDeletions.Inc()
	metrics.RecordUserOperation("delete", true)

	return nil
}

// ListUsers 用户列表（分页）
func (s *userService) ListUsers(ctx context.Context, limit, offset int32) ([]repository.User, int64, error) {
	// 1. 查询用户列表
	users, err := s.userRepo.ListUsers(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("service: list users: %w", err)
	}

	// 2. 查询总数
	total, err := s.userRepo.CountUsers(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("service: count users: %w", err)
	}

	return users, total, nil
}

// ============================================================================
// 事务方法示例
// ============================================================================

// TransferUserData 转移用户数据（示例：复杂的多步骤事务操作）
//
// 场景：将一个用户的数据转移到另一个用户（例如合并账户）
// 需要在一个事务中完成多个操作，确保数据一致性
func (s *userService) TransferUserData(ctx context.Context, fromUserID, toUserID int64) error {
	// 1. 验证源用户和目标用户都存在
	fromUser, err := s.userRepo.GetUserByID(ctx, fromUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("service: get source user: %w", err)
	}

	toUser, err := s.userRepo.GetUserByID(ctx, toUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("service: get target user: %w", err)
	}

	// 2. 在事务中执行数据转移
	return s.userRepo.WithTx(ctx, func(tx *sql.Tx) error {
		// 2.1 这里可以执行多个数据库操作
		// 例如：转移用户的订单、评论、收藏等数据到目标用户
		// txRepo := s.userRepo.WithTx(tx)
		
		// 示例：更新源用户状态为已合并
		// if err := txRepo.UpdateUserStatus(ctx, fromUserID, StatusMerged); err != nil {
		//     return fmt.Errorf("failed to update source user status: %w", err)
		// }

		// 示例：更新目标用户的最后更新时间
		// if err := txRepo.TouchUser(ctx, toUserID); err != nil {
		//     return fmt.Errorf("failed to update target user: %w", err)
		// }

		// 2.2 记录日志（示例）
		slog.InfoContext(ctx, "User data transferred successfully",
			"from_user_id", fromUser.ID,
			"from_email", fromUser.Email,
			"to_user_id", toUser.ID,
			"to_email", toUser.Email,
		)

		return nil
	})
}

// BatchUpdateUsers 批量更新用户（示例：批量操作事务）
//
// 场景：批量更新多个用户的信息
// 使用事务确保要么全部成功，要么全部失败
func (s *userService) BatchUpdateUsers(ctx context.Context, updates []BatchUpdateInput) error {
	if len(updates) == 0 {
		return ErrInvalidInput
	}

	// 使用批量事务操作
	ops := make([]func(ctx context.Context, tx *sql.Tx) error, 0, len(updates))

	for _, update := range updates {
		// 捕获循环变量
		u := update

		ops = append(ops, func(ctx context.Context, tx *sql.Tx) error {
			// 获取当前用户信息
			currentUser, err := s.userRepo.GetUserByID(ctx, u.UserID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return fmt.Errorf("user %d not found: %w", u.UserID, ErrUserNotFound)
				}
				return fmt.Errorf("failed to get user %d: %w", u.UserID, err)
			}

			// 准备更新参数
			username := currentUser.Username
			if u.Username != nil && *u.Username != "" {
				username = *u.Username
			}

			email := currentUser.Email
			if u.Email != nil && *u.Email != "" {
				email = *u.Email
			}

			avatar := currentUser.Avatar
			if u.Avatar != nil {
				if *u.Avatar != "" {
					avatar = sql.NullString{String: *u.Avatar, Valid: true}
				} else {
					avatar = sql.NullString{Valid: false}
				}
			}

			// 使用事务中的 Repository 执行更新
			txRepo := s.userRepo.WithTxRepo(tx)
			return txRepo.UpdateUser(ctx, repository.UpdateUserParams{
				ID:       u.UserID,
				Username: username,
				Email:    email,
				Avatar:   avatar,
			})
		})
	}

	// 在一个事务中执行所有操作
	return s.userRepo.BatchExecInTx(ctx, ops)
}
