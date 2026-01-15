package models

// CoinBalanceParam 定义了获取硬币余额接口的请求参数。
// UserID: 用户ID，对应数据库中的用户标识。
type CoinBalanceParam struct {
	UserID int `schema:"userid"`
}

// LoginParams 定义了登录接口的请求参数
type LoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse 定义了登录成功的响应
type LoginResponse struct {
	Code  int    `json:"code"`
	Token string `json:"token"`
}

// CoinBalanceResponse 定义了获取硬币余额接口的响应结构。
// Code: 业务状态码，通常 200 表示成功。
// Balance: 用户的硬币余额，使用 int64 防止溢出。
type CoinBalanceResponse struct {
	Code    int
	Balance int64
}
