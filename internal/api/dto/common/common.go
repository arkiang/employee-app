package common

import (
	common "employee-app/internal/common/model"
	"employee-app/internal/model"
	"time"
)

type CommonQueryParams struct {
	SortBy      string     `form:"sortBy"`                     // e.g. "created_at"
	Ascending   bool       `form:"ascending"`                  // ASC or DESC
	Limit       *int       `form:"limit"`                      // Optional: max items
	Page        int        `form:"page"`                       // Default: 1
}

type EmployeePeriodQueryParams struct {
	EmployeeIDs []string   `form:"employeeId[]"`               // Optional: filter by multiple emp IDs
	Start       *time.Time `form:"start"`                      // Optional: start date
	End         *time.Time `form:"end"`                        // Optional: end date
	SortBy      string     `form:"sortBy"`                     // e.g. "created_at"
	Ascending   bool       `form:"ascending"`                  // ASC or DESC
	Limit       *int       `form:"limit"`                      // Optional: max items
	Page        int        `form:"page"`                       // Default: 1
}

func (q CommonQueryParams) ToFilter() common.CommonFilter {
	// Apply defaults
	defaultSort := "created_at"
	defaultPage := 1

	sort := q.SortBy
	if sort == "" {
		sort = defaultSort
	}

	page := q.Page
	if page <= 0 {
		page = defaultPage
	}

	return common.CommonFilter{
		SortBy:    sort,
		Ascending: q.Ascending,
		Limit:     q.Limit,
		Page:      page,
	}
}

func (q EmployeePeriodQueryParams) ToFilter() model.EmployeePeriodFilter  {
	// Apply defaults
	defaultSort := "created_at"
	defaultPage := 1

	sort := q.SortBy
	if sort == "" {
		sort = defaultSort
	}

	page := q.Page
	if page <= 0 {
		page = defaultPage
	}

	return model.EmployeePeriodFilter{
		Base: common.CommonFilter{
			SortBy:    sort,
			Ascending: q.Ascending,
			Limit:     q.Limit,
			Page:      page,
		},
		EmpIds:  &q.EmployeeIDs,
		Start:   q.Start,
		End:     q.End,
	}
}