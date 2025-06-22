package dto

import (
	"employee-app/internal/api/dto/common"
	"time"

	"github.com/shopspring/decimal"
)

type SubmitReimbursementRequest struct {
	Date        common.DateOnly   `json:"date" binding:"required"`             // Reimbursement date (YYYY-MM-DD)
	Amount      string            `json:"amount" binding:"required,gt=0"`      // Must be greater than 0
	Description *string           `json:"description,omitempty"`               // Optional description
}

type ReimbursementResponse struct {
	ID          uint            `json:"id"`
	EmployeeID  uint            `json:"employeeId"`
	Date        string          `json:"date"`        // YYYY-MM-DD
	Amount      decimal.Decimal `json:"amount"`
	Description *string         `json:"description,omitempty"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}
