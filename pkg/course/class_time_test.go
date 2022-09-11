package course

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func FuzzClassTime(f *testing.F) {
	f.Add(uint8(123), uint16(24234), uint16(1249))
	f.Fuzz(func(t *testing.T, days uint8, startTime, endTime uint16) {
		startTime %= TimeOnlyMax
		endTime %= TimeOnlyMax
		var classTime ClassTime
		// Create the weekdays
		weekdays := make([]time.Weekday, 0, 7)
		weekdaysMap := make(map[time.Weekday]struct{}, 7)
		for i := 0; i < 7; i++ {
			if (days>>i)&1 == 1 {
				weekdays = append(weekdays, time.Weekday(i))
				weekdaysMap[time.Weekday(i)] = struct{}{}
			}
		}
		// Set the data
		classTime.Set(weekdays, NewTimeOnly(startTime), NewTimeOnly(endTime))
		// Get the data back and compare
		weekdaysGot, startTimeGot, endTimeGot := classTime.Get()
		weekdaysGotMap := make(map[time.Weekday]struct{}, len(weekdaysGot))
		for _, day := range weekdaysGot {
			weekdaysGotMap[day] = struct{}{}
		}
		assert.Equal(t, startTime, startTimeGot.t)
		assert.Equal(t, endTime, endTimeGot.t)
		assert.Equal(t, weekdaysMap, weekdaysGotMap)
	})
}

func TestClassTimeIntersects(t *testing.T) {
	tests := []struct {
		Name               string
		Weekdays           [2][]time.Weekday
		StartTime          [2]TimeOnly
		EndTime            [2]TimeOnly
		ExpectedIntersects bool
	}{
		{
			Name: "not intersecting time and day",
			Weekdays: [2][]time.Weekday{
				{
					time.Saturday,
				},
				{
					time.Wednesday,
				},
			},
			StartTime: [2]TimeOnly{
				NewTimeOnly(10*60 + 0),
				NewTimeOnly(15*60 + 0),
			},
			EndTime: [2]TimeOnly{
				NewTimeOnly(12*60 + 0),
				NewTimeOnly(18*60 + 0),
			},
			ExpectedIntersects: false,
		},
		{
			Name: "not intersecting time but day",
			Weekdays: [2][]time.Weekday{
				{
					time.Saturday,
				},
				{
					time.Saturday,
				},
			},
			StartTime: [2]TimeOnly{
				NewTimeOnly(10*60 + 0),
				NewTimeOnly(15*60 + 0),
			},
			EndTime: [2]TimeOnly{
				NewTimeOnly(12*60 + 0),
				NewTimeOnly(18*60 + 0),
			},
			ExpectedIntersects: false,
		},
		{
			Name: "not intersecting day but time",
			Weekdays: [2][]time.Weekday{
				{
					time.Saturday,
					time.Sunday,
				},
				{
					time.Wednesday,
					time.Thursday,
				},
			},
			StartTime: [2]TimeOnly{
				NewTimeOnly(10*60 + 0),
				NewTimeOnly(10*60 + 0),
			},
			EndTime: [2]TimeOnly{
				NewTimeOnly(12*60 + 0),
				NewTimeOnly(12*60 + 0),
			},
			ExpectedIntersects: false,
		},
		{
			Name: "intersecting complete overlap",
			Weekdays: [2][]time.Weekday{
				{
					time.Saturday,
					time.Sunday,
				},
				{
					time.Saturday,
					time.Sunday,
				},
			},
			StartTime: [2]TimeOnly{
				NewTimeOnly(10*60 + 0),
				NewTimeOnly(10*60 + 0),
			},
			EndTime: [2]TimeOnly{
				NewTimeOnly(12*60 + 0),
				NewTimeOnly(12*60 + 0),
			},
			ExpectedIntersects: true,
		},
		{
			Name: "intersecting partial time overlap",
			Weekdays: [2][]time.Weekday{
				{
					time.Saturday,
					time.Sunday,
				},
				{
					time.Saturday,
					time.Sunday,
				},
			},
			StartTime: [2]TimeOnly{
				NewTimeOnly(10*60 + 0),
				NewTimeOnly(11*60 + 0),
			},
			EndTime: [2]TimeOnly{
				NewTimeOnly(12*60 + 0),
				NewTimeOnly(13*60 + 0),
			},
			ExpectedIntersects: true,
		},
		{
			Name: "intersecting partial time and day overlap",
			Weekdays: [2][]time.Weekday{
				{
					time.Saturday,
					time.Sunday,
				},
				{
					time.Wednesday,
					time.Sunday,
				},
			},
			StartTime: [2]TimeOnly{
				NewTimeOnly(10*60 + 0),
				NewTimeOnly(11*60 + 0),
			},
			EndTime: [2]TimeOnly{
				NewTimeOnly(12*60 + 0),
				NewTimeOnly(13*60 + 0),
			},
			ExpectedIntersects: true,
		},
		{
			Name: "not intersecting next to each other",
			Weekdays: [2][]time.Weekday{
				{
					time.Saturday,
				},
				{
					time.Saturday,
				},
			},
			StartTime: [2]TimeOnly{
				NewTimeOnly(10*60 + 0),
				NewTimeOnly(12*60 + 0),
			},
			EndTime: [2]TimeOnly{
				NewTimeOnly(12*60 + 0),
				NewTimeOnly(13*60 + 0),
			},
			ExpectedIntersects: false,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var first, second ClassTime
			first.Set(test.Weekdays[0], test.StartTime[0], test.EndTime[0])
			second.Set(test.Weekdays[1], test.StartTime[1], test.EndTime[1])
			assert.Equal(t, test.ExpectedIntersects, first.Intersects(&second))
		})
	}
}

func BenchmarkClassTimeIntersects(b *testing.B) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var first, second ClassTime
	for i := 0; i < b.N; i++ {
		first.data.Store(rng.Uint32())
		second.data.Store(rng.Uint32())
		first.Intersects(&second)
	}
}
