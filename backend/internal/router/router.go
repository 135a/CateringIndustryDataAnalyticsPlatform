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

		// ================= 数据可视化统计路由 (服务于大屏 ECharts) =================
		statsGroup := apiGroup.Group("/statistics")
		{
			statsGroup.GET("/overview", api.Overview)                       // 指标卡
			statsGroup.GET("/category-pie", api.CategoryPie)                // 饼图
			statsGroup.GET("/district-bar", api.DistrictBar)                // 柱状图
			statsGroup.GET("/price-rating-scatter", api.PriceRatingScatter) // 散点图
			statsGroup.GET("/map-points", api.MapPoints)                    // 地图坐标点
		}

		// ================= 商户明细检索路由 =================
		// 分页查询商户列表 (支持排序与多条件筛选)
		apiGroup.GET("/restaurants", api.GetRestaurants)

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
