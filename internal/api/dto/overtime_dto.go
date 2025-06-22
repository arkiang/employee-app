package dto

import "time"

type SubmitOvertimeRequest struct {
	Date  time.Time `json:"date" binding:"required"`            // Date of overtime (YYYY-MM-DD)
	Hours uint8     `json:"hours" binding:"required,oneof=1 2 3"` // Only 1, 2, or 3 are valid
}

type OvertimeResponse struct {
	ID         uint      `json:"id"`
	EmployeeID uint      `json:"employeeId"`
	Date       string    `json:"date"`        // YYYY-MM-DD
	Hours      uint8     `json:"hours"`       // 1, 2, or 3
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}