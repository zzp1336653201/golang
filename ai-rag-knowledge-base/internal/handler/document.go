package handler

import (
	"ai-rag-knowledge-base/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DocumentHandler struct {
	svc *service.DocumentService
}

func NewDocumentHandler(svc *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{svc: svc}
}

func (h *DocumentHandler) Upload(c *gin.Context) {
	var req service.UploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doc, err := h.svc.Upload(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, doc)
}

func (h *DocumentHandler) List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "TODO: 实现列表查询"})
}

func (h *DocumentHandler) Get(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "TODO: 实现详情查询"})
}

func (h *DocumentHandler) Delete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "TODO: 实现删除"})
}
