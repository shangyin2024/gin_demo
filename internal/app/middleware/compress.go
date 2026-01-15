package middleware

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// CompressConfig 压缩配置
type CompressConfig struct {
	// 压缩级别: -1 (默认), 0 (不压缩), 1-9 (压缩级别)
	Level int
	// 最小压缩大小（字节）
	MinSize int
	// 排除的路径
	ExcludePaths []string
	// 排除的扩展名
	ExcludeExtensions []string
}

// DefaultCompressConfig 默认压缩配置
var DefaultCompressConfig = CompressConfig{
	Level:   gzip.DefaultCompression,
	MinSize: 1024, // 1KB
	ExcludePaths: []string{
		"/metrics", // Prometheus 指标
	},
	ExcludeExtensions: []string{
		".jpg", ".jpeg", ".png", ".gif", ".webp", // 图片已压缩
		".zip", ".gz", ".7z", ".rar",             // 压缩文件
		".mp4", ".avi", ".mov",                   // 视频已压缩
		".mp3", ".wav", ".flac",                  // 音频已压缩
	},
}

// Compress Gzip 压缩中间件
func Compress(config ...CompressConfig) gin.HandlerFunc {
	cfg := DefaultCompressConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	// 使用 gin-contrib/gzip
	return gzip.Gzip(cfg.Level, gzip.WithExcludedPaths(cfg.ExcludePaths))
}
