# WSL Ubuntu 开发环境配置

## 项目信息
- 项目路径：`/disk/www/golang`
- Windows 路径：`e:/xp.cn/www/golang`
- WSL 分发版：Ubuntu-24.04

## 已安装软件

| 软件 | 版本 | 安装路径 |
|------|------|----------|
| Go | 1.26.3 | /usr/local/go |
| GOPROXY | goproxy.cn,direct | 环境变量 |

## WSL 常用命令

```bash
# 进入 WSL
wsl -d Ubuntu-24.04

# 在 Windows 终端直接执行 WSL 命令
wsl bash -c "命令"

# WSL 内切换目录
cd /disk/www/golang

# 常用检查
go version
go env GOPROXY
```

## Go 相关命令

```bash
# 依赖管理
go mod tidy
go mod download
go get github.com/xxx/xxx

# 运行
go run main.go

# 构建
go build -o binary_name
```

## 注意事项
- 项目放在 WSL 本地文件系统 (/disk/www/)，不用 /mnt/e
- WSL 项目实际路径需在 WSL 终端用 `wsl -d Ubuntu-24.04` 进入后用 `cd ~` 查看
- GOPROXY 已配置为 goproxy.cn 加速

### ⚠️ 重要规则：WSL 命令执行
- **不要**直接执行 WSL 命令（如 wsl -d Ubuntu-24.04 -- bash -c "..."）
- 改为提供命令脚本，**由用户手动在 WSL 终端执行**
- 原因：WSL 路径映射复杂，直接执行可能失败

### 标准 WSL 操作方式
```
❌ 错误：直接用工具执行 wsl 命令
✅ 正确：提供命令让用户在 WSL 终端手动执行

示例：
# 提供给用户的命令
wsl -d Ubuntu-24.04
cd /disk/www/golang/gin-api
go mod tidy
go run main.go
```

## 项目工作模式与规则

### 统一工作模式：Windows + WSL 双环境协作

#### 1. 环境分工
- **Windows 环境**：代码编写、项目初始化、Git 管理
- **WSL 环境**：运行调试、依赖安装、数据库操作、服务启动

#### 2. 目录管理规则
- **每个项目或任务**必须单独创建目录
- 禁止在根目录直接创建项目文件
- 禁止多项目混用同一目录

#### 3. 项目目录命名规范
```
/disk/www/golang/
├── project-name-1/          # 任务1
├── project-name-2/          # 任务2
├── another-project/         # 其他项目
└── ...
```

#### 4. 标准项目结构
```
project-name/
├── go.mod                   # Go 模块配置
├── main.go                  # 主入口
├── README.md                # 项目说明
├── config/                  # 配置目录
│   └── config.go
├── models/                  # 数据模型
│   └── xxx.go
├── handlers/                # 业务处理
│   └── xxx.go
├── routes/                  # 路由配置
│   └── routes.go
└── docs/                    # 文档（可选）
```

#### 5. 开发流程规范

##### Windows 端（代码准备）
```powershell
# 1. 创建项目目录
mkdir project-name
cd project-name

# 2. 初始化 Go 模块
go mod init 项目名

# 3. 编写代码文件
# - main.go
# - config/config.go
# - models/*.go
# - handlers/*.go
# - routes/routes.go

# 4. Git 管理
git init
git add .
git commit -m "feat: 项目初始化"
git remote add origin <仓库地址>
git push
```

##### WSL 端（运行调试）
```bash
# 1. 进入项目目录
cd /disk/www/golang/project-name

# 2. 拉取代码（如有）
git pull

# 3. 安装依赖
go mod tidy

# 4. 运行项目
go run main.go

# 5. 测试接口
curl http://localhost:8080/api/xxx
```

#### 6. 作业/任务处理规范
- 每次新作业创建独立目录
- 在 skill 文件中记录项目关键信息
- 更新本 skill 文件的工作日志

#### 7. 代码同步机制
- Windows 编写 → Git 提交推送
- WSL 拉取 → 运行测试
- 避免交叉编辑导致冲突

#### 8. 文件创建规则 ⚠️
- **所有项目文件**（代码、SQL、配置、migrations）**必须在 Windows 工作区创建**
- **禁止**在 WSL/服务器上直接创建文件
- 原因：文件需通过 Git 版本管理，保持可追溯性
- 流程：Windows 创建/修改 → Git commit → WSL git pull

```
❌ 错误：在 WSL 直接创建 migrations 文件
✅ 正确：在 Windows 工作区创建 migrations 文件
```

#### 8. 依赖管理
- 所有依赖通过 `go mod tidy` 统一管理
- 禁止手动编辑 go.mod
- 第三方库通过 `go get` 安装

### 示例项目记录

| 项目名 | 创建时间 | 说明 | 状态 |
|--------|----------|------|------|
| gin-api | 2026-05-13 | Gin REST API 用户管理 | 待运行 |
| (后续项目...) | | | |
