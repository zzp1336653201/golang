package model

import (
	"time"
)

// Product 商品
type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Price       int64     `gorm:"type:decimal(10,2);not null" json:"price"` // 分为单位
	Stock       int       `gorm:"default:0" json:"stock"`                   // 库存
	Category    string    `gorm:"size:100" json:"category"`
	ImageURL    string    `gorm:"size:500" json:"image_url"`
	Status      int8      `gorm:"default:1" json:"status"` // 1:上架 0:下架
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Product) TableName() string {
	return "products"
}
