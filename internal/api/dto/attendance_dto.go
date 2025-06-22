package dto

import (
	"employee-app/internal/model"
	"time"
)

type SubmitAttendanceRequest struct {
	Time           int64            `json:"time" binding:"required"`                       // Timestamp of submission
	AttendanceType model.AttendanceType `json:"attendanceType" binding:"required,oneof=Check-In Check-Out"` // Enum (Check-In / Check-Out)
}

type AttendanceResponse struct {
	ID             uint                 `json:"id"`
	EmployeeID     uint                 `json:"employeeId"`
	AttendanceDate string               `json:"attendanceDate"` // Format: YYYY-MM-DD
	CheckInTime    string               `json:"checkInTime,omitempty"` // Format: HH:MM:SS
	CheckOutTime   *string               `json:"checkOutTime,omitempty"` // Format: HH:MM:SS
	CreatedAt      time.Time            `json:"createdAt"`
	UpdatedAt      time.Time            `json:"updatedAt"`
}