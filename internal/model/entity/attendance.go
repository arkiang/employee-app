package entity

import "time"

type Attendance struct {
	ID             uint       `gorm:"primaryKey"`
	EmployeeID     uint       `gorm:"not null;index:idx_employee_attendance_date,unique"`
	Employee       Employee   `gorm:"foreignKey:EmployeeID"`
	AttendanceDate time.Time  `gorm:"type:date;not null;index:idx_employee_attendance_date,unique"`
	CheckInTime    time.Time  `gorm:"type:timestamp;not null"`
	CheckOutTime   *time.Time `gorm:"type:timestamp"` // Optional if user only checks out
	CreatedBy      uint       `gorm:"column:created_by;not null"`
	UpdatedBy      uint       `gorm:"column:updated_by;not null"`
	CreatedAt      time.Time  `gorm:"column:created_at;autoCreateTime;not null"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;autoUpdateTime;not null"`
}