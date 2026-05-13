package models

import (
	"time"

	"gin-api/config"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email     string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Nickname  string         `gorm:"size:50" json:"nickname"`
	Status    int            `gorm:"default:1" json:"status"` // 1:正常 0:禁用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// AutoMigrate 自动迁移
func AutoMigrate() {
	err := config.DB.AutoMigrate(&User{})
	if err != nil {
		panic("数据库迁移失败: " + err.Error())
	}
}
