package api

import (
	"net/http"

	"catering-backend/internal/model"
	"catering-backend/pkg/database"

	"github.com/gin-gonic/gin"
)

// RestaurantQuery 接受前端的分页与条件检索参数
type RestaurantQuery struct {
	Page           int    `form:"page,default=1"`       // 当前页码
	PageSize       int    `form:"page_size,default=15"` // 每页条数
	DistrictID     int    `form:"district_id"`          // 行政区过滤
	CategoryID     int    `form:"category_id"`          // 分类过滤
	Keyword        string `form:"keyword"`              // 模糊搜索关键字(商户名)
	SortBy         string `form:"sort_by"`              // 排序策略
	HasFreeParking bool   `form:"has_free_parking"`     // 是否有免费停车
	IsReservable   bool   `form:"is_reservable"`        // 是否可订座
	HasBabyChair   bool   `form:"has_baby_chair"`       // 是否有宝宝椅
	HasPrivateRoom bool   `form:"has_private_room"`     // 是否有包厢
}

// GetRestaurants 分页条件查询商户明细列表
func GetRestaurants(c *gin.Context) {
	var query RestaurantQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "分页参数解析失败"})
		return
	}

	// 开启一个 Model 查询实例
	db := database.DB.Model(&model.Restaurant{})

	// 1. 构建多条件动态检索 Where 子句
	if query.DistrictID > 0 {
		db = db.Where("district_id = ?", query.DistrictID)
	}
	if query.CategoryID > 0 {
		db = db.Where("category_id = ?", query.CategoryID)
	}
	if query.Keyword != "" {
		db = db.Where("name LIKE ?", "%"+query.Keyword+"%")
	}
	if query.HasFreeParking {
		db = db.Where("has_free_parking = ?", 1)
	}
	if query.IsReservable {
		db = db.Where("is_reservable = ?", 1)
	}
	if query.HasBabyChair {
		db = db.Where("has_baby_chair = ?", 1)
	}
	if query.HasPrivateRoom {
		db = db.Where("has_private_room = ?", 1)
	}

	// 2. 统计满足条件的总条数 (必须在 Order 和 Limit 之前执行，专为前端 Pagination 组件提供 total 属性)
	var total int64
	db.Count(&total)

	// 3. 构建动态排序策略
	switch query.SortBy {
	case "rating": // 评分最高优先，同分看人气
		db = db.Order("rating DESC, review_count DESC")
	case "review": // 人气最旺优先 (评价最多)
		db = db.Order("review_count DESC")
	case "price_asc": // 价格从低到高
		db = db.Order("avg_price ASC")
	case "price_desc": // 价格从高到低
		db = db.Order("avg_price DESC")
	default:
		// 默认排序：按数据库入库最新抓取时间排序
		db = db.Order("id DESC")
	}

	// 4. 执行分页截断查询
	offset := (query.Page - 1) * query.PageSize
	var restaurants []model.Restaurant
	db.Limit(query.PageSize).Offset(offset).Find(&restaurants)

	// 5. 组装标准分页响应体下发给前端
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": gin.H{
			"total":     total,
			"page":      query.Page,
			"page_size": query.PageSize,
			"list":      restaurants,
		},
	})
}
