package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"

	"github.com/zhongjian-ready/goapi/internal/tools/response"
)

// UnAuthorizedError 定义了一个标准的未授权错误。
var UnAuthorizedError = errors.New("Invalid token or unauthorized.")

// ContextKey 定义上下文键的类型
type ContextKey string

// UserIDKey 是上下文中使用用户ID的键
const UserIDKey ContextKey = "userid"

// NewAuthorization 创建一个 Authorization 中间件，闭包持有 jwtSecret。
func NewAuthorization(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. 获取 Authorization Header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.WriteError(w, UnAuthorizedError.Error(), http.StatusUnauthorized)
				return
			}

			// 2. 解析 Bearer Token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				log.Error("Invalid authorization header format")
				response.WriteError(w, UnAuthorizedError.Error(), http.StatusUnauthorized)
				return
			}
			tokenString := parts[1]

			// 3. 解析并验证 Token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// 验证签名算法
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				// 使用传入的密钥
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				log.Error("Token validation failed:", err)
				response.WriteError(w, UnAuthorizedError.Error(), http.StatusUnauthorized)
				return
			}

			// 4. 提取 Claims 并放入 Context
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				// 安全获取 user_id
				userIDFloat, ok := claims["userid"].(float64)
				if !ok {
					log.Error("Token claims missing userid or invalid type")
					response.WriteError(w, UnAuthorizedError.Error(), http.StatusUnauthorized)
					return
				}

				userID := int(userIDFloat)
				ctx := context.WithValue(r.Context(), UserIDKey, userID)

				// 调用下一个处理函数，传入新的 Context
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				response.WriteError(w, UnAuthorizedError.Error(), http.StatusUnauthorized)
			}
		})
	}
}
