package utils

import (
	"strings"
	"time"
)

func calculateStartDate() time.Time {
	year, month, day := time.Now().Date()
	// Set start of day to this morning 00:00
	todayStart := time.Date(year, month, day, 0, 0, 0, 0, time.Local)

	today6hStart := time.Date(year, month, day, 6, 0, 0, 0, time.Local)

	timeNow := time.Now()

	if timeNow.After(todayStart) && timeNow.Before(today6hStart) {
		// 2 Days back 06:00
		return today6hStart.AddDate(0, 0, -2)
	}
	// yesterday 06:00
	return today6hStart.AddDate(0, 0, -1)
}

func GetStartDate() string {
	return strings.Replace(calculateStartDate().Format(time.RFC3339), "+", ".0+", 1)
}

func GetEndDate() string {
	return strings.Replace(calculateStartDate().AddDate(0, 0, 1).Format(time.RFC3339), "+", ".0+", 1)
}
