package course

import (
	"sync/atomic"
	"time"
)

// TimeOnlyMax is maximum internal value which can be stored in TimeOnly
const TimeOnlyMax = 24 * 60

// On which bit the weekday data starts
const classTimeWeekdayStartBit = 22

// ClassTime only holds the class start and end hours and days which
// the class is being held on.
//
// Internally, it's an atomic uint32. The start date and end date are stored as
// start_date + end_date * TimeOnlyMax. This means that the max bits used is 21. So we still have
// 11 bits left.
//
// Next 7 bits are used to determine which days the class is being held. 1 in bit means is
// held while 0 means not. Days order is the same as time.Weekday
type ClassTime struct {
	data atomic.Uint32
}

// NewClassTime is basically instantiation + set
//
//goland:noinspection GoVetCopyLock
func NewClassTime(days []time.Weekday, startTime, endTime TimeOnly) ClassTime {
	result := ClassTime{}
	result.Set(days, startTime, endTime)
	return result
}

// Get will extract the data from ClassTime
func (t *ClassTime) Get() (days []time.Weekday, startTime, endTime TimeOnly) {
	raw := t.data.Load()
	// Extract start time
	dateOnlyBits := raw & ((1 << classTimeWeekdayStartBit) - 1)
	startTime = NewTimeOnly(uint16(dateOnlyBits % TimeOnlyMax))
	endTime = NewTimeOnly(uint16((dateOnlyBits / TimeOnlyMax) % TimeOnlyMax))
	// Most classes have two days in week
	days = make([]time.Weekday, 0, 2)
	for i := classTimeWeekdayStartBit; i < classTimeWeekdayStartBit+7; i++ {
		if (raw>>i)&1 == 1 { // if that bit is set...
			days = append(days, time.Weekday(i-classTimeWeekdayStartBit))
		}
	}
	return
}

// Set will set the data of ClassTime based on its inputs
func (t *ClassTime) Set(days []time.Weekday, startTime, endTime TimeOnly) {
	var raw uint32
	// Store the time
	raw = uint32(startTime.t)
	raw += TimeOnlyMax * uint32(endTime.t)
	// Store dates
	for _, day := range days {
		raw |= 1 << (uint32(day) + classTimeWeekdayStartBit)
	}
	// Store
	t.data.Store(raw)
}

// Intersects checks if two ClassTimes intersects.
// This is useful to forbid the student to pick another class which has intersection problems
// with their other classes
func (t *ClassTime) Intersects(other *ClassTime) bool {
	// Get the info
	weekdays1, startTime1, endTime1 := t.Get()
	weekdays2, startTime2, endTime2 := other.Get()
	// I think it's better to just check time then date
	if startTime1.t < endTime2.t && startTime2.t < endTime1.t {
		// Here, we check if dates also intersect
		for _, weekday1 := range weekdays1 {
			for _, weekday2 := range weekdays2 {
				if weekday1 == weekday2 {
					return true
				}
			}
		}
	}
	return false
}

// TimeOnly only and only holds a time between 00:00 and 24:00.
//
// Internal representation is basically a number between 0 and 1440 which is calculated with
// minute + hour * 60
type TimeOnly struct {
	t uint16
}

// NewTimeOnly creates a TimeOnly with the initial as it's internal value
func NewTimeOnly(initial uint16) TimeOnly {
	return TimeOnly{initial}
}

func NewTimeOnlyFromTime(date time.Time) TimeOnly {
	return TimeOnly{uint16(date.Hour()*60 + date.Minute())}
}
