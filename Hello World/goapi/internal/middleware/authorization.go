package middleware

import (
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/zhongjian-ready/goapi/api"
	"github.com/zhongjian-ready/goapi/internal/tools"
)

// UnAuthorizedError 定义了一个标准的未授权错误。
// 当用户名或 Token 无效时使用。
var UnAuthorizedError = errors.New("Invalid username or token.")

// Authorization 是一个 HTTP 中间件函数，用于验证请求的合法性。
// 它接收一个 http.Handler (next)，并返回一个新的 http.Handler。
// 只有通过验证的请求才会调用 next.ServeHTTP 继续处理，否则直接返回错误。
func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. 获取认证信息
		// 从查询参数中获取 username。
		username := r.URL.Query().Get("username")
		// 从 HTTP Header 中获取 Authorization Token。
		token := r.Header.Get("Authorization")
		var err error

		// 2. 基础校验
		// 如果用户名或 Token 均为空，则认为请求未携带认证信息，记录错误并返回。
		if username == "" && token == "" {
			log.Error(UnAuthorizedError)
			api.RequestErrorHandler(w, UnAuthorizedError)
			return
		}

		// 3. 连接数据库
		// 创建数据库实例以查询用户信息。
		var database *tools.DatabaseInterface
		database, err = tools.NewDatabase()
		if err != nil {
			// 数据库连接失败通常是内部错误。
			api.InternalErrorHandler(w)
			return
		}

		// 4. 验证用户凭证
		// 从数据库中获取该用户的登录详情（包括正确的 AuthToken）。
		var loginDetails *tools.LoginDetails
		loginDetails = (*database).GetUserLoginDetails(username)

		// 5. 比对 Token
		// 如果用户不存在 (loginDetails == nil)
		// 或者请求头中的 Token 与数据库中存储的 Token 不匹配，则认证失败。
		if loginDetails == nil || (token != (*loginDetails).AuthToken) {
			log.Error(UnAuthorizedError)
			api.RequestErrorHandler(w, UnAuthorizedError)
			return
		}

		// 6. 放行请求
		// 认证通过，调用下一个处理函数 (Handler)。
		next.ServeHTTP(w, r)

	})
}
