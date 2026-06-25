package router

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter 初始化并配置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 配置跨域中间件 (开发环境允许所有来源跨域)
	r.Use(cors.Default())

	// 测试服务是否存活
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// 核心业务 API 路由组
	api := r.Group("/api/v1")
	{
		// TODO: 稍后在这里挂载字典、统计、商户明细等接口
		api.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "success",
				"data": "Welcome to Catering Analytics API!",
			})
		})
	}

	return r
}
