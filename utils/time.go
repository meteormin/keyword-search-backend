package utils

import "time"

const DefaultDateLayout = "2006-01-02 15:04:05"

func TimeIn(t time.Time, tz string) time.Time {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		panic(err)
	}

	return t.In(loc)
}
