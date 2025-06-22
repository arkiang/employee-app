package handler

import (
	"employee-app/internal/api/dto"
	"employee-app/internal/api/dto/common"
	"employee-app/internal/api/middleware"
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
	var queryParams common.CommonQueryParams
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := queryParams.ToFilter()

	employees, err := h.usecase.ListEmployees(c.Request.Context(), filter)
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

	employee, err := h.usecase.GetEmployeeByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
		return
	}

	c.JSON(http.StatusOK, employee)
}

// PUT /employees/:id
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	userID, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found in token"})
		return
	}
	
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee ID"})
		return
	}

	var input dto.UpdateEmployeeRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee := entity.Employee{
		ID:     uint(id),
		UserID: userID.(uint),
		Name:   input.Name,
		Salary: input.Salary,
	}

	updated, err := h.usecase.UpdateEmployee(c.Request.Context(), employee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update employee"})
		return
	}

	c.JSON(http.StatusOK, updated)
}