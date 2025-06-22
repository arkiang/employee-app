package entity

import "time"
import "github.com/shopspring/decimal"

type PayslipReimbursement struct {
	ID          uint            `gorm:"primaryKey"`
	PayslipID   uint            `gorm:"not null"`
	Date        time.Time       `gorm:"type:date;not null"`
	Amount      decimal.Decimal `gorm:"type:numeric(12,2);not null"`
	Description *string         `gorm:"type:text"`
}