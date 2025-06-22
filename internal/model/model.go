package model

import (
	common "employee-app/internal/common/model"
	"time"
)

type AttendanceType string
const (
	CHECK_IN AttendanceType = "Check-In"
	CHECK_OUT AttendanceType = "Check-Out"
)

type EmployeePeriodFilter struct {
	Base	common.CommonFilter `json:"base"`
	EmpIds  *[]string           `json:"employeeId,omitempty"`
	Start   *time.Time          `json:"start,omitempty"`
	End     *time.Time          `json:"end,omitempty"`
}

type PeriodFilter struct {
	Base     common.CommonFilter `json:"base"`
	PeriodID *string             `json:"periodId,omitempty"`
}

type contextKey string
const userIDKey contextKey = "userID"