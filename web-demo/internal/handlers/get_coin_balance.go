package handlers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/zhongjian-ready/goapi/api"
	"github.com/zhongjian-ready/goapi/internal/tools"
)

// GetCoinBalance 是处理获取用户硬币余额请求的 Handler 函数。
// 对应的路由是: GET /account/coins
func GetCoinBalance(w http.ResponseWriter, r *http.Request) {
	// 1. 从 Token 中获取 UserID
	// 鉴权中间件已经验证了 Token 并将 userid 放入了 Context
	userID, ok := r.Context().Value("userid").(int)
	if !ok {
		log.Error("User ID missing from context")
		api.InternalErrorHandler(w)
		return
	}

	// 2. 连接数据库
	// 使用 tools 包的 NewDatabase 工厂函数创建一个数据库接口实例。
	var database *tools.DatabaseInterface
	var err error
	database, err = tools.NewDatabase()

	if err != nil {
		// 如果数据库连接/初始化失败，返回内部服务器错误 (500 Internal Server Error)。
		api.InternalErrorHandler(w)
		return
	}

	// 3. 查询数据
	// 调用数据库接口的 GetUserCoins 方法，根据 UserID 查询硬币详情。
	var tokenDetails *tools.CoinDetails

	// 使用从 Token 提取的 userID
	tokenDetails = (*database).GetUserCoins(userID)

	if tokenDetails == nil {
		// 如果查询结果为空（用户不存在或没有数据），记录错误并返回。
		log.Error("Failed to get coin details for user:", userID)
		api.InternalErrorHandler(w)
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
