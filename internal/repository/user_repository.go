package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gin_demo/pkg/cache"
	dbContext "gin_demo/pkg/database"
)

// UserRepository 用户仓库层（结合缓存）
type UserRepository struct {
	*BaseRepository[User]
	queries *Queries
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *sql.DB, cacheManager *cache.Manager) *UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository[User](db, cacheManager),
		queries:        New(db),
	}
}

// ============================================================================
// 查询方法（带缓存）
// ============================================================================

// GetUserByID 通过 ID 查询用户（主键缓存）
func (r *UserRepository) GetUserByID(ctx context.Context, userID int64) (User, error) {
	return r.GetByIDWithCache(ctx, "user", userID, 5*time.Minute,
		func(ctx context.Context) (User, error) {
			row, err := r.queries.GetUserByID(ctx, userID)
			if err != nil {
				return User{}, err
			}
			return r.rowToUser(row.ID, row.Username, row.Email, "", row.Avatar, row.Status, row.CreatedAt, row.UpdatedAt), nil
		})
}

// rowToUser 转换查询结果为 User
func (r *UserRepository) rowToUser(id int64, username, email, password string, avatar sql.NullString, status int16, createdAt, updatedAt time.Time) User {
	return User{
		ID:        id,
		Username:  username,
		Email:     email,
		Password:  password,
		Avatar:    avatar,
		Status:    status,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

// GetUserByEmail 通过 Email 查询用户（包含密码，用于登录）
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (User, error) {
	ctx, cancel := dbContext.WithQueryTimeout(ctx)
	defer cancel()

	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return User{}, err
	}
	return r.rowToUser(row.ID, row.Username, row.Email, row.Password, row.Avatar, row.Status, row.CreatedAt, row.UpdatedAt), nil
}

// GetUserByUsername 通过 Username 查询用户（索引缓存）
func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (User, error) {
	ctx, cancel := dbContext.WithQueryTimeout(ctx)
	defer cancel()

	row, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return User{}, err
	}
	return r.rowToUser(row.ID, row.Username, row.Email, "", row.Avatar, row.Status, row.CreatedAt, row.UpdatedAt), nil
}

// ListUsers 列出用户（分页，不缓存）
func (r *UserRepository) ListUsers(ctx context.Context, limit, offset int32) ([]User, error) {
	return r.ListWithPagination(ctx, func(ctx context.Context) ([]User, error) {
		rows, err := r.queries.ListUsers(ctx, ListUsersParams{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			return nil, err
		}

		users := make([]User, 0, len(rows))
		for _, row := range rows {
			users = append(users, r.rowToUser(row.ID, row.Username, row.Email, "", row.Avatar, row.Status, row.CreatedAt, row.UpdatedAt))
		}
		return users, nil
	})
}

// CountUsers 统计用户总数（短期缓存）
func (r *UserRepository) CountUsers(ctx context.Context) (int64, error) {
	return r.CountWithCache(ctx, "user:count", 1*time.Minute,
		func(ctx context.Context) (int64, error) {
			return r.queries.CountUsers(ctx)
		})
}

// ============================================================================
// 写操作（自动清理缓存）
// ============================================================================

// CreateUser 创建用户（清理统计缓存）
func (r *UserRepository) CreateUser(ctx context.Context, params CreateUserParams) (User, error) {
	ctx, cancel := dbContext.WithQueryTimeout(ctx)
	defer cancel()

	// MySQL: 执行创建操作，返回 sql.Result
	result, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return User{}, err
	}

	// 获取自动生成的 ID
	userID, err := result.LastInsertId()
	if err != nil {
		return User{}, fmt.Errorf("failed to get last insert id: %w", err)
	}

	// 清理统计缓存
	_ = r.ExecWithCache(ctx, "user:count", "total", func(context.Context) error {
		return nil // 只删除缓存，不执行 DB 操作
	})

	// 查询新创建的用户以获取完整信息
	return r.GetUserByID(ctx, userID)
}

// UpdateUser 更新用户信息（清理主键和索引缓存）
func (r *UserRepository) UpdateUser(ctx context.Context, params UpdateUserParams) error {
	// 先获取旧数据（用于清理旧索引）
	oldUser, err := r.queries.GetUserByID(ctx, params.ID)
	if err != nil {
		return fmt.Errorf("repository: get old user: %w", err)
	}

	// 构建需要清理的索引 Key
	indexes := []string{
		r.Cache().BuildIndexKey("user", "email", oldUser.Email),
		r.Cache().BuildIndexKey("user", "email", params.Email),
		r.Cache().BuildIndexKey("user", "username", oldUser.Username),
		r.Cache().BuildIndexKey("user", "username", params.Username),
	}

	return r.ExecWithIndexCache(ctx, "user", params.ID, indexes,
		func(ctx context.Context) error {
			return r.queries.UpdateUser(ctx, params)
		})
}

// UpdateUserPassword 更新用户密码（只清理主键缓存）
func (r *UserRepository) UpdateUserPassword(ctx context.Context, userID int64, password string) error {
	return r.ExecWithCache(ctx, "user", userID, func(ctx context.Context) error {
		return r.queries.UpdateUserPassword(ctx, UpdateUserPasswordParams{
			ID:       userID,
			Password: password,
		})
	})
}

// DeleteUser 删除用户（软删除，清理所有相关缓存）
func (r *UserRepository) DeleteUser(ctx context.Context, userID int64) error {
	// 先获取用户数据（用于清理索引）
	user, err := r.queries.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("repository: get user: %w", err)
	}

	indexes := []string{
		r.Cache().BuildIndexKey("user", "email", user.Email),
		r.Cache().BuildIndexKey("user", "username", user.Username),
		r.Cache().BuildKey("user:count", "total"),
	}

	return r.ExecWithIndexCache(ctx, "user", userID, indexes,
		func(ctx context.Context) error {
			return r.queries.DeleteUser(ctx, userID)
		})
}

// ============================================================================
// 事务支持
// ============================================================================

// ============================================================================
// 事务支持
// ============================================================================

// WithTxRepo 返回使用事务的 Repository（实现接口）
func (r *UserRepository) WithTxRepo(tx *sql.Tx) UserRepositoryInterface {
	return &UserRepository{
		BaseRepository: r.BaseRepository,
		queries:        r.queries.WithTx(tx),
	}
}
