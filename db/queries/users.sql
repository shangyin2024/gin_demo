-- name: GetUserByID :one
-- 通过 ID 获取用户
SELECT id, username, email, avatar, status, created_at, updated_at
FROM users
WHERE id = ? AND status = 1
LIMIT 1;

-- name: GetUserByEmail :one
-- 通过 Email 获取用户（包含密码，用于登录验证）
SELECT id, username, email, password, avatar, status, created_at, updated_at
FROM users
WHERE email = ? AND status = 1
LIMIT 1;

-- name: GetUserByUsername :one
-- 通过 Username 获取用户
SELECT id, username, email, avatar, status, created_at, updated_at
FROM users
WHERE username = ? AND status = 1
LIMIT 1;

-- name: ListUsers :many
-- 列出用户（分页）
SELECT id, username, email, avatar, status, created_at, updated_at
FROM users
WHERE status = 1
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CreateUser :execresult
-- 创建用户（MySQL 使用 execresult 获取 LastInsertId）
INSERT INTO users (username, email, password, avatar)
VALUES (?, ?, ?, ?);

-- name: UpdateUser :exec
-- 更新用户信息
UPDATE users
SET username = ?,
    email = ?,
    avatar = ?
WHERE id = ?;

-- name: UpdateUserPassword :exec
-- 更新用户密码
UPDATE users
SET password = ?
WHERE id = ?;

-- name: DeleteUser :exec
-- 软删除用户（设置状态为禁用）
UPDATE users
SET status = 2
WHERE id = ?;

-- name: CountUsers :one
-- 统计用户总数
SELECT COUNT(*) as total
FROM users
WHERE status = 1;

-- name: GetUserIDByEmail :one
-- 通过 Email 获取用户 ID（用于缓存索引）
SELECT id
FROM users
WHERE email = ? AND status = 1
LIMIT 1;

-- name: GetUserIDByUsername :one
-- 通过 Username 获取用户 ID（用于缓存索引）
SELECT id
FROM users
WHERE username = ? AND status = 1
LIMIT 1;
