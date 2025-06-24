package handler

import (
	"employee-app/internal/api/dto"
	"employee-app/internal/api/dto/common"
	"employee-app/internal/api/handler/utils"
	"employee-app/internal/common/constant"
	commonFilter "employee-app/internal/common/model"
	"employee-app/internal/model"
	"employee-app/internal/usecase/attendance"
	"employee-app/internal/usecase/employee"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	usecase        attendance.AttendanceUsecase
	employeeUsecase employee.EmployeeUsecase
}

func NewAttendanceHandler(u attendance.AttendanceUsecase, e employee.EmployeeUsecase) *AttendanceHandler {
	return &AttendanceHandler{usecase: u, employeeUsecase: e}
}

// POST /attendance/{id}/submit
func (h *AttendanceHandler) SubmitAttendance(c *gin.Context) {
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

	var req dto.SubmitAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.usecase.SubmitAttendance(c, empID, req.Time, req.AttendanceType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "attendance submitted successfully"})
}


// GET /attendance
func (h *AttendanceHandler) GetAttendanceForPeriod(c *gin.Context) {
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

	attendances, err := h.usecase.GetAttendanceForPeriod(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve attendance"})
		return
	}

	// Map to DTO response
	var response []dto.AttendanceResponse
	for _, a := range attendances {
		resp := dto.AttendanceResponse{
			ID:             a.ID,
			EmployeeID:     a.EmployeeID,
			AttendanceDate: a.AttendanceDate.Format("2006-01-02"),
			CheckInTime:    a.CheckInTime.Format("15:04:05"),
			CreatedAt:      a.CreatedAt,
			UpdatedAt:      a.UpdatedAt,
		}

		if a.CheckOutTime != nil {
			formatted := a.CheckOutTime.Format("15:04:05")
			resp.CheckOutTime = &formatted
		}
		response = append(response, resp)
	}

	c.JSON(http.StatusOK, response)
}