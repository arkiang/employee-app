package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Payslip struct {
	ID            uint            `gorm:"primaryKey"`
	EmployeeID    uint            `gorm:"not null;index"` // FK to Employee
	Employee      Employee        `gorm:"foreignKey:EmployeeID"`

	PeriodStart   time.Time       `gorm:"not null;type:date;index"` // Payroll start date
	PeriodEnd     time.Time       `gorm:"not null;type:date;index"` // Payroll end date

	BaseSalary    decimal.Decimal `gorm:"type:numeric(12,2);not null"` // Monthly base salary
	AttendanceDays int             `gorm:"not null"` // Number of working days attended
	AttendancePay  decimal.Decimal `gorm:"type:numeric(12,2);not null"`
	OvertimeHours  int             `gorm:"not null"` // Total overtime hours
	OvertimePay    decimal.Decimal `gorm:"type:numeric(12,2);not null"`
	ReimbursementTotal decimal.Decimal `gorm:"type:numeric(12,2);not null"`
	TakeHomePay   decimal.Decimal `gorm:"type:numeric(12,2);not null"`

	Attendances    []PayslipAttendance    `gorm:"foreignKey:PayslipID"`
	Overtimes      []PayslipOvertime      `gorm:"foreignKey:PayslipID"`
	Reimbursements []PayslipReimbursement `gorm:"foreignKey:PayslipID"`

	CreatedBy     uint       `gorm:"column:created_by;not null"`
	UpdatedBy     uint       `gorm:"column:updated_by;not null"`
	CreatedAt     time.Time      `gorm:"not null;autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"not null;autoUpdateTime"`
}