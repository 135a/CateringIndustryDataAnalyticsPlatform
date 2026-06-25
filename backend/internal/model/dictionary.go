package model

import "time"

// District 行政区映射模型
type District struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	CityName     string    `gorm:"column:city_name" json:"city_name"`
	DistrictName string    `gorm:"column:district_name" json:"district_name"`
	CreatedAt    time.Time `json:"-"` // 返回前端时忽略该字段
}

// Category 餐饮分类映射模型
type Category struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	SortOrder int       `gorm:"column:sort_order" json:"sort_order"` // 供前端下拉框排序使用
	CreatedAt time.Time `json:"-"`
}
