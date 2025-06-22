package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type PayslipAttendanceDTO struct {
	Date         string    `json:"date"`         // YYYY-MM-DD
	CheckInTime  time.Time `json:"checkInTime"`
	CheckOutTime *time.Time `json:"checkOutTime"`
}

type PayslipOvertimeDTO struct {
	Date  string          `json:"date"`  // YYYY-MM-DD
	Hours uint8           `json:"hours"` // 1–3
}

type PayslipReimbursementDTO struct {
	Date        string          `json:"date"`        // YYYY-MM-DD
	Amount      decimal.Decimal `json:"amount"`
	Description *string         `json:"description,omitempty"`
}

type PayslipResponse struct {
	PayslipID       uint                      `json:"payslipId"`
	EmployeeID      uint                      `json:"employeeId"`
	EmployeeName    string                    `json:"employeeName"`
	PeriodID        uint                      `json:"periodId"`
	PeriodStartDate string                    `json:"periodStartDate"` // YYYY-MM-DD
	PeriodEndDate   string                    `json:"periodEndDate"`   // YYYY-MM-DD

	// Breakdown
	AttendanceDays  int                       `json:"attendanceDays"`
	BaseSalary      decimal.Decimal           `json:"baseSalary"`
	ProratedSalary  decimal.Decimal           `json:"proratedSalary"`
	OvertimePay     decimal.Decimal           `json:"overtimePay"`
	Reimbursement   decimal.Decimal           `json:"reimbursement"`
	TotalTakeHome   decimal.Decimal           `json:"totalTakeHome"`

	// Breakdown details
	Attendances     []PayslipAttendanceDTO    `json:"attendances"`
	Overtimes       []PayslipOvertimeDTO      `json:"overtimes"`
	Reimbursements  []PayslipReimbursementDTO `json:"reimbursements"`

	GeneratedAt     string                 `json:"generatedAt"`
}

type PayslipSummaryDTO struct {
	EmployeeID     uint            `json:"employeeId"`
	EmployeeName   string          `json:"employeeName"`
	TotalTakeHome  decimal.Decimal `json:"totalTakeHome"`
}

type PayslipSummaryListDTO struct {
	Summaries     []PayslipSummaryDTO `json:"summaries"`
	TotalPayroll  decimal.Decimal     `json:"totalPayroll"`
}