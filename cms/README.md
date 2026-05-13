# CMS 电商系统

基于 Go + Gin 的企业级电商后端系统，包含商品管理、订单流程、用户认证等核心功能。

## 核心特性

| 模块 | 功能 |
|------|------|
| **商品管理** | CRUD、分类、上下架、库存管理 |
| **订单系统** | 创建→支付→发货→收货→完成 全流程 |
| **用户认证** | JWT Token、注册/登录 |
| **缓存策略** | Redis 缓存、防缓存穿透 |
| **库存控制** | 分布式锁防超卖、事务保证一致性 |
| **限流保护** | 令牌桶限流 |

## 面试亮点（必问）

### 1. 缓存穿透/击穿/雪崩

```go
// 缓存穿透：空值缓存
if product == nil {
    cache.Set(key, "__nil__", 60*time.Second) // 短时间缓存空值
}

// 缓存击穿：分布式锁
locked, _ := cache.LockStock(productID)
if !locked {
    return errors.New("系统繁忙")
}
defer cache.UnlockStock(productID)

// 雪崩：过期时间加随机值
ttl := 10*time.Minute + time.Duration(rand.Int63n(60))*time.Second
```

### 2. 超卖问题解决

```go
// 分布式锁 + 行锁双重保障
tx := db.Begin()

// 1. SELECT ... FOR UPDATE 加行锁
tx.Set("gorm:query_option", "FOR UPDATE").First(&product, id)

// 2. 检查库存
if product.Stock < quantity {
    tx.Rollback()
    return errors.New("库存不足")
}

// 3. 扣减库存
tx.Model(&Product{}).Where("id = ?", id).Update("stock", newStock)

// 4. 提交事务
tx.Commit()
```

### 3. 订单事务保证

```go
// 创建订单的事务保证
tx.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

// 1. 扣减多个商品库存
for _, item := range items {
    tx.Model(&Product{}).Where("id = ?", item.ProductID).
        Update("stock", gorm.Expr("stock - ?", item.Quantity))
}

// 2. 创建订单
tx.Create(order)

// 3. 创建订单项
tx.Create(orderItems)

// 4. 提交
tx.Commit()
```

### 4. 限流算法

```go
// 令牌桶限流
type TokenBucket struct {
    rate   int      // 每秒补充令牌数
    bucket int      // 桶容量
    tokens int      // 当前令牌
}

func (tb *TokenBucket) Allow() bool {
    // 补充令牌
    tb.tokens += int(time.Since(lastUpdate).Seconds() * float64(tb.rate))
    tb.tokens = min(tb.tokens, tb.bucket)
    
    if tb.tokens > 0 {
        tb.tokens--
        return true
    }
    return false
}
```

## 技术栈

| 层级 | 技术 |
|------|------|
| 语言 | Go 1.22+ |
| Web 框架 | Gin |
| ORM | GORM |
| 数据库 | MySQL |
| 缓存 | Redis |
| 认证 | JWT |
| 日志 | Zap |

## 项目结构

```
cms/
├── cmd/server/main.go        # 入口
├── internal/
│   ├── config/               # 配置加载
│   ├── handler/              # HTTP 处理器
│   ├── middleware/           # 中间件（认证、限流）
│   ├── model/                # 数据模型
│   └── service/              # 业务逻辑
├── pkg/
│   ├── cache/                # Redis 封装
│   ├── database/             # MySQL 连接
│   └── ratelimit/            # 限流器
└── config.yaml
```

## 快速开始

### 1. 环境准备

```bash
# MySQL
mysql -u root -p < CREATE DATABASE cms;

# Redis
redis-server
```

### 2. 配置

```yaml
# config.yaml
database:
  host: localhost
  port: 3306
  user: root
  password: your_password
  dbname: cms

redis:
  host: localhost
  port: 6379
```

### 3. 运行

```bash
go mod tidy
go run cmd/server/main.go
```

## API 接口

### 认证

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/v1/auth/register | 注册 |
| POST | /api/v1/auth/login | 登录 |

### 商品

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | /api/v1/products | 商品列表 | 否 |
| GET | /api/v1/products/:id | 商品详情 | 否 |
| POST | /api/v1/products | 创建商品 | 是 |

### 订单

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | /api/v1/orders | 创建订单 | 是 |
| GET | /api/v1/orders | 订单列表 | 是 |
| POST | /api/v1/orders/:order_no/pay | 支付 | 是 |
| DELETE | /api/v1/orders/:order_no | 取消订单 | 是 |

## 运行要求

- MySQL 5.7+
- Redis 6.0+
- Go 1.22+

## License

MIT
