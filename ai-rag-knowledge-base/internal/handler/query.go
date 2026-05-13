package handler

import (
	"ai-rag-knowledge-base/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type QueryHandler struct {
	svc *service.RAGService
}

func NewQueryHandler(svc *service.RAGService) *QueryHandler {
	return &QueryHandler{svc: svc}
}

func (h *QueryHandler) Query(c *gin.Context) {
	var req service.QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.TopK == 0 {
		req.TopK = 5
	}
	if req.Collection == "" {
		req.Collection = "documents"
	}

	resp, err := h.svc.Query(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *QueryHandler) Search(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "TODO: 实现混合检索"})
}
