package common

import "math"

type CommonFilter struct {
	SortBy   string `json:"sortBy,omitempty"`   // e.g., "name", "created_at"
	Ascending bool  `json:"ascending"`          // true = ASC, false = DESC
	Limit    *int   `json:"limit,omitempty"`    // Optional limit, defaults to 20
	Page     int    `json:"page,omitempty"`     // 1-based page number
}

// NewCommonFilter creates a CommonFilter with default values if missing.
func NewCommonFilter() *CommonFilter {
	defaultLimit := 20
	return &CommonFilter{
		Ascending: true,
		Limit:     &defaultLimit,
		Page:      1,
	}
}

func (f *CommonFilter) GetSortByOrDefault(defaultSortBy string) string {
	if f.SortBy != "" {
		return f.SortBy
	}
	return defaultSortBy
}

func (f *CommonFilter) GetSortBySQL() string {
	if f.Ascending {
		return "ASC"
	}
	return "DESC"
}

func (f *CommonFilter) GetLimitOrDefault(defaultLimit int) *int {
	if f.Limit == nil {
		return nil
	}
	limit := int(math.Min(float64(abs(*f.Limit)), float64(defaultLimit)))
	return &limit
}

func (f *CommonFilter) GetOffset() int {
	if f.Limit == nil {
		return 0
	}
	page := f.Page
	if page < 1 {
		page = 1
	}
	return (page - 1) * abs(*f.Limit)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}