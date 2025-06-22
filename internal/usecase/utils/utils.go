package utils

import (
	"employee-app/internal/model/entity"

	"github.com/shopspring/decimal"
)

func CountUniqueDates(attendances []*entity.Attendance) int {
	dates := make(map[string]bool)
	for _, a := range attendances {
		dates[a.AttendanceDate.Format("2006-01-02")] = true
	}
	return len(dates)
}

func SumOvertimeHours(overtimes []*entity.Overtime) int {
	total := int(0)
	for _, o := range overtimes {
		total += int(o.Hours)
	}
	return total
}

func SumReimbursementAmounts(reimbursements []*entity.Reimbursement) decimal.Decimal {
	total := decimal.NewFromInt(0)
	for _, r := range reimbursements {
		total = total.Add(r.Amount)
	}
	return total
}