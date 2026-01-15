package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// ============================================================================
// 业务指标
// ============================================================================

var (
	// 用户注册指标
	UserRegistrations = promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_registrations_total",
		Help: "Total number of user registrations",
	})

	// 用户登录指标
	UserLogins = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "user_logins_total",
		Help: "Total number of user logins",
	}, []string{"status"}) // status: success, failed

	// 活跃用户指标（Gauge）
	ActiveUsers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "active_users_current",
		Help: "Current number of active users",
	})

	// 在线用户指标
	OnlineUsers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "online_users_current",
		Help: "Current number of online users",
	})

	// 用户操作指标
	UserOperations = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "user_operations_total",
		Help: "Total number of user operations",
	}, []string{"operation", "status"}) // operation: create, update, delete; status: success, failed

	// 密码修改指标
	PasswordChanges = promauto.NewCounter(prometheus.CounterOpts{
		Name: "password_changes_total",
		Help: "Total number of password changes",
	})

	// 用户删除指标
	UserDeletions = promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_deletions_total",
		Help: "Total number of user deletions",
	})
)

// ============================================================================
// 认证指标
// ============================================================================

var (
	// Token 生成指标
	TokenGenerations = promauto.NewCounter(prometheus.CounterOpts{
		Name: "token_generations_total",
		Help: "Total number of token generations",
	})

	// Token 验证指标
	TokenValidations = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "token_validations_total",
		Help: "Total number of token validations",
	}, []string{"status"}) // status: success, expired, invalid

	// Token 刷新指标
	TokenRefreshes = promauto.NewCounter(prometheus.CounterOpts{
		Name: "token_refreshes_total",
		Help: "Total number of token refreshes",
	})

	// 认证失败指标
	AuthFailures = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "auth_failures_total",
		Help: "Total number of authentication failures",
	}, []string{"reason"}) // reason: invalid_credentials, expired_token, invalid_token
)

// ============================================================================
// 权限指标
// ============================================================================

var (
	// 权限检查指标
	PermissionChecks = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "permission_checks_total",
		Help: "Total number of permission checks",
	}, []string{"role", "permission", "result"}) // result: allowed, denied

	// 权限拒绝指标
	PermissionDenials = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "permission_denials_total",
		Help: "Total number of permission denials",
	}, []string{"role", "operation"})
)

// ============================================================================
// 业务错误指标
// ============================================================================

var (
	// 业务错误计数
	BusinessErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "business_errors_total",
		Help: "Total number of business errors",
	}, []string{"error_code", "operation"})

	// 服务端错误计数
	InternalErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "internal_errors_total",
		Help: "Total number of internal errors",
	}, []string{"component", "error_type"}) // component: service, repository, cache
)

// ============================================================================
// 辅助函数
// ============================================================================

// RecordUserRegistration 记录用户注册
func RecordUserRegistration() {
	UserRegistrations.Inc()
}

// RecordUserLogin 记录用户登录
func RecordUserLogin(success bool) {
	status := "success"
	if !success {
		status = "failed"
	}
	UserLogins.WithLabelValues(status).Inc()
}

// RecordUserOperation 记录用户操作
func RecordUserOperation(operation string, success bool) {
	status := "success"
	if !success {
		status = "failed"
	}
	UserOperations.WithLabelValues(operation, status).Inc()
}

// RecordTokenValidation 记录 Token 验证
func RecordTokenValidation(status string) {
	TokenValidations.WithLabelValues(status).Inc()
}

// RecordPermissionCheck 记录权限检查
func RecordPermissionCheck(role, permission string, allowed bool) {
	result := "allowed"
	if !allowed {
		result = "denied"
	}
	PermissionChecks.WithLabelValues(role, permission, result).Inc()
	
	if !allowed {
		PermissionDenials.WithLabelValues(role, permission).Inc()
	}
}

// RecordBusinessError 记录业务错误
func RecordBusinessError(errorCode, operation string) {
	BusinessErrors.WithLabelValues(errorCode, operation).Inc()
}

// RecordInternalError 记录内部错误
func RecordInternalError(component, errorType string) {
	InternalErrors.WithLabelValues(component, errorType).Inc()
}

// UpdateActiveUsers 更新活跃用户数
func UpdateActiveUsers(count float64) {
	ActiveUsers.Set(count)
}

// UpdateOnlineUsers 更新在线用户数
func UpdateOnlineUsers(count float64) {
	OnlineUsers.Set(count)
}
