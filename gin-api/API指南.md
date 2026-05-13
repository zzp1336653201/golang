# 逐文件介绍 gin-api 项目

---

## 📁 文件 1：`go.mod` - 模块配置文件

### 文件路径
```
gin-api/go.mod
```

### 文件内容
```go
module gin-api

go 1.26.3

require (
	github.com/gin-gonic/gin v1.10.0
	gorm.io/driver/postgres v1.5.9
	gorm.io/gorm v1.25.12
)
```

### 作用说明
- **module**：定义模块名为 `gin-api`，其他文件引用时用这个名字
- **go 1.26.3**：指定 Go 语言版本
- **require**：声明项目依赖的第三方库
  - `gin`：Web 框架，用于处理 HTTP 请求
  - `gorm`：ORM 库，用于操作数据库
  - `postgres`：Gorm 的 PostgreSQL 驱动

### 依赖版本说明
| 库 | 版本 | 用途 |
|---|---|---|
| gin | v1.10.0 | HTTP 路由和中间件 |
| gorm | v1.25.12 | ORM 对象关系映射 |
| postgres | v1.5.9 | PostgreSQL 数据库驱动 |

---

## 📁 文件 2：`main.go` - 程序入口

### 文件路径
```
gin-api/main.go
```

### 文件内容
```go
package main

import (
	"log"

	"gin-api/config"
	"gin-api/models"
	"gin-api/routes"
)

func main() {
	// 1. 初始化数据库连接
	config.InitDB()

	// 2. 自动迁移数据库表结构
	models.AutoMigrate()
	log.Println("数据库迁移完成")

	// 3. 配置路由
	r := routes.SetupRouter()

	// 4. 启动服务监听 8080 端口
	log.Println("服务启动: http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
```

### 作用说明
- **程序启动顺序**：
  1. 连接数据库
  2. 自动创建/更新数据表
  3. 配置路由规则
  4. 启动 HTTP 服务

### 启动后访问地址
- 服务地址：http://localhost:8080
- 健康检查：http://localhost:8080/health

---

## 📁 文件 3：`config/config.go` - 数据库配置

### 文件路径
```
gin-api/config/config.go
```

### 文件内容
```go
package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Config 数据库配置结构体
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// InitDB 初始化数据库连接
func InitDB() {
	cfg := Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "zzp"),
		DBName:   getEnv("DB_NAME", "gin_api"),
	}

	// 拼接 PostgreSQL 连接字符串
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port)

	// 建立连接
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	log.Println("数据库连接成功")
}

// getEnv 获取环境变量，支持默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
```

### 作用说明
- 管理数据库连接
- 支持通过环境变量配置数据库参数
- 默认配置适用于本地开发环境

### 数据库配置参数
| 参数 | 默认值 | 说明 |
|---|---|---|
| DB_HOST | localhost | 数据库地址 |
| DB_PORT | 5432 | PostgreSQL 端口 |
| DB_USER | postgres | 用户名 |
| DB_PASSWORD | zzp | 密码 |
| DB_NAME | gin_api | 数据库名 |

### 环境变量配置方法（可选）
```bash
# 在 WSL 中运行前设置环境变量
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=zzp
export DB_NAME=gin_api
```

---

## 📁 文件 4：`models/user.go` - 用户数据模型

### 文件路径
```
gin-api/models/user.go
```

### 文件内容
```go
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

// AutoMigrate 自动迁移（创建/更新数据表）
func AutoMigrate() {
	err := config.DB.AutoMigrate(&User{})
	if err != nil {
		panic("数据库迁移失败: " + err.Error())
	}
}
```

### 字段说明
| 字段 | 类型 | 约束 | JSON 标签 | 说明 |
|---|---|---|---|---|
| ID | uint | primaryKey | id | 主键自增 |
| Username | string | unique, not null | username | 用户名（唯一） |
| Email | string | unique, not null | email | 邮箱（唯一） |
| Password | string | not null | - | 密码（不返回给前端） |
| Nickname | string | - | nickname | 昵称 |
| Status | int | default:1 | status | 状态：1正常 0禁用 |
| CreatedAt | time.Time | - | created_at | 创建时间 |
| UpdatedAt | time.Time | - | updated_at | 更新时间 |
| DeletedAt | gorm.DeletedAt | index | - | 软删除标记 |

### Gorm 标签说明
- `primaryKey`：设为主键
- `uniqueIndex`：唯一索引
- `not null`：非空约束
- `size:50`：字段长度
- `default:1`：默认值
- `json:"-"`：JSON 序列化时忽略此字段

---

## 📁 文件 5：`handlers/user.go` - 接口处理函数

### 文件路径
```
gin-api/handlers/user.go
```

### 接口列表
| 方法 | 路由 | 功能 |
|---|---|---|
| GET | /api/users | 获取用户列表（分页） |
| GET | /api/users/:id | 获取单个用户 |
| POST | /api/users | 创建用户 |
| PUT | /api/users/:id | 更新用户 |
| DELETE | /api/users/:id | 删除用户 |

### 响应格式（统一）
```json
{
	"code": 0,
	"message": "success",
	"data": { ... }
}
```

### 各接口详细说明

#### 1️⃣ GET /api/users - 获取用户列表
```go
// 请求参数（Query）
page=1        // 页码，默认1
page_size=10  // 每页数量，默认10

// 响应示例
{
	"code": 0,
	"message": "success",
	"data": {
		"list": [
			{"id": 1, "username": "alice", "email": "alice@example.com", ...}
		],
		"total": 50,
		"page": 1,
		"page_size": 10
	}
}
```

#### 2️⃣ GET /api/users/:id - 获取单个用户
```go
// 请求参数（Path）
id=1  // 用户ID

// 响应示例
{
	"code": 0,
	"message": "success",
	"data": {
		"id": 1,
		"username": "alice",
		"email": "alice@example.com",
		"nickname": "Alice",
		"status": 1,
		"created_at": "2026-05-13T10:00:00Z"
	}
}
```

#### 3️⃣ POST /api/users - 创建用户
```go
// 请求体（JSON）
{
	"username": "bob",
	"email": "bob@example.com",
	"password": "123456",
	"nickname": "Bob"
}

// 响应示例
{
	"code": 0,
	"message": "创建成功",
	"data": {
		"id": 2,
		"username": "bob",
		"email": "bob@example.com",
		"nickname": "Bob",
		"status": 1
	}
}
```

#### 4️⃣ PUT /api/users/:id - 更新用户
```go
// 请求参数（Path）
id=1

// 请求体（JSON，可只传需要更新的字段）
{
	"email": "newemail@example.com",
	"nickname": "NewAlice"
}

// 响应示例
{
	"code": 0,
	"message": "更新成功",
	"data": { ... }
}
```

#### 5️⃣ DELETE /api/users/:id - 删除用户
```go
// 请求参数（Path）
id=1

// 响应示例
{
	"code": 0,
	"message": "删除成功"
}
```

### 错误码说明
| code | 说明 |
|---|---|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 📁 文件 6：`routes/routes.go` - 路由配置

### 文件路径
```
gin-api/routes/routes.go
```

### 文件内容
```go
package routes

import (
	"gin-api/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API 路由组
	api := r.Group("/api")
	{
		// 用户管理路由组
		users := api.Group("/users")
		{
			users.GET("", handlers.GetUsers)          // GET /api/users
			users.GET("/:id", handlers.GetUser)      // GET /api/users/:id
			users.POST("", handlers.CreateUser)       // POST /api/users
			users.PUT("/:id", handlers.UpdateUser)    // PUT /api/users/:id
			users.DELETE("/:id", handlers.DeleteUser) // DELETE /api/users/:id
		}
	}

	return r
}
```

### 路由结构图
```
/
├── /health                    # 健康检查
└── /api
    └── /users
        ├── GET    /           # 获取列表
        ├── GET    /:id        # 获取单个
        ├── POST   /           # 创建
        ├── PUT    /:id        # 更新
        └── DELETE /:id       # 删除
```

---

## 🛠️ 接口测试工具推荐

### 推荐工具对比

| 工具 | 类型 | 适用场景 | 学习曲线 |
|---|---|---|---|
| **Postman** | 桌面应用 | 专业 API 开发测试 | 中等 |
| **Thunder Client** | VS Code 插件 | VS Code 用户 | 简单 |
| **Bruno** | 桌面应用 | 开源替代 Postman | 简单 |
| **curl** | 命令行 | 快速测试 | 需要学习 |

### 1️⃣ Thunder Client（推荐新手）

**特点**：轻量级，VS Code 内直接使用，无需切换应用

**安装步骤**：
1. 打开 VS Code
2. 按 `Ctrl+Shift+X` 打开扩展
3. 搜索 "Thunder Client"
4. 点击安装

**使用示例**：
```bash
# 创建用户
POST http://localhost:8080/api/users
Content-Type: application/json

{
	"username": "testuser",
	"email": "test@example.com",
	"password": "123456",
	"nickname": "测试用户"
}
```

### 2️⃣ Postman（功能最全）

**下载地址**：https://www.postman.com/downloads/

**特点**：
- 功能最全面
- 支持 Collections（接口集合）
- 支持环境变量
- 团队协作

### 3️⃣ Bruno（开源替代）

**下载地址**：https://www.bruno.so/

**特点**：
- 开源免费
- 接口文件用纯文本保存（Git 友好）
- 界面简洁

### 4️⃣ curl（命令行快速测试）

Windows 10+ 自带 curl，直接在 PowerShell 使用：

```powershell
# 健康检查
curl http://localhost:8080/health

# 获取用户列表
curl http://localhost:8080/api/users

# 获取单个用户
curl http://localhost:8080/api/users/1

# 创建用户
curl -X POST http://localhost:8080/api/users `
  -H "Content-Type: application/json" `
  -d '{"username":"alice","email":"alice@example.com","password":"123456"}'

# 更新用户
curl -X PUT http://localhost:8080/api/users/1 `
  -H "Content-Type: application/json" `
  -d '{"nickname":"新昵称"}'

# 删除用户
curl -X DELETE http://localhost:8080/api/users/1
```

---

## 🚀 下一步操作

### 在 WSL 中启动服务
```bash
# 1. 进入项目目录
cd /disk/www/golang/gin-api

# 2. 安装依赖
go mod tidy

# 3. 创建数据库（如果还没有）
# 先用 postgres 用户登录
sudo -u postgres psql

# 在 psql 中执行：
CREATE DATABASE gin_api;

# 退出 psql
\q

# 4. 启动服务
go run main.go
```

### 测试接口
服务启动后，打开 Thunder Client 或 Postman，测试以下接口：

1. `GET http://localhost:8080/health` - 健康检查
2. `POST http://localhost:8080/api/users` - 创建用户
3. `GET http://localhost:8080/api/users` - 获取列表
4. `GET http://localhost:8080/api/users/1` - 获取单个
5. `PUT http://localhost:8080/api/users/1` - 更新用户
6. `DELETE http://localhost:8080/api/users/1` - 删除用户

---

你有任何问题或需要进一步解释的地方吗？