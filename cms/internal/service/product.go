package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cms/internal/model"
	"cms/pkg/cache"

	"gorm.io/gorm"
)

// ProductService 商品服务
type ProductService struct {
	db *gorm.DB
}

// NewProductService 创建商品服务
func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{db: db}
}

// Create 创建商品
func (s *ProductService) Create(ctx context.Context, product *model.Product) error {
	return s.db.WithContext(ctx).Create(product).Error
}

// GetByID 根据ID获取商品（带缓存）
// 【面试重点】缓存穿透、击穿、雪崩处理
func (s *ProductService) GetByID(ctx context.Context, id uint) (*model.Product, error) {
	key := fmt.Sprintf(cache.ProductCache, id)
	
	// 1. 先查缓存
	data, err := cache.Get(ctx, key)
	if err == nil && data != "" && data != "__nil__" {
		var product model.Product
		if json.Unmarshal([]byte(data), &product) == nil {
			return &product, nil
		}
	}
	
	// 2. 缓存未命中，查数据库
	var product model.Product
	if err := s.db.WithContext(ctx).First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 【防止缓存穿透】- 空值缓存
			cache.Set(ctx, key, "__nil__", 60*time.Second)
			return nil, nil
		}
		return nil, err
	}
	
	// 3. 写入缓存
	if jsonData, err := json.Marshal(product); err == nil {
		cache.Set(ctx, key, string(jsonData), 10*time.Minute)
	}
	
	return &product, nil
}

// List 分页查询商品
func (s *ProductService) List(ctx context.Context, page, pageSize int, category string) ([]*model.Product, int64, error) {
	var products []*model.Product
	var total int64
	
	query := s.db.WithContext(ctx).Model(&model.Product{})
	if category != "" {
		query = query.Where("category = ?", category)
	}
	
	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id desc").Find(&products).Error; err != nil {
		return nil, 0, err
	}
	
	return products, total, nil
}

// Update 更新商品
func (s *ProductService) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	result := s.db.WithContext(ctx).Model(&model.Product{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	
	// 删除缓存
	cache.Delete(ctx, fmt.Sprintf(cache.ProductCache, id))
	
	return nil
}

// Delete 删除商品
func (s *ProductService) Delete(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&model.Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	
	// 删除缓存
	cache.Delete(ctx, fmt.Sprintf(cache.ProductCache, id))
	
	return nil
}

// GetStock 获取库存
func (s *ProductService) GetStock(ctx context.Context, productID uint) (int, error) {
	var product model.Product
	if err := s.db.WithContext(ctx).Select("stock").First(&product, productID).Error; err != nil {
		return 0, err
	}
	return product.Stock, nil
}

// DecreaseStock 扣减库存
// 【面试重点】库存扣减需要考虑：
// 1. 分布式锁防止超卖
// 2. 乐观锁/悲观锁
// 3. 库存不足时回滚
func (s *ProductService) DecreaseStock(ctx context.Context, productID uint, quantity int) error {
	// 1. 先加分布式锁
	locked, err := cache.LockStock(ctx, productID)
	if err != nil {
		return fmt.Errorf("获取库存锁失败: %w", err)
	}
	if !locked {
		return fmt.Errorf("系统繁忙，请稍后重试")
	}
	defer cache.UnlockStock(ctx, productID) // 确保释放锁
	
	// 2. 开启事务
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	
	// 3. 查询当前库存（带行锁）
	var product model.Product
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&product, productID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("商品不存在")
	}
	
	// 4. 检查库存
	if product.Stock < quantity {
		tx.Rollback()
		return fmt.Errorf("库存不足，当前库存: %d", product.Stock)
	}
	
	// 5. 扣减库存
	newStock := product.Stock - quantity
	if err := tx.Model(&model.Product{}).Where("id = ?", productID).Update("stock", newStock).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("库存扣减失败: %w", err)
	}
	
	// 6. 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("事务提交失败: %w", err)
	}
	
	// 7. 删除缓存
	cache.Delete(ctx, fmt.Sprintf(cache.ProductCache, productID))
	
	return nil
}

// IncreaseStock 回补库存（退款、退货时使用）
func (s *ProductService) IncreaseStock(ctx context.Context, productID uint, quantity int) error {
	result := s.db.WithContext(ctx).Model(&model.Product{}).
		Where("id = ?", productID).
		Update("stock", gorm.Expr("stock + ?", quantity))
	
	if result.Error != nil {
		return result.Error
	}
	
	// 删除缓存
	cache.Delete(ctx, fmt.Sprintf(cache.ProductCache, productID))
	
	return nil
}
