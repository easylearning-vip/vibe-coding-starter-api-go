package mcp

import (
	"net/http"
	"strings"
)

// ExtractToken 从HTTP请求中提取用户token
// 支持以下格式:
// - X-User-Token header
// - Authorization: Token <token>
// - Authorization: Bearer <token>
func ExtractToken(r *http.Request) string {
	// 优先检查 X-User-Token header
	if v := strings.TrimSpace(r.Header.Get("X-User-Token")); v != "" {
		return v
	}

	// 检查 Authorization header
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	if auth == "" {
		return ""
	}

	// 支持 "Token <token>" 格式
	lower := strings.ToLower(auth)
	if strings.HasPrefix(lower, "token ") {
		return strings.TrimSpace(auth[len("Token "):])
	}

	// 支持 "Bearer <token>" 格式
	if strings.HasPrefix(lower, "bearer ") {
		return strings.TrimSpace(auth[len("Bearer "):])
	}

	return ""
}

