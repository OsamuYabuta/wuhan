package utils

import (
	"strings"
	"time"
)

const layout string = "Mon Jan 2 15:04:05 MST 2006"
const layout2 string = "2006-01-02 03:04:05"
const layout3 string = "2006-01-02"
const layout4 string = "Mon Jan 2 15:04:05 -0700 2006"

//	var date string = "Sun May 5 18:01:29 JST 2019"

func ParseTweetedTime(dateTimeString string) string {
	parsedTime, _ := time.Parse(layout, strings.ReplaceAll(dateTimeString, "+0000", "JST"))

	return parsedTime.Format(layout2)
}

func FormatTime(t time.Time) string {
	return t.Format(layout2)
}

func FormatDate(t time.Time) string {
	return t.Format(layout3)
}

func ParseStringDate(date string) time.Time {
	parsedDate, err := time.Parse(layout2, date)

	if err != nil {
		return time.Now()
	} else {
		return parsedDate
	}
}
