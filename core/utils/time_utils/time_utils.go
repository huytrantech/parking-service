package time_utils

import "time"

func ConvertFromTimeToSecond(time *time.Time) int {
	if time == nil {
		return 0
	}
	return time.Hour() * 60 * 60 + time.Minute()*60 + 59
}

