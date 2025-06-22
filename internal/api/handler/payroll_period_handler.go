package handler

import (
	"employee-app/internal/common/constant"
	commonFilter "employee-app/internal/common/model"
	"employee-app/internal/usecase/payroll"

	"employee-app/internal/api/dto"
	"employee-app/internal/api/dto/common"
	"employee-app/internal/api/handler/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PayrollPeriodHandler struct {
	payrollUC payroll.PayrollPeriodUsecase
}

func NewPayrollPeriodHandler(payrollUC payroll.PayrollPeriodUsecase) *PayrollPeriodHandler {
	return &PayrollPeriodHandler{payrollUC: payrollUC}
}

// POST /payroll-periods/{adminID}
func (h *PayrollPeriodHandler) CreatePeriod(c *gin.Context) {
	_, status, errMsg := utils.GetUserID(c, constant.EnumRoleAdmin)
	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	var req dto.CreatePayrollPeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminIDStr := c.Param("adminId")
	adminIDInt, err := strconv.Atoi(adminIDStr)
	if err != nil || adminIDInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid adminId path parameter"})
		return
	}
	adminID := uint(adminIDInt)

	err = h.payrollUC.CreatePeriod(c.Request.Context(), req.StartDate, req.EndDate, adminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// GET /payroll-periods
func (h *PayrollPeriodHandler) ListPeriods(c *gin.Context) {
	_, status, errMsg := utils.GetUserID(c, constant.EnumRoleAdmin)
	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	var query common.CommonQueryParams
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := commonFilter.CommonFilter{
		SortBy:    query.SortBy,
		Ascending: query.Ascending,
		Limit:     query.Limit,
		Page:      query.Page,
	}

	periods, err := h.payrollUC.ListPeriod(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []dto.PayrollPeriodResponse
	for _, p := range periods {
		responses = append(responses, dto.PayrollPeriodResponse{
			ID:        p.ID,
			StartDate: p.StartDate.Format("2006-01-02"),
			EndDate:   p.EndDate.Format("2006-01-02"),
			CreatedBy: p.CreatedBy,
			UpdatedBy: p.UpdatedBy,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, responses)
}

func (h *PayrollPeriodHandler) GetPeriodByID(c *gin.Context) {
	_, status, errMsg := utils.GetUserID(c, constant.EnumRoleAdmin)
	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": errMsg})
		return
	}
	
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payroll period ID"})
		return
	}

	period, err := h.payrollUC.GetPeriodByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payroll period not found"})
		return
	}

	resp := dto.PayrollPeriodResponse{
		ID:        period.ID,
		StartDate: period.StartDate.Format("2006-01-02"),
		EndDate:   period.EndDate.Format("2006-01-02"),
		CreatedBy: period.CreatedBy,
		UpdatedBy: period.UpdatedBy,
		CreatedAt: period.CreatedAt,
		UpdatedAt: period.UpdatedAt,
	}

	c.JSON(http.StatusOK, resp)
}