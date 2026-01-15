package repository

import (
	"context"
	"database/sql"
)

// UserRepositoryInterface 用户仓库接口（用于依赖注入和测试）
type UserRepositoryInterface interface {
	// ========================================
	// 查询方法
	// ========================================
	
	// GetUserByID 通过 ID 查询用户
	GetUserByID(ctx context.Context, userID int64) (User, error)

	// GetUserByEmail 通过 Email 查询用户（包含密码）
	GetUserByEmail(ctx context.Context, email string) (User, error)

	// GetUserByUsername 通过 Username 查询用户
	GetUserByUsername(ctx context.Context, username string) (User, error)

	// ListUsers 列出用户（分页）
	ListUsers(ctx context.Context, limit, offset int32) ([]User, error)

	// CountUsers 统计用户总数
	CountUsers(ctx context.Context) (int64, error)

	// ========================================
	// 写操作方法
	// ========================================
	
	// CreateUser 创建用户
	CreateUser(ctx context.Context, params CreateUserParams) (User, error)

	// UpdateUser 更新用户信息
	UpdateUser(ctx context.Context, params UpdateUserParams) error

	// UpdateUserPassword 更新用户密码
	UpdateUserPassword(ctx context.Context, userID int64, password string) error

	// DeleteUser 删除用户
	DeleteUser(ctx context.Context, userID int64) error

	// ========================================
	// 事务方法
	// ========================================
	
	// WithTx 在事务中执行
	WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error

	// BatchExecInTx 批量操作（在一个事务中执行多个写操作）
	BatchExecInTx(ctx context.Context, ops []func(ctx context.Context, tx *sql.Tx) error) error

	// WithTx 返回使用事务的 Repository（用于在事务中执行多个操作）
	WithTxRepo(tx *sql.Tx) UserRepositoryInterface
}

// 确保 UserRepository 实现了接口
var _ UserRepositoryInterface = (*UserRepository)(nil)
