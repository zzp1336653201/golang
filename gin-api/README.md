# Gin-API 用户管理接口

基于 Gin + GORM + PostgreSQL 的用户管理 API 项目

## 技术栈

| 技术 | 说明 |
|------|------|
| Go | 1.26.3 |
| Gin | Web 框架 |
| GORM | ORM |
| PostgreSQL | 数据库 |

## 项目结构

```
gin-api/
├── main.go           # 入口文件
├── go.mod            # 模块定义
├── config/
│   └── config.go     # 数据库配置
├── models/
│   └── user.go       # 用户模型
├── handlers/
│   └── user.go       # 用户处理器(CRUD)
└── routes/
    └── routes.go     # 路由配置
```

## API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /health | 健康检查 |
| GET | /api/users | 获取用户列表(分页) |
| GET | /api/users/:id | 获取单个用户 |
| POST | /api/users | 创建用户 |
| PUT | /api/users/:id | 更新用户 |
| DELETE | /api/users/:id | 删除用户 |

## 使用方法

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置数据库

设置环境变量或修改 `config/config.go`：

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=zzp
export DB_NAME=gin_api
```

### 3. 创建数据库

```sql
CREATE DATABASE gin_api;
```

### 4. 运行

```bash
go run main.go
```

### 5. 测试接口

```bash
# 创建用户
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","nickname":"测试用户"}'

# 获取用户列表
curl http://localhost:8080/api/users

# 获取单个用户
curl http://localhost:8080/api/users/1

# 更新用户
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{"nickname":"新昵称"}'

# 删除用户
curl -X DELETE http://localhost:8080/api/users/1
```

## 数据库连接信息

| 配置项 | 值 |
|--------|-----|
| Host | localhost |
| Port | 5432 |
| User | postgres |
| Password | zzp |
| Database | gin_api |
