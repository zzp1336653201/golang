# 数据库脚本

## 初始化数据库

### 方式一：执行 SQL 脚本

```bash
mysql -u root -p < scripts/init.sql
```

### 方式二：GORM AutoMigrate（代码自动创建）

项目启动时会自动调用 `AutoMigrate`，无需手动执行 SQL。

```go
// pkg/database/mysql.go
func autoMigrate() error {
    return DB.AutoMigrate(
        &model.User{},
        &model.Product{},
        &model.Order{},
        &model.OrderItem{},
        &model.CartItem{},
    )
}
```

## 测试账号

| 用户名 | 密码 | 说明 |
|--------|------|------|
| admin | admin123 | 管理员 |

> 注意：生产环境请修改默认密码！
