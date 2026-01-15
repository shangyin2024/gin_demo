package repository

import (
	"context"
	"database/sql"
	"fmt"
	"gin_demo/pkg/cache"
	"strconv"
	"time"
)

// BaseRepository 基础仓库（提供通用的泛型方法）
type BaseRepository[T any] struct {
	db    *sql.DB
	cache *cache.Manager
}

// NewBaseRepository 创建基础仓库
func NewBaseRepository[T any](db *sql.DB, cacheManager *cache.Manager) *BaseRepository[T] {
	return &BaseRepository[T]{
		db:    db,
		cache: cacheManager,
	}
}

// ListWithPagination 通用的分页查询（泛型）
func (r *BaseRepository[T]) ListWithPagination(
	ctx context.Context,
	queryFn func(ctx context.Context) ([]T, error),
) ([]T, error) {
	return queryFn(ctx)
}

// CountWithCache 通用的计数查询（带缓存）
func (r *BaseRepository[T]) CountWithCache(
	ctx context.Context,
	entity string,
	ttl time.Duration,
	countFn func(ctx context.Context) (int64, error),
) (int64, error) {
	return cache.TakeByID(ctx, r.cache, entity, "count", ttl, countFn)
}

// GetByIDWithCache 通过 ID 获取数据（带缓存）
func (r *BaseRepository[T]) GetByIDWithCache(
	ctx context.Context,
	entity string,
	id any,
	ttl time.Duration,
	queryFn func(context.Context) (T, error),
) (T, error) {
	return cache.TakeByID(ctx, r.cache, entity, id, ttl, queryFn)
}

// GetByIndexWithCache 通过索引获取数据（带缓存）
func (r *BaseRepository[T]) GetByIndexWithCache(
	ctx context.Context,
	entity, field string,
	value any,
	ttl time.Duration,
	indexQueryFn func(context.Context) (int64, error),
	dataQueryFn func(context.Context, int64) (T, error),
) (T, error) {
	return cache.TakeByIndex(
		ctx, r.cache, entity, field, value, ttl,
		indexQueryFn,
		dataQueryFn,
		func(idStr string) (int64, error) {
			// 字符串转 int64
			return strconv.ParseInt(idStr, 10, 64)
		},
	)
}

// ExecWithCache 执行写操作并清理缓存
func (r *BaseRepository[T]) ExecWithCache(
	ctx context.Context,
	entity string,
	id any,
	execFn func(context.Context) error,
) error {
	return r.cache.ExecByID(ctx, entity, id, execFn)
}

// ExecWithIndexCache 执行写操作并清理索引缓存
func (r *BaseRepository[T]) ExecWithIndexCache(
	ctx context.Context,
	entity string,
	id any,
	indexes []string,
	execFn func(context.Context) error,
) error {
	return r.cache.ExecByIDWithIndexes(ctx, entity, id, indexes, execFn)
}

// ============================================================================
// 事务管理
// ============================================================================

// TxOptions 事务选项
type TxOptions struct {
	Isolation sql.IsolationLevel
	ReadOnly  bool
}

// DefaultTxOptions 默认事务选项
var DefaultTxOptions = &sql.TxOptions{
	Isolation: sql.LevelDefault,
	ReadOnly:  false,
}

// ReadOnlyTxOptions 只读事务选项（用于查询）
var ReadOnlyTxOptions = &sql.TxOptions{
	Isolation: sql.LevelReadCommitted,
	ReadOnly:  true,
}

// WithTx 在事务中执行（使用默认选项）
func (r *BaseRepository[T]) WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	return r.WithTxOptions(ctx, nil, fn)
}

// WithTxOptions 在事务中执行（自定义选项）
func (r *BaseRepository[T]) WithTxOptions(ctx context.Context, opts *sql.TxOptions, fn func(tx *sql.Tx) error) error {
	tx, err := r.db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	
	// 使用 defer 确保事务回滚（如果未提交）
	defer func() {
		if p := recover(); p != nil {
			// 发生 panic 时回滚
			_ = tx.Rollback()
			panic(p) // 重新抛出 panic
		} else if err != nil {
			// 发生错误时回滚
			_ = tx.Rollback()
		}
	}()

	// 执行事务函数
	err = fn(tx)
	if err != nil {
		return err
	}

	// 提交事务
	return tx.Commit()
}

// WithReadOnlyTx 在只读事务中执行（用于需要一致性读的查询）
func (r *BaseRepository[T]) WithReadOnlyTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	return r.WithTxOptions(ctx, ReadOnlyTxOptions, fn)
}

// ExecInTx 在事务中执行写操作并清理缓存
func (r *BaseRepository[T]) ExecInTx(
	ctx context.Context,
	entity string,
	id any,
	fn func(ctx context.Context, tx *sql.Tx) error,
) error {
	return r.WithTx(ctx, func(tx *sql.Tx) error {
		// 执行数据库操作
		if err := fn(ctx, tx); err != nil {
			return err
		}
		
		// 清理缓存（在事务提交前）
		return r.cache.ExecByID(ctx, entity, id, func(context.Context) error {
			return nil // 只删除缓存，不执行额外操作
		})
	})
}

// BatchExecInTx 批量操作（在一个事务中执行多个写操作）
func (r *BaseRepository[T]) BatchExecInTx(
	ctx context.Context,
	ops []func(ctx context.Context, tx *sql.Tx) error,
) error {
	return r.WithTx(ctx, func(tx *sql.Tx) error {
		for i, op := range ops {
			if err := op(ctx, tx); err != nil {
				return fmt.Errorf("batch operation %d failed: %w", i, err)
			}
		}
		return nil
	})
}

// ============================================================================
// 辅助方法
// ============================================================================

// Cache 获取缓存管理器
func (r *BaseRepository[T]) Cache() *cache.Manager {
	return r.cache
}

// DB 获取数据库连接
func (r *BaseRepository[T]) DB() *sql.DB {
	return r.db
}
