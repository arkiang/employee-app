package migration

import (
	"employee-app/internal/model/entity"
	"log"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&entity.User{},
		&entity.Employee{},
		&entity.Attendance{},
		&entity.Overtime{},
		&entity.Reimbursement{},
		&entity.PayrollPeriod{},
		&entity.Payslip{},
		&entity.PayslipAttendance{},
		&entity.PayslipOvertime{},
		&entity.PayslipReimbursement{},
	)

	if err != nil {
		log.Fatalf("Failed to run AutoMigrate: %v", err)
	}

	log.Println("Database migration completed.")
}