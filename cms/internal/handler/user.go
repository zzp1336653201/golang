package handler

import (
	"net/http"

	"cms/internal/model"
	"cms/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Register 注册
// POST /api/v1/auth/register
func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Password string `json:"password" binding:"required,min=6"`
		Email    string `json:"email"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(err.Error()))
		return
	}
	
	user, err := h.userService.Register(c.Request.Context(), req.Username, req.Password, req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(err.Error()))
		return
	}
	
	c.JSON(http.StatusCreated, model.Success(gin.H{
		"user_id":  user.ID,
		"username": user.Username,
	}))
}

// Login 登录
// POST /api/v1/auth/login
func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(err.Error()))
		return
	}
	
	token, user, err := h.userService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.Error(err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, model.Success(gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
		},
	}))
}

// GetProfile 获取当前用户信息
// GET /api/v1/user/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.Error("未登录"))
		return
	}
	
	u := user.(*model.User)
	c.JSON(http.StatusOK, model.Success(gin.H{
		"id":       u.ID,
		"username": u.Username,
		"email":    u.Email,
		"phone":    u.Phone,
		"nickname": u.Nickname,
	}))
}
