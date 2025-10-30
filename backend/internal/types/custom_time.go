package types

import (
	"encoding/json"
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		return nil
	}

	formats := []string{
		time.RFC3339,          // 2025-10-30T23:25:00Z
		"2006-01-02T15:04",    // 2025-10-30T23:25
		"2006-01-02 15:04:05", // 2025-10-30 23:25:00
		"2006-01-02",          // 2025-10-30
	}

	var parsed time.Time
	var err error
	for _, f := range formats {
		parsed, err = time.Parse(f, s)
		if err == nil {
			ct.Time = parsed
			return nil
		}
	}
	return err
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(ct.Time.Format(time.RFC3339))
}
