# 数据库初始化 SQL

## 前置条件

在执行以下 SQL 之前，请确保：

1. **PostgreSQL 已安装并运行**
2. **数据库已创建**
3. **连接配置正确**

---

## 1. 创建数据库

```sql
-- 创建 gin_api 数据库（如果不存在）
CREATE DATABASE gin_api;

-- 或在 psql 命令行中：
-- createdb gin_api
```

### 连接数据库（psql）
```bash
# 本地连接
psql -U postgres -d gin_api

# 或指定主机
psql -h localhost -U postgres -d gin_api
```

---

## 2. 创建 users 表

```sql
-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,                    -- 主键自增
    username VARCHAR(50) NOT NULL UNIQUE,     -- 用户名，唯一
    email VARCHAR(100) NOT NULL UNIQUE,      -- 邮箱，唯一
    password VARCHAR(255) NOT NULL,           -- 密码
    nickname VARCHAR(50),                     -- 昵称
    status INTEGER DEFAULT 1,                 -- 状态：1正常 0禁用
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 更新时间
    deleted_at TIMESTAMP                      -- 软删除时间
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- 添加注释
COMMENT ON TABLE users IS '用户表';
COMMENT ON COLUMN users.username IS '用户名';
COMMENT ON COLUMN users.email IS '邮箱';
COMMENT ON COLUMN users.password IS '密码（加密存储）';
COMMENT ON COLUMN users.nickname IS '昵称';
COMMENT ON COLUMN users.status IS '状态：1正常 0禁用';
COMMENT ON COLUMN users.deleted_at IS '软删除时间，为NULL表示未删除';
```

---

## 3. 插入测试数据

```sql
-- 插入测试用户
INSERT INTO users (username, email, password, nickname, status) VALUES
    ('alice', 'alice@example.com', '$2a$10$...', 'Alice', 1),
    ('bob', 'bob@example.com', '$2a$10$...', 'Bob', 1),
    ('charlie', 'charlie@example.com', '$2a$10$...', 'Charlie', 1);

-- 注意：实际项目中密码应该用 bcrypt 加密，这里用 ... 表示
```

---

## 4. 常用查询

```sql
-- 查看所有用户（排除已删除）
SELECT id, username, email, nickname, status, created_at 
FROM users 
WHERE deleted_at IS NULL;

-- 根据 ID 查询用户
SELECT * FROM users WHERE id = 1 AND deleted_at IS NULL;

-- 分页查询
SELECT * FROM users 
WHERE deleted_at IS NULL 
ORDER BY id DESC 
LIMIT 10 OFFSET 0;

-- 统计用户数量
SELECT COUNT(*) FROM users WHERE deleted_at IS NULL;
```

---

## 5. 查看表结构

```sql
-- 查看表结构
\d users

-- 或更详细的结构
SELECT 
    column_name, 
    data_type, 
    character_maximum_length,
    is_nullable,
    column_default
FROM information_schema.columns 
WHERE table_name = 'users';

-- 查看所有表
\dt

-- 查看索引
\di
```

---

## 6. 完整操作流程（推荐）

在 WSL 终端中执行：

```bash
# 1. 切换到 postgres 用户
sudo -u postgres psql

# 2. 在 psql 中执行：
CREATE DATABASE gin_api;
\c gin_api

# 3. 创建表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    nickname VARCHAR(50),
    status INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

# 4. 创建索引
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

# 5. 插入测试数据
INSERT INTO users (username, email, password, nickname, status) VALUES
    ('alice', 'alice@example.com', 'password123', 'Alice', 1),
    ('bob', 'bob@example.com', 'password123', 'Bob', 1);

# 6. 验证
SELECT * FROM users;

# 7. 退出
\q
```

---

## 7. 关闭 AutoMigrate（推荐）

由于我们已经手动创建了表，建议关闭 AutoMigrate 以避免意外修改表结构。修改 `main.go`：

```go
// 注释掉这一行
// models.AutoMigrate()
log.Println("数据库迁移完成（已跳过，使用手动建表）")
```

---

## 8. 字段说明

| 字段名 | 类型 | 约束 | 说明 |
|---|---|---|---|
| id | SERIAL | PRIMARY KEY | 主键，自增 |
| username | VARCHAR(50) | NOT NULL, UNIQUE | 用户名 |
| email | VARCHAR(100) | NOT NULL, UNIQUE | 邮箱 |
| password | VARCHAR(255) | NOT NULL | 密码（实际应加密） |
| nickname | VARCHAR(50) | - | 昵称 |
| status | INTEGER | DEFAULT 1 | 状态：1正常 0禁用 |
| created_at | TIMESTAMP | DEFAULT NOW() | 创建时间 |
| updated_at | TIMESTAMP | DEFAULT NOW() | 更新时间 |
| deleted_at | TIMESTAMP | - | 软删除标记 |

---

## 总结

### 快速执行步骤

```bash
# 1. 进入 psql
sudo -u postgres psql

# 2. 执行以下 SQL：
CREATE DATABASE gin_api;
\c gin_api
CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, username VARCHAR(50) NOT NULL UNIQUE, email VARCHAR(100) NOT NULL UNIQUE, password VARCHAR(255) NOT NULL, nickname VARCHAR(50), status INTEGER DEFAULT 1, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, deleted_at TIMESTAMP);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

# 3. 插入测试数据
INSERT INTO users (username, email, password, nickname, status) VALUES ('test', 'test@example.com', '123456', '测试用户', 1);

# 4. 验证
SELECT * FROM users;
```

执行完这些后，你就可以启动项目测试接口了！
