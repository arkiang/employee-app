package entity

import "time"

type PayrollPeriod struct {
	ID         uint      `gorm:"primaryKey"`
	Name	   string    `gorm:"type:varchar(50);not null;uniqueIndex"`
	StartDate  time.Time `gorm:"type:date;not null"`
	EndDate    time.Time `gorm:"type:date;not null"`
	Processed  bool      `gorm:"default:false"`
	CreatedBy  uint      `gorm:"column:created_by"`
	UpdatedBy  uint      `gorm:"column:updated_by"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime;not null"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime;not null"`
}