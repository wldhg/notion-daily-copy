package main

import "time"

func get_time() time.Time {
	return time.Now().AddDate(0, 0, int(dateOffset))
}
