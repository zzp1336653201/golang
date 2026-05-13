package cache

import (
	"context"
	"fmt"
	"time"

	"cms/internal/config"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func InitRedis(cfg *config.RedisConfig) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis 连接失败: %w", err)
	}

	return nil
}

func CloseRedis() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}

// ============ 缓存辅助函数 ============

// Set 设置缓存
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return Client.Set(ctx, key, value, expiration).Err()
}

// Get 获取缓存
func Get(ctx context.Context, key string) (string, error) {
	return Client.Get(ctx, key).Result()
}

// Delete 删除缓存
func Delete(ctx context.Context, keys ...string) error {
	return Client.Del(ctx, keys...).Err()
}

// Incr 自增
func Incr(ctx context.Context, key string) (int64, error) {
	return Client.Incr(ctx, key).Result()
}

// SetNX 设置仅当不存在时
func SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return Client.SetNX(ctx, key, value, expiration).Result()
}

// ============ 缓存 Key 模板 ============

const (
	// ProductCache 商品缓存
	ProductCache   = "product:%d"
	ProductListCache = "product:list:%s"
	
	// StockLock 库存锁
	StockLock    = "stock:lock:%d"
	StockLockTTL = 10 * time.Second // 锁超时时间
	
	// UserToken 用户Token
	UserToken = "user:token:%s"
	
	// OrderNoSeq 订单号序列
	OrderNoSeq = "order:seq"
)

// ============ 库存操作（面试重点）============

// LockStock 加库存锁（防止超卖）
// 【面试重点】分布式锁实现：SETNX + TTL
func LockStock(ctx context.Context, productID uint) (bool, error) {
	key := fmt.Sprintf(StockLock, productID)
	return SetNX(ctx, key, "1", StockLockTTL)
}

// UnlockStock 解锁
func UnlockStock(ctx context.Context, productID uint) error {
	key := fmt.Sprintf(StockLock, productID)
	return Delete(ctx, key)
}
