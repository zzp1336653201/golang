# AI + RAG 知识库系统

基于 Go + Ollama/ChromaDB 的智能知识库系统，支持多格式文档处理、混合检索和 RAG 增强问答。

## 核心特性

- **多格式文档处理**：支持 PDF、Word、Markdown、TXT 等文档解析
- **混合检索**：向量检索 + 关键词检索融合
- **RAG 引擎**：基于检索增强的生成，支持上下文注入
- **本地部署**：使用 Ollama 本地运行大模型，保护数据隐私
- **RESTful API**：清晰易用的 HTTP 接口

## 技术栈

| 层级 | 技术 |
|------|------|
| 语言 | Go 1.22+ |
| 向量数据库 | ChromaDB / pgvector |
| LLM | Ollama (本地部署) |
| 文档解析 | gooxml, PDF 解析库 |
| Web 框架 | Gin |
| 配置管理 | Viper |
| 日志 | Zap |

## 项目结构

```
ai-rag-knowledge-base/
├── cmd/server/          # 应用入口
├── internal/
│   ├── config/          # 配置加载
│   ├── handler/         # HTTP 处理器
│   ├── model/           # 数据模型
│   ├── repository/       # 数据访问层
│   ├── service/         # 业务逻辑层
│   └── rag/              # RAG 核心引擎
├── pkg/
│   ├── vector/           # 向量处理封装
│   └── llm/              # LLM 调用封装
└── docs/                 # 文档
```

## 快速开始

### 1. 环境准备

```bash
# 安装 Ollama
curl -fsSL https://ollama.com/install.sh | sh

# 拉取模型
ollama pull llama3.2

# 启动 Ollama 服务
ollama serve
```

### 2. 配置

```bash
cp config.example.yaml config.yaml
# 编辑 config.yaml 配置数据库和模型参数
```

### 3. 运行

```bash
go mod tidy
go run cmd/server/main.go
```

## API 接口

### 文档管理

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/v1/documents | 上传文档 |
| GET | /api/v1/documents | 列出文档 |
| GET | /api/v1/documents/:id | 获取文档详情 |
| DELETE | /api/v1/documents/:id | 删除文档 |

### 知识库查询

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/v1/query | RAG 问答 |
| GET | /api/v1/search | 混合检索 |

## 求职亮点

1. **RAG 架构设计**：完整实现检索-增强-生成流程
2. **混合检索实现**：向量相似度 + BM25 融合排序
3. **本地 LLM 集成**：Ollama API 对接实践
4. **可生产级代码**：依赖注入、错误处理、日志规范

## License

MIT
