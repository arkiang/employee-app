package entity

import "time"
import "github.com/shopspring/decimal"

type Reimbursement struct {
	ID                uint            `gorm:"primaryKey"`
	EmployeeID        uint            `gorm:"not null"`
	Employee          Employee        `gorm:"foreignKey:EmployeeID"`
	ReimbursementDate time.Time       `gorm:"type:date;not null"`
	Amount            decimal.Decimal `gorm:"type:numeric(10,2);not null"`
	Description       *string         `gorm:"type:text"`
	CreatedBy         uint            `gorm:"column:created_by;not null"`
	UpdatedBy         uint            `gorm:"column:updated_by;not null"`
	CreatedAt         time.Time       `gorm:"column:created_at;autoCreateTime;not null"`
	UpdatedAt         time.Time       `gorm:"column:updated_at;autoUpdateTime;not null"`
}