package router

import (
	"net/http"

	"catering-backend/internal/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter 初始化并配置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 配置跨域中间件 (开发环境允许所有来源跨域)
	r.Use(cors.Default())

	// 核心业务 API 路由组
	apiGroup := r.Group("/api/v1")
	{
		// 测试接口
		apiGroup.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": "Welcome to Catering Analytics API!"})
		})

		// ================= 字典与基础路由 =================
		// 获取行政区列表下拉框
		apiGroup.GET("/districts", api.GetDistricts)
		// 获取餐饮分类下拉框
		apiGroup.GET("/categories", api.GetCategories)

		// ================= 用户相关路由 =================
		userGroup := apiGroup.Group("/user")
		{
			// POST /api/v1/user/register
			userGroup.POST("/register", api.Register)
			// POST /api/v1/user/login
			userGroup.POST("/login", api.Login)
		}
	}

	return r
}
