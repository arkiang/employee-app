package handler

import (
	"employee-app/internal/api/dto"
	"employee-app/internal/api/dto/common"
	"employee-app/internal/api/handler/utils"
	"employee-app/internal/common/constant"
	commonFilter "employee-app/internal/common/model"
	"employee-app/internal/model"
	"employee-app/internal/usecase/employee"
	"employee-app/internal/usecase/reimbursement"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type ReimbursementHandler struct {
	reimbursementUC reimbursement.ReimbursementUsecase
	employeeUsecase  employee.EmployeeUsecase
}

func NewReimbursementHandler(reimbursementUC reimbursement.ReimbursementUsecase, employeeUsecase employee.EmployeeUsecase) *ReimbursementHandler {
	return &ReimbursementHandler{reimbursementUC: reimbursementUC, employeeUsecase: employeeUsecase}
}

// POST /reimbursements
func (h *ReimbursementHandler) SubmitReimbursement(c *gin.Context) {
	userID, status, errMsg := utils.GetUserID(c, constant.EnumRoleEmployee)
	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	empIDStr := c.Param("empId")
	empIDInt, err := strconv.Atoi(empIDStr)
	if err != nil || empIDInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid empId path parameter"})
		return
	}
	empID := uint(empIDInt)

	employee, err := h.employeeUsecase.GetEmployeeByID(c, empID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
		return
	}

	_, status, errMsg = utils.CheckAccess(c, userID, employee.UserID)
	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	var req dto.SubmitReimbursementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount"})
		return
	}
	if err := h.reimbursementUC.SubmitReimbursement(
		c,
		empID,
		req.Date,
		amount,
		req.Description,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reimbursement submitted successfully"})
}

// GET /reimbursements
func (h *ReimbursementHandler) GetReimbursementsForPeriod(c *gin.Context) {
	_, status, errMsg := utils.GetUserID(c, constant.EnumRoleEmployee)
	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	var query common.EmployeePeriodQueryParams
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := model.EmployeePeriodFilter{
		EmpIds: &query.EmployeeIDs,
		Start:  query.Start,
		End:    query.End,
		Base: commonFilter.CommonFilter{
			SortBy:   query.SortBy,
			Ascending: query.Ascending,
			Limit:    query.Limit,
			Page:     query.Page,
		},
	}

	reimbursements, err := h.reimbursementUC.GetReimbursementsForPeriod(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	var response []dto.ReimbursementResponse
	for _, r := range reimbursements {
		response = append(response, dto.ReimbursementResponse{
			ID:          r.ID,
			EmployeeID:  r.EmployeeID,
			Date:        r.ReimbursementDate.Format("2006-01-02"),
			Amount:      r.Amount,
			Description: r.Description,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}