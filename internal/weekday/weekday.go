package weekday

import "time"

func Parse(weekday []int) []time.Weekday {
	weekdays := make([]time.Weekday, 0)
	for _, w := range weekday {
		weekdays = append(weekdays, time.Weekday(w))
	}
	return weekdays
}

func Next(startDate time.Time, weekday time.Weekday) time.Time {
	for startDate.Weekday() != weekday {
		startDate = startDate.AddDate(0, 0, 1)
	}
	return startDate
}
