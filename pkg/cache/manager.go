package cache

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"gin_demo/pkg/metrics"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

var sfGroup singleflight.Group

const (
	NotFoundPlaceholder = "*"             // 数据库不存在记录时的占位符
	DefaultNotFoundTTL  = 5 * time.Minute // 占位符的默认过期时间
)

type Manager struct {
	rdb redis.UniversalClient
}

func NewManager(rdb redis.UniversalClient) *Manager {
	return &Manager{rdb: rdb}
}

// ----------------------------------------------------------------------------
// Key 构造工具
// ----------------------------------------------------------------------------

// BuildKey 构造主键缓存 Key: cache:user:1
func (m *Manager) BuildKey(entity string, id any) string {
	return fmt.Sprintf("cache:%s:%v", entity, id)
}

// BuildIndexKey 构造索引缓存 Key: cache:user:email:abc@example.com
func (m *Manager) BuildIndexKey(entity, field string, value any) string {
	return fmt.Sprintf("cache:%s:%s:%v", entity, field, value)
}

// getJitterTTL 在基础时间上增加随机扰动，防止缓存雪崩
func (m *Manager) getJitterTTL(baseTTL time.Duration) time.Duration {
	if baseTTL <= 0 {
		return baseTTL
	}
	// 20% 范围波动 + 0~30秒噪声
	jitter := rand.Int63n(int64(baseTTL) / 5)
	noise := time.Duration(rand.Int63n(30)) * time.Second
	return baseTTL + time.Duration(jitter) + noise
}

// ----------------------------------------------------------------------------
// 核心方法 1：TakeByID (通过主键获取)
// ----------------------------------------------------------------------------

func TakeByID[T any](ctx context.Context, m *Manager, entity string, id any, baseTTL time.Duration, queryFn func(context.Context) (T, error)) (T, error) {
	key := m.BuildKey(entity, id)
	var data T

	// 1. 查缓存
	val, err := m.rdb.Get(ctx, key).Result()
	if err == nil {
		// 缓存命中
		metrics.RecordCacheHit(entity)
		
		if val == NotFoundPlaceholder {
			return data, sql.ErrNoRows
		}
		if err := json.Unmarshal([]byte(val), &data); err == nil {
			return data, nil
		}
		_ = m.rdb.Del(ctx, key) // 数据损坏则删除
		metrics.RecordCacheError("get", "deserialization_error")
	} else {
		// 缓存未命中
		metrics.RecordCacheMiss(entity)
	}

	// 2. 缓存未命中，使用 singleflight 防击穿
	raw, err, _ := sfGroup.Do(key, func() (any, error) {
		// Double check
		if v, e := m.rdb.Get(ctx, key).Result(); e == nil {
			return v, nil
		}

		res, dbErr := queryFn(ctx)
		if dbErr != nil {
			if errors.Is(dbErr, sql.ErrNoRows) {
				_ = m.rdb.Set(ctx, key, NotFoundPlaceholder, DefaultNotFoundTTL).Err()
			}
			return nil, dbErr
		}

		bs, _ := json.Marshal(res)
		err := m.rdb.Set(ctx, key, string(bs), m.getJitterTTL(baseTTL)).Err()
		if err == nil {
			metrics.RecordCacheSet(entity)
		} else {
			metrics.RecordCacheError("set", "write_error")
		}
		return string(bs), nil
	})

	if err != nil {
		return data, err
	}

	s := raw.(string)
	if s == NotFoundPlaceholder {
		return data, sql.ErrNoRows
	}
	_ = json.Unmarshal([]byte(s), &data)
	return data, nil
}

// ----------------------------------------------------------------------------
// 核心方法 2：TakeByIndex (通过索引回填 ID 模式)
// ----------------------------------------------------------------------------

func TakeByIndex[T any, ID any](ctx context.Context, m *Manager, entity, field string, value any, baseTTL time.Duration,
	indexQueryFn func(context.Context) (ID, error),
	dataQueryFn func(context.Context, ID) (T, error),
	idConverter func(string) (ID, error)) (T, error) {

	indexKey := m.BuildIndexKey(entity, field, value)
	var data T

	// 1. 尝试获取 ID 映射
	idRaw, err, _ := sfGroup.Do(indexKey, func() (any, error) {
		if v, e := m.rdb.Get(ctx, indexKey).Result(); e == nil {
			return v, nil
		}

		id, dbErr := indexQueryFn(ctx)
		if dbErr != nil {
			if errors.Is(dbErr, sql.ErrNoRows) {
				_ = m.rdb.Set(ctx, indexKey, NotFoundPlaceholder, DefaultNotFoundTTL).Err()
			}
			return nil, dbErr
		}

		idStr := fmt.Sprintf("%v", id)
		_ = m.rdb.Set(ctx, indexKey, idStr, m.getJitterTTL(baseTTL*2)).Err()
		return idStr, nil
	})

	if err != nil {
		return data, err
	}

	idStr := idRaw.(string)
	if idStr == NotFoundPlaceholder {
		return data, sql.ErrNoRows
	}

	id, convertErr := idConverter(idStr)
	if convertErr != nil {
		return data, fmt.Errorf("ID conversion failed: %w", convertErr)
	}

	// 2. 拿到 ID 后走主键缓存逻辑
	return TakeByID(ctx, m, entity, id, baseTTL, func(ctx context.Context) (T, error) {
		return dataQueryFn(ctx, id)
	})
}

// ----------------------------------------------------------------------------
// 核心方法 3：Exec (写操作清理)
// ----------------------------------------------------------------------------

// ExecByID 清理主键缓存
func (m *Manager) ExecByID(ctx context.Context, entity string, id any, execFn func(context.Context) error) error {
	if err := execFn(ctx); err != nil {
		return err
	}
	err := m.rdb.Del(ctx, m.BuildKey(entity, id)).Err()
	if err == nil {
		metrics.RecordCacheDelete(entity)
	} else {
		metrics.RecordCacheError("delete", "delete_error")
	}
	return err
}

// ExecByIDWithIndexes 清理主键及相关索引缓存
func (m *Manager) ExecByIDWithIndexes(ctx context.Context, entity string, id any, indexes []string, execFn func(context.Context) error) error {
	if err := execFn(ctx); err != nil {
		return err
	}
	keys := append([]string{m.BuildKey(entity, id)}, indexes...)
	err := m.rdb.Del(ctx, keys...).Err()
	if err == nil {
		// 记录批量删除
		for range keys {
			metrics.RecordCacheDelete(entity)
		}
	} else {
		metrics.RecordCacheError("delete", "batch_delete_error")
	}
	return err
}
