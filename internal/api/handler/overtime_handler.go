package handler

import (
	"employee-app/internal/api/dto"
	"employee-app/internal/api/dto/common"
	"employee-app/internal/api/handler/utils"
	"employee-app/internal/common/constant"
	commonFilter "employee-app/internal/common/model"
	"employee-app/internal/model"
	"employee-app/internal/usecase/overtime"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OvertimeHandler struct {
	overtimeUC overtime.OvertimeUsecase
}

func NewOvertimeHandler(overtimeUC overtime.OvertimeUsecase) *OvertimeHandler {
	return &OvertimeHandler{overtimeUC: overtimeUC}
}

// POST /overtime
func (h *OvertimeHandler) SubmitOvertime(c *gin.Context) {
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

	var req dto.SubmitOvertimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call usecase
	if err := h.overtimeUC.SubmitOvertime(c.Request.Context(), empID, req.Date, req.Hours); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// GET /overtime
func (h *OvertimeHandler) GetOvertimeForPeriod(c *gin.Context) {
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

	overtimes, err := h.overtimeUC.GetOvertimeForPeriod(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map to DTO
	var response []dto.OvertimeResponse
	for _, o := range overtimes {
		response = append(response, dto.OvertimeResponse{
			ID:         o.ID,
			EmployeeID: o.EmployeeID,
			Date:       o.OvertimeDate.Format("2006-01-02"),
			Hours:      o.Hours,
			CreatedAt:  o.CreatedAt,
			UpdatedAt:  o.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}