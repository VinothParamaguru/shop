package utilities

import "time"

// GetCurrentTimeStampString
// Get the current timestamp as string
func GetCurrentTimeStampString() string {
	now := time.Now().Format(time.RFC3339)
	return now
}

// GetEpochTime
// Get the epoch time from the timestamp string specified
func GetEpochTime(timestamp string) int64 {
	givenTime, _ := time.Parse(time.RFC3339, timestamp)
	return givenTime.Unix()
}
