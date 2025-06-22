package entity

import "time"

type Overtime struct {
	ID             uint      `gorm:"primaryKey"`
	EmployeeID     uint      `gorm:"not null;index:idx_employee_overtime_date,unique"`
	Employee       Employee  `gorm:"foreignKey:EmployeeID"`
	OvertimeDate   time.Time `gorm:"type:date;not null;index:idx_employee_overtime_date,unique"`
	Hours          uint8      `gorm:"not null"`
	CreatedBy      uint      `gorm:"column:created_by;not null"`
	UpdatedBy      uint      `gorm:"column:updated_by;not null"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime;not null"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime;not null"`
}