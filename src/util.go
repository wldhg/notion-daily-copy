package main

import "time"

func get_time() time.Time {
	return time.Now().AddDate(0, 0, int(dateOffset))
}

func get_0am_time() time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).AddDate(0, 0, int(dateOffset))
}
