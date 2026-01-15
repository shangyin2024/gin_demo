package app

import (
	"gin_demo/internal/app/middleware"
	"gin_demo/pkg/auth"

	"github.com/gin-gonic/gin"
)

// setupAPIv1Routes 配置 API v1 路由
func setupAPIv1Routes(rg *gin.RouterGroup, handlers *Handlers) {
	// 用户路由
	setupUserRoutes(rg, handlers)

	// 可以在这里添加更多 v1 模块路由
	// setupArticleRoutes(rg, handlers)
	// setupCommentRoutes(rg, handlers)
	// setupOrderRoutes(rg, handlers)
}

// setupUserRoutes 配置用户路由（包含 RBAC 权限控制）
func setupUserRoutes(rg *gin.RouterGroup, handlers *Handlers) {
	users := rg.Group("/users")
	{
		// ========================================
		// 公开路由（无需认证）
		// ========================================
		users.POST("/register", handlers.User.Register) // 用户注册
		users.POST("/login", handlers.User.Login)       // 用户登录

		// ========================================
		// 个人资料路由（需要认证，操作当前用户）
		// ========================================
		profile := users.Group("")
		profile.Use(handlers.Auth.Handle()) // 基础认证即可
		{
			profile.GET("/me", handlers.User.GetProfile)              // 获取当前用户信息
			profile.PUT("/me", handlers.User.UpdateProfile)           // 更新当前用户信息
			profile.PUT("/me/password", handlers.User.ChangePassword) // 修改密码
		}

		// ========================================
		// 管理员路由（需要管理员权限）
		// ========================================
		// 方式 1: 使用 RequireRole 中间件（推荐）
		admin := users.Group("")
		admin.Use(handlers.Auth.Handle())                                            // 先认证
		admin.Use(middleware.RequireRole(auth.RoleAdmin, auth.RoleSuperAdmin))     // 再检查角色
		{
			admin.GET("", handlers.User.ListUsers)         // 用户列表（需要 admin 或 super_admin 角色）
			admin.GET("/:id", handlers.User.GetUser)       // 获取指定用户
			admin.PUT("/:id", handlers.User.UpdateUser)    // 更新指定用户
		}

		// ========================================
		// 超级管理员路由（需要超级管理员权限）
		// ========================================
		// 方式 2: 使用 RequireSuperAdmin 快捷中间件
		superAdmin := users.Group("")
		superAdmin.Use(handlers.Auth.Handle())              // 先认证
		superAdmin.Use(middleware.RequireSuperAdmin())      // 超级管理员专用
		{
			superAdmin.DELETE("/:id", handlers.User.DeleteUser) // 删除用户（仅超级管理员）
		}

		// ========================================
		// 基于权限的路由示例（细粒度控制）
		// ========================================
		// 方式 3: 使用 RequirePermission 中间件
		// management := users.Group("/management")
		// management.Use(handlers.Auth.Handle())
		// management.Use(middleware.RequirePermission(auth.PermissionUserWrite))
		// {
		//     management.POST("/batch", handlers.User.BatchUpdate)  // 批量更新（需要 user:write 权限）
		// }
	}
}

// ========================================
// RBAC 使用示例和最佳实践
// ========================================
//
// 1. 角色层级（从低到高）：
//    - guest: 游客（未登录）
//    - user: 普通用户
//    - moderator: 版主/审核员
//    - admin: 管理员
//    - super_admin: 超级管理员
//
// 2. 权限继承：
//    - 超级管理员拥有所有权限
//    - 管理员拥有大部分权限（除了系统配置）
//    - 版主拥有内容审核权限
//    - 普通用户只能操作自己的数据
//
// 3. 路由保护方式：
//
//    方式 A: 角色检查（适用于粗粒度控制）
//    ```go
//    admin := router.Group("/admin")
//    admin.Use(handlers.Auth.Handle())                    // 认证
//    admin.Use(middleware.RequireAdmin())                 // 需要管理员角色
//    ```
//
//    方式 B: 权限检查（适用于细粒度控制）
//    ```go
//    users := router.Group("/users")
//    users.Use(handlers.Auth.Handle())                    // 认证
//    users.Use(middleware.RequirePermission(              // 需要特定权限
//        auth.PermissionUserWrite,
//        auth.PermissionUserDelete,
//    ))
//    ```
//
//    方式 C: 混合使用（最灵活）
//    ```go
//    sensitive := router.Group("/sensitive")
//    sensitive.Use(handlers.Auth.Handle())                // 认证
//    sensitive.Use(middleware.RequireRole(                // 角色检查
//        auth.RoleAdmin,
//        auth.RoleSuperAdmin,
//    ))
//    sensitive.Use(middleware.RequirePermission(          // 权限检查
//        auth.PermissionSystemConfig,
//    ))
//    ```
//
// 4. Handler 内部权限检查：
//    ```go
//    func (h *Handler) SomeAction(c *gin.Context) {
//        claims := middleware.GetRBACClaims(c)
//        if claims == nil {
//            response.Error(c, response.ErrUnauthorized)
//            return
//        }
//
//        // 检查是否有权限操作目标资源
//        if !claims.HasPermission(auth.PermissionUserDelete) {
//            response.Error(c, response.ErrForbidden)
//            return
//        }
//
//        // 执行操作...
//    }
//    ```
//
// 5. 注意事项：
//    - 先使用 handlers.Auth.Handle() 进行认证
//    - 再使用 RequireRole/RequirePermission 进行授权
//    - 顺序不能颠倒
//    - 可以多层叠加（先角色再权限）
//
