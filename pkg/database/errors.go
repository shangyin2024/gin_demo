package database

import "errors"

var (
	// ErrUnsupportedDatabaseType 不支持的数据库类型
	ErrUnsupportedDatabaseType = errors.New("unsupported database type")
)
