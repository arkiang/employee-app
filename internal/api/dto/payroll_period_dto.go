package dto

import (
	"employee-app/internal/api/dto/common"
	"time"
)

type CreatePayrollPeriodRequest struct {
	StartDate common.DateOnly `json:"startDate" binding:"required"` // Start of payroll period
	EndDate   common.DateOnly `json:"endDate" binding:"required"`   // End of payroll period
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