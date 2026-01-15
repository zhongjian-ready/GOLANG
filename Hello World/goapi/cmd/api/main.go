package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"github.com/zhongjian-ready/goapi/internal/handlers"
	"github.com/zhongjian-ready/goapi/internal/tools"
)

func main() {
	// 0. 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Warn("Error loading .env file")
	}

	// 解析命令行参数
	// 使用 -migrate 参数来触发数据库表结构初始化
	migrate := flag.Bool("migrate", false, "Run database migration (create tables)")
	flag.Parse()

	// 1. 设置日志配置
	// 配置 logrus 日志库，让它在打印日志时，额外显示是哪个函数、哪个文件调用的。
	// 这对于排查问题非常有用（比如会显示: main.main at /path/to/main.go:14）。
	log.SetReportCaller(true)

	// 如果指定了 -migrate 参数，执行建表逻辑
	if *migrate {
		log.Info("Running database migration...")
		db, err := tools.NewDatabase()
		if err != nil {
			log.Fatal("Failed to connect to database for migration:", err)
		}
		// 注意：NewDatabase 返回的是 *DatabaseInterface，这是一个指向接口的指针。
		// 在 Go 中，接口本身已经包含了底层数据的类型和值，通常不需要使用指向接口的指针。
		// 虽然这里定义稍微有点冗余（返回了 *DatabaseInterface），但我们需要解引用来调用接口方法。
		if err := (*db).SetupSchema(); err != nil {
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
	// 这里是你定义业务逻辑的地方，比如 "/account/coins" 这个 URL 应该由谁来处理，
	// 都需要在这个函数里“挂载”到路由器 r 上。
	handlers.Handler(r)

	// 4. 打印启动欢迎语
	// 简单的控制台输出，告诉开发者服务正在启动。
	fmt.Println("Start go api service...")

	// 打印刚才我们生成的那个帅气的 ASCII Art 图案。
	fmt.Println(`
   ______  ____    ___    ____  ____
  / ____/ / __ \  /   |  / __ \/  _/
 / / __  / / / / / /| | / /_/ // /  
/ /_/ / / /_/ / / ___ |/ ____// /   
\____/  \____/ /_/  |_/_/   /___/   
`)

	// 5. 启动 HTTP 服务 (核心)
	// 这是一个阻塞调用，程序运行到这里就会一直“卡住”监听。
	// "localhost:8000"：指定监听本地的 8000 端口。
	// r：把刚才配置好规则的路由器传进去，作为请求处理器。
	// 如果启动失败（比如端口被占用），它会返回一个 error。
	err := http.ListenAndServe("localhost:8000", r)

	// ... 后面的代码用于处理启动失败的情况
	if err != nil {
		log.Error(err)
	}
}
