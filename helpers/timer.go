package helpers

import (
	"strconv"
	"time"
)

type Timer struct {
	current time.Time
	last    time.Time
	records []time.Time
}

func CreateTimer() Timer {
	var timer Timer
	timer.records = make([]time.Time, 0)
	timer.Reset()
	return timer
}
func (timer *Timer) PrintDiff() {
	println(strconv.FormatFloat(time.Since(timer.current).Seconds(), 'f', 2, 32))
}
func (timer *Timer) Reset() {
	timer.current = time.Now()
	timer.last = timer.current
}
func (timer *Timer) RecordTime() {
	timer.records = append(timer.records, time.Now())
}
