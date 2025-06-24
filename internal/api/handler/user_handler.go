package handler

import (
	"employee-app/internal/api/dto"
	"employee-app/internal/api/handler/utils"
	"employee-app/internal/common/constant"
	"employee-app/internal/usecase/user"
	"employee-app/pkg/security"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	usecase user.UserUsecase
}

func NewUserHandler(usecase user.UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

// Login authenticates a user and returns basic user info
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.usecase.Login(c, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := security.GenerateToken(user.ID, user.Role, time.Hour*24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
		"token":    token,
	})
}

// GetByID returns user info by ID
func (h *UserHandler) GetByID(c *gin.Context) {
	userID, status, errMsg := utils.GetUserID(c, constant.EnumRoleEmployee)
	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": errMsg})
		return
	}
	
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.usecase.GetByID(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if userID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not allowed to access this user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}