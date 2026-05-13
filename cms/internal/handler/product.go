package handler

import (
	"net/http"
	"strconv"

	"cms/internal/middleware"
	"cms/internal/model"
	"cms/internal/service"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// Create 创建商品
// POST /api/v1/products
func (h *ProductHandler) Create(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Price       float64 `json:"price" binding:"required,gt=0"`
		Stock       int     `json:"stock" binding:"gte=0"`
		Category    string  `json:"category"`
		ImageURL    string  `json:"image_url"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(err.Error()))
		return
	}
	
	product := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       int64(req.Price * 100), // 元转分
		Stock:       req.Stock,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
		Status:      1,
	}
	
	if err := h.productService.Create(c.Request.Context(), product); err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(err.Error()))
		return
	}
	
	c.JSON(http.StatusCreated, model.Success(product))
}

// Get 获取商品详情
// GET /api/v1/products/:id
func (h *ProductHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error("无效的商品ID"))
		return
	}
	
	product, err := h.productService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, model.Error("商品不存在"))
		return
	}
	
	c.JSON(http.StatusOK, model.Success(product))
}

// List 商品列表
// GET /api/v1/products
func (h *ProductHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	category := c.Query("category")
	
	products, total, err := h.productService.List(c.Request.Context(), page, pageSize, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, model.Success(model.NewPageResponse(total, page, pageSize, products)))
}

// Update 更新商品
// PUT /api/v1/products/:id
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error("无效的商品ID"))
		return
	}
	
	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
		Category    string  `json:"category"`
		ImageURL    string  `json:"image_url"`
		Status      int8    `json:"status"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(err.Error()))
		return
	}
	
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Price > 0 {
		updates["price"] = int64(req.Price * 100)
	}
	if req.Stock >= 0 {
		updates["stock"] = req.Stock
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.ImageURL != "" {
		updates["image_url"] = req.ImageURL
	}
	if req.Status > 0 {
		updates["status"] = req.Status
	}
	
	if err := h.productService.Update(c.Request.Context(), uint(id), updates); err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, model.Success(nil))
}

// Delete 删除商品
// DELETE /api/v1/products/:id
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error("无效的商品ID"))
		return
	}
	
	if err := h.productService.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, model.Success(nil))
}

// GetStock 获取库存
// GET /api/v1/products/:id/stock
func (h *ProductHandler) GetStock(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error("无效的商品ID"))
		return
	}
	
	stock, err := h.productService.GetStock(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, model.Success(gin.H{"stock": stock}))
}

// ============ 订单相关 ============

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// CreateOrder 创建订单
// POST /api/v1/orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, model.Error("请先登录"))
		return
	}
	
	var req struct {
		Items    []struct {
			ProductID uint `json:"product_id" binding:"required"`
			Quantity  int  `json:"quantity" binding:"required,gt=0"`
		} `json:"items" binding:"required,min=1"`
		Receiver service.ReceiverInfo `json:"receiver" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(err.Error()))
		return
	}
	
	items := make([]struct {
		ProductID uint
		Quantity  int
	}, len(req.Items))
	for i, item := range req.Items {
		items[i] = struct {
			ProductID uint
			Quantity  int
		}{item.ProductID, item.Quantity}
	}
	
	order, err := h.orderService.CreateOrder(c.Request.Context(), userID, items, req.Receiver)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(err.Error()))
		return
	}
	
	c.JSON(http.StatusCreated, model.Success(order))
}

// GetOrders 获取订单列表
// GET /api/v1/orders
func (h *OrderHandler) GetOrders(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, model.Error("请先登录"))
		return
	}
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	
	orders, total, err := h.orderService.GetUserOrders(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, model.Success(model.NewPageResponse(total, page, pageSize, orders)))
}

// PayOrder 支付订单
// POST /api/v1/orders/:order_no/pay
func (h *OrderHandler) PayOrder(c *gin.Context) {
	orderNo := c.Param("order_no")
	if orderNo == "" {
		c.JSON(http.StatusBadRequest, model.Error("订单号不能为空"))
		return
	}
	
	if err := h.orderService.PayOrder(c.Request.Context(), orderNo); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, model.Success(nil))
}

// CancelOrder 取消订单
// DELETE /api/v1/orders/:order_no
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, model.Error("请先登录"))
		return
	}
	
	orderNo := c.Param("order_no")
	if orderNo == "" {
		c.JSON(http.StatusBadRequest, model.Error("订单号不能为空"))
		return
	}
	
	if err := h.orderService.CancelOrder(c.Request.Context(), orderNo, userID); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, model.Success(nil))
}
