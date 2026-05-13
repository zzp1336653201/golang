package service

import (
	"context"
	"fmt"
	"time"

	"cms/internal/model"
	"cms/pkg/database"

	"gorm.io/gorm"
)

// OrderService 订单服务
type OrderService struct {
	db             *gorm.DB
	productService *ProductService
}

// NewOrderService 创建订单服务
func NewOrderService(db *gorm.DB, productService *ProductService) *OrderService {
	return &OrderService{
		db:             db,
		productService: productService,
	}
}

// CreateOrder 创建订单
// 【面试重点】订单创建流程：
// 1. 校验商品和库存
// 2. 扣减库存（分布式锁）
// 3. 创建订单和订单项
// 4. 整个过程需要事务保证原子性
func (s *OrderService) CreateOrder(ctx context.Context, userID uint, items []struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}, receiver model.ReceiverInfo) (*model.Order, error) {
	
	// 开启事务
	tx := database.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	
	// 计算订单总额并构建订单项
	var totalAmount int64
	orderItems := make([]model.OrderItem, 0, len(items))
	
	for _, item := range items {
		// 查询商品信息
		var product model.Product
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&product, item.ProductID).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("商品不存在: %d", item.ProductID)
		}
		
		// 检查库存（行锁保证一致性）
		if product.Stock < item.Quantity {
			tx.Rollback()
			return nil, fmt.Errorf("商品 %s 库存不足，需要: %d, 剩余: %d", product.Name, item.Quantity, product.Stock)
		}
		
		// 扣减库存
		if err := tx.Model(&model.Product{}).
			Where("id = ?", item.ProductID).
			Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("库存扣减失败")
		}
		
		// 计算小计
		subTotal := product.Price * int64(item.Quantity)
		totalAmount += subTotal
		
		orderItems = append(orderItems, model.OrderItem{
			ProductID:   product.ID,
			ProductName: product.Name,
			Price:       product.Price,
			Quantity:    item.Quantity,
			SubTotal:    subTotal,
		})
	}
	
	// 创建订单
	order := &model.Order{
		OrderNo:       model.GenerateOrderNo(),
		UserID:        userID,
		TotalAmount:   totalAmount,
		Status:        model.OrderStatusPending, // 待支付
		ReceiverName:  receiver.Name,
		ReceiverPhone: receiver.Phone,
		ReceiverAddr:  receiver.Address,
		Remark:        receiver.Remark,
	}
	
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("订单创建失败: %w", err)
	}
	
	// 创建订单项
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
	}
	if err := tx.Create(&orderItems).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("订单项创建失败: %w", err)
	}
	
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("事务提交失败: %w", err)
	}
	
	// 查询完整订单
	s.db.Preload("Items").First(order, order.ID)
	
	return order, nil
}

// PayOrder 支付订单
func (s *OrderService) PayOrder(ctx context.Context, orderNo string) error {
	now := time.Now()
	result := s.db.WithContext(ctx).Model(&model.Order{}).
		Where("order_no = ? AND status = ?", orderNo, model.OrderStatusPending).
		Updates(map[string]interface{}{
			"status":   model.OrderStatusPaid,
			"paid_at":  now,
		})
	
	if result.RowsAffected == 0 {
		return fmt.Errorf("订单不存在或状态异常")
	}
	
	return result.Error
}

// CancelOrder 取消订单
func (s *OrderService) CancelOrder(ctx context.Context, orderNo string, userID uint) error {
	// 查询订单
	var order model.Order
	if err := s.db.WithContext(ctx).
		Where("order_no = ? AND user_id = ?", orderNo, userID).
		First(&order).Error; err != nil {
		return fmt.Errorf("订单不存在")
	}
	
	// 只有待支付订单可以取消
	if order.Status != model.OrderStatusPending {
		return fmt.Errorf("订单状态不允许取消")
	}
	
	// 开启事务回补库存
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	
	// 查询订单项
	var items []model.OrderItem
	if err := tx.Where("order_id = ?", order.ID).Find(&items).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// 回补库存
	for _, item := range items {
		if err := tx.Model(&model.Product{}).
			Where("id = ?", item.ProductID).
			Update("stock", gorm.Expr("stock + ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	
	// 更新订单状态
	now := time.Now()
	if err := tx.Model(&order).Updates(map[string]interface{}{
		"status":      model.OrderStatusCanceled,
		"canceled_at": now,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	return tx.Commit().Error
}

// GetUserOrders 获取用户订单列表
func (s *OrderService) GetUserOrders(ctx context.Context, userID uint, page, pageSize int) ([]*model.Order, int64, error) {
	var orders []*model.Order
	var total int64
	
	query := s.db.WithContext(ctx).Model(&model.Order{}).Where("user_id = ?", userID)
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	offset := (page - 1) * pageSize
	if err := query.Preload("Items").
		Offset(offset).Limit(pageSize).
		Order("id desc").
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	
	return orders, total, nil
}

// OrderStatusXXX 订单状态常量
const (
	OrderStatusPending   = 1 // 待支付
	OrderStatusPaid      = 2 // 已支付
	OrderStatusShipped   = 3 // 已发货
	OrderStatusReceived  = 4 // 已收货
	OrderStatusCompleted = 5 // 已完成
	OrderStatusCanceled  = 6 // 已取消
)

// ReceiverInfo 收货信息
type ReceiverInfo struct {
	Name    string `json:"name" binding:"required"`
	Phone   string `json:"phone" binding:"required"`
	Address string `json:"address" binding:"required"`
	Remark  string `json:"remark"`
}
