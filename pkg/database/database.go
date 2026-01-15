package database

import (
	"database/sql"
	"time"
)

// Type 数据库类型
type Type string

const (
	// TypePostgreSQL PostgreSQL 数据库
	TypePostgreSQL Type = "postgres"
	// TypeMySQL MySQL 数据库
	TypeMySQL Type = "mysql"
)

// CommonConfig 通用数据库配置
type CommonConfig struct {
	Type            Type          // 数据库类型
	Host            string        // 主机地址
	Port            int           // 端口
	User            string        // 用户名
	Password        string        // 密码
	DBName          string        // 数据库名
	MaxOpenConns    int           // 最大打开连接数
	MaxIdleConns    int           // 最大空闲连接数
	ConnMaxLifetime time.Duration // 连接最大生命周期
	ConnMaxIdleTime time.Duration // 连接最大空闲时间

	// PostgreSQL 特定配置
	SSLMode string // SSL 模式

	// MySQL 特定配置
	Charset   string // 字符集
	ParseTime bool   // 是否解析时间
	Loc       string // 时区
}

// New 创建数据库连接（根据类型自动选择）
func New(cfg CommonConfig) (*sql.DB, error) {
	switch cfg.Type {
	case TypePostgreSQL:
		return NewPostgres(Config{
			Host:            cfg.Host,
			Port:            cfg.Port,
			User:            cfg.User,
			Password:        cfg.Password,
			DBName:          cfg.DBName,
			SSLMode:         cfg.SSLMode,
			MaxOpenConns:    cfg.MaxOpenConns,
			MaxIdleConns:    cfg.MaxIdleConns,
			ConnMaxLifetime: int(cfg.ConnMaxLifetime.Seconds()),
			ConnMaxIdleTime: int(cfg.ConnMaxIdleTime.Seconds()),
		})
	case TypeMySQL:
		return NewMySQL(Config{
			Host:            cfg.Host,
			Port:            cfg.Port,
			User:            cfg.User,
			Password:        cfg.Password,
			DBName:          cfg.DBName,
			SSLMode:         "", // MySQL 不使用
			MaxOpenConns:    cfg.MaxOpenConns,
			MaxIdleConns:    cfg.MaxIdleConns,
			ConnMaxLifetime: int(cfg.ConnMaxLifetime.Seconds()),
			ConnMaxIdleTime: int(cfg.ConnMaxIdleTime.Seconds()),
		})
	default:
		return nil, ErrUnsupportedDatabaseType
	}
}

// Close 关闭数据库连接
func Close(db *sql.DB) error {
	if db == nil {
		return nil
	}
	return db.Close()
}
