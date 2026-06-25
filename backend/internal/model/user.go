package model

import "time"

// User 与 MySQL 中的 users 表对应
type User struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string    `gorm:"type:varchar(50);not null;unique" json:"username"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"` // json:"-" 保证密码哈希永远不会泄漏给前端
	Role         int8      `gorm:"type:tinyint;default:0" json:"role"`  // 0: 普通用户, 1: 管理员
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
