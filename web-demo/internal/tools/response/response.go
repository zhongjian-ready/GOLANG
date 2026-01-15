package response

import (
	"encoding/json"
	"net/http"
)

// Error 定义了通用的错误响应结构。
type Error struct {
	Code    int
	Message string
}

// WriteError 是一个辅助函数，用于统一向客户端写入错误响应。
func WriteError(w http.ResponseWriter, message string, code int) {
	resp := Error{
		Code:    code,
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}

// RequestErrorHandler 用于处理请求相关的错误。
var RequestErrorHandler = func(w http.ResponseWriter, err error) {
	WriteError(w, err.Error(), http.StatusBadRequest)
}

// InternalErrorHandler 用于处理服务器内部错误。
var InternalErrorHandler = func(w http.ResponseWriter) {
	WriteError(w, "An Unexpected Error Occurred.", http.StatusInternalServerError)
}
