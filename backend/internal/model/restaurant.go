package model

import "time"

// Restaurant 餐饮商户实体，与 MySQL 中的 restaurants 表对应
type Restaurant struct {
	ID               int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ShopID           string    `gorm:"type:varchar(100);unique" json:"shop_id"`
	Name             string    `gorm:"type:varchar(255);not null" json:"name"`
	CategoryID       int       `gorm:"not null" json:"category_id"`
	DistrictID       int       `gorm:"not null" json:"district_id"`
	Address          string    `gorm:"type:varchar(500)" json:"address"`
	AvgPrice         float64   `gorm:"type:decimal(10,2)" json:"avg_price"`
	Rating           float64   `gorm:"type:decimal(3,1)" json:"rating"`
	ReviewCount      int       `gorm:"default:0" json:"review_count"`
	OpeningHours     string    `gorm:"type:varchar(255)" json:"opening_hours"`
	TasteScore       float64   `gorm:"type:decimal(3,1)" json:"taste_score"`
	EnvironmentScore float64   `gorm:"type:decimal(3,1)" json:"environment_score"`
	ServiceScore     float64   `gorm:"type:decimal(3,1)" json:"service_score"`
	HasFreeParking   bool      `gorm:"type:tinyint(1);default:0" json:"has_free_parking"`
	IsReservable     bool      `gorm:"type:tinyint(1);default:0" json:"is_reservable"`
	HasBabyChair     bool      `gorm:"type:tinyint(1);default:0" json:"has_baby_chair"`
	HasPrivateRoom   bool      `gorm:"type:tinyint(1);default:0" json:"has_private_room"`
	CrawledAt        time.Time `json:"crawled_at"`
}
