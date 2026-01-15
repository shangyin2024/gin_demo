package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"gin_demo/pkg/cache"

	"github.com/spf13/viper"
)

// Config 全局配置
type Config struct {
	// 服务配置
	Server ServerConfig

	// 数据库配置
	Database DatabaseConfig

	// Redis 配置
	Redis RedisConfig

	// 日志配置
	Logger LoggerConfig

	// JWT 配置
	JWT JWTConfig

	// CORS 配置
	CORS CORSConfig

	// 安全配置
	Security SecurityConfig

	// 缓存配置
	Cache CacheConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host               string
	Port               int
	Mode               string // debug, release, test
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	IdleTimeout        time.Duration
	MaxRequestBodySize int64 // 最大请求体大小（字节）
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string        // mysql 或 postgres
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string        // 仅 PostgreSQL 使用
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// RedisConfig Redis 配置
type RedisConfig struct {
	// 单机模式
	Host         string
	Port         int
	Password     string
	DB           int
	MaxRetries   int
	PoolSize     int
	MinIdleConns int
	
	// 哨兵模式
	SentinelEnabled bool     // 是否启用哨兵模式
	SentinelMaster  string   // 哨兵主节点名称
	SentinelAddrs   []string // 哨兵地址列表
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level        string // debug, info, warn, error
	Format       string // json, text
	IsJSON       bool   // 是否使用 JSON 格式
	AddSource    bool   // 是否添加源码位置
	RequestIDKey string // Request ID 的键名
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret     string        // JWT 密钥
	Expiration time.Duration // Token 过期时间
}

// CORSConfig CORS 跨域配置
type CORSConfig struct {
	AllowedOrigins   []string // 允许的来源
	AllowCredentials bool     // 是否允许携带凭证
	MaxAge           int      // 预检请求缓存时间（秒）
}

// CacheConfig 缓存配置（直接使用 pkg/cache 的配置）
type CacheConfig = cache.CacheConfig

// Load 加载配置（支持多环境）
func Load() (*Config, error) {
	// 1. 设置默认值
	setDefaults()

	// 2. 读取环境变量（支持覆盖配置文件）
	// 替换点号为下划线以兼容环境变量，例如 DATABASE_PASSWORD
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// 3. 确定环境（优先级: APP_ENV > ENV > 默认 dev）
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = os.Getenv("ENV")
	}
	if env == "" {
		env = "dev"
	}

	// 4. 加载配置文件（分层加载策略）
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// 4.1 首先尝试加载基础配置文件 config.yaml
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// 配置文件存在但读取失败
			return nil, fmt.Errorf("config: failed to read base config: %w", err)
		}
		// 基础配置文件不存在也可以（完全依赖环境变量）
		slog.Warn("Base config file not found, using defaults and environment variables only")
	} else {
		slog.Info("Loaded base configuration", "file", viper.ConfigFileUsed())
	}

	// 4.2 然后尝试合并环境特定配置文件 config.{env}.yaml
	if env != "" {
		envConfigName := fmt.Sprintf("config.%s", env)
		viper.SetConfigName(envConfigName)
		
		// 使用 MergeInConfig 而不是 ReadInConfig，以便覆盖基础配置
		if err := viper.MergeInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				// 环境配置文件存在但读取失败
				return nil, fmt.Errorf("config: failed to read environment config: %w", err)
			}
			// 环境特定配置文件不存在也可以（使用基础配置）
			slog.Debug("Environment-specific config file not found", "env", env, "file", envConfigName+".yaml")
		} else {
			slog.Info("Loaded environment-specific configuration", 
				"env", env, 
				"file", viper.ConfigFileUsed(),
			)
		}
	}

	// 5. 记录最终使用的环境
	slog.Info("Configuration loaded", 
		"env", env,
		"config_files", viper.ConfigFileUsed(),
	)

	// 6. 解析配置
	cfg := &Config{
		Server: ServerConfig{
			Host:               viper.GetString("server.host"),
			Port:               viper.GetInt("server.port"),
			Mode:               viper.GetString("server.mode"),
			ReadTimeout:        viper.GetDuration("server.read_timeout"),
			WriteTimeout:       viper.GetDuration("server.write_timeout"),
			IdleTimeout:        viper.GetDuration("server.idle_timeout"),
			MaxRequestBodySize: viper.GetInt64("server.max_request_body_size"),
		},
		Database: DatabaseConfig{
			Host:            viper.GetString("database.host"),
			Port:            viper.GetInt("database.port"),
			User:            viper.GetString("database.user"),
			Password:        viper.GetString("database.password"),
			DBName:          viper.GetString("database.dbname"),
			SSLMode:         viper.GetString("database.sslmode"),
			MaxOpenConns:    viper.GetInt("database.max_open_conns"),
			MaxIdleConns:    viper.GetInt("database.max_idle_conns"),
			ConnMaxLifetime: viper.GetDuration("database.conn_max_lifetime"),
			ConnMaxIdleTime: viper.GetDuration("database.conn_max_idle_time"),
		},
		Redis: RedisConfig{
			Host:         viper.GetString("redis.host"),
			Port:         viper.GetInt("redis.port"),
			Password:     viper.GetString("redis.password"),
			DB:           viper.GetInt("redis.db"),
			MaxRetries:   viper.GetInt("redis.max_retries"),
			PoolSize:     viper.GetInt("redis.pool_size"),
			MinIdleConns: viper.GetInt("redis.min_idle_conns"),
		},
		Logger: LoggerConfig{
			Level:        viper.GetString("logger.level"),
			Format:       viper.GetString("logger.format"),
			IsJSON:       viper.GetString("logger.format") == "json",
			AddSource:    viper.GetBool("logger.add_source"),
			RequestIDKey: viper.GetString("logger.request_id_key"),
		},
		JWT: JWTConfig{
			Secret:     viper.GetString("jwt.secret"),
			Expiration: viper.GetDuration("jwt.expiration"),
		},
		CORS: CORSConfig{
			AllowedOrigins:   viper.GetStringSlice("cors.allowed_origins"),
			AllowCredentials: viper.GetBool("cors.allow_credentials"),
			MaxAge:           viper.GetInt("cors.max_age"),
		},
		Security: SecurityConfig{
			Headers: SecurityHeadersConfig{
				Enabled:               viper.GetBool("security.headers.enabled"),
				EnableHSTS:            viper.GetBool("security.headers.enable_hsts"),
				HSTSMaxAge:            viper.GetInt("security.headers.hsts_max_age"),
				HSTSIncludeSubdomains: viper.GetBool("security.headers.hsts_include_subdomains"),
				EnableCSP:             viper.GetBool("security.headers.enable_csp"),
				CSPPolicy:             viper.GetString("security.headers.csp_policy"),
				EnableFrameOptions:    viper.GetBool("security.headers.enable_frame_options"),
				FrameOptions:          viper.GetString("security.headers.frame_options"),
			},
			EnableCompression: viper.GetBool("security.enable_compression"),
			CompressionLevel:  viper.GetInt("security.compression_level"),
			TLS: TLSConfig{
				Enabled:    viper.GetBool("security.tls.enabled"),
				CertFile:   viper.GetString("security.tls.cert_file"),
				KeyFile:    viper.GetString("security.tls.key_file"),
				MinVersion: viper.GetString("security.tls.min_version"),
			},
		},
		Cache: CacheConfig{
			DefaultTTL:     viper.GetDuration("cache.default_ttl"),
			UserTTL:        viper.GetDuration("cache.user_ttl"),
			UserIndexTTL:   viper.GetDuration("cache.user_index_ttl"),
			UserCountTTL:   viper.GetDuration("cache.user_count_ttl"),
			UserSessionTTL: viper.GetDuration("cache.user_session_ttl"),
			ContentTTL:     viper.GetDuration("cache.content_ttl"),
			ContentListTTL: viper.GetDuration("cache.content_list_ttl"),
			StatsTTL:       viper.GetDuration("cache.stats_ttl"),
			NotFoundTTL:    viper.GetDuration("cache.not_found_ttl"),
			EnableJitter:   viper.GetBool("cache.enable_jitter"),
			JitterPercent:  viper.GetInt("cache.jitter_percent"),
		},
	}

	// 7. 验证配置
	if err := cfg.Validate(env); err != nil {
		return nil, fmt.Errorf("config: validate: %w", err)
	}

	// 8. 记录关键配置（不包含敏感信息）
	slog.Debug("Configuration details",
		"server.mode", cfg.Server.Mode,
		"server.port", cfg.Server.Port,
		"database.host", cfg.Database.Host,
		"redis.host", cfg.Redis.Host,
		"cors.allowed_origins", cfg.CORS.AllowedOrigins,
		"security.headers.enabled", cfg.Security.Headers.Enabled,
		"security.tls.enabled", cfg.Security.TLS.Enabled,
	)

	return cfg, nil
}

// setDefaults 设置默认值
func setDefaults() {
	// 服务器默认值
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", 10*time.Second)
	viper.SetDefault("server.write_timeout", 10*time.Second)
	viper.SetDefault("server.idle_timeout", 60*time.Second)
	viper.SetDefault("server.max_request_body_size", int64(10*1024*1024)) // 10MB

	// 数据库默认值
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.dbname", "gin_demo")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", 5*time.Minute)
	viper.SetDefault("database.conn_max_idle_time", 10*time.Minute)

	// Redis 默认值
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.max_retries", 3)
	viper.SetDefault("redis.pool_size", 10)
	viper.SetDefault("redis.min_idle_conns", 5)

	// 日志默认值
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.format", "json")
	viper.SetDefault("logger.add_source", false)
	viper.SetDefault("logger.request_id_key", "request_id")

	// JWT 默认值
	viper.SetDefault("jwt.secret", "your-secret-key-change-in-production")
	viper.SetDefault("jwt.expiration", 24*time.Hour)

	// CORS 默认值
	viper.SetDefault("cors.allowed_origins", []string{"http://localhost:3000", "http://localhost:8080"})
	viper.SetDefault("cors.allow_credentials", true)
	viper.SetDefault("cors.max_age", 43200)

	// 安全默认值
	viper.SetDefault("security.headers.enabled", true)
	viper.SetDefault("security.headers.enable_hsts", false)
	viper.SetDefault("security.headers.hsts_max_age", 31536000)
	viper.SetDefault("security.headers.hsts_include_subdomains", true)
	viper.SetDefault("security.headers.enable_csp", true)
	viper.SetDefault("security.headers.csp_policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline';")
	viper.SetDefault("security.headers.enable_frame_options", true)
	viper.SetDefault("security.headers.frame_options", "DENY")
	viper.SetDefault("security.enable_compression", true)
	viper.SetDefault("security.compression_level", 5)
	viper.SetDefault("security.tls.enabled", false)
	viper.SetDefault("security.tls.cert_file", "")
	viper.SetDefault("security.tls.key_file", "")
	viper.SetDefault("security.tls.min_version", "1.2")

	// 缓存默认值
	viper.SetDefault("cache.default_ttl", 5*time.Minute)
	viper.SetDefault("cache.user_ttl", 5*time.Minute)
	viper.SetDefault("cache.user_index_ttl", 10*time.Minute)
	viper.SetDefault("cache.user_count_ttl", 1*time.Minute)
	viper.SetDefault("cache.user_session_ttl", 30*time.Minute)
	viper.SetDefault("cache.content_ttl", 10*time.Minute)
	viper.SetDefault("cache.content_list_ttl", 2*time.Minute)
	viper.SetDefault("cache.stats_ttl", 1*time.Minute)
	viper.SetDefault("cache.not_found_ttl", 5*time.Minute)
	viper.SetDefault("cache.enable_jitter", true)
	viper.SetDefault("cache.jitter_percent", 20)
}

// Validate 验证配置（根据环境进行不同级别的校验）
func (c *Config) Validate(env string) error {
	// 验证服务器配置
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	// 验证数据库配置
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}

	// 验证日志配置
	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[c.Logger.Level] {
		return fmt.Errorf("invalid log level: %s", c.Logger.Level)
	}

	// 验证 JWT 配置
	if c.JWT.Secret == "" {
		return fmt.Errorf("jwt.secret is required")
	}
	
	// 生产环境额外校验
	if env == "prod" || env == "production" || c.Server.Mode == "release" {
		// 强制要求自定义 JWT 密钥
		if c.JWT.Secret == "your-secret-key-change-in-production" {
			return fmt.Errorf("jwt.secret must be customized in production (current: %s)", c.JWT.Secret)
		}
		
		// 警告：生产环境未启用 TLS
		if !c.Security.TLS.Enabled {
			slog.Warn("⚠️  TLS is disabled in production environment - this is insecure!")
		}
		
		// 警告：生产环境使用 debug 模式
		if c.Server.Mode == "debug" {
			return fmt.Errorf("server.mode must not be 'debug' in production")
		}
		
		// 警告：生产环境日志级别过低
		if c.Logger.Level == "debug" {
			slog.Warn("⚠️  Log level is 'debug' in production - consider using 'info' or 'warn'")
		}
	}
	
	// 开发/测试环境放宽 JWT 校验
	if env == "dev" || env == "development" || env == "test" {
		if c.JWT.Secret == "your-secret-key-change-in-production" {
			slog.Warn("Using default JWT secret in development mode")
		}
	}
	
	if c.JWT.Expiration <= 0 {
		return fmt.Errorf("jwt.expiration must be positive")
	}

	if c.Server.MaxRequestBodySize <= 0 {
		return fmt.Errorf("server.max_request_body_size must be positive")
	}

	return nil
}

// GetDSN 获取数据库连接字符串
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// GetRedisAddr 获取 Redis 地址
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}
