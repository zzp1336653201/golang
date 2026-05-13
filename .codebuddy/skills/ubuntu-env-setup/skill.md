# WSL Ubuntu-24.04 环境搭建记录

## 环境信息
- **系统**: WSL Ubuntu-24.04
- **用途**: Go 后端开发学习 + 生产部署流程实践
- **状态**: Go 环境已安装完成

## 搭建进度追踪

### Go 开发环境
- [x] Go 语言安装 (1.26.3)
- [x] Go 环境变量配置 (GOROOT, GOPATH)
- [x] Go Modules 配置 (goproxy.cn)
- [ ] Gin 框架验证运行

### 数据库环境
- [x] PostgreSQL 安装 (16.13)
- [x] PostgreSQL 基础配置
- [ ] 创建开发数据库
- [ ] 配置远程访问（如需）

### 缓存环境
- [x] Redis 安装 (7.0.15)
- [x] Redis 基础配置
- [ ] 配置持久化

### 部署工具
- [x] Git 配置 (已有)
- [x] Docker 安装 (29.4.3)
- [x] Docker Compose (随 Docker 一起安装)
- [x] Nginx 安装 (1.24.0)

### 开发工具
- [ ] 代码编辑器/IDE 配置
- [ ] 终端工具配置
- [ ] SSH 密钥配置

## 环境搭建记录

### 2026-05-13 Go 环境安装（版本更新）

**安装版本**: Go 1.26.3 linux/amd64

**安装路径**: `/usr/local/go`

**环境变量配置** (`~/.bashrc`):
```bash
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

**关键环境变量**:
| 变量 | 值 | 说明 |
|------|-----|------|
| GOROOT | `/usr/local/go` | Go 安装目录 |
| GOPATH | `/home/zzp/go` | 工作目录/依赖缓存 |
| GOMODCACHE | `/home/zzp/go/pkg/mod` | 模块缓存路径 |
| GOVERSION | `go1.26.3` | 版本号 |
| GOOS | `linux` | 目标系统 |
| GOARCH | `amd64` | 目标架构 |

**Go Modules 配置**:
- 默认使用 `https://proxy.golang.org` 作为代理
- 国内网络慢时可配置国内镜像

**常用命令**:
```bash
go version                    # 查看版本
go env                        # 查看环境变量
go mod init <name>            # 初始化模块
go get <package>              # 安装依赖
go run <file>                 # 运行程序
go build                      # 编译程序
```

## 生产部署流程笔记

<!-- 记录从开发到部署的完整流程 -->

## 常用命令备忘

```bash
# WSL 相关
wsl -l -v                    # 查看 WSL 状态
wsl --shutdown               # 关闭 WSL
wsl -d Ubuntu-24.04          # 进入 Ubuntu

# Ubuntu 基础
sudo apt update              # 更新包列表
sudo apt upgrade             # 升级包
```
