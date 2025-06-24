package main

import (
	"context"
	"employee-app/configs"
	persistent "employee-app/infrastructure/peristent"
	"employee-app/internal/api/handler"
	http "employee-app/internal/api/v1"
	"employee-app/internal/migration"
	"employee-app/internal/seed"
	"employee-app/internal/usecase/attendance"
	"employee-app/internal/usecase/employee"
	"employee-app/internal/usecase/overtime"
	"employee-app/internal/usecase/payroll"
	"employee-app/internal/usecase/payslip"
	"employee-app/internal/usecase/registration"
	"employee-app/internal/usecase/reimbursement"
	"employee-app/internal/usecase/user"
	"log"
)

func main() {
	err := configs.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables ", err)
	}

	db, err := configs.ConnectDB(&configs.AppConfig)
	if err != nil {
		log.Fatal("Could not connect to database", err)
	}

	// AutoMigrate database
	migration.AutoMigrate(db)

	// Repository layer
	userRepo := persistent.NewUserRepository(db)
	employeeRepo := persistent.NewEmployee(db)
	attendanceRepo := persistent.NewAttendance(db)
	overtimeRepo := persistent.NewOvertime(db)
	reimbursementRepo := persistent.NewReimbursement(db)
	periodRepo := persistent.NewPayrollPeriod(db)
	payslipRepo := persistent.NewPayslip(db)

	// Usecase layer
	userUsecase := user.New(userRepo, employeeRepo)
	registrationUsecase := registration.New(userRepo, employeeRepo)
	employeeUsecase := employee.New(employeeRepo)
	attendanceUsecase := attendance.New(attendanceRepo)
	overtimeUsecase := overtime.New(overtimeRepo, attendanceRepo)
	reimbursementUsecase := reimbursement.New(reimbursementRepo)
	periodUsecase := payroll.New(periodRepo)
	payslipUsecase := payslip.New(employeeRepo, attendanceRepo, overtimeRepo, reimbursementRepo, periodRepo, payslipRepo)

	// Seed initial data
	seed.SeedAdmin(context.Background(), userRepo, periodRepo)
	seed.SeedEmployeesFromCSV(context.Background(), userRepo, employeeRepo, attendanceRepo, overtimeRepo, reimbursementRepo, registrationUsecase)

	// Handler layer
	userHandler := handler.NewUserHandler(userUsecase)
	registrationHandler := handler.NewRegistrationHandler(registrationUsecase)
	employeeHandler := handler.NewEmployeeHandler(employeeUsecase)
	attendanceHandler := handler.NewAttendanceHandler(attendanceUsecase, employeeUsecase)
	overtimeHandler := handler.NewOvertimeHandler(overtimeUsecase, employeeUsecase)
	reimbursementHandler := handler.NewReimbursementHandler(reimbursementUsecase, employeeUsecase)
	periodHandler := handler.NewPayrollPeriodHandler(periodUsecase)
	payslipHandler := handler.NewPayslipHandler(payslipUsecase)

	router := http.SetupRouter(
		userHandler,
		attendanceHandler,
		overtimeHandler,
		reimbursementHandler,
		payslipHandler,
		employeeHandler,
		periodHandler,
		registrationHandler,
	)

	log.Printf("Server is running on port %s...", configs.AppConfig.ServerPort)
	if err := router.Run("0.0.0.0:" + configs.AppConfig.ServerPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

