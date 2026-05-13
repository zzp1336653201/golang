package service

import (
	"context"
	"fmt"

	"ai-rag-knowledge-base/internal/repository"
	"ai-rag-knowledge-base/pkg/llm"
	"ai-rag-knowledge-base/pkg/vector"
)

type RAGService struct {
	vectorDB *vector.ChromaDB
	llm      *llm.OllamaClient
	docRepo  *repository.DocumentRepository
}

func NewRAGService(
	vectorDB *vector.ChromaDB,
	llm *llm.OllamaClient,
	docRepo *repository.DocumentRepository,
) *RAGService {
	return &RAGService{
		vectorDB: vectorDB,
		llm:      llm,
		docRepo:  docRepo,
	}
}

type QueryRequest struct {
	Query      string   `json:"query"`
	TopK       int      `json:"top_k"`
	Collection string   `json:"collection"`
}

type QueryResponse struct {
	Answer    string   `json:"answer"`
	Sources   []string `json:"sources"`
	Score     float64  `json:"score"`
}

func (s *RAGService) Query(ctx context.Context, req *QueryRequest) (*QueryResponse, error) {
	// 1. 向量检索
	results, err := s.vectorDB.Search(ctx, req.Query, req.TopK, req.Collection)
	if err != nil {
		return nil, err
	}

	// 2. 构建上下文
	context := s.buildContext(results)

	// 3. 调用 LLM 生成答案
	prompt := s.buildPrompt(req.Query, context)
	answer, err := s.llm.Generate(ctx, prompt)
	if err != nil {
		return nil, err
	}

	return &QueryResponse{
		Answer:  answer,
		Sources: s.extractSources(results),
		Score:   results[0].Score,
	}, nil
}

func (s *RAGService) buildContext(results []vector.SearchResult) string {
	context := "相关文档片段：\n"
	for i, r := range results {
		context += fmt.Sprintf("%d. %s\n\n", i+1, r.Content)
	}
	return context
}

func (s *RAGService) buildPrompt(query, context string) string {
	return f"""基于以下上下文回答问题。如果上下文中没有相关信息，请说明无法从提供的内容中找到答案。

上下文：
{context}

问题：{query}

答案："""
}

func (s *RAGService) extractSources(results []vector.SearchResult) []string {
	sources := make([]string, 0, len(results))
	for _, r := range results {
		sources = append(sources, r.DocumentID)
	}
	return sources
}
