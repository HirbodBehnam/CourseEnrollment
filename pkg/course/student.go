package course

import (
	"CourseEnrollment/pkg/proto"
	"context"
	"fmt"
	"github.com/benbjohnson/clock"
	"sync"
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
	// This value must not change... So no atomic.
	// This is in unix milliseconds.
	EnrollmentStartTime int64
	// Remaining actions such as removing a course or changing group
	RemainingActions uint8
	// Maximum units user can pick up
	MaxUnits uint8
	// How many units user has registered in
	RegisteredUnits uint8
	StudentSex      Sex
	// List of courses which the student has enrolled in. The key is the course ID and the value is
	// the group ID
	RegisteredCourses map[CourseID]GroupID
	// A simple locker for this user
	mu sync.RWMutex
}

// EnrollCourse tries to enroll the student in a course.
// It does all the checks and then enrolls the student if possible.
func (s *Student) EnrollCourse(ctx context.Context, courses *Courses, courseID CourseID, groupID GroupID, batcher Batcher) error {
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
		if examTimesIntersect(registeredCourse.ExamTime.Load(), course.ExamTime.Load()) {
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
	registered, err := course.EnrollStudent(ctx, s.ID, batcher)
	if err != nil {
		return err
	}
	if !registered {
		return NoCapacityLeftErr
	}
	// We are good!
	s.RegisteredUnits += course.Units
	s.RegisteredCourses[courseID] = groupID
	return nil
}

// DisenrollCourse will remove student from a course from
func (s *Student) DisenrollCourse(ctx context.Context, courses *Courses, courseID CourseID, batcher Batcher) error {
	// We check the start time at very first
	if !s.IsEnrollTimeOK() {
		return NotEnrollmentTimeErr
	}
	// Lock the user to do stuff with them
	s.mu.Lock()
	defer s.mu.Unlock()
	// Check the actions
	if s.RemainingActions == 0 {
		return NoRemainingActionsErr
	}
	// Get the course
	groupID, exists := s.RegisteredCourses[courseID]
	if !exists {
		return NotExistsErr
	}
	course := courses.GetCourse(courseID, groupID)
	if course == nil {
		panic(fmt.Sprintf("invalid registered lesson %d-%d for user %d", courseID, groupID, s.ID))
	}
	// Disenroll
	err := course.DisenrollStudent(ctx, s.ID, batcher)
	if err != nil {
		return err
	}
	// Remove from map
	delete(s.RegisteredCourses, courseID)
	s.RegisteredUnits -= course.Units
	s.RemainingActions--
	return nil
}

// ChangeGroup will atomically change group of a user in a course
func (s *Student) ChangeGroup(ctx context.Context, courses *Courses, courseID CourseID, destinationGroupID GroupID, batcher Batcher) error {
	// We check the start time at very first
	if !s.IsEnrollTimeOK() {
		return NotEnrollmentTimeErr
	}
	// Lock the user to do stuff with them
	s.mu.Lock()
	defer s.mu.Unlock()
	// Check the actions
	if s.RemainingActions == 0 {
		return NoRemainingActionsErr
	}
	// Get the course
	sourceGroupID, exists := s.RegisteredCourses[courseID]
	if !exists {
		return NotExistsErr
	}
	// Check same group ID
	if sourceGroupID == destinationGroupID {
		return PlayedYourselfErr
	}
	// Now, get the courses
	sourceCourse := courses.GetCourse(courseID, sourceGroupID)
	if sourceCourse == nil {
		panic(fmt.Sprintf("invalid registered lesson %d-%d for user %d", courseID, sourceGroupID, s.ID))
	}
	destinationCourse := courses.GetCourse(courseID, destinationGroupID)
	if destinationCourse == nil {
		return NotExistsErr
	}
	// Check the time of the course with registered courses (except the source)
	for registeredCourseID, registeredGroupID := range s.RegisteredCourses {
		// We are going to remove this, so we don't check it
		if registeredCourseID == courseID {
			continue
		}
		// Get the course
		registeredCourse := courses.GetCourse(registeredCourseID, registeredGroupID)
		if registeredCourse == nil {
			panic(fmt.Sprintf("inconsistent user state: course %d group %d is registered but not found", registeredCourseID, registeredGroupID))
		}
		// Check exam time
		if examTimesIntersect(registeredCourse.ExamTime.Load(), destinationCourse.ExamTime.Load()) {
			return ExamConflictErr{
				CourseID: registeredCourse.ID,
				GroupID:  registeredCourse.GroupID,
			}
		}
		// Check time
		if registeredCourse.ClassHeldTime.Intersects(&destinationCourse.ClassHeldTime) {
			return ClassTimeConflictErr{
				CourseID: registeredCourse.ID,
				GroupID:  registeredCourse.GroupID,
			}
		}
	}
	// Change the group
	changed, err := sourceCourse.ChangeGroupOfStudent(ctx, s.ID, destinationCourse, batcher)
	if err != nil {
		return err
	}
	if !changed {
		return NoCapacityLeftErr
	}
	// Done!
	s.RemainingActions--
	s.RegisteredCourses[courseID] = destinationGroupID
	return nil
}

// ForceEnrollCourse will forcibly enroll the student in a course.
// This will add capacity to course if needed.
// This function will return error if user is already registered in the course
func (s *Student) ForceEnrollCourse(ctx context.Context, courses *Courses, courseID CourseID, groupID GroupID, batcher Batcher) error {
	// We get the course which is basically lock-free. (we are all reading from this map)
	course := courses.GetCourse(courseID, groupID)
	if course == nil {
		return NotExistsErr
	}
	// Lock the user to do stuff with them
	s.mu.Lock()
	defer s.mu.Unlock()
	// Check if user has already registered in this course
	if _, alreadyRegistered := s.RegisteredCourses[courseID]; alreadyRegistered {
		return AlreadyRegisteredErr
	}
	// Register in course
	err := course.ForceEnroll(ctx, s.ID, batcher)
	if err != nil {
		return err
	}
	// Update values
	s.RegisteredCourses[courseID] = groupID
	s.RegisteredUnits += course.Units
	return nil
}

// IsEnrollTimeOK checks if the user can enroll based on its Student.EnrollmentStartTime.
// Students have one hour to do their enrollment
func (s *Student) IsEnrollTimeOK() bool {
	const enrollmentDurationMilliseconds = 60 * 60 * 1000
	now := studentClock.Now().UnixMilli()
	return now > s.EnrollmentStartTime && now < s.EnrollmentStartTime+enrollmentDurationMilliseconds
}

// GetEnrolledCoursesProto gets all the enrolled courses of user as a protobuf message
func (s *Student) GetEnrolledCoursesProto(courses *Courses) *proto.StudentCourseDataArray {
	// Lock student
	s.mu.RLock()
	defer s.mu.RUnlock()
	// Create the result and populate it
	result := &proto.StudentCourseDataArray{
		Data: make([]*proto.StudentCourseData, 0, len(s.RegisteredCourses)),
	}
	for courseID, groupID := range s.RegisteredCourses {
		course := courses.GetCourse(courseID, groupID)
		if course == nil {
			panic(fmt.Sprintf("invalid registered lesson %d-%d for user %d", courseID, groupID, s.ID))
		}
		result.Data = append(result.Data, course.ToStudentCourseDataProto(s.ID))
	}
	// Done
	return result
}

// examTimesIntersect checks if two exam times intersect.
// As a side note that why this is a separate function, 0 as time means no exam.
func examTimesIntersect(a, b int64) bool {
	// If at least one of them doesn't have exam, it's fine
	if a == 0 || b == 0 {
		return false
	}
	return a == b
}
