package seed

import (
	"context"
	"employee-app/internal/model/entity"
	"employee-app/internal/repository"
	"employee-app/internal/usecase/registration"
	"employee-app/pkg/security"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func SeedAdmin(ctx context.Context, userRepo repository.UserRepository) {
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

	err = userRepo.WithTransaction(ctx, func(tx *gorm.DB) error {
		user := entity.User{
			Username:     "admin",
			Role:         "admin",
			Salt:         salt,
			PasswordHash: hashedPassword,
		}

		_, err := userRepo.CreateTx(ctx, tx, &user)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		panic(fmt.Sprintf("Failed to create admin user: %v", err))
	}

	fmt.Println("Admin user seeded successfully.")
}

func SeedEmployeesFromCSV(ctx context.Context, userRepo repository.UserRepository, regUC registration.RegistrationUsecase) {
	file, err := os.Open("internal/seed/user.csv")
	if err != nil {
		panic(fmt.Sprintf("Failed to open CSV: %v", err))
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic(fmt.Sprintf("Failed to read CSV: %v", err))
	}

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
	}
}

func SeedInitialUsers(ctx context.Context, regUC registration.RegistrationUsecase) error {
	user := entity.User{
		Username: "user1",
		Role:     "employee",
	}

	employee := entity.Employee{
		Name:   "John Doe",
		Salary: decimal.NewFromFloat(5000000), // dummy salary
	}
	// Seed Employee
	err := regUC.RegisterEmployee(ctx, user, employee, "12345")
	if err != nil {
		return fmt.Errorf("failed to seed employee user: %w", err)
	}

	fmt.Println("✅ Employee user seeded.")
	return nil
}