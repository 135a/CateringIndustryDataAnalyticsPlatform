package main

import (
	"log"
	"os"

	"catering-backend/internal/router"
	"catering-backend/pkg/database"

	"github.com/joho/godotenv"
)

func main() {
	// 1. 加载环境变量 (读取项目根目录下的 .env)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using default environment variables.")
	}

	// 2. 初始化数据库连接
	database.InitDB()

	// 3. 初始化并装载 Gin 路由
	r := router.SetupRouter()

	// 4. 启动 HTTP 服务
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // 默认端口
	}

	log.Printf("🚀 Server is starting on http://localhost:%s ...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}
