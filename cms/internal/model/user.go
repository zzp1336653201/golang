package model

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User 用户
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Email     string    `gorm:"size:100" json:"email"`
	Phone     string    `gorm:"size:20" json:"phone"`
	Nickname  string    `gorm:"size:100" json:"nickname"`
	Avatar    string    `gorm:"size:500" json:"avatar"`
	Status    int8      `gorm:"default:1" json:"status"` // 1:正常 0:禁用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

// SetPassword 加密密码
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// Order 订单
type Order struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	OrderNo      string    `gorm:"size:32;uniqueIndex;not null" json:"order_no"`
	UserID       uint      `gorm:"index;not null" json:"user_id"`
	TotalAmount  int64     `gorm:"type:decimal(10,2);not null" json:"total_amount"` // 订单总额（分）
	Status       int8      `gorm:"default:1" json:"status"`                        // 1:待支付 2:已支付 3:已发货 4:已收货 5:已完成 6:已取消
	ReceiverName string    `gorm:"size:100" json:"receiver_name"`
	ReceiverPhone string   `gorm:"size:20" json:"receiver_phone"`
	ReceiverAddr  string   `gorm:"type:text" json:"receiver_addr"`
	Remark       string    `gorm:"type:text" json:"remark"`
	PaidAt       *time.Time `json:"paid_at"`
	ShippedAt    *time.Time `json:"shipped_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	CanceledAt   *time.Time `json:"canceled_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
	// 关联
	Items []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	User  *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}

// OrderItem 订单项
type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `gorm:"index;not null" json:"order_id"`
	ProductID uint    `gorm:"index;not null" json:"product_id"`
	ProductName string `gorm:"size:255" json:"product_name"`
	Price     int64   `gorm:"type:decimal(10,2);not null" json:"price"` // 购买时价格
	Quantity  int     `gorm:"not null" json:"quantity"`
	SubTotal  int64   `gorm:"type:decimal(10,2);not null" json:"sub_total"` // 小计
	
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (OrderItem) TableName() string {
	return "order_items"
}

// CartItem 购物车项
type CartItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	ProductID uint      `gorm:"index;not null" json:"product_id"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (CartItem) TableName() string {
	return "cart_items"
}

// GenerateOrderNo 生成订单号
func GenerateOrderNo() string {
	return time.Now().Format("20060102150405") + uuid.New().String()[:8]
}
