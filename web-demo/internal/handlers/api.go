package handlers

import (
	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"

	"github.com/zhongjian-ready/goapi/internal/middleware"
)

// Handler 负责注册整个应用的所有路由。
// 它接收一个 *chi.Mux 指针，并在其上定义 URL 路径与处理函数的映射关系。
func Handler(r *chi.Mux) {
	// 中间件：StripSlashes
	// 这个中间件会自动处理 URL 末尾的斜杠。
	// 例如，如果用户访问 "/account/coins/"，它会自动重定向或处理为 "/account/coins"。
	// 这使得路由匹配更加灵活。
	r.Use(chimiddle.StripSlashes)

	// 定义路由组 "/account"
	// 所有以 "/account" 开头的请求都会进入这个闭包内配置的规则。
	r.Route("/account", func(router chi.Router) {
		// 1. 公开接口 (不需要鉴权)
		router.Post("/login", Login)

		// 2. 受保护接口 (需要鉴权)
		router.Group(func(r chi.Router) {
			// 路由组中间件：Authorization
			// 在这个路由组下的所有请求，都会先经过 middleware.Authorization 中间件。
			r.Use(middleware.Authorization)

			// 路由规则：GET /account/coins
			// 当收到 GET 方法的 "/account/coins" 请求时，调用 GetCoinBalance 函数处理。
			r.Get("/coins", GetCoinBalance)
		})
	})
}
