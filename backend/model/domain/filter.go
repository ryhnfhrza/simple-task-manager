package domain

import "time"

type TaskFilter struct {
	Completed *int
	DueBefore *time.Time
	DueAfter  *time.Time
	SortBy    string
	Order     string
	Limit     int
	Offset    int
}
