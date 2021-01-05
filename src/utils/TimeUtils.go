package utils

import (
	"fmt"
	"time"
)

func GetFormattedTimeAsString(timestamp time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d.%02d",
		timestamp.Year(), timestamp.Month(), timestamp.Day(),
		timestamp.Hour(), timestamp.Minute(), timestamp.Second(),
		padNanosecondDigits(timestamp.Nanosecond()/1000, 6))
}

func countNanosecondDigits(number int) int {
	digitCount := 0
	for number != 0 {
		number /= 10
		digitCount += 1
	}
	return digitCount
}

func padNanosecondDigits(number int, desiredLength int) int {
	outDigits := number
	for i := countNanosecondDigits(number); i < desiredLength; i++ {
		outDigits = outDigits * 10
	}
	return outDigits
}
