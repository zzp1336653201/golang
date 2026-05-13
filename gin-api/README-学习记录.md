# 学习记录

## 2026-05-13

- 11:31 **问**：GORM AutoMigrate 会自动创建数据库/表吗？
  → 指向：`docs/PostgreSQL完全指南.html`
  → AI总结：AutoMigrate 只会管理表结构，不会创建数据库本身，必须先手动创建数据库
  → Z的总结：GORM 不负责建库，只管表结构

- 11:31 **问**：PostgreSQL 相比 MySQL 有什么特点？
  → 指向：`docs/PostgreSQL完全指南.html`
  → AI总结：PostgreSQL 有数组、JSONB、范围类型等特有功能，并发使用 MVCC 机制
  → Z的总结：数组、JSONB 可建索引、MVCC 并发优秀，语法差异在自增ID和Upsert

- 11:35 **问**：Migration 是什么？生产环境如何管理数据库变更？
  → 指向：`docs/Migration详解.html`
  → AI总结：Migration = 数据库的 Git，记录每次表结构变更，有 up.sql（升级）和 down.sql（降级），与代码一起提交
  → Z的总结：Migration 是独立 CLI 工具，不是 Go 依赖，生产环境通过命令执行

- 14:08 **问**：Migration 是 Go 插件还是依赖？生产环境怎么用？
  → 指向：`docs/Migration详解.html`
  → AI总结：Migration 是独立 CLI 工具（如 golang-migrate、goose），不是 Go 依赖，在部署服务器用命令行执行，如 `migrate -path migrations -database "postgres://..." up`
  → Z的总结：独立工具，生产用命令 `migrate up/down`，也可集成 CI/CD

---

## 规则

- 每行不超过 150 字
- 格式：时间 + 问题 + HTML链接 + AI总结 + Z的总结
- 只记录核心内容，语句简短
