package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrInvalidToken Token 无效
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken Token 已过期
	ErrExpiredToken = errors.New("token expired")
)

// Claims JWT 声明（泛型版本，支持不同类型的 UserID）
type Claims[T any] struct {
	UserID T `json:"user_id"`
	jwt.RegisteredClaims
}

// DefaultClaims 默认 Claims（UserID 为 int64）
type DefaultClaims = Claims[int64]

// JWTManager JWT 管理器（泛型版本）
type JWTManager[T any] struct {
	secret     string
	expiration time.Duration
}

// NewJWTManager 创建 JWT 管理器
func NewJWTManager[T any](secret string, expiration time.Duration) *JWTManager[T] {
	return &JWTManager[T]{
		secret:     secret,
		expiration: expiration,
	}
}

// GenerateToken 生成 JWT Token
func (m *JWTManager[T]) GenerateToken(userID T) (string, error) {
	now := time.Now()
	claims := Claims[T]{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken 验证 JWT Token
func (m *JWTManager[T]) ValidateToken(tokenString string) (*Claims[T], error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims[T]{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*Claims[T])
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// RefreshToken 刷新 Token（验证旧 Token 并生成新 Token）
func (m *JWTManager[T]) RefreshToken(tokenString string) (string, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil && !errors.Is(err, ErrExpiredToken) {
		return "", err
	}

	// 即使 Token 过期，只要签名有效，就允许刷新
	return m.GenerateToken(claims.UserID)
}

// DefaultJWTManager 默认的 JWT 管理器（UserID 为 int64）
type DefaultJWTManager = JWTManager[int64]

// NewDefaultJWTManager 创建默认的 JWT 管理器（UserID 为 int64）
func NewDefaultJWTManager(secret string, expiration time.Duration) *DefaultJWTManager {
	return NewJWTManager[int64](secret, expiration)
}
