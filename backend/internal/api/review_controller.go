package api

import (
	"net/http"

	"catering-backend/internal/model"
	"catering-backend/pkg/database"

	"github.com/gin-gonic/gin"
)

// ReviewQuery 接受前端传来的分页参数
type ReviewQuery struct {
	Page         int   `form:"page,default=1"`       // 当前页码
	PageSize     int   `form:"page_size,default=15"` // 每页条数
	RestaurantID int64 `form:"restaurant_id"`        // 可选：只看特定商户的评价
}

// GetReviews 分页获取评价明细列表（支持附带所属餐饮店信息）
func GetReviews(c *gin.Context) {
	var query ReviewQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "分页参数解析失败"})
		return
	}

	db := database.DB.Model(&model.Review{})

	// 如果传了商户ID，就查某个商户的评价；否则查全市的评价
	if query.RestaurantID > 0 {
		db = db.Where("restaurant_id = ?", query.RestaurantID)
	}

	// 统计这批条件下一共有多少条评价
	var total int64
	db.Count(&total)

	// 排序：最新的评价永远在最上面
	db = db.Order("created_at DESC")

	// 【重头戏】: 执行分页查询，并利用 Preload("Restaurant") 将 MySQL 中的主外键关系直接映射为嵌套的 JSON
	offset := (query.Page - 1) * query.PageSize
	var reviews []model.Review
	db.Preload("Restaurant").Limit(query.PageSize).Offset(offset).Find(&reviews)

	// 组装标准分页响应体下发给前端
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": gin.H{
			"total":     total,
			"page":      query.Page,
			"page_size": query.PageSize,
			"list":      reviews,
		},
	})
}
