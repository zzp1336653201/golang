package main

import (
	"context"
	"fmt"
	"os"

	"ai-rag-knowledge-base/internal/config"
	"ai-rag-knowledge-base/internal/handler"
	"ai-rag-knowledge-base/internal/repository"
	"ai-rag-knowledge-base/internal/service"
	"ai-rag-knowledge-base/pkg/llm"
	"ai-rag-knowledge-base/pkg/vector"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// 初始化日志
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	// 加载配置
	if err := config.Load(); err != nil {
		sugar.Fatalf("配置加载失败: %v", err)
	}

	// 初始化组件
	ctx := context.Background()

	// 向量存储
	vectorDB, err := vector.NewChromaDB(viper.GetString("chroma.endpoint"))
	if err != nil {
		sugar.Fatalf("向量数据库初始化失败: %v", err)
	}

	// LLM 客户端
	llmClient := llm.NewOllamaClient(
		viper.GetString("ollama.endpoint"),
		viper.GetString("ollama.model"),
	)

	// 仓储层
	docRepo := repository.NewDocumentRepository()

	// 服务层
	ragService := service.NewRAGService(vectorDB, llmClient, docRepo)
	docService := service.NewDocumentService(docRepo, vectorDB)

	// 处理器
	docHandler := handler.NewDocumentHandler(docService)
	queryHandler := handler.NewQueryHandler(ragService)

	// 启动服务
	port := viper.GetInt("server.port")
	fmt.Printf("🚀 RAG 知识库服务启动中... 端口: %d\n", port)

	// TODO: 启动 HTTP 服务器
	_ = ctx
	_ = docHandler
	_ = queryHandler
	_ = sugar
}
