package entity

import "time"

type PayslipAttendance struct {
	ID         uint      `gorm:"primaryKey"`
	PayslipID  uint      `gorm:"not null"`
	Date       time.Time `gorm:"not null"`
	CheckIn    time.Time `gorm:"type:timestamp;not null"`
	CheckOut   *time.Time `gorm:"type:timestamp"` // Optional if user only checks out
}