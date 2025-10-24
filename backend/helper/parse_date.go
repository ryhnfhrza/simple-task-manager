package helper

import (
	"errors"
	"time"
)

func ParseFlexibleDate(dateStr string) (time.Time, error) {
	layouts := []string{
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("invalid date format (use YYYY-MM-DD or YYYY-MM-DDTHH:MM:SS)")
}
