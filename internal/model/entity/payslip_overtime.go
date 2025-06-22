package entity

import "time"

type PayslipOvertime struct {
	ID         uint      `gorm:"primaryKey"`
	PayslipID  uint      `gorm:"not null"`
	Date       time.Time `gorm:"not null"`
	Hours      uint8     `gorm:"not null"`
	Amount     float64   `gorm:"not null"`
}