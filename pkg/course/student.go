package course

import (
	"sync"
	"time"
)

// StudentID is the type of each student's ID
type StudentID uint64

// Student holds the info needed for a student
type Student struct {
	// The id of student
	ID StudentID
	// When does the course enrollment start for this student?
	EnrollmentStartTime time.Time
	// Maximum units user can pick up
	MaxUnits uint16
	// Remaining actions such as removing a course or changing group
	RemainingActions uint16
	DepartmentID     DepartmentID
	StudentSex       Sex
	// List of courses which the student has enrolled in. The key is the course ID and the value is
	// the group ID
	RegisteredCourses map[CourseID]GroupID
	// A simple locker for this user
	mu sync.Mutex
}
