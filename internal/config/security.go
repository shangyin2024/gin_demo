package config

// SecurityConfig 安全配置
type SecurityConfig struct {
	// HTTP 安全头
	Headers SecurityHeadersConfig `mapstructure:"headers"`
	
	// 是否启用压缩
	EnableCompression bool `mapstructure:"enable_compression"`
	
	// 压缩级别 (-1 到 9)
	CompressionLevel int `mapstructure:"compression_level"`
	
	// TLS/HTTPS 配置
	TLS TLSConfig `mapstructure:"tls"`
}

// SecurityHeadersConfig HTTP 安全头配置
type SecurityHeadersConfig struct {
	// 是否启用安全头
	Enabled bool `mapstructure:"enabled"`
	
	// HSTS 配置
	EnableHSTS            bool   `mapstructure:"enable_hsts"`
	HSTSMaxAge            int    `mapstructure:"hsts_max_age"`
	HSTSIncludeSubdomains bool   `mapstructure:"hsts_include_subdomains"`
	
	// CSP 配置
	EnableCSP bool   `mapstructure:"enable_csp"`
	CSPPolicy string `mapstructure:"csp_policy"`
	
	// Frame Options
	EnableFrameOptions bool   `mapstructure:"enable_frame_options"`
	FrameOptions       string `mapstructure:"frame_options"`
}

// TLSConfig TLS/HTTPS 配置
type TLSConfig struct {
	// 是否启用 TLS
	Enabled bool `mapstructure:"enabled"`
	
	// 证书文件路径
	CertFile string `mapstructure:"cert_file"`
	
	// 密钥文件路径
	KeyFile string `mapstructure:"key_file"`
	
	// 最小 TLS 版本: "1.2", "1.3"
	MinVersion string `mapstructure:"min_version"`
}
