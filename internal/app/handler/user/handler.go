package user

import (
	"log/slog"

	"gin_demo/internal/app/middleware"
	"gin_demo/internal/domain/service"
	"gin_demo/internal/repository"
	"gin_demo/internal/response"
	"gin_demo/pkg/auth"

	"github.com/gin-gonic/gin"
)

// Handler 用户处理器
type Handler struct {
	userService service.UserService
	jwtManager  *auth.DefaultJWTManager
}

// NewHandler 创建用户处理器
func NewHandler(userService service.UserService, jwtManager *auth.DefaultJWTManager) *Handler {
	return &Handler{
		userService: userService,
		jwtManager:  jwtManager,
	}
}

// Register 用户注册
//
// @Summary 用户注册
// @Description 创建新用户账户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册信息"
// @Success 200 {object} response.Response{data=Response} "注册成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 409 {object} response.Response "用户已存在"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /users/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.NewWithError(response.CodeInvalidParams, "参数错误", err))
		return
	}

	user, err := h.userService.Register(c.Request.Context(), service.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "Register failed", "error", err)
		response.Error(c, err)
		return
	}

	response.Success(c, toResponse(user))
}

// Login 用户登录
//
// @Summary 用户登录
// @Description 验证用户凭证并返回 JWT Token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=LoginResponse} "登录成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "认证失败"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /users/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.NewWithError(response.CodeInvalidParams, "参数错误", err))
		return
	}

	user, err := h.userService.Login(c.Request.Context(), service.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		slog.WarnContext(c.Request.Context(), "Login failed", "email", req.Email, "error", err)
		response.Error(c, err)
		return
	}

	// 生成 JWT Token（只包含 UserID）
	token, err := h.jwtManager.GenerateToken(user.ID)
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "Generate token failed", "user_id", user.ID, "error", err)
		response.Error(c, response.NewWithError(response.CodeInternalError, "生成 Token 失败", err))
		return
	}

	response.Success(c, LoginResponse{
		User:  toResponse(user),
		Token: token,
	})
}

// GetProfile 获取当前用户信息
//
// @Summary 获取当前用户信息
// @Description 获取已登录用户的个人资料
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=Response} "获取成功"
// @Failure 401 {object} response.Response "未认证"
// @Failure 404 {object} response.Response "用户不存在"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /users/me [get]
func (h *Handler) GetProfile(c *gin.Context) {
	// 从认证中间件获取当前用户 ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, response.NewWithError(response.CodeUnauthorized, "未认证", nil))
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, toResponse(user))
}

// GetUser 获取指定用户信息（管理员）
//
// @Summary 获取指定用户信息
// @Description 管理员获取任意用户的详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response{data=Response} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未认证"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 404 {object} response.Response "用户不存在"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /users/{id} [get]
func (h *Handler) GetUser(c *gin.Context) {
	var req IDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.Error(c, response.NewWithError(response.CodeInvalidParams, "无效的用户ID", err))
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), req.ID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, toResponse(user))
}

// UpdateProfile 更新当前用户信息
//
// @Summary 更新当前用户信息
// @Description 更新已登录用户的个人资料
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateUserRequest true "更新信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未认证"
// @Failure 409 {object} response.Response "邮箱/用户名已存在"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /users/me [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	// 从认证中间件获取当前用户 ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, response.NewWithError(response.CodeUnauthorized, "未认证", nil))
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.NewWithError(response.CodeInvalidParams, "参数错误", err))
		return
	}

	err := h.userService.UpdateUser(c.Request.Context(), service.UpdateUserInput{
		UserID:   userID,
		Username: &req.Username,
		Email:    &req.Email,
		Avatar:   &req.Avatar,
	})
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "Update user failed", "user_id", userID, "error", err)
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// UpdateUser 更新指定用户信息（管理员）
//
// @Summary 更新指定用户信息
// @Description 管理员更新任意用户的信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param request body UpdateUserRequest true "更新信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未认证"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 404 {object} response.Response "用户不存在"
// @Failure 409 {object} response.Response "邮箱/用户名已存在"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /users/{id} [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	var idReq IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.Error(c, response.NewWithError(response.CodeInvalidParams, "无效的用户ID", err))
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.NewWithError(response.CodeInvalidParams, "参数错误", err))
		return
	}

	err := h.userService.UpdateUser(c.Request.Context(), service.UpdateUserInput{
		UserID:   idReq.ID,
		Username: &req.Username,
		Email:    &req.Email,
		Avatar:   &req.Avatar,
	})
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "Update user failed", "user_id", idReq.ID, "error", err)
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// ChangePassword 修改当前用户密码
//
// @Summary 修改密码
// @Description 修改当前登录用户的密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "密码信息"
// @Success 200 {object} response.Response "修改成功"
// @Failure 400 {object} response.Response "参数错误或密码不正确"
// @Failure 401 {object} response.Response "未认证"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /users/me/password [put]
func (h *Handler) ChangePassword(c *gin.Context) {
	// 从认证中间件获取当前用户 ID
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, response.NewWithError(response.CodeUnauthorized, "未认证", nil))
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.NewWithError(response.CodeInvalidParams, "参数错误", err))
		return
	}

	err := h.userService.ChangePassword(c.Request.Context(), service.ChangePasswordInput{
		UserID:      userID,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		slog.WarnContext(c.Request.Context(), "Change password failed", "user_id", userID, "error", err)
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "密码修改成功"})
}

// DeleteUser 删除用户
//
// @Summary 删除用户
// @Description 管理员删除指定用户（软删除）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未认证"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 404 {object} response.Response "用户不存在"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	var req IDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.Error(c, response.NewWithError(response.CodeInvalidParams, "无效的用户ID", err))
		return
	}

	err := h.userService.DeleteUser(c.Request.Context(), req.ID)
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "Delete user failed", "user_id", req.ID, "error", err)
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// ListUsers 用户列表
//
// @Summary 获取用户列表
// @Description 管理员获取用户列表（分页）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=response.ListResponse} "获取成功"
// @Failure 401 {object} response.Response "未认证"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /users [get]
func (h *Handler) ListUsers(c *gin.Context) {
	// 获取分页参数
	pagination := response.GetPagination(c)

	users, total, err := h.userService.ListUsers(c.Request.Context(), pagination.GetLimit(), pagination.GetOffset())
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "List users failed", "error", err)
		response.Error(c, response.Wrap(err, response.CodeInternalError, "获取用户列表失败"))
		return
	}

	// 转换为响应 DTO
	responses := make([]Response, 0, len(users))
	for _, user := range users {
		responses = append(responses, toResponse(user))
	}

	// 构建分页响应
	paginationResp := response.NewPaginationResponse(pagination.Page, pagination.PageSize, total)

	response.Success(c, response.NewListResponse(responses, paginationResp))
}

// toResponse 转换为响应 DTO
func toResponse(user repository.User) Response {
	avatar := ""
	if user.Avatar.Valid {
		avatar = user.Avatar.String
	}

	return Response{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Avatar:    avatar,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
