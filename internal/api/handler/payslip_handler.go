package handler

import (
	"employee-app/internal/api/dto"
	"employee-app/internal/api/dto/common"
	"employee-app/internal/api/handler/utils"
	"employee-app/internal/common/constant"
	commonFilter "employee-app/internal/common/model"
	"employee-app/internal/model/entity"
	"employee-app/internal/usecase/payslip"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type PayslipHandler struct {
	payslipUC payslip.PayslipUsecase
}

func NewPayslipHandler(payslipUC payslip.PayslipUsecase) *PayslipHandler {
	return &PayslipHandler{payslipUC: payslipUC}
}

// POST /payroll-periods/:id/run
func (h *PayslipHandler) RunPayroll(c *gin.Context) {
	_, status, errMsg := utils.GetUserID(c, constant.EnumRoleAdmin)
	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	periodID, err := strconv.Atoi(c.Param("id"))
	if err != nil || periodID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period ID"})
		return
	}

	if err := h.payslipUC.RunPayroll(c, uint(periodID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// GET /payslips/:periodId/me
func (h *PayslipHandler) GetPayslipForMe(c *gin.Context) {
		_, status, errMsg := utils.GetUserID(c, constant.EnumRoleEmployee)
	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	periodID, err := strconv.Atoi(c.Param("periodId"))
	if err != nil || periodID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period ID"})
		return
	}

	payslip, err := h.payslipUC.GetPayslipForEmployee(c, uint(periodID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payslip not found" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, payslip)
}

func (h *PayslipHandler) GetPayslipSummary(c *gin.Context) {
	_, status, errMsg := utils.GetUserID(c, constant.EnumRoleAdmin)
	if status != http.StatusOK {
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	periodID, err := strconv.Atoi(c.Param("periodId"))
	if err != nil || periodID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period ID"})
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

	payslips, err := h.payslipUC.GetPayslipSummary(c, uint(periodID), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var summaries []dto.PayslipSummaryDTO
	total := decimal.Zero

	for _, p := range payslips {
		summaries = append(summaries, dto.PayslipSummaryDTO{
			EmployeeID:    p.EmployeeID,
			EmployeeName:  p.Employee.Name,
			TotalTakeHome: p.TakeHomePay,
		})
		total = total.Add(p.TakeHomePay)
	}

	c.JSON(http.StatusOK, dto.PayslipSummaryListDTO{
		Summaries:    summaries,
		TotalPayroll: total,
	})
}

func toPayslipResponse(p *entity.Payslip) dto.PayslipResponse {
	attendances := make([]dto.PayslipAttendanceDTO, 0)
	for _, a := range p.Attendances {
		if a.CheckOut == nil {
			continue // atau log.Println("nil attendance entry")
		}
		attendances = append(attendances, dto.PayslipAttendanceDTO{
			Date:         a.Date.Format("2006-01-02"),
			CheckInTime:  a.CheckIn,
			CheckOutTime: a.CheckOut,
		})
	}

	overtimes := make([]dto.PayslipOvertimeDTO, 0)
	for _, o := range p.Overtimes {
		overtimes = append(overtimes, dto.PayslipOvertimeDTO{
			Date:  o.Date.Format("2006-01-02"),
			Hours: o.Hours,
		})
	}

	reimbursements := make([]dto.PayslipReimbursementDTO, 0)
	for _, r := range p.Reimbursements {
		reimbursements = append(reimbursements, dto.PayslipReimbursementDTO{
			Date:        r.Date.Format("2006-01-02"),
			Amount:      r.Amount,
			Description: r.Description,
		})
	}

	return dto.PayslipResponse{
		PayslipID:       p.ID,
		EmployeeID:      p.EmployeeID,
		EmployeeName:    p.Employee.Name,
		PeriodStartDate: p.PeriodStart.Format("2006-01-02"),
		PeriodEndDate:   p.PeriodEnd.Format("2006-01-02"),
		AttendanceDays:  p.AttendanceDays,
		BaseSalary:      p.BaseSalary,
		OvertimePay:     p.OvertimePay,
		Reimbursement:   p.ReimbursementTotal,
		TotalTakeHome:   p.TakeHomePay,
		Attendances:     attendances,
		Overtimes:       overtimes,
		Reimbursements:  reimbursements,
		GeneratedAt:     p.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}