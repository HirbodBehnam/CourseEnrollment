package course

import (
	"CourseEnrollment/pkg/util"
	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/assert"
	"math"
	"sync/atomic"
	"testing"
	"time"
)

func TestStudentIsEnrollTimeOK(t *testing.T) {
	clk := clock.NewMock()
	studentClock = clk
	tests := []struct {
		Name                  string
		CurrentTime           time.Time
		StudentEnrollmentTime time.Time
		ExpectedAllowed       bool
	}{
		{
			Name:                  "before enrollment time",
			CurrentTime:           time.Date(2022, 9, 12, 8, 0, 0, 0, time.UTC),
			StudentEnrollmentTime: time.Date(2022, 9, 12, 9, 0, 0, 0, time.UTC),
			ExpectedAllowed:       false,
		},
		{
			Name:                  "edge enrollment time (exact time)",
			CurrentTime:           time.Date(2022, 9, 12, 9, 0, 0, 0, time.UTC),
			StudentEnrollmentTime: time.Date(2022, 9, 12, 9, 0, 0, 0, time.UTC),
			ExpectedAllowed:       false, // expected
		},
		{
			Name:                  "edge enrollment time (after)",
			CurrentTime:           time.Date(2022, 9, 12, 9, 0, 0, 1e7, time.UTC),
			StudentEnrollmentTime: time.Date(2022, 9, 12, 9, 0, 0, 0, time.UTC),
			ExpectedAllowed:       true,
		},
		{
			Name:                  "after enrollment time",
			CurrentTime:           time.Date(2022, 9, 12, 9, 30, 0, 0, time.UTC),
			StudentEnrollmentTime: time.Date(2022, 9, 12, 9, 0, 0, 0, time.UTC),
			ExpectedAllowed:       true,
		},
		{
			Name:                  "edge enrollment time (end)",
			CurrentTime:           time.Date(2022, 9, 12, 10, 0, 0, 0, time.UTC),
			StudentEnrollmentTime: time.Date(2022, 9, 12, 9, 0, 0, 0, time.UTC),
			ExpectedAllowed:       false,
		},
		{
			Name:                  "end of enrollment time",
			CurrentTime:           time.Date(2022, 9, 12, 11, 0, 0, 0, time.UTC),
			StudentEnrollmentTime: time.Date(2022, 9, 12, 9, 0, 0, 0, time.UTC),
			ExpectedAllowed:       false,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			clk.Set(test.CurrentTime)
			std := Student{}
			std.EnrollmentStartTime.Store(test.StudentEnrollmentTime.UnixMilli())
			assert.Equal(t, test.ExpectedAllowed, std.IsEnrollTimeOK())
		})
	}
}

func TestStudentEnrollCourse(t *testing.T) {
	clk := clock.NewMock()
	studentClock = clk
	courses := Courses{
		courses: map[CourseID][]*Course{
			CourseID(1): {
				{
					ID:                 CourseID(1),
					GroupID:            GroupID(1),
					Units:              1,
					Capacity:           5,
					RegisteredStudents: map[StudentID]struct{}{},
					ReserveCapacity:    5,
					ReserveQueue:       util.NewQueue[StudentID](),
					ExamTime:           newAtomicTimeUnix(time.Date(2022, 9, 12, 1, 0, 0, 0, time.UTC)),
					ClassHeldTime: NewClassTime(
						[]time.Weekday{time.Wednesday},
						TimeOnly{13 * 60},
						TimeOnly{15 * 60},
					),
				},
				{
					ID:                 CourseID(1),
					GroupID:            GroupID(2),
					Units:              1,
					Capacity:           5,
					RegisteredStudents: map[StudentID]struct{}{},
					ReserveCapacity:    5,
					ReserveQueue:       util.NewQueue[StudentID](),
					ExamTime:           newAtomicTimeUnix(time.Date(2022, 9, 12, 1, 0, 0, 0, time.UTC)),
					ClassHeldTime: NewClassTime(
						[]time.Weekday{time.Wednesday},
						TimeOnly{15 * 60},
						TimeOnly{17 * 60},
					),
				},
			},
			CourseID(2): {
				{
					ID:                 CourseID(2),
					GroupID:            GroupID(1),
					Units:              3,
					Capacity:           5,
					RegisteredStudents: map[StudentID]struct{}{},
					ReserveCapacity:    5,
					ReserveQueue:       util.NewQueue[StudentID](),
					ExamTime:           newAtomicTimeUnix(time.Date(2022, 9, 12, 1, 0, 0, 0, time.UTC)),
					ClassHeldTime: NewClassTime(
						[]time.Weekday{time.Saturday},
						TimeOnly{13*60 + 30},
						TimeOnly{15 * 60},
					),
				},
			},
			CourseID(3): {
				{
					ID:                 CourseID(3),
					GroupID:            GroupID(1),
					Units:              2,
					Capacity:           5,
					RegisteredStudents: map[StudentID]struct{}{},
					ReserveCapacity:    5,
					ReserveQueue:       util.NewQueue[StudentID](),
					ExamTime:           newAtomicTimeUnix(time.Date(2022, 9, 12, 3, 0, 0, 0, time.UTC)),
					ClassHeldTime: NewClassTime(
						[]time.Weekday{time.Saturday, time.Monday},
						TimeOnly{12 * 60},
						TimeOnly{14 * 60},
					),
				},
			},
			CourseID(4): {
				{
					ID:                 CourseID(4),
					GroupID:            GroupID(1),
					Units:              2,
					Capacity:           5,
					RegisteredStudents: map[StudentID]struct{}{},
					ReserveCapacity:    5,
					ReserveQueue:       util.NewQueue[StudentID](),
					ExamTime:           newAtomicTimeUnix(time.Date(2022, 9, 12, 4, 0, 0, 0, time.UTC)),
					ClassHeldTime: NewClassTime(
						[]time.Weekday{time.Friday},
						TimeOnly{0},
						TimeOnly{1},
					),
				},
			},
			CourseID(5): {
				{
					ID:                 CourseID(5),
					GroupID:            GroupID(1),
					Units:              2,
					Capacity:           5,
					RegisteredStudents: map[StudentID]struct{}{},
					ReserveCapacity:    5,
					ReserveQueue:       util.NewQueue[StudentID](),
					ExamTime:           newAtomicTimeUnix(time.Date(2022, 9, 12, 5, 0, 0, 0, time.UTC)),
					ClassHeldTime: NewClassTime(
						[]time.Weekday{time.Friday},
						TimeOnly{5},
						TimeOnly{10},
					),
					SexLock: SexLockMaleOnly,
				},
				{
					ID:                 CourseID(5),
					GroupID:            GroupID(2),
					Units:              2,
					Capacity:           5,
					RegisteredStudents: map[StudentID]struct{}{},
					ReserveCapacity:    5,
					ReserveQueue:       util.NewQueue[StudentID](),
					ExamTime:           newAtomicTimeUnix(time.Date(2022, 9, 12, 5, 0, 0, 0, time.UTC)),
					ClassHeldTime: NewClassTime(
						[]time.Weekday{time.Friday},
						TimeOnly{5},
						TimeOnly{10},
					),
					SexLock: SexLockFemaleOnly,
				},
			},
		},
	}
	t.Run("enrollment time check", func(t *testing.T) {
		std := Student{RegisteredCourses: map[CourseID]GroupID{}}
		std.MaxUnits = math.MaxUint8
		clk.Set(time.Unix(10, 0))
		std.EnrollmentStartTime.Store(time.Unix(15, 0).UnixMilli())
		assert.ErrorIs(t, std.EnrollCourse(&courses, CourseID(1), GroupID(1)), NotEnrollmentTimeErr)
		std.EnrollmentStartTime.Store(time.Unix(5, 0).UnixMilli())
		assert.NoError(t, std.EnrollCourse(&courses, CourseID(1), GroupID(1)))
	})
	// We allow all other register times without setting them
	clk.Set(time.Unix(1, 0))
	// A non-existent course must throw an error
	t.Run("non existent course", func(t *testing.T) {
		std := Student{RegisteredCourses: map[CourseID]GroupID{}}
		assert.ErrorIs(t, std.EnrollCourse(&courses, CourseID(1), GroupID(10)), NotExistsErr)
		assert.ErrorIs(t, std.EnrollCourse(&courses, CourseID(100), GroupID(1)), NotExistsErr)
	})
	// Unlocked "sex lock" lessons are checked before. We just check the locked ones
	t.Run("sex lock", func(t *testing.T) {
		std := Student{
			StudentSex:        SexFemale,
			MaxUnits:          math.MaxUint8,
			RegisteredCourses: map[CourseID]GroupID{},
		}
		assert.ErrorIs(t, std.EnrollCourse(&courses, CourseID(5), GroupID(1)), SexLockErr)
		assert.NoError(t, std.EnrollCourse(&courses, CourseID(5), GroupID(2)))
		std.StudentSex = SexMale
		std.RegisteredCourses = map[CourseID]GroupID{}
		assert.ErrorIs(t, std.EnrollCourse(&courses, CourseID(5), GroupID(2)), SexLockErr)
		assert.NoError(t, std.EnrollCourse(&courses, CourseID(5), GroupID(1)))
	})
	// Check max registered courses
	t.Run("unit limit", func(t *testing.T) {
		std := Student{
			MaxUnits:          3,
			RegisteredCourses: map[CourseID]GroupID{},
		}
		assert.NoError(t, std.EnrollCourse(&courses, CourseID(2), GroupID(1)))
		assert.Equal(t, uint8(3), std.RegisteredUnits)
		assert.ErrorIs(t, std.EnrollCourse(&courses, CourseID(1), GroupID(1)), UnitLimitReachedErr)
	})
	// Do not allow already registered courses
	t.Run("already registered", func(t *testing.T) {
		std := Student{
			MaxUnits:          math.MaxUint8,
			RegisteredCourses: map[CourseID]GroupID{},
		}
		assert.NoError(t, std.EnrollCourse(&courses, CourseID(1), GroupID(1)))
		assert.ErrorIs(t, std.EnrollCourse(&courses, CourseID(1), GroupID(1)), AlreadyRegisteredErr)
		assert.ErrorIs(t, std.EnrollCourse(&courses, CourseID(1), GroupID(2)), AlreadyRegisteredErr)
	})
	// Exam times must not overlap
	t.Run("exam time", func(t *testing.T) {
		std := Student{
			MaxUnits: math.MaxUint8,
			RegisteredCourses: map[CourseID]GroupID{
				CourseID(1): GroupID(1),
			},
		}
		assert.ErrorIs(t, std.EnrollCourse(&courses, CourseID(2), GroupID(1)), ExamConflictErr{CourseID(1), GroupID(1)})
	})
	// Class times must not overlap as well
	t.Run("class time", func(t *testing.T) {
		std := Student{
			MaxUnits: math.MaxUint8,
			RegisteredCourses: map[CourseID]GroupID{
				CourseID(2): GroupID(1),
			},
		}
		assert.ErrorIs(t, std.EnrollCourse(&courses, CourseID(3), GroupID(1)), ClassTimeConflictErr{CourseID(2), GroupID(1)})
	})
	// Must never happen. We just test for panic
	t.Run("inconsistent state panic", func(t *testing.T) {
		std := Student{
			MaxUnits: math.MaxUint8,
			RegisteredCourses: map[CourseID]GroupID{
				CourseID(1000): GroupID(100),
			},
		}
		assert.PanicsWithValue(t, "inconsistent user state: course 1000 group 100 is registered but not found",
			func() {
				_ = std.EnrollCourse(&courses, CourseID(3), GroupID(1))
			})
	})
	// Capacity error
	t.Run("capacity", func(t *testing.T) {
		courses.courses[CourseID(1)][0].Capacity = 1
		courses.courses[CourseID(1)][0].RegisteredStudents = map[StudentID]struct{}{StudentID(100): {}}
		courses.courses[CourseID(1)][0].ReserveCapacity = 1
		courses.courses[CourseID(1)][0].ReserveQueue = util.NewQueue[StudentID]()
		courses.courses[CourseID(1)][0].ReserveQueue.Enqueue(StudentID(101))
		std := Student{
			MaxUnits:          math.MaxUint8,
			RegisteredCourses: map[CourseID]GroupID{},
		}
		assert.ErrorIs(t, std.EnrollCourse(&courses, CourseID(1), GroupID(1)), NoCapacityLeftErr)
	})
}

//goland:noinspection GoVetCopyLock
func newAtomicTimeUnix(t time.Time) atomic.Int64 {
	result := atomic.Int64{}
	result.Store(t.Unix())
	return result
}
