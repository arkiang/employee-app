package handler

import (
	"employee-app/internal/api/dto"
	"employee-app/internal/api/handler/utils"
	"employee-app/internal/model/entity"
	"employee-app/internal/usecase/registration"
	"net/http"

	"employee-app/internal/common/constant"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type RegistrationHandler struct {
	usecase registration.RegistrationUsecase
}

func NewRegistrationHandler(u registration.RegistrationUsecase) *RegistrationHandler {
	return &RegistrationHandler{
		usecase: u,
	}
}

func (h *RegistrationHandler) Register(c *gin.Context) {
	_, status, errMsg := utils.GetUserID(c, constant.EnumRoleAdmin)
	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	var req dto.RegisterEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to entity
	user := entity.User{
		Username: req.Username,
		Role:     "employee",
	}

	salary, err := decimal.NewFromString(req.Salary)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount"})
		return
	}
	employee := entity.Employee{
		Name:   req.Name,
		Salary: salary,
	}

	if err := h.usecase.RegisterEmployee(c, user, employee, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "registration successful"})
}