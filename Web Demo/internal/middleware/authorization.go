package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"

	"github.com/zhongjian-ready/goapi/api"
)

// UnAuthorizedError 定义了一个标准的未授权错误。
var UnAuthorizedError = errors.New("Invalid token or unauthorized.")

// Authorization 是一个 HTTP 中间件函数，用于验证 JWT Token。
func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. 获取 Authorization Header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			api.WriteError(w, UnAuthorizedError.Error(), http.StatusUnauthorized)
			return
		}

		// 2. 解析 Bearer Token
		// 格式应该是 "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Error("Invalid authorization header format")
			api.WriteError(w, UnAuthorizedError.Error(), http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

		// 3. 解析并验证 Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// 获取密钥 (生产环境应从环境变量获取)
			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				secret = "default_secret_key"
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			log.Error("Token validation failed:", err)
			api.WriteError(w, UnAuthorizedError.Error(), http.StatusUnauthorized)
			return
		}

		// 4. 提取 Claims 并放入 Context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// 将 user_id 放入 context 中，供后续 Handler 使用
			// 注意：JWT 解析出来的数字通常是 float64
			userID := int(claims["userid"].(float64))
			ctx := context.WithValue(r.Context(), "userid", userID)

			// 调用下一个处理函数，传入新的 Context
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			api.WriteError(w, UnAuthorizedError.Error(), http.StatusUnauthorized)
		}
	})
}
