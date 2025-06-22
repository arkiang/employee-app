package handler

import (
	"employee-app/internal/api/dto"
	"employee-app/internal/api/dto/common"
	"employee-app/internal/api/handler/utils"
	"employee-app/internal/common/constant"
	commonFilter "employee-app/internal/common/model"
	"employee-app/internal/model"
	"employee-app/internal/usecase/reimbursement"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReimbursementHandler struct {
	reimbursementUC reimbursement.ReimbursementUsecase
}

func NewReimbursementHandler(reimbursementUC reimbursement.ReimbursementUsecase) *ReimbursementHandler {
	return &ReimbursementHandler{reimbursementUC: reimbursementUC}
}

// POST /reimbursements
func (h *ReimbursementHandler) SubmitReimbursement(c *gin.Context) {
	_, status, errMsg := utils.GetUserID(c, constant.EnumRoleEmployee)
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

	var req dto.SubmitReimbursementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.reimbursementUC.SubmitReimbursement(
		c.Request.Context(),
		empID,
		req.Date,
		req.Amount,
		req.Description,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
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

	reimbursements, err := h.reimbursementUC.GetReimbursementsForPeriod(c.Request.Context(), filter)
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