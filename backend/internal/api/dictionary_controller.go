package api

import (
	"net/http"

	"catering-backend/internal/model"
	"catering-backend/pkg/database"

	"github.com/gin-gonic/gin"
)

// GetDistricts 获取所有行政区列表 (用于前端下拉框过滤)
func GetDistricts(c *gin.Context) {
	var districts []model.District

	// 从数据库查询所有行政区，按 ID 升序排列
	if err := database.DB.Order("id asc").Find(&districts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取行政区列表失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": districts,
	})
}

// GetCategories 获取所有餐饮分类列表 (用于前端下拉框过滤)
func GetCategories(c *gin.Context) {
	var categories []model.Category

	// 从数据库查询所有分类，优先按 sort_order 排序，其次按 ID 排序
	if err := database.DB.Order("sort_order asc, id asc").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取分类列表失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": categories,
	})
}
