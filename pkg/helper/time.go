package helper

import (
	"fmt"
	"time"
)

func ParseToDate(value string) *time.Time {
	layout := "2006-01-02"

	time, err := time.Parse(layout, value)
	if err != nil {
		panic(fmt.Sprintf("error parse data from %s", value))
	}

	return &time
}
