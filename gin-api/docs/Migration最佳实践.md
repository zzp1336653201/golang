# 数据库 Migration 管理

## 什么是 Migration？

Migration 是一种数据库版本管理方式，记录数据库结构的每次变更。

### 为什么需要 Migration？

| 对比项 | AutoMigrate | Migration |
|---|---|---|
| **安全性** | ❌ 可能丢数据 | ✅ 可控、可回滚 |
| **可追溯** | ❌ 无记录 | ✅ 清晰的变更历史 |
| **多人协作** | ❌ 冲突风险大 | ✅ 合并 SQL 文件即可 |
| **生产部署** | ❌ 禁止使用 | ✅ DBA 审核执行 |

---

## 推荐 Migration 工具

### 1. golang-migrate（推荐）

**安装：**
```bash
# macOS
brew install golang-migrate

# Windows (下载二进制)
# https://github.com/golang-migrate/migrate/releases

# 或用 Go 安装
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

**项目结构：**
```
gin-api/
├── migrations/
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_add_status_column.up.sql
│   └── 000002_add_status_column.down.sql
└── ...
```

**创建 Migration 文件：**
```bash
# 自动生成带序号的文件
migrate create -ext sql -dir migrations -seq create_users_table
```

**执行 Migration：**
```bash
# 升级（执行 up.sql）
migrate -path migrations -database "postgres://postgres:zzp@localhost:5432/gin_api?sslmode=disable" up

# 降级（执行 down.sql）
migrate -path migrations -database "postgres://postgres:zzp@localhost:5432/gin_api?sslmode=disable" down 1

# 查看状态
migrate -path migrations -database "postgres://postgres:zzp@localhost:5432/gin_api?sslmode=disable" version
```

---

### 2. goose（更简单）

**安装：**
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

**创建 Migration：**
```bash
# 创建 users 表
goose postgres "postgres://postgres:zzp@localhost:5432/gin_api?sslmode=disable" create add users sql

# 编辑生成的文件
```

**执行：**
```bash
# 升级
goose postgres "postgres://postgres:zzp@localhost:5432/gin_api?sslmode=disable" up

# 降级
goose postgres "postgres://postgres:zzp@localhost:5432/gin_api?sslmode=disable" down

# 查看状态
goose postgres "postgres://postgres:zzp@localhost:5432/gin_api?sslmode=disable" status
```

---

## 实际项目建议

### 1. 开发流程

```
┌─────────────────────────────────────────────────────────┐
│  开发流程                                               │
├─────────────────────────────────────────────────────────┤
│  1. 本地开发 → 测试功能                                 │
│           ↓                                            │
│  2. 确定需要的表变更                                    │
│           ↓                                            │
│  3. 创建 migration 文件（up.sql + down.sql）           │
│           ↓                                            │
│  4. 提交代码 + migration 文件                          │
│           ↓                                            │
│  5. CI/CD 或 DBA 执行 migration                        │
│           ↓                                            │
│  6. 部署应用代码                                        │
└─────────────────────────────────────────────────────────┘
```

### 2. 关闭 AutoMigrate

在 `main.go` 中注释掉或删除 AutoMigrate：

```go
func main() {
    config.InitDB()

    // ❌ 生产环境关闭
    // models.AutoMigrate()

    // ✅ 改用 migration 管理
    log.Println("使用 migration 管理数据库结构")

    r := routes.SetupRouter()
    log.Println("服务启动: http://localhost:8080")
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("服务启动失败: %v", err)
    }
}
```

### 3. 手动建表 SQL（当前项目）

使用 `database.sql` 中的 SQL 手动建表：

```bash
sudo -u postgres psql -d gin_api -f database.sql
```

---

## 总结

### 生产环境必须做到

1. ✅ **关闭 AutoMigrate** - 禁止代码自动修改表
2. ✅ **使用 Migration 工具** - 版本化管理
3. ✅ **DBA 审核** - 重大变更需要 DBA 检查
4. ✅ **先执行 Migration** - 再部署代码
5. ✅ **保留 down.sql** - 方便回滚

### 开发环境可以

- 使用 AutoMigrate 快速原型
- 但最终生产前必须生成 migration 文件

---

你说得对：**生产环境绝对不能开 AutoMigrate**，这是一个非常专业的认识！👍