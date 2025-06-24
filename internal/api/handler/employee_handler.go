package handler

import (
	"employee-app/internal/api/dto"
	"employee-app/internal/api/dto/common"
	"employee-app/internal/api/handler/utils"
	"employee-app/internal/common/constant"
	"employee-app/internal/model/entity"
	"employee-app/internal/usecase/employee"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	usecase employee.EmployeeUsecase
}

func NewEmployeeHandler(u employee.EmployeeUsecase) *EmployeeHandler {
	return &EmployeeHandler{usecase: u}
}

// GET /employees
func (h *EmployeeHandler) ListEmployees(c *gin.Context) {
	_, status, errMsg := utils.GetUserID(c, constant.EnumRoleAdmin)
	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	var queryParams common.CommonQueryParams
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := queryParams.ToFilter()

	employees, err := h.usecase.ListEmployees(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list employees"})
		return
	}

	c.JSON(http.StatusOK, employees)
}

// GET /employees/:id
func (h *EmployeeHandler) GetEmployeeByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee ID"})
		return
	}

	employee, err := h.usecase.GetEmployeeByID(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
		return
	}

	userID, status, errMsg := utils.GetUserID(c, constant.EnumRoleEmployee)
	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	_, status, errMsg = utils.CheckAccess(c, userID, employee.UserID)
	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, employee)
}

// PUT /employees/:id
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	userID, status, errMsg := utils.GetUserID(c, constant.EnumRoleAdmin)
	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": errMsg})
		return
	}
	
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee ID"})
		return
	}

	employee, err := h.usecase.GetEmployeeByID(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
		return
	}

	_, status, errMsg = utils.CheckAccess(c, userID, employee.ID)
	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	var input dto.UpdateEmployeeRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedEmployee := entity.Employee{
		ID:     uint(id),
		Name:   input.Name,
		Salary: input.Salary,
	}

	updated, err := h.usecase.UpdateEmployee(c, updatedEmployee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update employee"})
		return
	}

	c.JSON(http.StatusOK, updated)
}