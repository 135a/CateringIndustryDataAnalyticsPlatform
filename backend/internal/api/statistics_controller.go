package api

import (
	"net/http"

	"catering-backend/internal/model"
	"catering-backend/pkg/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StatsQuery 通用筛选参数
type StatsQuery struct {
	DistrictID int `form:"district_id"`
	CategoryID int `form:"category_id"`
}

// buildStatsQuery 抽取公共方法：根据前端传来的条件构造 GORM 的 Where 子句
func buildStatsQuery(c *gin.Context) *gorm.DB {
	var query StatsQuery
	c.ShouldBindQuery(&query) // 从 URL 解析 ?district_id=1&category_id=2

	db := database.DB.Model(&model.Restaurant{})
	if query.DistrictID > 0 {
		db = db.Where("restaurants.district_id = ?", query.DistrictID)
	}
	if query.CategoryID > 0 {
		db = db.Where("restaurants.category_id = ?", query.CategoryID)
	}
	return db
}

// Overview 全局数据指标卡 (总店数、均价、总人气)
func Overview(c *gin.Context) {
	db := buildStatsQuery(c)

	var result struct {
		TotalShops   int64   `json:"total_shops"`
		AvgPrice     float64 `json:"avg_price"`
		TotalReviews int64   `json:"total_reviews"`
		AvgRating    float64 `json:"avg_rating"`
	}

	// 聚合查询
	db.Select("COUNT(*) as total_shops, COALESCE(AVG(avg_price),0) as avg_price, COALESCE(SUM(review_count),0) as total_reviews, COALESCE(AVG(rating),0) as avg_rating").Scan(&result)

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": result})
}

// CategoryPie 品类占比分析 (供 ECharts 饼图使用)
func CategoryPie(c *gin.Context) {
	db := buildStatsQuery(c)

	type PieData struct {
		Name  string `json:"name"`
		Value int64  `json:"value"`
	}
	var results []PieData

	// 连表查询分类名，并根据分类分组 COUNT
	db.Select("categories.name as name, COUNT(restaurants.id) as value").
		Joins("JOIN categories ON restaurants.category_id = categories.id").
		Group("categories.id, categories.name").
		Order("value DESC").
		Limit(15). // 为保证饼图可读性，只取前15大分类
		Scan(&results)

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": results})
}

// DistrictBar 区域消费对比 (供 ECharts 柱状图使用)
func DistrictBar(c *gin.Context) {
	db := buildStatsQuery(c)

	type BarData struct {
		DistrictName string  `json:"district_name"`
		ShopCount    int64   `json:"shop_count"`
		AvgPrice     float64 `json:"avg_price"`
	}
	var results []BarData

	db.Select("districts.district_name as district_name, COUNT(restaurants.id) as shop_count, COALESCE(AVG(restaurants.avg_price),0) as avg_price").
		Joins("JOIN districts ON restaurants.district_id = districts.id").
		Group("districts.id, districts.district_name").
		Order("shop_count DESC").
		Scan(&results)

	// 为了完美契合 ECharts Bar 图表格式，将对象数组转换成分离的 Array
	districtsList := make([]string, 0)
	shopCountsList := make([]int64, 0)
	avgPricesList := make([]float64, 0)
	for _, r := range results {
		districtsList = append(districtsList, r.DistrictName)
		shopCountsList = append(shopCountsList, r.ShopCount)
		avgPricesList = append(avgPricesList, r.AvgPrice)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": gin.H{
		"xAxis":       districtsList,
		"shop_counts": shopCountsList,
		"avg_prices":  avgPricesList,
	}})
}

// PriceRatingScatter 性价比散点图 (供 ECharts Scatter 使用)
func PriceRatingScatter(c *gin.Context) {
	db := buildStatsQuery(c)

	type ScatterData struct {
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Rating   float64 `json:"rating"`
		Category string  `json:"category"`
	}
	var results []ScatterData

	db.Select("restaurants.name, restaurants.avg_price as price, restaurants.rating, categories.name as category").
		Joins("JOIN categories ON restaurants.category_id = categories.id").
		Where("restaurants.avg_price > 0 AND restaurants.rating > 0").
		Limit(3000). // 限制点数防止前端渲染卡顿
		Scan(&results)

	// ECharts 散点图的 dataset 数据格式要求通常是二维数组 [[x, y, 附加数据], ...]
	var echartsData [][]interface{}
	for _, r := range results {
		// [价格(X轴), 评分(Y轴), 店名, 分类]
		echartsData = append(echartsData, []interface{}{r.Price, r.Rating, r.Name, r.Category})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": echartsData})
}

// MapPoints 地理空间分布打点 (供 ECharts 百度/高德地图散点图使用)
func MapPoints(c *gin.Context) {
	db := buildStatsQuery(c)

	type MapPoint struct {
		Name        string  `json:"name"`
		Longitude   float64 `json:"longitude"`
		Latitude    float64 `json:"latitude"`
		ReviewCount int     `json:"review_count"`
	}
	var results []MapPoint

	db.Select("name, longitude, latitude, review_count").
		Where("longitude > 0 AND latitude > 0").
		Scan(&results)

	// 转换为 ECharts geo / scatter 所需的特定 JSON 结构
	var echartsData []gin.H
	for _, r := range results {
		echartsData = append(echartsData, gin.H{
			"name": r.Name,
			// ECharts value 数组一般第一项是经度，第二项纬度，第三项是权重(用于决定气泡大小)
			"value": []interface{}{r.Longitude, r.Latitude, r.ReviewCount},
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": echartsData})
}
