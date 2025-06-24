package seed

import (
	_ "embed"
	"strconv"
	"time"

	"context"
	"employee-app/internal/model/entity"
	"employee-app/internal/repository"
	"employee-app/internal/usecase/registration"
	"employee-app/pkg/security"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func SeedAdmin(ctx context.Context, userRepo repository.UserRepository, payrollRepo repository.PayrollPeriodRepository) {
	admin, _ := userRepo.GetByUsername(ctx, "admin");
	if admin != nil {
		fmt.Println("Admin already exist")
		return
	}
	// Generate salt and hash
	salt, err := security.GenerateSalt()
	if err != nil {
		panic(fmt.Sprintf("Failed to generate salt: %v", err))
	}

	hashedPassword, err := security.HashPassword("admin", salt)
	if err != nil {
		panic(fmt.Sprintf("Failed to hash password: %v", err))
	}

	var adminId uint

	err = userRepo.WithTransaction(ctx, func(tx *gorm.DB) error {
		user := entity.User{
			Username:     "admin",
			Role:         "admin",
			Salt:         salt,
			PasswordHash: hashedPassword,
		}

		admin, err := userRepo.CreateTx(ctx, tx, &user)
		if err != nil {
			return err
		}

		adminId = admin.ID
		return nil
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to create admin user: %v", err))
	}

	payroll := &entity.PayrollPeriod{
		Name:       "First Half June 2025 Payroll",
		StartDate:  time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate:    time.Date(2025, 6, 15, 23, 59, 59, 0, time.UTC),
		CreatedBy:  adminId,
		UpdatedBy:  adminId,
	}

	_, err = payrollRepo.Create(ctx, payroll)
	if err != nil {
		panic(fmt.Sprintf("Failed to create payroll period: %v", err))
	}

	fmt.Println("Admin and Payroll Period user seeded successfully.")
}

//go:embed user.csv
var userCSV string
func SeedEmployeesFromCSV(ctx context.Context, 
		userRepo repository.UserRepository, 
		employeeRepo repository.EmployeeRepository,
		attendanceRepo repository.AttendanceRepository,
		overtimeRepo repository.OvertimeRepository,
		reimbursementRepo repository.ReimbursementRepository,
		regUC registration.RegistrationUsecase) {
	reader := csv.NewReader(strings.NewReader(userCSV))
	records, err := reader.ReadAll()
	if err != nil {
		panic(fmt.Sprintf("Failed to read CSV: %v", err))
	}

	// Map to keep username -> employeeID
	usernameToEmpID := make(map[string]uint)

	// Skip header
	for i, row := range records {
		if i == 0 {
			continue
		}

		username := strings.TrimSpace(row[0])
		name := strings.TrimSpace(row[1])
		salaryStr := strings.TrimSpace(row[2])

		// Convert salary
		salary, err := decimal.NewFromString(salaryStr)
		if err != nil {
			fmt.Printf("Invalid salary for '%s': %v\n", username, err)
			continue
		}

		user := entity.User{
			Username: username,
			Role:     "employee",
		}

		employee := entity.Employee{
			Name:   name,
			Salary: salary,
		}

		userExist, _ := userRepo.GetByUsername(ctx, username);
		if userExist != nil {
			fmt.Printf("User %s already exists\n", username)
			continue
		}

		if err := regUC.RegisterEmployee(ctx, user, employee, "12345"); err != nil {
			fmt.Printf("Failed to register '%s': %v\n", username, err)
		} else {
			fmt.Printf("Registered employee '%s'\n", username)
		}

		insertedUser, _ := userRepo.GetByUsername(ctx, username);

		createdEmp, _ := employeeRepo.GetByUserID(ctx, insertedUser.ID)
		usernameToEmpID[username] = createdEmp.ID
	}
	
	seedAttendances(ctx, usernameToEmpID, attendanceRepo)
	seedOvertimes(ctx, usernameToEmpID, overtimeRepo)
	seedReimbursements(ctx, usernameToEmpID, reimbursementRepo)
}

//go:embed attendances.csv
var attendancesCSV string
func seedAttendances(ctx context.Context, usernameToEmpID map[string]uint, repo repository.AttendanceRepository) {
	reader := csv.NewReader(strings.NewReader(attendancesCSV))
	records, _ := reader.ReadAll()

	for i, row := range records {
		if i == 0 {
			continue
		}
		username := row[0]
		date, _ := time.Parse("2006-01-02", row[1])
		checkIn, _ := time.Parse("15:04:00", row[2])
		checkOut, _ := time.Parse("15:04:00", row[3])

		employeeID, ok := usernameToEmpID[username]
		if !ok {
			continue
		}

		_, err := repo.Create(ctx, &entity.Attendance{
			EmployeeID:     employeeID,
			AttendanceDate: date,
			CheckInTime:    time.Date(date.Year(), date.Month(), date.Day(), checkIn.Hour(), checkIn.Minute(), checkIn.Second(), 0, time.Local),
			CheckOutTime:   ptrTime(time.Date(date.Year(), date.Month(), date.Day(), checkOut.Hour(), checkOut.Minute(), checkOut.Second(), 0, time.Local)),
		})
		if err != nil {
			fmt.Printf("Failed to create attendance for %s: %v\n", username, err)
		}
	}
}

//go:embed overtimes.csv
var overtimesCSV string
func seedOvertimes(ctx context.Context, usernameToEmpID map[string]uint, repo repository.OvertimeRepository) {
	reader := csv.NewReader(strings.NewReader(overtimesCSV))
	records, _ := reader.ReadAll()

	for i, row := range records {
		if i == 0 {
			continue
		}
		username := row[0]
		date, _ := time.Parse("2006-01-02", row[1])
		hours, _ := strconv.ParseUint(row[2], 10, 8)

		employeeID, ok := usernameToEmpID[username]
		if !ok {
			continue
		}

		_, err := repo.Create(ctx, &entity.Overtime{
			EmployeeID:   employeeID,
			OvertimeDate: date,
			Hours:        uint8(hours),
		})
		if err != nil {
			fmt.Printf("Failed to create overtime for %s: %v\n", username, err)
		}
	}
}

//go:embed reimbursements.csv
var reimbursementsCSV string
func seedReimbursements(ctx context.Context, usernameToEmpID map[string]uint, repo repository.ReimbursementRepository) {
	reader := csv.NewReader(strings.NewReader(reimbursementsCSV))
	records, _ := reader.ReadAll()

	for i, row := range records {
		if i == 0 {
			continue
		}
		username := row[0]
		date, _ := time.Parse("2006-01-02", row[1])
		amount, _ := decimal.NewFromString(row[2])
		desc := row[3]

		employeeID, ok := usernameToEmpID[username]
		if !ok {
			continue
		}

		_, err := repo.Create(ctx, &entity.Reimbursement{
			EmployeeID:        employeeID,
			ReimbursementDate: date,
			Amount:            amount,
			Description:       &desc,
		})
		if err != nil {
			fmt.Printf("Failed to create reimbursement for %s: %v\n", username, err)
		}
	}
}

func ptrTime(t time.Time) *time.Time {
	return &t
}