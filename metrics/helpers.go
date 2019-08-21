package metrics

import "time"

// Timer returns a function which returns an int64 that represents the difference
// in time between Time being called and the returning function being called
// in milliseconds.
func Timer() func() int64 {
	start := MSTime()
	return func() int64 {
		return Duration(start)
	}
}

// MSTime returns the current time in milliseconds.
func MSTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// Duration provided a start time in milliseconds, returns the difference /
// duration between that time and the time of Duration being called as an int64.
// The returned calue is the duration in time in milliseconds.
func Duration(start int64) int64 {
	return MSTime() - start
}
