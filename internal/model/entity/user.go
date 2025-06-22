package entity

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Username     string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	Salt         string    `gorm:"type:text;not null"`
	PasswordHash string    `gorm:"type:text;not null"`
	Role         string    `gorm:"type:varchar(10);not null;check:role IN ('admin','employee')"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime;not null"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime;not null"`
}