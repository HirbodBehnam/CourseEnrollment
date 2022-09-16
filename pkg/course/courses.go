package course

import (
	"CourseEnrollment/pkg/util"
	"fmt"
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

// NewCourses creates Courses from it's map
func NewCourses(courses map[CourseID][]*Course) *Courses {
	return &Courses{courses: courses}
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
	// Could not enroll the course
	return c.threadUnsafeEnrollStudent(studentID)
}

// threadUnsafeEnrollStudent does EnrollStudent but without locking the course
func (c *Course) threadUnsafeEnrollStudent(studentID StudentID) bool {
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

// DisenrollStudent will remove the student from course.
//
// Will panic if the student is not enrolled in course.
func (c *Course) DisenrollStudent(studentID StudentID) {
	// Lock the course and unlock it when we are leaving
	c.mu.Lock()
	defer c.mu.Unlock()
	c.threadUnsafeDisenrollStudent(studentID)
}

// threadUnsafeDisenrollStudent is basically DisenrollStudent but without locking
// the course
func (c *Course) threadUnsafeDisenrollStudent(studentID StudentID) {
	// Check registered list
	if _, registered := c.RegisteredStudents[studentID]; registered {
		delete(c.RegisteredStudents, studentID)
		// Now put first person from reserve queue into registered users
		// (if exists)
		if c.ReserveQueue.Len() != 0 {
			c.RegisteredStudents[c.ReserveQueue.Dequeue()] = struct{}{}
		}
		// Done
		return
	}
	// Otherwise remove from queue
	if !c.ReserveQueue.Remove(studentID) {
		panic(fmt.Sprintf("user %d has lesson %d-%d in their registered courses but lesson map does not have this user", studentID, c.ID, c.GroupID))
	}
}

// ChangeGroupOfStudent tries to change the group of a student between two courses
func (c *Course) ChangeGroupOfStudent(studentID StudentID, other *Course) bool {
	// Check courses
	if c.ID != other.ID {
		panic("different courses provided")
	}
	if c.GroupID == other.GroupID {
		panic("same group provided")
	}
	// For locking, we at first lock the smaller group ID
	// to avoid deadlocks
	if c.GroupID < other.GroupID {
		c.mu.Lock()
		other.mu.Lock()
	} else {
		other.mu.Lock()
		c.mu.Lock()
	}
	defer c.mu.Unlock()
	defer other.mu.Unlock()
	// Check if user exists in this course
	if _, existsInRegistered := c.RegisteredStudents[studentID]; !existsInRegistered {
		if c.ReserveQueue.Exists(studentID) == -1 {
			panic(fmt.Sprintf("student %d does not exists in course %d-%d", studentID, c.ID, c.GroupID))
		}
	}
	// Now try to add it to other course
	if !other.threadUnsafeEnrollStudent(studentID) {
		return false // we could not enroll user due to capacity
	}
	// Now remove the user from this course
	c.threadUnsafeDisenrollStudent(studentID)
	return true
}
