package course

import (
	"CourseEnrollment/pkg/util"
	"sync"
	"sync/atomic"
)

// CourseID is the type of the course's ID
type CourseID int32

// GroupID is the type of the group's ID used in Course
type GroupID uint8

// Courses held a list of all courses
type Courses struct {
	courses map[CourseID][]*Course
	mu      sync.RWMutex
}

// Course represents a single course
type Course struct {
	// The course ID
	ID CourseID
	// The group ID
	GroupID GroupID
	// What is the department of this course?
	Department DepartmentID
	// Lecturer name
	LecturerName string
	// Number of units
	Units uint8
	// The total capacity of this course
	Capacity int
	// List of students which have registered in this course
	RegisteredStudents map[StudentID]struct{}
	// Total number of students which can be in reserved queue
	ReserveCapacity int
	// The queue of reserved students
	ReserveQueue util.Queue[StudentID]
	// When is the exam of this course? In unix epoch (seconds)
	ExamTime atomic.Int64
	// The time and days which class is held on
	ClassHeldTime ClassTime
	// Does this class has a sex lock?
	SexLock SexLock
	// The mutex to work with this course
	mu sync.Mutex
}

// GetCourse will get the course based on group ID and course ID. If the course does not exist,
// it returns nil
func (c *Courses) GetCourse(courseID CourseID, groupID GroupID) *Course {
	var result *Course
	c.mu.RLock()
	for _, course := range c.courses[courseID] {
		if course.GroupID == groupID {
			result = course
			break
		}
	}
	c.mu.RUnlock()
	return result
}

// EnrollStudent will enroll the student in this course.
//
// NOTE: This method does not check for class conflicts or etc. It just adds a student to a course
// if possible. At first, it tries to add student to Course.RegisteredStudents and if it's not possible,
// (due to limitations), it ties to add they in Course.ReserveQueue. If that's also not possible,
// it returns false.
// It also does not check if the student is enrolled in this course or not
func (c *Course) EnrollStudent(studentID StudentID) bool {
	// Lock the course and unlock it when we are leaving
	c.mu.Lock()
	defer c.mu.Unlock()
	// At first check the registered count
	if len(c.RegisteredStudents) < c.Capacity {
		c.RegisteredStudents[studentID] = struct{}{}
		return true
	}
	// Next check the reserve queue
	if c.ReserveQueue.Len() < c.ReserveCapacity {
		c.ReserveQueue.Enqueue(studentID)
		return true
	}
	// Could not enroll the course
	return false
}

// UnrollStudent will remove the student from course.
//
// Returns true if the student was enrolled before; Otherwise false
func (c *Course) UnrollStudent(studentID StudentID) bool {
	// Lock the course and unlock it when we are leaving
	c.mu.Lock()
	defer c.mu.Unlock()
	// Check registered list
	if _, registered := c.RegisteredStudents[studentID]; registered {
		delete(c.RegisteredStudents, studentID)
		// Not put someone from reserve queue into registered users
		if c.ReserveQueue.Len() != 0 {
			c.RegisteredStudents[c.ReserveQueue.Dequeue()] = struct{}{}
		}
		// Done
		return true
	}
	// Otherwise remove from queue
	return c.ReserveQueue.Remove(studentID)
}
