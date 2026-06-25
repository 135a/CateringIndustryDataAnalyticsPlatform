package model

import "time"

// Review 对应 MySQL 中的 reviews 表
type Review struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	RestaurantID int64     `gorm:"not null" json:"restaurant_id"`
	UserName     string    `gorm:"type:varchar(50);default:匿名用户" json:"user_name"`
	Rating       float64   `gorm:"type:decimal(3,1)" json:"rating"`
	Content      string    `gorm:"type:text" json:"content"`
	CreatedAt    time.Time `json:"created_at"`

	// GORM 连表外键配置：一个 Review 属于一个 Restaurant。
	// 这使得我们在查询 Review 时，可以自动查出对应的那家店的名字等信息。
	Restaurant Restaurant `gorm:"foreignKey:RestaurantID" json:"restaurant"`
}
