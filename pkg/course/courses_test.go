package course

import (
	"CourseEnrollment/pkg/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCourseEnrollStudent(t *testing.T) {
	assertion := assert.New(t)
	course := Course{
		Capacity:           10,
		RegisteredStudents: make(map[StudentID]struct{}),
		ReserveCapacity:    4,
		ReserveQueue:       util.NewQueue[StudentID](),
	}
	// Add to registered courses
	for i := 0; i < course.Capacity; i++ {
		assertion.True(course.EnrollStudent(StudentID(i)))
	}
	// Check
	assertion.Len(course.RegisteredStudents, course.Capacity)
	for i := 0; i < course.Capacity; i++ {
		_, exists := course.RegisteredStudents[StudentID(i)]
		assertion.True(exists)
	}
	assertion.Equal(0, course.ReserveQueue.Len())
	// Add to reserve queue
	for i := 0; i < course.ReserveCapacity; i++ {
		assertion.True(course.EnrollStudent(StudentID(course.ReserveCapacity + i)))
	}
	assertion.Equal(course.ReserveCapacity, course.ReserveQueue.Len())
	for i := 0; i < course.ReserveCapacity; i++ {
		index := course.ReserveQueue.Exists(StudentID(course.ReserveCapacity + i))
		assertion.Equal(i, index)
	}
	// Add over capacity
	for i := 0; i < 100; i++ {
		assertion.False(course.EnrollStudent(StudentID(course.ReserveCapacity + course.Capacity + i)))
	}
}

func TestCourseUnrollStudent(t *testing.T) {
	assertion := assert.New(t)
	course := Course{
		ID:                 CourseID(1),
		GroupID:            GroupID(1),
		Capacity:           10,
		RegisteredStudents: make(map[StudentID]struct{}),
		ReserveCapacity:    4,
		ReserveQueue:       util.NewQueue[StudentID](),
	}
	// At first, we just run unroll on empty course
	assertion.PanicsWithValue("user 1 has lesson 1-1 in their registered courses but lesson map does not have this user", func() {
		course.DisenrollStudent(StudentID(1))
	})
	// Then, we add some users to registered user. We don't go to reserved capacity
	for i := 0; i < 5; i++ {
		assertion.True(course.EnrollStudent(StudentID(i)))
	}
	// Then we unroll the first student
	assertion.NotPanics(func() {
		course.DisenrollStudent(StudentID(0))
	})
	assertion.Len(course.RegisteredStudents, 4) // zero must be removed
	for i := 1; i < 5; i++ {
		_, exists := course.RegisteredStudents[StudentID(i)]
		assertion.True(exists)
	}
	// Then we add them again
	course.RegisteredStudents = make(map[StudentID]struct{})
	for i := 0; i < course.Capacity; i++ {
		assertion.True(course.EnrollStudent(StudentID(i)))
	}
	for i := 0; i < course.ReserveCapacity; i++ {
		assertion.True(course.EnrollStudent(StudentID(course.Capacity + i)))
	}
	assertion.NotPanics(func() {
		course.DisenrollStudent(StudentID(0))
	})
	// Check it
	for i := 1; i < course.Capacity+1; i++ {
		_, exists := course.RegisteredStudents[StudentID(i)]
		assertion.True(exists)
	}
	// Check reserved
	{
		expectedReserved := make([]StudentID, 0, course.ReserveCapacity)
		for i := 1; i < course.ReserveCapacity; i++ {
			expectedReserved = append(expectedReserved, StudentID(course.Capacity+i))
		}
		assertion.Equal(expectedReserved, course.ReserveQueue.CopyAsArray())
	}
}
