package http

import (
	"employee-app/internal/api/handler"
	"employee-app/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userHandler *handler.UserHandler,
	attendanceHandler *handler.AttendanceHandler,
	overtimeHandler *handler.OvertimeHandler,
	reimbursementHandler *handler.ReimbursementHandler,
	payslipHandler *handler.PayslipHandler,
	employeeHandler *handler.EmployeeHandler,
	periodHandler *handler.PayrollPeriodHandler,
	registrationHandler *handler.RegistrationHandler,
) *gin.Engine {
	router := gin.Default()

	// CORS, logger, recovery, etc. can be added here
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.CORSMiddleware())

	// Public routes
	public := router.Group("/api")
	public.POST("/login", userHandler.Login)
	public.POST("/register", registrationHandler.Register)

	// Auth-protected routes
	protected := public.Group("/")
	protected.Use(middleware.AuthMiddleware())

	// User
	protected.GET("/users/me", userHandler.GetByID)

	// Employee
	protected.GET("/employees", employeeHandler.ListEmployees)
	protected.GET("/employees/:id", employeeHandler.GetEmployeeByID)
	protected.PUT("/employees/:id", employeeHandler.UpdateEmployee)

	// Attendance
	protected.POST("/attendance/:empId", attendanceHandler.SubmitAttendance)
	protected.GET("/attendance", attendanceHandler.GetAttendanceForPeriod)

	// Overtime
	protected.POST("/overtime/:empId", overtimeHandler.SubmitOvertime)
	protected.GET("/overtime", overtimeHandler.GetOvertimeForPeriod)

	// Reimbursement
	protected.POST("/reimbursements/:empId", reimbursementHandler.SubmitReimbursement)
	protected.GET("/reimbursements", reimbursementHandler.GetReimbursementsForPeriod)

	// Payroll Periods
	protected.POST("/payroll-periods", periodHandler.CreatePeriod)
	protected.GET("/payroll-periods", periodHandler.ListPeriods)
	protected.GET("/payroll-periods/:id", periodHandler.GetPeriodByID)

	// Payslips
	protected.POST("/payroll/:id/run", payslipHandler.RunPayroll)
	protected.GET("/payslips/:periodId/me", payslipHandler.GetPayslipForMe)
	protected.GET("/payslips/:periodId/summary", payslipHandler.GetPayslipSummary)

	return router
}