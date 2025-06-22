package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type UserLoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required,min=8,max=128"`
}

type RegisterEmployeeRequest struct {
	// User fields
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=employee"` // Only employee allowed

	// Employee fields
	Name   string          `json:"name" binding:"required,min=3,max=100"`
	Salary decimal.Decimal `json:"salary" binding:"required,gt=0"`
}

type RegisterEmployeeResponse struct {
	UserID     uint            `json:"userId"`
	Username   string          `json:"username"`
	EmployeeID uint            `json:"employeeId"`
	Name       string          `json:"name"`
	Salary     decimal.Decimal `json:"salary"`
	Role       string          `json:"role"`
	CreatedAt  time.Time       `json:"createdAt"`
}

type UpdateEmployeeRequest struct {
	Name   string          `json:"name,omitempty"`
	Salary decimal.Decimal `json:"salary,omitempty"`
	Role   string          `json:"role,omitempty"`
}

type ListEmployeeRequest struct {
	SortBy    string `form:"sortBy"`             // e.g., name or created_at
	Ascending bool   `form:"ascending"`          // true = ASC, false = DESC
	Limit     *int   `form:"limit"`              // optional
	Page      int    `form:"page"`               // default to 1
}

type EmployeeResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Salary    float64   `json:"salary"`
	UserID    uint      `json:"userId"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}