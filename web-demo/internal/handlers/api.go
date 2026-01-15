package handlers

import (
	"net/http"
	"time" // 引入 time 包

	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"

	"github.com/zhongjian-ready/goapi/internal/config"
	"github.com/zhongjian-ready/goapi/internal/database"
	"github.com/zhongjian-ready/goapi/internal/middleware"
)

// API 结构体用于持有数据库连接等依赖
// 之前命名为 Repository 有歧义，Repository 通常指数据访问层。
// 而这里是 Handler 层，所以改名为 API 或 ServiceContainer 更合适。
type API struct {
	DB     database.DatabaseInterface
	Config *config.Config
}

// NewAPI 创建一个新的 API 实例
func NewAPI(db database.DatabaseInterface, cfg *config.Config) *API {
	return &API{
		DB:     db,
		Config: cfg,
	}
}

// Handler 负责注册整个应用的所有路由。
// 它接收一个 *chi.Mux 指针，并在其上定义 URL 路径与处理函数的映射关系。
func Handler(r *chi.Mux, db database.DatabaseInterface, cfg *config.Config) {
	// 初始化 API 实例
	api := NewAPI(db, cfg)

	// === 全局中间件配置 ===
	// 1. Recoverer: 捕获 Panic，防止服务因 handler 内部错误而崩溃
	// 还可以返回一个友好的 500 页面而不是空白连接断开
	r.Use(chimiddle.Recoverer)

	// 2. Logger: 记录每个 HTTP 请求的详细日志 (Method, URL, Status, Duration)
	r.Use(chimiddle.Logger)

	// 3. StripSlashes: 自动移除 URL 末尾的斜杠
	r.Use(chimiddle.StripSlashes)

	// 健康检查接口 (Health Check)
	// 用于 K8s Liveness Probe 或负载均衡器心跳检测
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 定义路由组 "/account"
	// 所有以 "/account" 开头的请求都会进入这个闭包内配置的规则。
	r.Route("/account", func(router chi.Router) {
		// 1. 公开接口 (不需要鉴权)
		router.Post("/login", api.Login)

		// [测试用] 模拟耗时请求验证优雅关闭
		// 访问此接口后立即按 Ctrl+C，你会发现服务器会等待。
		router.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
			// 模拟一个需要 3 秒才能处理完的任务
			time.Sleep(3 * time.Second)
			w.Write([]byte("Task finished!"))
		})

		// 2. 受保护接口 (需要鉴权)
		router.Group(func(r chi.Router) {
			// 路由组中间件：Authorization
			// 在这个路由组下的所有请求，都会先经过 middleware.Authorization 中间件。
			// 使用闭包注入配置中的 JWT Secret
			r.Use(middleware.NewAuthorization(cfg.JWTSecret))

			// 路由规则：GET /account/coins
			// 当收到 GET 方法的 "/account/coins" 请求时，调用 GetCoinBalance 函数处理。
			r.Get("/coins", api.GetCoinBalance)
		})
	})
}
