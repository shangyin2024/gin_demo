package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"gin_demo/pkg/cache"
)

// BenchmarkGetUserByID 性能测试 - 通过 ID 获取用户
func BenchmarkGetUserByID(b *testing.B) {
	db, cleanup := setupBenchDB(b)
	defer cleanup()

	rdb, cleanupRedis := setupBenchRedis(b)
	defer cleanupRedis()

	repo := NewUserRepository(db, cache.NewManager(rdb))

	// 准备测试数据
	ctx := context.Background()
	result, err := repo.CreateUser(ctx, CreateUserParams{
		Username: "benchuser",
		Email:    "bench@example.com",
		Password: "hashedpassword",
		Avatar:   sql.NullString{String: "avatar.jpg", Valid: true},
	})
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	// 性能测试
	for i := 0; i < b.N; i++ {
		_, err := repo.GetUserByID(ctx, result.ID)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetUserByEmail 性能测试 - 通过 Email 获取用户
func BenchmarkGetUserByEmail(b *testing.B) {
	db, cleanup := setupBenchDB(b)
	defer cleanup()

	rdb, cleanupRedis := setupBenchRedis(b)
	defer cleanupRedis()

	repo := NewUserRepository(db, cache.NewManager(rdb))

	// 准备测试数据
	ctx := context.Background()
	email := "bench2@example.com"
	_, err := repo.CreateUser(ctx, CreateUserParams{
		Username: "benchuser2",
		Email:    email,
		Password: "hashedpassword",
		Avatar:   sql.NullString{String: "avatar.jpg", Valid: true},
	})
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	// 性能测试
	for i := 0; i < b.N; i++ {
		_, err := repo.GetUserByEmail(ctx, email)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkListUsers 性能测试 - 列出用户（分页）
func BenchmarkListUsers(b *testing.B) {
	db, cleanup := setupBenchDB(b)
	defer cleanup()

	rdb, cleanupRedis := setupBenchRedis(b)
	defer cleanupRedis()

	repo := NewUserRepository(db, cache.NewManager(rdb))

	// 准备批量测试数据
	ctx := context.Background()
	for i := 0; i < 100; i++ {
		username := fmt.Sprintf("benchuser%d", i)
		email := fmt.Sprintf("bench%d@example.com", i)
		_, err := repo.CreateUser(ctx, CreateUserParams{
			Username: username,
			Email:    email,
			Password: "hashedpassword",
			Avatar:   sql.NullString{String: "avatar.jpg", Valid: true},
		})
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	b.ReportAllocs()

	// 性能测试
	for i := 0; i < b.N; i++ {
		_, err := repo.ListUsers(ctx, 10, 0)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCreateUser 性能测试 - 创建用户
func BenchmarkCreateUser(b *testing.B) {
	db, cleanup := setupBenchDB(b)
	defer cleanup()

	rdb, cleanupRedis := setupBenchRedis(b)
	defer cleanupRedis()

	repo := NewUserRepository(db, cache.NewManager(rdb))
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	// 性能测试
	for i := 0; i < b.N; i++ {
		username := fmt.Sprintf("createuser%d", i)
		email := fmt.Sprintf("create%d@example.com", i)
		_, err := repo.CreateUser(ctx, CreateUserParams{
			Username: username,
			Email:    email,
			Password: "hashedpassword",
			Avatar:   sql.NullString{String: "avatar.jpg", Valid: true},
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUpdateUser 性能测试 - 更新用户
func BenchmarkUpdateUser(b *testing.B) {
	db, cleanup := setupBenchDB(b)
	defer cleanup()

	rdb, cleanupRedis := setupBenchRedis(b)
	defer cleanupRedis()

	repo := NewUserRepository(db, cache.NewManager(rdb))

	// 准备测试数据
	ctx := context.Background()
	user, err := repo.CreateUser(ctx, CreateUserParams{
		Username: "benchuser",
		Email:    "bench@example.com",
		Password: "hashedpassword",
		Avatar:   sql.NullString{String: "avatar.jpg", Valid: true},
	})
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	// 性能测试
	for i := 0; i < b.N; i++ {
		username := fmt.Sprintf("updateuser%d", i)
		email := fmt.Sprintf("update%d@example.com", i)
		err := repo.UpdateUser(ctx, UpdateUserParams{
			ID:       user.ID,
			Username: username,
			Email:    email,
			Avatar:   sql.NullString{String: "new-avatar.jpg", Valid: true},
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCountUsers 性能测试 - 统计用户总数
func BenchmarkCountUsers(b *testing.B) {
	db, cleanup := setupBenchDB(b)
	defer cleanup()

	rdb, cleanupRedis := setupBenchRedis(b)
	defer cleanupRedis()

	repo := NewUserRepository(db, cache.NewManager(rdb))

	// 准备测试数据
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		username := fmt.Sprintf("countuser%d", i)
		email := fmt.Sprintf("count%d@example.com", i)
		_, err := repo.CreateUser(ctx, CreateUserParams{
			Username: username,
			Email:    email,
			Password: "hashedpassword",
			Avatar:   sql.NullString{String: "avatar.jpg", Valid: true},
		})
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	b.ReportAllocs()

	// 性能测试
	for i := 0; i < b.N; i++ {
		_, err := repo.CountUsers(ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCacheHit 性能测试 - 缓存命中场景
func BenchmarkCacheHit(b *testing.B) {
	db, cleanup := setupBenchDB(b)
	defer cleanup()

	rdb, cleanupRedis := setupBenchRedis(b)
	defer cleanupRedis()

	repo := NewUserRepository(db, cache.NewManager(rdb))

	// 准备测试数据
	ctx := context.Background()
	user, err := repo.CreateUser(ctx, CreateUserParams{
		Username: "cacheuser",
		Email:    "cache@example.com",
		Password: "hashedpassword",
		Avatar:   sql.NullString{String: "avatar.jpg", Valid: true},
	})
	if err != nil {
		b.Fatal(err)
	}

	// 预热缓存
	_, _ = repo.GetUserByID(ctx, user.ID)

	b.ResetTimer()
	b.ReportAllocs()

	// 性能测试（应该全部命中缓存）
	for i := 0; i < b.N; i++ {
		_, err := repo.GetUserByID(ctx, user.ID)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCacheMiss 性能测试 - 缓存未命中场景
func BenchmarkCacheMiss(b *testing.B) {
	db, cleanup := setupBenchDB(b)
	defer cleanup()

	rdb, cleanupRedis := setupBenchRedis(b)
	defer cleanupRedis()

	repo := NewUserRepository(db, cache.NewManager(rdb))

	// 准备测试数据
	ctx := context.Background()
	var userIDs []int64
	for i := 0; i < 100; i++ {
		username := fmt.Sprintf("missuser%d", i)
		email := fmt.Sprintf("miss%d@example.com", i)
		user, err := repo.CreateUser(ctx, CreateUserParams{
			Username: username,
			Email:    email,
			Password: "hashedpassword",
			Avatar:   sql.NullString{String: "avatar.jpg", Valid: true},
		})
		if err != nil {
			b.Fatal(err)
		}
		userIDs = append(userIDs, user.ID)
	}

	b.ResetTimer()
	b.ReportAllocs()

	// 性能测试（每次查询不同 ID，缓存未命中）
	for i := 0; i < b.N; i++ {
		userID := userIDs[i%len(userIDs)]
		// 先清除缓存
		_ = rdb.Del(ctx, cache.NewManager(rdb).BuildKey("user", userID))
		// 查询（缓存未命中）
		_, err := repo.GetUserByID(ctx, userID)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// 注意：setupBenchDB 和 setupBenchRedis 函数已在 user_repository_test.go 中定义
// 这里不需要重复定义
