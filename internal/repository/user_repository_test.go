package repository

import (
	"context"
	"database/sql"
	"testing"

	"gin_demo/pkg/cache"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 注意：这些是集成测试，需要真实的数据库和 Redis
// 运行前需要启动测试环境：docker-compose up -d

// setupTestDB 设置测试数据库
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	// 使用测试数据库配置
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=gin_demo_test sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Skipf("Skipping integration test: cannot connect to test database: %v", err)
		return nil, func() {}
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		t.Skipf("Skipping integration test: database not available: %v", err)
		return nil, func() {}
	}

	// 清理函数
	cleanup := func() {
		// 清理测试数据
		_, _ = db.Exec("DELETE FROM users WHERE email LIKE 'test%@example.com'")
		db.Close()
	}

	return db, cleanup
}

// setupTestRedis 设置测试 Redis
func setupTestRedis(t *testing.T) (*redis.Client, func()) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1, // 使用不同的数据库编号
	})

	// 测试连接
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skipf("Skipping integration test: redis not available: %v", err)
		return nil, func() {}
	}

	// 清理函数
	cleanup := func() {
		// 清理测试缓存
		_ = rdb.FlushDB(ctx).Err()
		_ = rdb.Close()
	}

	return rdb, cleanup
}

// TestUserRepository_Integration 集成测试
func TestUserRepository_Integration(t *testing.T) {
	// 跳过短测试
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// 设置测试环境
	db, dbCleanup := setupTestDB(t)
	if db == nil {
		return
	}
	defer dbCleanup()

	rdb, redisCleanup := setupTestRedis(t)
	if rdb == nil {
		return
	}
	defer redisCleanup()

	// 创建 Repository
	cacheManager := cache.NewManager(rdb)
	repo := NewUserRepository(db, cacheManager)
	ctx := context.Background()

	t.Run("创建和查询用户", func(t *testing.T) {
		// 创建用户
		user, err := repo.CreateUser(ctx, CreateUserParams{
			Username: "testuser1",
			Email:    "test1@example.com",
			Password: "hashed_password",
			Avatar:   sql.NullString{Valid: false},
		})
		require.NoError(t, err)
		assert.NotZero(t, user.ID)
		assert.Equal(t, "testuser1", user.Username)
		assert.Equal(t, "test1@example.com", user.Email)

		// 通过 ID 查询（第一次从数据库）
		user1, err := repo.GetUserByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.ID, user1.ID)
		assert.Equal(t, user.Username, user1.Username)

		// 再次查询（应该从缓存读取）
		user2, err := repo.GetUserByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.ID, user2.ID)

		// 通过 Email 查询
		user3, err := repo.GetUserByEmail(ctx, user.Email)
		require.NoError(t, err)
		assert.Equal(t, user.ID, user3.ID)

		// 通过 Username 查询
		user4, err := repo.GetUserByUsername(ctx, user.Username)
		require.NoError(t, err)
		assert.Equal(t, user.ID, user4.ID)
	})

	t.Run("更新用户并清理缓存", func(t *testing.T) {
		// 创建用户
		user, err := repo.CreateUser(ctx, CreateUserParams{
			Username: "testuser2",
			Email:    "test2@example.com",
			Password: "hashed_password",
			Avatar:   sql.NullString{Valid: false},
		})
		require.NoError(t, err)

		// 查询用户（缓存）
		_, err = repo.GetUserByID(ctx, user.ID)
		require.NoError(t, err)

		// 更新用户
		err = repo.UpdateUser(ctx, UpdateUserParams{
			ID:       user.ID,
			Username: "testuser2_updated",
			Email:    "test2_updated@example.com",
			Avatar:   sql.NullString{String: "https://example.com/avatar.jpg", Valid: true},
		})
		require.NoError(t, err)

		// 再次查询（缓存应该已清理）
		updatedUser, err := repo.GetUserByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, "testuser2_updated", updatedUser.Username)
		assert.Equal(t, "test2_updated@example.com", updatedUser.Email)
		assert.True(t, updatedUser.Avatar.Valid)
		assert.Equal(t, "https://example.com/avatar.jpg", updatedUser.Avatar.String)
	})

	t.Run("更新密码", func(t *testing.T) {
		// 创建用户
		user, err := repo.CreateUser(ctx, CreateUserParams{
			Username: "testuser3",
			Email:    "test3@example.com",
			Password: "old_password",
			Avatar:   sql.NullString{Valid: false},
		})
		require.NoError(t, err)

		// 更新密码
		err = repo.UpdateUserPassword(ctx, user.ID, "new_password")
		require.NoError(t, err)

		// 通过 Email 查询（包含密码）
		updatedUser, err := repo.GetUserByEmail(ctx, user.Email)
		require.NoError(t, err)
		assert.Equal(t, "new_password", updatedUser.Password)
	})

	t.Run("删除用户（软删除）", func(t *testing.T) {
		// 创建用户
		user, err := repo.CreateUser(ctx, CreateUserParams{
			Username: "testuser4",
			Email:    "test4@example.com",
			Password: "password",
			Avatar:   sql.NullString{Valid: false},
		})
		require.NoError(t, err)

		// 删除用户
		err = repo.DeleteUser(ctx, user.ID)
		require.NoError(t, err)

		// 查询应该返回 not found（软删除）
		_, err = repo.GetUserByID(ctx, user.ID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("用户列表和统计", func(t *testing.T) {
		// 创建多个用户
		for i := 0; i < 5; i++ {
			_, err := repo.CreateUser(ctx, CreateUserParams{
				Username: "listuser" + string(rune('0'+i)),
				Email:    "listuser" + string(rune('0'+i)) + "@example.com",
				Password: "password",
				Avatar:   sql.NullString{Valid: false},
			})
			require.NoError(t, err)
		}

		// 查询列表
		users, err := repo.ListUsers(ctx, 10, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 5)

		// 统计总数
		count, err := repo.CountUsers(ctx)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(5))

		// 再次统计（应该从缓存读取）
		count2, err := repo.CountUsers(ctx)
		require.NoError(t, err)
		assert.Equal(t, count, count2)
	})

	t.Run("缓存穿透防护", func(t *testing.T) {
		// 查询不存在的用户
		_, err := repo.GetUserByID(ctx, 999999)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)

		// 再次查询（应该从缓存返回占位符）
		_, err = repo.GetUserByID(ctx, 999999)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)

		// 验证缓存中存在占位符
		key := cacheManager.BuildKey("user", 999999)
		val, err := rdb.Get(ctx, key).Result()
		require.NoError(t, err)
		assert.Equal(t, cache.NotFoundPlaceholder, val)
	})

	t.Run("并发安全测试", func(t *testing.T) {
		// 创建测试用户
		user, err := repo.CreateUser(ctx, CreateUserParams{
			Username: "concurrent_user",
			Email:    "concurrent@example.com",
			Password: "password",
			Avatar:   sql.NullString{Valid: false},
		})
		require.NoError(t, err)

		// 清理缓存
		key := cacheManager.BuildKey("user", user.ID)
		_ = rdb.Del(ctx, key)

		// 并发查询同一用户（测试 singleflight）
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				_, err := repo.GetUserByID(ctx, user.ID)
				assert.NoError(t, err)
				done <- true
			}()
		}

		// 等待所有协程完成
		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

// BenchmarkUserRepository_GetUserByID 性能基准测试
func BenchmarkUserRepository_GetUserByID(b *testing.B) {
	// 设置测试环境
	db, dbCleanup := setupBenchDB(b)
	if db == nil {
		return
	}
	defer dbCleanup()

	rdb, redisCleanup := setupBenchRedis(b)
	if rdb == nil {
		return
	}
	defer redisCleanup()

	// 创建 Repository
	cacheManager := cache.NewManager(rdb)
	repo := NewUserRepository(db, cacheManager)
	ctx := context.Background()

	// 创建测试用户
	user, err := repo.CreateUser(ctx, CreateUserParams{
		Username: "benchuser",
		Email:    "bench@example.com",
		Password: "password",
		Avatar:   sql.NullString{Valid: false},
	})
	if err != nil {
		b.Fatalf("Failed to create test user: %v", err)
	}

	b.ResetTimer()

	// 基准测试
	b.Run("WithCache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = repo.GetUserByID(ctx, user.ID)
		}
	})

	b.Run("CacheMiss", func(b *testing.B) {
		key := cacheManager.BuildKey("user", user.ID)
		for i := 0; i < b.N; i++ {
			// 每次都清理缓存，模拟缓存未命中
			_ = rdb.Del(ctx, key)
			_, _ = repo.GetUserByID(ctx, user.ID)
		}
	})
}

func setupBenchDB(b *testing.B) (*sql.DB, func()) {
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=gin_demo_test sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		b.Skipf("Skipping benchmark: cannot connect to test database: %v", err)
		return nil, func() {}
	}

	if err := db.Ping(); err != nil {
		b.Skipf("Skipping benchmark: database not available: %v", err)
		return nil, func() {}
	}

	cleanup := func() {
		_, _ = db.Exec("DELETE FROM users WHERE email = 'bench@example.com'")
		db.Close()
	}

	return db, cleanup
}

func setupBenchRedis(b *testing.B) (*redis.Client, func()) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		b.Skipf("Skipping benchmark: redis not available: %v", err)
		return nil, func() {}
	}

	cleanup := func() {
		_ = rdb.FlushDB(ctx).Err()
		_ = rdb.Close()
	}

	return rdb, cleanup
}
