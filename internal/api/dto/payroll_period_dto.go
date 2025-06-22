package dto

import "time"

type CreatePayrollPeriodRequest struct {
	StartDate time.Time `json:"startDate" binding:"required"` // Start of payroll period
	EndDate   time.Time `json:"endDate" binding:"required"`   // End of payroll period
}

type PayrollPeriodResponse struct {
	ID        uint      `json:"id"`
	StartDate string    `json:"startDate"` // Format: YYYY-MM-DD
	EndDate   string    `json:"endDate"`   // Format: YYYY-MM-DD
	CreatedBy uint      `json:"createdBy"`
	UpdatedBy uint      `json:"updatedBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}