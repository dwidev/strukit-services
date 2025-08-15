package helper

import (
	"fmt"
	"time"
)

func ParseToDate(value *string) *time.Time {
	if value == nil {
		return nil
	}

	layout := "2006-01-02"
	time, err := time.Parse(layout, *value)
	if err != nil {
		panic(fmt.Sprintf("error parse data from %s", *value))
	}

	return &time
}

// Helper function untuk parse time-only dengan multiple format
func ParseTimeOnly(timeStr string) *time.Time {
	formats := []string{
		"15:04:05",   // HH:MM:SS
		"15:04",      // HH:MM
		"3:04:05 PM", // H:MM:SS AM/PM
		"3:04 PM",    // H:MM AM/PM
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return &t
		}
	}

	// Fallback ke zero time jika gagal parse
	return nil
}
