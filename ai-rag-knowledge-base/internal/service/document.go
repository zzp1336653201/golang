package service

import (
	"ai-rag-knowledge-base/internal/model"
	"ai-rag-knowledge-base/internal/repository"
	"ai-rag-knowledge-base/pkg/vector"
	"context"
)

type DocumentService struct {
	repo      *repository.DocumentRepository
	vectorDB  *vector.ChromaDB
}

func NewDocumentService(repo *repository.DocumentRepository, vectorDB *vector.ChromaDB) *DocumentService {
	return &DocumentService{
		repo:     repo,
		vectorDB: vectorDB,
	}
}

type UploadRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Type      string `json:"type"`
	FilePath  string `json:"file_path"`
}

func (s *DocumentService) Upload(ctx context.Context, req *UploadRequest) (*model.Document, error) {
	// 1. 创建文档记录
	doc := model.NewDocument(req.Title, req.Content, req.Type, req.FilePath)
	if err := s.repo.Create(doc); err != nil {
		return nil, err
	}

	// 2. 分块并向量化
	chunks := s.chunkText(req.Content)
	if err := s.vectorDB.Insert(ctx, doc.ID, chunks, req.Title); err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *DocumentService) chunkText(text string) []string {
	// 简单分块策略：按固定长度分块
	const chunkSize = 500
	const overlap = 50

	chunks := make([]string, 0)
	runes := []rune(text)

	for i := 0; i < len(runes); i += chunkSize - overlap {
		end := i + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
		if end == len(runes) {
			break
		}
	}

	return chunks
}
