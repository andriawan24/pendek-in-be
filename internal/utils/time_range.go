package utils

import "time"

type TimeRange string

const (
	TimeRange7Days  TimeRange = "7d"
	TimeRange30Days TimeRange = "30d"
	TimeRange90Days TimeRange = "90d"
	TimeRangeAll    TimeRange = "all"
)

func (t TimeRange) IsValid() bool {
	switch t {
	case TimeRange7Days, TimeRange30Days, TimeRange90Days, TimeRangeAll:
		return true
	}
	return false
}

func (t TimeRange) GetFromDate() time.Time {
	now := time.Now()
	switch t {
	case TimeRange7Days:
		return now.AddDate(0, 0, -7)
	case TimeRange30Days:
		return now.AddDate(0, 0, -30)
	case TimeRange90Days:
		return now.AddDate(0, 0, -90)
	case TimeRangeAll:
		return time.Time{}
	default:
		return now.AddDate(0, 0, -30)
	}
}

func ParseTimeRange(s string) TimeRange {
	switch s {
	case "7d":
		return TimeRange7Days
	case "30d":
		return TimeRange30Days
	case "90d":
		return TimeRange90Days
	case "all":
		return TimeRangeAll
	default:
		return TimeRange30Days
	}
}
