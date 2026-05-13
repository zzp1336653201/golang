package vector

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ChromaDB struct {
	endpoint string
	client   *http.Client
}

type SearchResult struct {
	ID        string
	Content   string
	DocumentID string
	Score     float64
}

func NewChromaDB(endpoint string) (*ChromaDB, error) {
	return &ChromaDB{
		endpoint: endpoint,
		client:   &http.Client{},
	}, nil
}

func (c *ChromaDB) Search(ctx context.Context, query string, topK int, collection string) ([]SearchResult, error) {
	// 构建查询请求
	reqBody := map[string]interface{}{
		"query_embeddings": query, // 简化：实际需要先 embedding
		"n_results":        topK,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/collections/%s/query", c.endpoint, collection)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ChromaDB 查询失败: %w", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var result struct {
		IDs       []string    `json:"ids"`
		Distances [][]float64 `json:"distances"`
		Metadatas [][]map[string]interface{} `json:"metadatas"`
		Documents [][]string   `json:"documents"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// 转换结果
	results := make([]SearchResult, 0, len(result.IDs))
	for i, id := range result.IDs {
		if len(result.Documents) > 0 && len(result.Documents[i]) > 0 {
			results = append(results, SearchResult{
				ID:         id,
				Content:    result.Documents[i][0],
				DocumentID: id,
				Score:      1 - result.Distances[i][0], // 转换距离为相似度
			})
		}
	}

	return results, nil
}

func (c *ChromaDB) Insert(ctx context.Context, docID string, chunks []string, metadata string) error {
	reqBody := map[string]interface{}{
		"ids":          []string{docID},
		"documents":    chunks,
		"metadatas":    []map[string]interface{}{{"source": metadata}},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/collections/documents/add", c.endpoint)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ChromaDB 插入失败: status %d", resp.StatusCode)
	}

	return nil
}
