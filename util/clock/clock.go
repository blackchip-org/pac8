package clock

import "time"

var mockNow time.Time

func Now() time.Time {
	if !mockNow.IsZero() {
		return mockNow
	}
	return time.Now()
}

func SetNow(t time.Time) {
	mockNow = t
}
