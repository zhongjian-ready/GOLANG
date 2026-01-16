package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"github.com/zhongjian-ready/goapi/internal/config"
	"github.com/zhongjian-ready/goapi/internal/database"
	"github.com/zhongjian-ready/goapi/internal/handlers"
)

func main() {
	// 0. 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Warn("Error loading .env file")
	}

	// 1. 加载配置 (Fail Fast)
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// 解析命令行参数
	// 使用 -migrate 参数来触发数据库表结构初始化
	migrate := flag.Bool("migrate", false, "Run database migration (create tables)")
	flag.Parse()

	// 2. 设置日志配置
	// 配置 logrus 日志库，让它在打印日志时，额外显示是哪个函数、哪个文件调用的。
	// 这对于排查问题非常有用（比如会显示: main.main at /path/to/main.go:14）。
	log.SetReportCaller(true)

	// 初始化数据库连接 (全局单例)
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 如果指定了 -migrate 参数，执行建表逻辑
	if *migrate {
		log.Info("Running database migration...")
		// 之前是 (*db).SetupSchema()，现在可以直接调用接口方法
		if err := db.SetupSchema(); err != nil {
			log.Fatal("Failed to setup schema:", err)
		}
		log.Info("Database migration completed successfully.")
	}

	// 2. 初始化路由 (Router)
	// 创建一个新的 Chi 路由器实例。Chi 是一个轻量级、高性能的路由库。
	// *chi.Mux 是路由器的指针类型，它负责把收到的 HTTP 请求分发给对应的处理函数。
	var r *chi.Mux = chi.NewRouter()

	// 3. 注册路由规则
	// 调用 internal/handlers 包里的 Handler 函数。
	// 将数据库连接和配置注入到 Handler 中
	handlers.Handler(r, db, cfg)

	// 4. 打印启动欢迎语
	// 简单的控制台输出，告诉开发者服务正在启动。
	fmt.Println("Start go api service...")

	// 打印刚才我们生成的那个帅气的 ASCII Art 图案。
	fmt.Print(`
   ______  ____    ___    ____  ____
  / ____/ / __ \  /   |  / __ \/  _/
 / / __  / / / / / /| | / /_/ // /  
/ /_/ / / /_/ / / ___ |/ ____// /   
\____/  \____/ /_/  |_/_/   /___/   
`)

	// 5. 启动 HTTP 服务 (核心)
	// 使用自定义的 Server 结构体来支持优雅关闭
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	addr := ":" + port

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 启动一个 goroutine 来监听服务
	go func() {
		fmt.Printf("Starting server on port %s...\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 监听中断信号 (Ctrl+C, SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// 创建一个 5 秒超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 尝试优雅关闭
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Info("Server exiting")
}
