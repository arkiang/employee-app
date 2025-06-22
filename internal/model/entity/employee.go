package entity

import "time"
import "github.com/shopspring/decimal"

type Employee struct {
	ID        uint              `gorm:"primaryKey"`
	UserID    uint              `gorm:"uniqueIndex;not null"`
	User      User              `gorm:"foreignKey:UserID"`
	Name      string            `gorm:"type:varchar(100);not null"`
	Salary    decimal.Decimal   `gorm:"type:numeric(12,2);not null"`
	CreatedAt time.Time         `gorm:"column:created_at;autoCreateTime;not null"`
	UpdatedAt time.Time         `gorm:"column:updated_at;autoUpdateTime;not null"`
}
