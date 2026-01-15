package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
	"github.com/zhongjian-ready/goapi/api"
	"github.com/zhongjian-ready/goapi/internal/tools"
)

// GetCoinBalance 是处理获取用户硬币余额请求的 Handler 函数。
// 对应的路由是: GET /account/coins
func GetCoinBalance(w http.ResponseWriter, r *http.Request) {
	// 1. 定义请求参数结构体
	// 使用 api 包中定义的 CoinBalanceParam 结构体来接收参数。
	var params = api.CoinBalanceParam{}

	// 2. 解析请求参数
	// 创建一个 gorilla/schema 的解码器实例。
	var decoder *schema.Decoder = schema.NewDecoder()
	var err error

	// 尝试将 URL 查询参数（Query Parameters）解码到 params 结构体中。
	// 例如：?username=alex 会被映射到 params.UserName。
	err = decoder.Decode(&params, r.URL.Query())

	if err != nil {
		// 如果解码失败，记录错误日志，并返回客户端错误响应 (400 Bad Request)。
		log.Error("Failed to decode request parameters:", err)
		api.RequestErrorHandler(w, err)
		return
	}

	// 3. 连接数据库
	// 使用 tools 包的 NewDatabase 工厂函数创建一个数据库接口实例。
	var database *tools.DatabaseInterface
	database, err = tools.NewDatabase()

	if err != nil {
		// 如果数据库连接/初始化失败，返回内部服务器错误 (500 Internal Server Error)。
		// 注意：具体的错误信息在 NewDatabase 内部可能已经被记录了，或者这里可以补充记录。
		api.InternalErrorHandler(w)
		return
	}

	// 4. 查询数据
	// 调用数据库接口的 GetUserCoins 方法，根据用户名查询硬币详情。
	var tokenDetails *tools.CoinDetails
	tokenDetails = (*database).GetUserCoins(params.UserName)

	if tokenDetails == nil {
		// 如果查询结果为空（用户不存在或没有数据），记录错误并返回。
		// 这里虽然叫 RequestErrorHandler，但可能在该上下文中也用于表示资源未找到或逻辑错误。
		log.Error("Failed to get coin details for user:", params.UserName)
		api.RequestErrorHandler(w, err)
		return
	}

	// 5. 构建响应
	// 创建 api.CoinBalanceResponse 结构体，填充查询到的余额信息。
	var response = api.CoinBalanceResponse{
		Code:    http.StatusOK, // HTTP 200
		Balance: (*tokenDetails).Balance,
	}

	// 6. 发送响应
	// 设置 Content-Type 头部为 application/json。
	w.Header().Set("Content-Type", "application/json")

	// 将响应结构体编码为 JSON 并写入 ResponseWriter。
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		// 如果发送响应失败（例如网络断开），记录错误并尝试发送 500 错误。
		log.Error("Failed to encode response:", err)
		api.InternalErrorHandler(w)
		return
	}

}
