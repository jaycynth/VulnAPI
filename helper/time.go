package helper

import (
	"fmt"
	"time"
)

func TimeToString(timeStr string) time.Time {
	layout := "2006-01-02 15:04:05"

	parsedTime, err := time.Parse(layout, timeStr)
	if err != nil {
		fmt.Println("Error parsing time:", err)
	}

	return parsedTime

}
