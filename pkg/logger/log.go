package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

// AttrExtractor 定义如何从 context 中提取属性
type AttrExtractor func(ctx context.Context) []slog.Attr

// Config 日志初始化配置
type Config struct {
	Level        slog.Level      // 日志级别
	IsJSON       bool            // 是否开启 JSON 格式
	AddSource    bool            // 是否添加行号源码信息
	RequestIDKey string          // 从 Context 中读取 RequestID 的 Key (默认 requestId)
	Extractors   []AttrExtractor // 额外的自定义提取器
}

// ContextHandler 包装器：在处理日志记录前从 Context 提取属性
type ContextHandler struct {
	slog.Handler
	extractors []AttrExtractor
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx != nil {
		for _, fn := range h.extractors {
			if attrs := fn(ctx); len(attrs) > 0 {
				r.AddAttrs(attrs...)
			}
		}
	}
	return h.Handler.Handle(ctx, r)
}

// Setup 初始化全局 slog 配置
func Setup(cfg Config) {
	// 获取当前工作目录，用于计算相对路径
	wd, _ := os.Getwd()

	// 1. 设置默认 RequestID 提取逻辑
	if cfg.RequestIDKey == "" {
		cfg.RequestIDKey = "requestId" // 默认匹配大部分中间件的 Key
	}

	defaultExtractor := func(ctx context.Context) []slog.Attr {
		if rid, ok := ctx.Value(cfg.RequestIDKey).(string); ok {
			return []slog.Attr{slog.String("request_id", rid)}
		}
		return nil
	}

	// 合并内置和自定义提取器
	allExtractors := append([]AttrExtractor{defaultExtractor}, cfg.Extractors...)

	// 2. 配置 Handler 选项
	opts := &slog.HandlerOptions{
		AddSource: cfg.AddSource,
		Level:     cfg.Level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// 格式化时间
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format(time.DateTime))
			}
			// 2. 优化 Source 字段：由对象改为 "文件名:行号" 字符串
			if a.Key == slog.SourceKey && a.Value.Kind() == slog.KindAny {
				if source, ok := a.Value.Any().(*slog.Source); ok {
					// 💡 优化点：将绝对路径转换为相对路径
					file := source.File
					if rel, err := strings.CutPrefix(file, wd+"/"); err {
						file = rel
					} else {
						// 如果不在当前工作目录下（比如引用的第三方库），则只取最后两级
						// 避免输出冗长的 /Users/xxx/go/pkg/mod/...
						parts := strings.Split(file, "/")
						if len(parts) > 2 {
							file = strings.Join(parts[len(parts)-2:], "/")
						}
					}
					shortPath := fmt.Sprintf("%s:%d", file, source.Line)
					return slog.String(slog.SourceKey, shortPath)
				}
			}
			return a
		},
	}

	// 3. 构造 Handler
	var baseHandler slog.Handler
	if cfg.IsJSON {
		baseHandler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		baseHandler = slog.NewTextHandler(os.Stdout, opts)
	}

	// 4. 设置为全局默认日志实例
	slog.SetDefault(slog.New(&ContextHandler{
		Handler:    baseHandler,
		extractors: allExtractors,
	}))
}
