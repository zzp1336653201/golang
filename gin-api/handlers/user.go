package handlers

import (
	"net/http"
	"strconv"

	"gin-api/config"
	"gin-api/models"

	"github.com/gin-gonic/gin"
)

// UserResponse 统一响应结构
type UserResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// GetUsers 获取用户列表
// @Summary 获取用户列表
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} UserResponse
// @Router /api/users [get]
func GetUsers(c *gin.Context) {
	var users []models.User
	var total int64

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 查询总数
	config.DB.Model(&models.User{}).Count(&total)

	// 查询列表
	config.DB.Offset(offset).Limit(pageSize).Order("id DESC").Find(&users)

	c.JSON(http.StatusOK, UserResponse{
		Code:    0,
		Message: "success",
		Data: gin.H{
			"list":      users,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetUser 获取单个用户
// @Summary 获取单个用户
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} UserResponse
// @Router /api/users/{id} [get]
func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, UserResponse{
			Code:    404,
			Message: "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		Code:    0,
		Message: "success",
		Data:    user,
	})
}

// CreateUser 创建用户
// @Summary 创建用户
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "用户信息"
// @Success 200 {object} UserResponse
// @Router /api/users [post]
func CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, UserResponse{
			Code:    400,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	// 设置默认昵称
	if user.Nickname == "" {
		user.Nickname = user.Username
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, UserResponse{
			Code:    500,
			Message: "创建失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		Code:    0,
		Message: "创建成功",
		Data:    user,
	})
}

// UpdateUser 更新用户
// @Summary 更新用户
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param user body models.User true "用户信息"
// @Success 200 {object} UserResponse
// @Router /api/users/{id} [put]
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	var updateData models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, UserResponse{
			Code:    404,
			Message: "用户不存在",
		})
		return
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, UserResponse{
			Code:    400,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	// 只更新允许的字段
	updates := map[string]interface{}{}
	if updateData.Username != "" {
		updates["username"] = updateData.Username
	}
	if updateData.Email != "" {
		updates["email"] = updateData.Email
	}
	if updateData.Nickname != "" {
		updates["nickname"] = updateData.Nickname
	}
	if updateData.Status != 0 {
		updates["status"] = updateData.Status
	}

	if err := config.DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, UserResponse{
			Code:    500,
			Message: "更新失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		Code:    0,
		Message: "更新成功",
		Data:    user,
	})
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} UserResponse
// @Router /api/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, UserResponse{
			Code:    404,
			Message: "用户不存在",
		})
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, UserResponse{
			Code:    500,
			Message: "删除失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		Code:    0,
		Message: "删除成功",
	})
}
