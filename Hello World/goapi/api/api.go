package api

import (
	"encoding/json"
	"net/http"
)

// CoinBalanceParam 定义了获取硬币余额接口的请求参数。
// 虽然这个接口主要通过 Authorization 头获取用户信息，但这里预留了查询参数结构。
// UserID: 用户ID，对应数据库中的用户标识。
type CoinBalanceParam struct {
	UserID int `schema:"userid"`
}

// CoinBalanceResponse 定义了获取硬币余额接口的响应结构。
// Code: 业务状态码，通常 200 表示成功。
// Balance: 用户的硬币余额，使用 int64 防止溢出。
type CoinBalanceResponse struct {
	Code    int
	Balance int64
}

// Error 定义了通用的错误响应结构。
// 当 API 发生错误时，会返回这个 JSON 结构给客户端。
// Code: 错误状态码。
// Message: 具体的错误描述信息。
type Error struct {
	Code    int
	Message string
}

// writeError 是一个辅助函数，用于统一向客户端写入错误响应。
// 它可以避免在每个 handler 里重复写如下代码：设置 Header、设置 Status Code、编码 JSON。
// w: HTTP 响应写入器。
// message: 错误提示信息。
// code: HTTP 状态码（如 400, 500）。
func writeError(w http.ResponseWriter, message string, code int) {
	resp := Error{
		Code:    code,
		Message: message,
	}
	// 设置 Content-Type 为 JSON，因为我们返回的是 JSON 格式的错误信息。
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	// 将 Error 结构体编码为 JSON 并写入响应体。
	json.NewEncoder(w).Encode(resp)
}

var (
	// RequestErrorHandler 用于处理请求相关的错误（如参数解析失败）。
	// 它会返回 400 Bad Request 状态码，并附带具体的错误信息。
	RequestErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusBadRequest)
	}

	// InternalErrorHandler 用于处理服务器内部的 unexpected 错误（如数据库连接失败）。
	// 为了安全起见，通常不向客户端暴露具体的内部错误细节，而是返回一个通用的提示。
	// 这里返回 500 Internal Server Error。
	InternalErrorHandler = func(w http.ResponseWriter) {
		writeError(w, "An Unexpected Error Occurred.", http.StatusInternalServerError)
	}
)
