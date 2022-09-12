package course

import (
	"fmt"
	"github.com/benbjohnson/clock"
	"sync"
	"sync/atomic"
)

// studentClock is the clock which we use to check the student enrollment time
var studentClock = clock.New()

// StudentID is the type of each student's ID
type StudentID uint64

// Student holds the info needed for a student
type Student struct {
	// The id of student
	ID StudentID
	// When does the course enrollment start for this student?
	// In unix milliseconds epoch because atomic operations.
	EnrollmentStartTime atomic.Int64
	// Remaining actions such as removing a course or changing group
	RemainingActions uint16
	// Maximum units user can pick up
	MaxUnits uint8
	// How many units user has registered in
	RegisteredUnits uint8
	// What department this user is for
	DepartmentID DepartmentID
	StudentSex   Sex
	// List of courses which the student has enrolled in. The key is the course ID and the value is
	// the group ID
	RegisteredCourses map[CourseID]GroupID
	// A simple locker for this user
	mu sync.Mutex
}

// EnrollCourse tries to enroll the student in a course.
// It does all the checks and then enrolls the student if possible.
func (s *Student) EnrollCourse(courses *Courses, courseID CourseID, groupID GroupID) error {
	// We check the start time at very first
	if !s.IsEnrollTimeOK() {
		return NotEnrollmentTimeErr
	}
	// We get the course which is basically lock-free. (we are all reading from this map)
	course := courses.GetCourse(courseID, groupID)
	if course == nil {
		return NotExistsErr
	}
	// We check the sex without lock because it doesn't change!
	if !sexLockCompatible(course.SexLock, s.StudentSex) {
		return SexLockErr
	}
	// Then we lock the user to do stuff with him/her
	s.mu.Lock()
	defer s.mu.Unlock()
	// We check the max units
	if s.RegisteredUnits+course.Units > s.MaxUnits {
		return UnitLimitReachedErr
	}
	// Check if user has already registered in this course
	if _, alreadyRegistered := s.RegisteredCourses[courseID]; alreadyRegistered {
		return AlreadyRegisteredErr
	}
	// Check the time of the course with registered courses
	for registeredCourseID, registeredGroupID := range s.RegisteredCourses {
		// Get the course
		registeredCourse := courses.GetCourse(registeredCourseID, registeredGroupID)
		if registeredCourse == nil {
			panic(fmt.Sprintf("inconsistent user state: course %d group %d is registered but not found", registeredCourseID, registeredGroupID))
		}
		// Check exam time
		if registeredCourse.ExamTime.Load() == course.ExamTime.Load() {
			return ExamConflictErr{
				CourseID: registeredCourse.ID,
				GroupID:  registeredCourse.GroupID,
			}
		}
		// Check time
		if registeredCourse.ClassHeldTime.Intersects(&course.ClassHeldTime) {
			return ClassTimeConflictErr{
				CourseID: registeredCourse.ID,
				GroupID:  registeredCourse.GroupID,
			}
		}
	}
	// At last, we register the course
	if !course.EnrollStudent(s.ID) {
		return NoCapacityLeftErr
	}
	// We are good!
	s.RegisteredUnits += course.Units
	s.RegisteredCourses[courseID] = groupID
	return nil
}

// IsEnrollTimeOK checks if the user can enroll based on its Student.EnrollmentStartTime.
// Students have one hour to do their enrollment
func (s *Student) IsEnrollTimeOK() bool {
	const enrollmentDurationMilliseconds = 60 * 60 * 1000
	now := studentClock.Now().UnixMilli()
	enrollmentTime := s.EnrollmentStartTime.Load()
	return now > enrollmentTime && now < enrollmentTime+enrollmentDurationMilliseconds
}
