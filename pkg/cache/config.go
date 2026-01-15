package cache

import "time"

// CacheConfig 缓存配置
type CacheConfig struct {
	// 默认 TTL
	DefaultTTL time.Duration `mapstructure:"default_ttl"`

	// 用户相关缓存 TTL
	UserTTL          time.Duration `mapstructure:"user_ttl"`
	UserIndexTTL     time.Duration `mapstructure:"user_index_ttl"`
	UserCountTTL     time.Duration `mapstructure:"user_count_ttl"`
	UserSessionTTL   time.Duration `mapstructure:"user_session_ttl"`

	// 内容相关缓存 TTL
	ContentTTL       time.Duration `mapstructure:"content_ttl"`
	ContentListTTL   time.Duration `mapstructure:"content_list_ttl"`

	// 统计数据缓存 TTL
	StatsTTL         time.Duration `mapstructure:"stats_ttl"`

	// NotFound 占位符 TTL
	NotFoundTTL      time.Duration `mapstructure:"not_found_ttl"`

	// 是否启用 Jitter（随机扰动）
	EnableJitter     bool          `mapstructure:"enable_jitter"`
	
	// Jitter 范围（占基础 TTL 的百分比）
	JitterPercent    int           `mapstructure:"jitter_percent"`
}

// DefaultCacheConfig 默认缓存配置
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		DefaultTTL:       5 * time.Minute,
		UserTTL:          5 * time.Minute,
		UserIndexTTL:     10 * time.Minute, // 索引缓存时间更长
		UserCountTTL:     1 * time.Minute,  // 统计数据变化频繁，TTL较短
		UserSessionTTL:   30 * time.Minute,
		ContentTTL:       10 * time.Minute,
		ContentListTTL:   2 * time.Minute,
		StatsTTL:         1 * time.Minute,
		NotFoundTTL:      5 * time.Minute,
		EnableJitter:     true,
		JitterPercent:    20, // 20% 的随机扰动
	}
}

// GetUserTTL 获取用户缓存 TTL
func (c *CacheConfig) GetUserTTL() time.Duration {
	return c.applyJitter(c.UserTTL)
}

// GetUserIndexTTL 获取用户索引缓存 TTL
func (c *CacheConfig) GetUserIndexTTL() time.Duration {
	return c.applyJitter(c.UserIndexTTL)
}

// GetUserCountTTL 获取用户统计缓存 TTL
func (c *CacheConfig) GetUserCountTTL() time.Duration {
	return c.applyJitter(c.UserCountTTL)
}

// GetContentTTL 获取内容缓存 TTL
func (c *CacheConfig) GetContentTTL() time.Duration {
	return c.applyJitter(c.ContentTTL)
}

// GetStatsTTL 获取统计数据缓存 TTL
func (c *CacheConfig) GetStatsTTL() time.Duration {
	return c.applyJitter(c.StatsTTL)
}

// GetNotFoundTTL 获取 NotFound 占位符 TTL
func (c *CacheConfig) GetNotFoundTTL() time.Duration {
	return c.applyJitter(c.NotFoundTTL)
}

// applyJitter 应用随机扰动（防止缓存雪崩）
func (c *CacheConfig) applyJitter(baseTTL time.Duration) time.Duration {
	if !c.EnableJitter || c.JitterPercent <= 0 {
		return baseTTL
	}

	// 计算 jitter 范围
	jitterRange := int64(float64(baseTTL) * float64(c.JitterPercent) / 100.0)
	
	// 添加随机扰动 [-jitterRange/2, +jitterRange/2]
	jitter := int64(0)
	if jitterRange > 0 {
		jitter = int64(float64(jitterRange) * (0.5 - float64(time.Now().UnixNano()%1000)/1000.0))
	}

	return baseTTL + time.Duration(jitter)
}

// Validate 验证配置
func (c *CacheConfig) Validate() error {
	if c.DefaultTTL <= 0 {
		c.DefaultTTL = 5 * time.Minute
	}
	if c.UserTTL <= 0 {
		c.UserTTL = c.DefaultTTL
	}
	if c.UserIndexTTL <= 0 {
		c.UserIndexTTL = c.UserTTL * 2
	}
	if c.NotFoundTTL <= 0 {
		c.NotFoundTTL = c.DefaultTTL
	}
	if c.JitterPercent < 0 || c.JitterPercent > 50 {
		c.JitterPercent = 20 // 默认 20%
	}
	return nil
}
