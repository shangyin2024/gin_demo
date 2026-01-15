package wire

import (
	"database/sql"
	"gin_demo/internal/config"
	internalHealth "gin_demo/internal/health"
	"gin_demo/pkg/auth"
	"gin_demo/pkg/cache"
	"gin_demo/pkg/database"
	"gin_demo/pkg/health"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

// InfrastructureSet 基础设施层 Provider 集合
var InfrastructureSet = wire.NewSet(
	provideDatabase,
	provideRedis,
	provideCacheManager,
	provideJWTManager,
	provideHealthChecker,
)

// provideDatabase 提供数据库连接（支持 MySQL 和 PostgreSQL）
func provideDatabase(cfg *config.Config) (*sql.DB, error) {
	dbConfig := database.Config{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		DBName:          cfg.Database.DBName,
		SSLMode:         cfg.Database.SSLMode,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: int(cfg.Database.ConnMaxLifetime.Seconds()),
		ConnMaxIdleTime: int(cfg.Database.ConnMaxIdleTime.Seconds()),
	}

	switch cfg.Database.Driver {
	case "mysql":
		return database.NewMySQL(dbConfig)
	case "postgres", "postgresql":
		return database.NewPostgres(dbConfig)
	default:
		return database.NewMySQL(dbConfig) // 默认使用 MySQL
	}
}

// provideRedis 提供 Redis 连接（支持哨兵模式）
func provideRedis(cfg *config.Config) redis.UniversalClient {
	// 哨兵模式
	if cfg.Redis.SentinelEnabled {
		return redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    cfg.Redis.SentinelMaster,
			SentinelAddrs: cfg.Redis.SentinelAddrs,
			Password:      cfg.Redis.Password,
			DB:            cfg.Redis.DB,
			MaxRetries:    cfg.Redis.MaxRetries,
			PoolSize:      cfg.Redis.PoolSize,
			MinIdleConns:  cfg.Redis.MinIdleConns,
		})
	}
	
	// 单机模式
	return redis.NewClient(&redis.Options{
		Addr:         cfg.GetRedisAddr(),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		MaxRetries:   cfg.Redis.MaxRetries,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
	})
}

// provideCacheManager 提供缓存管理器
func provideCacheManager(rdb redis.UniversalClient) *cache.Manager {
	return cache.NewManager(rdb)
}

// provideJWTManager 提供 JWT 管理器（默认 int64 类型）
func provideJWTManager(cfg *config.Config) *auth.DefaultJWTManager {
	return auth.NewDefaultJWTManager(cfg.JWT.Secret, cfg.JWT.Expiration)
}

// provideHealthChecker 提供健康检查器
func provideHealthChecker(db *sql.DB, rdb redis.UniversalClient) health.Checker {
	// 创建组件检查器
	dbChecker := internalHealth.NewDatabaseChecker(db)
	redisChecker := internalHealth.NewRedisChecker(rdb)

	// 创建多组件检查器
	return health.NewMultiChecker("3.0.0", dbChecker, redisChecker)
}
