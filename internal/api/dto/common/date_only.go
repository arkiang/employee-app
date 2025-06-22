package common

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type DateOnly struct {
	time.Time
}

const layoutDateOnly = "02-01-2006" // 👈 DD-MM-YYYY

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse(layoutDateOnly, s)
	if err != nil {
		return fmt.Errorf("invalid date format: must be DD-MM-YYYY")
	}
	d.Time = t
	return nil
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format(layoutDateOnly))
}

func (d DateOnly) ToTime() time.Time {
	return d.Time
}