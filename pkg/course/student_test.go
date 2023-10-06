package course

import (
	"CourseEnrollment/pkg/util"
	"context"
	"errors"
	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestExamTimesIntersect(t *testing.T) {
	tests := []struct {
		Name           string
		A              int64
		B              int64
		ExpectedResult bool
	}{
		{
			Name:           "not equal",
			A:              1,
			B:              2,
			ExpectedResult: false,
		},
		{
			Name:           "equal",
			A:              2,
			B:              2,
			ExpectedResult: true,
		},
		{
			Name:           "zero",
			A:              0,
			B:              2,
			ExpectedResult: false,
		},
		{
			Name:           "zero",
			A:              2,
			B:              0,
			ExpectedResult: false,
		},
		{
			Name:           "both zero",
			A:              0,
			B:              0,
			ExpectedResult: false,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t, test.ExpectedResult, examTimesIntersect(test.A, test.B))
		})
	}
}

func TestStudentIsEnrollTime(t *testing.T) {
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
			std.EnrollmentStartTime = test.StudentEnrollmentTime.UnixMilli()
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
		std.EnrollmentStartTime = time.Unix(15, 0).UnixMilli()
		assert.ErrorIs(t, std.EnrollCourse(context.Background(), &courses, CourseID(1), GroupID(1), noOpBatcher{}), NotEnrollmentTimeErr)
		std.EnrollmentStartTime = time.Unix(5, 0).UnixMilli()
		assert.NoError(t, std.EnrollCourse(context.Background(), &courses, CourseID(1), GroupID(1), noOpBatcher{}))
	})
	// We allow all other register times without setting them
	clk.Set(time.Unix(1, 0))
	// A non-existent course must throw an error
	t.Run("non existent course", func(t *testing.T) {
		std := Student{RegisteredCourses: map[CourseID]GroupID{}}
		assert.ErrorIs(t, std.EnrollCourse(context.Background(), &courses, CourseID(1), GroupID(10), noOpBatcher{}), NotExistsErr)
		assert.ErrorIs(t, std.EnrollCourse(context.Background(), &courses, CourseID(100), GroupID(1), noOpBatcher{}), NotExistsErr)
	})
	// Unlocked "sex lock" lessons are checked before. We just check the locked ones
	t.Run("sex lock", func(t *testing.T) {
		std := Student{
			StudentSex:        SexFemale,
			MaxUnits:          math.MaxUint8,
			RegisteredCourses: map[CourseID]GroupID{},
		}
		assert.ErrorIs(t, std.EnrollCourse(context.Background(), &courses, CourseID(5), GroupID(1), noOpBatcher{}), SexLockErr)
		assert.NoError(t, std.EnrollCourse(context.Background(), &courses, CourseID(5), GroupID(2), noOpBatcher{}))
		std.StudentSex = SexMale
		std.RegisteredCourses = map[CourseID]GroupID{}
		assert.ErrorIs(t, std.EnrollCourse(context.Background(), &courses, CourseID(5), GroupID(2), noOpBatcher{}), SexLockErr)
		assert.NoError(t, std.EnrollCourse(context.Background(), &courses, CourseID(5), GroupID(1), noOpBatcher{}))
	})
	// Check max registered courses
	t.Run("unit limit", func(t *testing.T) {
		std := Student{
			MaxUnits:          3,
			RegisteredCourses: map[CourseID]GroupID{},
		}
		assert.NoError(t, std.EnrollCourse(context.Background(), &courses, CourseID(2), GroupID(1), noOpBatcher{}))
		assert.Equal(t, uint8(3), std.RegisteredUnits)
		assert.ErrorIs(t, std.EnrollCourse(context.Background(), &courses, CourseID(1), GroupID(1), noOpBatcher{}), UnitLimitReachedErr)
	})
	// Do not allow already registered courses
	t.Run("already registered", func(t *testing.T) {
		std := Student{
			MaxUnits:          math.MaxUint8,
			RegisteredCourses: map[CourseID]GroupID{},
		}
		assert.NoError(t, std.EnrollCourse(context.Background(), &courses, CourseID(1), GroupID(1), noOpBatcher{}))
		assert.ErrorIs(t, std.EnrollCourse(context.Background(), &courses, CourseID(1), GroupID(1), noOpBatcher{}), AlreadyRegisteredErr)
		assert.ErrorIs(t, std.EnrollCourse(context.Background(), &courses, CourseID(1), GroupID(2), noOpBatcher{}), AlreadyRegisteredErr)
	})
	// Exam times must not overlap
	t.Run("exam time", func(t *testing.T) {
		std := Student{
			MaxUnits: math.MaxUint8,
			RegisteredCourses: map[CourseID]GroupID{
				CourseID(1): GroupID(1),
			},
		}
		assert.ErrorIs(t, std.EnrollCourse(context.Background(), &courses, CourseID(2), GroupID(1), noOpBatcher{}), ExamConflictErr{CourseID(1), GroupID(1)})
	})
	// Class times must not overlap as well
	t.Run("class time", func(t *testing.T) {
		std := Student{
			MaxUnits: math.MaxUint8,
			RegisteredCourses: map[CourseID]GroupID{
				CourseID(2): GroupID(1),
			},
		}
		assert.ErrorIs(t, std.EnrollCourse(context.Background(), &courses, CourseID(3), GroupID(1), noOpBatcher{}), ClassTimeConflictErr{CourseID(2), GroupID(1)})
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
				_ = std.EnrollCourse(context.Background(), &courses, CourseID(3), GroupID(1), noOpBatcher{})
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
		assert.ErrorIs(t, std.EnrollCourse(context.Background(), &courses, CourseID(1), GroupID(1), noOpBatcher{}), NoCapacityLeftErr)
	})
}

func TestStudentDisenrollCourse(t *testing.T) {
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
		},
	}
	t.Run("enrollment time check", func(t *testing.T) {
		std := Student{
			ID: StudentID(1),
			RegisteredCourses: map[CourseID]GroupID{
				CourseID(1): GroupID(1),
			},
			RemainingActions: 10,
		}
		courses.courses[CourseID(1)][0].RegisteredStudents[StudentID(1)] = struct{}{}
		clk.Set(time.Unix(10, 0))
		std.EnrollmentStartTime = time.Unix(15, 0).UnixMilli()
		assert.ErrorIs(t, std.DisenrollCourse(context.Background(), &courses, CourseID(1), noOpBatcher{}), NotEnrollmentTimeErr)
		std.EnrollmentStartTime = time.Unix(5, 0).UnixMilli()
		assert.NoError(t, std.DisenrollCourse(context.Background(), &courses, CourseID(1), noOpBatcher{}))
		assert.Len(t, courses.courses[CourseID(1)][0].RegisteredStudents, 0)
	})
	// We allow all other register times without setting them
	clk.Set(time.Unix(1, 0))
	// Remaining actions must not be zero and must be reduced each time
	t.Run("remaining actions", func(t *testing.T) {
		std := Student{
			ID: StudentID(1),
			RegisteredCourses: map[CourseID]GroupID{
				CourseID(1): GroupID(1),
				CourseID(2): GroupID(1),
			},
			RemainingActions: 1,
		}
		courses.courses[CourseID(1)][0].RegisteredStudents[StudentID(1)] = struct{}{}
		courses.courses[CourseID(2)][0].RegisteredStudents[StudentID(1)] = struct{}{}
		assert.NoError(t, std.DisenrollCourse(context.Background(), &courses, CourseID(1), noOpBatcher{}))
		assert.Equal(t, uint8(0), std.RemainingActions)
		assert.ErrorIs(t, std.DisenrollCourse(context.Background(), &courses, CourseID(2), noOpBatcher{}), NoRemainingActionsErr)
	})
	// A course which user is not registered in or does not exist
	t.Run("invalid course", func(t *testing.T) {
		std := Student{
			ID: StudentID(1),
			RegisteredCourses: map[CourseID]GroupID{
				CourseID(100): GroupID(1),
				CourseID(2):   GroupID(2),
			},
			RemainingActions: 10,
		}
		assert.ErrorIs(t, std.DisenrollCourse(context.Background(), &courses, CourseID(10), noOpBatcher{}), NotExistsErr)
		assert.PanicsWithValue(t, "invalid registered lesson 2-2 for user 1", func() {
			_ = std.DisenrollCourse(context.Background(), &courses, CourseID(2), noOpBatcher{})
		})
		assert.PanicsWithValue(t, "invalid registered lesson 100-1 for user 1", func() {
			_ = std.DisenrollCourse(context.Background(), &courses, CourseID(100), noOpBatcher{})
		})
	})
	// Inconsistency of course and registered map in student struct
	t.Run("registered inconsistency", func(t *testing.T) {
		std := Student{
			ID: StudentID(1),
			RegisteredCourses: map[CourseID]GroupID{
				CourseID(1): GroupID(1),
			},
			RemainingActions: 10,
		}
		courses.courses[CourseID(1)][0].RegisteredStudents = map[StudentID]struct{}{}
		assert.PanicsWithValue(t, "user 1 has lesson 1-1 in their registered courses but lesson map does not have this user", func() {
			_ = std.DisenrollCourse(context.Background(), &courses, CourseID(1), noOpBatcher{})
		})
	})
	// Just remove a student from course
	t.Run("general test", func(t *testing.T) {
		const addedRegisteredUnits = 2
		const initialRemainingActions = 10
		std := Student{
			ID: StudentID(1),
			RegisteredCourses: map[CourseID]GroupID{
				CourseID(1): GroupID(1),
			},
			RemainingActions: initialRemainingActions,
			RegisteredUnits:  courses.courses[CourseID(1)][0].Units + addedRegisteredUnits,
		}
		courses.courses[CourseID(1)][0].RegisteredStudents = map[StudentID]struct{}{
			StudentID(1): {},
		}
		assert.NoError(t, std.DisenrollCourse(context.Background(), &courses, CourseID(1), noOpBatcher{}))
		assert.Equal(t, uint8(addedRegisteredUnits), std.RegisteredUnits)
		assert.Equal(t, uint8(initialRemainingActions-1), std.RemainingActions)
		assert.Len(t, std.RegisteredCourses, 0)
		assert.Len(t, courses.courses[CourseID(1)][0].RegisteredStudents, 0)
	})
}

func TestStudentChangeGroup(t *testing.T) {
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
						[]time.Weekday{time.Friday},
						TimeOnly{13*60 + 30},
						TimeOnly{15 * 60},
					),
				},
				{
					ID:                 CourseID(2),
					GroupID:            GroupID(2),
					Units:              3,
					Capacity:           5,
					RegisteredStudents: map[StudentID]struct{}{},
					ReserveCapacity:    5,
					ReserveQueue:       util.NewQueue[StudentID](),
					ExamTime:           newAtomicTimeUnix(time.Date(2022, 9, 12, 5, 0, 0, 0, time.UTC)),
					ClassHeldTime: NewClassTime(
						[]time.Weekday{time.Saturday, time.Monday},
						TimeOnly{12 * 60},
						TimeOnly{14 * 60},
					),
				},
				{
					ID:                 CourseID(2),
					GroupID:            GroupID(3),
					Units:              3,
					Capacity:           5,
					RegisteredStudents: map[StudentID]struct{}{},
					ReserveCapacity:    5,
					ReserveQueue:       util.NewQueue[StudentID](),
					ExamTime:           newAtomicTimeUnix(time.Date(2022, 9, 12, 1, 0, 0, 0, time.UTC)),
					ClassHeldTime: NewClassTime(
						[]time.Weekday{time.Thursday},
						TimeOnly{13*60 + 30},
						TimeOnly{15 * 60},
					),
				},
			},
			CourseID(3): {
				{
					ID:                 CourseID(3),
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
			CourseID(4): {
				{
					ID:                 CourseID(4),
					GroupID:            GroupID(1),
					Units:              3,
					Capacity:           1,
					RegisteredStudents: map[StudentID]struct{}{},
					ReserveCapacity:    0,
					ReserveQueue:       util.NewQueue[StudentID](),
					ExamTime:           newAtomicTimeUnix(time.Date(2022, 9, 12, 1, 0, 0, 0, time.UTC)),
					ClassHeldTime: NewClassTime(
						[]time.Weekday{time.Saturday},
						TimeOnly{13*60 + 30},
						TimeOnly{15 * 60},
					),
				},
				{
					ID:                 CourseID(4),
					GroupID:            GroupID(2),
					Units:              3,
					Capacity:           1,
					RegisteredStudents: map[StudentID]struct{}{},
					ReserveCapacity:    0,
					ReserveQueue:       util.NewQueue[StudentID](),
					ExamTime:           newAtomicTimeUnix(time.Date(2022, 9, 12, 1, 0, 0, 0, time.UTC)),
					ClassHeldTime: NewClassTime(
						[]time.Weekday{time.Saturday},
						TimeOnly{13*60 + 30},
						TimeOnly{15 * 60},
					),
				},
			},
		},
	}
	t.Run("enrollment time check", func(t *testing.T) {
		std := Student{
			ID: StudentID(1),
			RegisteredCourses: map[CourseID]GroupID{
				CourseID(1): GroupID(1),
			},
			RemainingActions: 10,
		}
		courses.courses[CourseID(1)][0].RegisteredStudents[StudentID(1)] = struct{}{}
		clk.Set(time.Unix(10, 0))
		std.EnrollmentStartTime = time.Unix(15, 0).UnixMilli()
		assert.ErrorIs(t, std.ChangeGroup(context.Background(), &courses, CourseID(1), GroupID(2), noOpBatcher{}), NotEnrollmentTimeErr)
		std.EnrollmentStartTime = time.Unix(5, 0).UnixMilli()
		assert.NoError(t, std.ChangeGroup(context.Background(), &courses, CourseID(1), GroupID(2), noOpBatcher{}))
		assert.Len(t, courses.courses[CourseID(1)][0].RegisteredStudents, 0)
		assert.Len(t, courses.courses[CourseID(1)][1].RegisteredStudents, 1)
	})
	// We allow all other register times without setting them
	clk.Set(time.Unix(1, 0))
	// Remaining actions must not be zero and must be reduced each time
	t.Run("remaining actions", func(t *testing.T) {
		std := Student{
			ID: StudentID(1),
			RegisteredCourses: map[CourseID]GroupID{
				CourseID(1): GroupID(1),
			},
			RemainingActions: 1,
		}
		courses.courses[CourseID(1)][0].RegisteredStudents = map[StudentID]struct{}{
			StudentID(1): {},
		}
		assert.NoError(t, std.ChangeGroup(context.Background(), &courses, CourseID(1), GroupID(2), noOpBatcher{}))
		assert.Equal(t, uint8(0), std.RemainingActions)
		assert.ErrorIs(t, std.ChangeGroup(context.Background(), &courses, CourseID(1), GroupID(2), noOpBatcher{}), NoRemainingActionsErr)
	})
	// A course which user is not registered in or does not exist
	t.Run("invalid course", func(t *testing.T) {
		std := Student{
			ID: StudentID(1),
			RegisteredCourses: map[CourseID]GroupID{
				CourseID(100): GroupID(1),
			},
			RemainingActions: 10,
		}
		assert.ErrorIs(t, std.ChangeGroup(context.Background(), &courses, CourseID(1), GroupID(2), noOpBatcher{}), NotExistsErr)
		std.RegisteredCourses[CourseID(1)] = GroupID(1)
		assert.ErrorIs(t, std.ChangeGroup(context.Background(), &courses, CourseID(1), GroupID(3), noOpBatcher{}), NotExistsErr)
		assert.ErrorIs(t, std.ChangeGroup(context.Background(), &courses, CourseID(1), GroupID(1), noOpBatcher{}), PlayedYourselfErr)
		assert.PanicsWithValue(t, "invalid registered lesson 100-1 for user 1", func() {
			_ = std.ChangeGroup(context.Background(), &courses, CourseID(100), GroupID(2), noOpBatcher{})
		})
	})
	// Time conflict checks
	t.Run("time conflicts", func(t *testing.T) {
		std := Student{
			ID: StudentID(1),
			RegisteredCourses: map[CourseID]GroupID{
				CourseID(2): GroupID(1),
				CourseID(3): GroupID(1),
			},
			RemainingActions: 10,
		}
		courses.courses[CourseID(2)][0].RegisteredStudents = map[StudentID]struct{}{
			StudentID(1): {},
		}
		assert.ErrorIs(t, std.ChangeGroup(context.Background(), &courses, CourseID(2), GroupID(2), noOpBatcher{}), ClassTimeConflictErr{CourseID: 3, GroupID: 1})
		assert.ErrorIs(t, std.ChangeGroup(context.Background(), &courses, CourseID(2), GroupID(3), noOpBatcher{}), ExamConflictErr{CourseID: 3, GroupID: 1})
	})
	// Capacity test
	t.Run("capacity", func(t *testing.T) {
		std := Student{
			ID: StudentID(1),
			RegisteredCourses: map[CourseID]GroupID{
				CourseID(4): GroupID(1),
			},
			RemainingActions: 10,
		}
		courses.courses[CourseID(4)][0].RegisteredStudents[StudentID(1)] = struct{}{}
		courses.courses[CourseID(4)][1].RegisteredStudents[StudentID(2)] = struct{}{}
		assert.ErrorIs(t, std.ChangeGroup(context.Background(), &courses, CourseID(4), GroupID(2), noOpBatcher{}), NoCapacityLeftErr)
	})
	// Normal test
	t.Run("normal test", func(t *testing.T) {
		std := Student{
			ID: StudentID(1),
			RegisteredCourses: map[CourseID]GroupID{
				CourseID(1): GroupID(1),
			},
			RemainingActions: 1,
		}
		courses.courses[CourseID(1)][0].RegisteredStudents = map[StudentID]struct{}{
			StudentID(1): {},
		}
		courses.courses[CourseID(1)][1].RegisteredStudents = map[StudentID]struct{}{}
		assert.NoError(t, std.ChangeGroup(context.Background(), &courses, CourseID(1), GroupID(2), noOpBatcher{}))
		assert.Len(t, courses.courses[CourseID(1)][0].RegisteredStudents, 0)
		assert.Equal(t, map[StudentID]struct{}{
			StudentID(1): {},
		}, courses.courses[CourseID(1)][1].RegisteredStudents)
		assert.Equal(t, uint8(0), std.RemainingActions)
		assert.Equal(t, map[CourseID]GroupID{
			CourseID(1): GroupID(2),
		}, std.RegisteredCourses)
	})
}

func TestStudentChangeGroupLock1(t *testing.T) {
	clk := clock.NewMock()
	clk.Set(time.Unix(1, 0))
	studentClock = clk
	const numberOfGroups = 20
	for numberOfRotations := 1; numberOfRotations <= numberOfGroups; numberOfRotations++ {
		for numberOfSteps := 1; numberOfSteps < numberOfRotations; numberOfSteps++ {
			// Create the courses
			courses := Courses{courses: map[CourseID][]*Course{
				CourseID(1): make([]*Course, numberOfGroups),
			}}
			students := make([]Student, numberOfGroups)
			for i := 0; i < numberOfGroups; i++ {
				courses.courses[CourseID(1)][i] = &Course{
					ID:                 CourseID(1),
					GroupID:            GroupID(i),
					Units:              1,
					Capacity:           1,
					RegisteredStudents: map[StudentID]struct{}{StudentID(i): {}},
					ReserveCapacity:    1,
					ReserveQueue:       util.NewQueue[StudentID](),
				}
				students[i] = Student{
					ID:                StudentID(i),
					RemainingActions:  uint8(numberOfRotations),
					MaxUnits:          1,
					RegisteredUnits:   1,
					RegisteredCourses: map[CourseID]GroupID{CourseID(1): GroupID(i)},
				}
			}
			// Spawn goroutines and spin students
			wg := new(sync.WaitGroup)
			wg.Add(numberOfGroups)
			for i := 0; i < numberOfGroups; i++ {
				go func(index int) {
				OuterLoop:
					for i := 0; i < numberOfRotations; i++ {
						nextGroupID := (int(students[index].RegisteredCourses[CourseID(1)]) + numberOfSteps) % numberOfGroups
						for {
							err := students[index].ChangeGroup(context.Background(), &courses, CourseID(1), GroupID(nextGroupID), noOpBatcher{})
							// If there is no capacity, we just try again
							if errors.Is(err, NoCapacityLeftErr) {
								runtime.Gosched() // Let other threads do stuff
								continue
							}
							// Otherwise we check for errors
							if err != nil {
								t.Errorf("unexcpeceted error when rotating: %s", err)
								break OuterLoop
							}
							// We are good so continue
							break
						}
					}
					wg.Done()
				}(i)
			}
			wg.Wait()
			// Check the position
			positionOffset := (numberOfRotations * numberOfSteps) % numberOfGroups
			for i := range students {
				assert.Equalf(t, map[CourseID]GroupID{
					CourseID(1): GroupID((i + positionOffset) % numberOfGroups),
				}, students[i].RegisteredCourses, "error on %d %d student %d",
					numberOfRotations, numberOfSteps, students[i].ID)
			}
		}
	}
}

func TestStudentChangeGroupLock2(t *testing.T) {
	clk := clock.NewMock()
	clk.Set(time.Unix(1, 0))
	studentClock = clk
	const numberOfGroups = 20
	for numberOfRotations := 1; numberOfRotations <= numberOfGroups; numberOfRotations++ {
		const numberOfSteps = 1
		// Create the courses
		courses := Courses{courses: map[CourseID][]*Course{
			CourseID(1): make([]*Course, numberOfGroups),
		}}
		students := make([]Student, numberOfGroups)
		for i := 0; i < numberOfGroups; i++ {
			courses.courses[CourseID(1)][i] = &Course{
				ID:                 CourseID(1),
				GroupID:            GroupID(i),
				Capacity:           1,
				Units:              1,
				RegisteredStudents: map[StudentID]struct{}{StudentID(i): {}},
				ReserveCapacity:    0,
				ReserveQueue:       util.NewQueue[StudentID](),
			}
			if i%5 == 0 {
				courses.courses[CourseID(1)][i].Capacity++
			}
			students[i] = Student{
				ID:                StudentID(i),
				RemainingActions:  uint8(numberOfRotations),
				MaxUnits:          1,
				RegisteredUnits:   1,
				RegisteredCourses: map[CourseID]GroupID{CourseID(1): GroupID(i)},
			}
		}
		// Spawn goroutines and spin students
		barriers := make([]sync.WaitGroup, numberOfRotations)
		for i := range barriers {
			barriers[i].Add(numberOfGroups)
		}
		wg := new(sync.WaitGroup)
		wg.Add(numberOfGroups)
		for i := 0; i < numberOfGroups; i++ {
			go func(index int) {
			OuterLoop:
				for j := 0; j < numberOfRotations; j++ {
					// Wait for workers before to finish up
					if j != 0 {
						barriers[j-1].Wait()
					}
					nextGroupID := (int(students[index].RegisteredCourses[CourseID(1)]) + numberOfSteps) % numberOfGroups
					for {
						err := students[index].ChangeGroup(context.Background(), &courses, CourseID(1), GroupID(nextGroupID), noOpBatcher{})
						// If there is no capacity, we just try again
						if errors.Is(err, NoCapacityLeftErr) {
							runtime.Gosched() // Let other threads do stuff
							continue
						}
						// Otherwise we check for errors
						if err != nil {
							t.Errorf("unexcpeceted error when rotating: %s", err)
							break OuterLoop
						}
						// We are good so continue
						break
					}
					barriers[j].Done()
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
		// Check the position
		positionOffset := (numberOfRotations * numberOfSteps) % numberOfGroups
		for i := range students {
			assert.Equalf(t, map[CourseID]GroupID{
				CourseID(1): GroupID((i + positionOffset) % numberOfGroups),
			}, students[i].RegisteredCourses, "error on %d %d student %d",
				numberOfRotations, numberOfSteps, students[i].ID)
		}
	}
}

func TestStudentForceEnrollCourse(t *testing.T) {
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
		emptyCourses(&courses)
		std := Student{RegisteredCourses: map[CourseID]GroupID{}}
		std.MaxUnits = math.MaxUint8
		clk.Set(time.Unix(10, 0))
		std.EnrollmentStartTime = time.Unix(15, 0).UnixMilli()
		assert.NoError(t, std.ForceEnrollCourse(context.Background(), &courses, CourseID(1), GroupID(1), noOpBatcher{}))
	})
	// We allow all other register times without setting them
	clk.Set(time.Unix(1, 0))
	// A non-existent course must throw an error
	t.Run("non existent course", func(t *testing.T) {
		emptyCourses(&courses)
		std := Student{RegisteredCourses: map[CourseID]GroupID{}}
		assert.ErrorIs(t, std.ForceEnrollCourse(context.Background(), &courses, CourseID(1), GroupID(10), noOpBatcher{}), NotExistsErr)
		assert.ErrorIs(t, std.ForceEnrollCourse(context.Background(), &courses, CourseID(100), GroupID(1), noOpBatcher{}), NotExistsErr)
	})
	t.Run("sex lock", func(t *testing.T) {
		emptyCourses(&courses)
		std := Student{
			StudentSex:        SexFemale,
			MaxUnits:          math.MaxUint8,
			RegisteredCourses: map[CourseID]GroupID{},
		}
		assert.NoError(t, std.ForceEnrollCourse(context.Background(), &courses, CourseID(5), GroupID(1), noOpBatcher{}))
		std.RegisteredCourses = map[CourseID]GroupID{}
		assert.NoError(t, std.ForceEnrollCourse(context.Background(), &courses, CourseID(5), GroupID(2), noOpBatcher{}))
		std.StudentSex = SexMale
		std.RegisteredCourses = map[CourseID]GroupID{}
		assert.NoError(t, std.ForceEnrollCourse(context.Background(), &courses, CourseID(5), GroupID(2), noOpBatcher{}))
		std.RegisteredCourses = map[CourseID]GroupID{}
		assert.NoError(t, std.ForceEnrollCourse(context.Background(), &courses, CourseID(5), GroupID(1), noOpBatcher{}))
	})
	// Check max registered courses
	t.Run("unit limit", func(t *testing.T) {
		emptyCourses(&courses)
		std := Student{
			MaxUnits:          3,
			RegisteredCourses: map[CourseID]GroupID{},
		}
		assert.NoError(t, std.ForceEnrollCourse(context.Background(), &courses, CourseID(2), GroupID(1), noOpBatcher{}))
		assert.Equal(t, uint8(3), std.RegisteredUnits)
		assert.NoError(t, std.ForceEnrollCourse(context.Background(), &courses, CourseID(1), GroupID(1), noOpBatcher{}))
	})
	// Do not allow already registered courses
	t.Run("already registered", func(t *testing.T) {
		emptyCourses(&courses)
		std := Student{
			MaxUnits:          math.MaxUint8,
			RegisteredCourses: map[CourseID]GroupID{},
		}
		assert.NoError(t, std.ForceEnrollCourse(context.Background(), &courses, CourseID(1), GroupID(1), noOpBatcher{}))
		assert.ErrorIs(t, std.ForceEnrollCourse(context.Background(), &courses, CourseID(1), GroupID(1), noOpBatcher{}), AlreadyRegisteredErr)
		assert.ErrorIs(t, std.ForceEnrollCourse(context.Background(), &courses, CourseID(1), GroupID(2), noOpBatcher{}), AlreadyRegisteredErr)
	})
	t.Run("exam time", func(t *testing.T) {
		emptyCourses(&courses)
		std := Student{
			MaxUnits: math.MaxUint8,
			RegisteredCourses: map[CourseID]GroupID{
				CourseID(1): GroupID(1),
			},
		}
		assert.NoError(t, std.ForceEnrollCourse(context.Background(), &courses, CourseID(2), GroupID(1), noOpBatcher{}))
	})
	// Class times must not overlap as well
	t.Run("class time", func(t *testing.T) {
		emptyCourses(&courses)
		std := Student{
			MaxUnits: math.MaxUint8,
			RegisteredCourses: map[CourseID]GroupID{
				CourseID(2): GroupID(1),
			},
		}
		assert.NoError(t, std.ForceEnrollCourse(context.Background(), &courses, CourseID(3), GroupID(1), noOpBatcher{}))
	})
	// Capacity error
	t.Run("capacity", func(t *testing.T) {
		emptyCourses(&courses)
		courses.courses[CourseID(1)][0].Capacity = 1
		courses.courses[CourseID(1)][0].RegisteredStudents = map[StudentID]struct{}{StudentID(100): {}}
		courses.courses[CourseID(1)][0].ReserveCapacity = 1
		courses.courses[CourseID(1)][0].ReserveQueue = util.NewQueue[StudentID]()
		courses.courses[CourseID(1)][0].ReserveQueue.Enqueue(StudentID(101))
		std := Student{
			ID:                1,
			MaxUnits:          math.MaxUint8,
			RegisteredCourses: map[CourseID]GroupID{},
		}
		assert.NoError(t, std.ForceEnrollCourse(context.Background(), &courses, CourseID(1), GroupID(1), noOpBatcher{}))
		assert.Equal(t, 2, courses.courses[CourseID(1)][0].Capacity)
		assert.Equal(t, map[StudentID]struct{}{StudentID(100): {}, StudentID(1): {}}, courses.courses[CourseID(1)][0].RegisteredStudents)
	})
}

func BenchmarkStudentAverage(b *testing.B) {
	// Setup time and rng
	clk := clock.NewMock()
	clk.Set(time.Unix(1, 0))
	studentClock = clk
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Setup courses and students
	const numberOfStudents = 10000
	const numberOfCourses = 2000
	const numberOfGroups = 5
	students := make(map[StudentID]*Student, numberOfStudents)
	courses := Courses{courses: make(map[CourseID][]*Course, numberOfCourses)}
	for i := 0; i < numberOfCourses; i++ {
		courses.courses[CourseID(i)] = make([]*Course, numberOfGroups)
		for j := 0; j < numberOfGroups; j++ {
			courses.courses[CourseID(i)][j] = &Course{
				ID:                 CourseID(i),
				GroupID:            GroupID(j),
				Units:              0, // don't account anything
				Capacity:           5,
				RegisteredStudents: map[StudentID]struct{}{},
				ReserveCapacity:    5,
				ReserveQueue:       util.NewQueue[StudentID](),
				ExamTime:           newAtomicFromValue(int64(i*numberOfGroups + j)),
				ClassHeldTime:      NewClassTime([]time.Weekday{time.Wednesday}, NewTimeOnly(uint16(i*numberOfGroups+j)), NewTimeOnly(uint16(i*numberOfGroups+j+1))),
			}
		}
	}
	for i := 0; i < numberOfStudents; i++ {
		students[StudentID(i)] = &Student{
			ID:                StudentID(i),
			RegisteredCourses: map[CourseID]GroupID{},
		}
	}
	// We benchmark like this:
	// We choose a student then an action from 3 possible actions
	// Then we randomly choose a course and hope for the best
	totalActions := 0
	successfulActions := 0
	for i := 0; i < b.N; i++ {
		totalActions++
		stdNumber := StudentID(rng.Intn(numberOfStudents))
		students[stdNumber].RemainingActions = 1
		action := rng.Intn(10)
		var err error
		if action <= 7 {
			err = students[stdNumber].EnrollCourse(context.Background(), &courses, CourseID(rng.Intn(numberOfCourses)), GroupID(rng.Intn(numberOfGroups)), noOpBatcher{})
		} else if action == 8 {
			var course CourseID
			if len(students[stdNumber].RegisteredCourses) != 0 {
				course = randomKeyFromMap(students[stdNumber].RegisteredCourses)
			}
			err = students[stdNumber].ChangeGroup(context.Background(), &courses, course, GroupID(rng.Intn(numberOfGroups)), noOpBatcher{})
		} else {
			var course CourseID
			if len(students[stdNumber].RegisteredCourses) != 0 {
				course = randomKeyFromMap(students[stdNumber].RegisteredCourses)
			}
			err = students[stdNumber].DisenrollCourse(context.Background(), &courses, course, noOpBatcher{})
		}
		if err == nil {
			successfulActions++
		}
	}
	b.Logf("done %d actions and %d of them were ok", totalActions, successfulActions)
}

func BenchmarkStudentEnrollOnly(b *testing.B) {
	// Setup time and rng
	clk := clock.NewMock()
	clk.Set(time.Unix(1, 0))
	studentClock = clk
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Setup courses and students
	const numberOfStudents = 10000
	const numberOfCourses = 2000
	const numberOfGroups = 5
	students := make(map[StudentID]*Student, numberOfStudents)
	courses := Courses{courses: make(map[CourseID][]*Course, numberOfCourses)}
	for i := 0; i < numberOfCourses; i++ {
		courses.courses[CourseID(i)] = make([]*Course, numberOfGroups)
		for j := 0; j < numberOfGroups; j++ {
			courses.courses[CourseID(i)][j] = &Course{
				ID:                 CourseID(i),
				GroupID:            GroupID(j),
				Units:              0, // don't account anything
				Capacity:           5,
				RegisteredStudents: map[StudentID]struct{}{},
				ReserveCapacity:    5,
				ReserveQueue:       util.NewQueue[StudentID](),
				ExamTime:           newAtomicFromValue(int64(i*numberOfGroups + j)),
				ClassHeldTime:      NewClassTime([]time.Weekday{time.Wednesday}, NewTimeOnly(uint16(i*numberOfGroups+j)), NewTimeOnly(uint16(i*numberOfGroups+j+1))),
			}
		}
	}
	for i := 0; i < numberOfStudents; i++ {
		students[StudentID(i)] = &Student{
			ID:                StudentID(i),
			RegisteredCourses: map[CourseID]GroupID{},
		}
	}
	// Benchmark
	totalActions := 0
	successfulActions := 0
	for i := 0; i < b.N; i++ {
		totalActions++
		stdNumber := rng.Intn(numberOfStudents)
		err := students[StudentID(stdNumber)].EnrollCourse(context.Background(), &courses, CourseID(rng.Intn(numberOfCourses)), GroupID(rng.Intn(numberOfGroups)), noOpBatcher{})
		if err == nil {
			successfulActions++
		}
	}
	b.Logf("done %d actions and %d of them were ok", totalActions, successfulActions)
}

//goland:noinspection GoVetCopyLock
func newAtomicTimeUnix(t time.Time) atomic.Int64 {
	result := atomic.Int64{}
	result.Store(t.Unix())
	return result
}

//goland:noinspection GoVetCopyLock
func newAtomicFromValue(in int64) atomic.Int64 {
	result := atomic.Int64{}
	result.Store(in)
	return result
}

// from the very nice https://stackoverflow.com/q/23482786/4213397
func randomKeyFromMap[K comparable, V any](m map[K]V) K {
	for k := range m {
		return k
	}
	panic("empty map")
}

// emptyCourses will empty all courses in a Courses object
func emptyCourses(courses *Courses) {
	for _, course := range courses.courses {
		for _, group := range course {
			group.RegisteredStudents = map[StudentID]struct{}{}
			group.ReserveQueue = util.NewQueue[StudentID]()
		}
	}
}
