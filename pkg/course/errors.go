package course

import (
	"errors"
	"fmt"
)

// NotExistsErr means that the requested course does not exist
var NotExistsErr = errors.New("course does not exist")

// SexLockErr happens when student tries to pick a course with sex lock
var SexLockErr = errors.New("you cannot pick up this course due to sex lock")

// NotEnrollmentTimeErr means that user cannot do anything because it's not their enrollment
// time
var NotEnrollmentTimeErr = errors.New("it's not your enrollment time")

// UnitLimitReachedErr is when user cannot register anymore because the unit limit has been
// reached
var UnitLimitReachedErr = errors.New("unit limit has been reached")

// AlreadyRegisteredErr means that user is trying to register in a course
// which they have already registered in
var AlreadyRegisteredErr = errors.New("you are already registered in this course")

// ExamConflictErr is returned when the courses of user have exam conflicts
type ExamConflictErr struct {
	CourseID CourseID
	GroupID  GroupID
}

func (e ExamConflictErr) Error() string {
	return fmt.Sprintf("exam conflict with course %d-%d", e.CourseID, e.GroupID)
}

// ClassTimeConflictErr is returned when the courses of user have class time conflicts
type ClassTimeConflictErr struct {
	CourseID CourseID
	GroupID  GroupID
}

func (e ClassTimeConflictErr) Error() string {
	return fmt.Sprintf("class time conflict with course %d-%d", e.CourseID, e.GroupID)
}

// NoCapacityLeftErr means that we cannot register user because the capacity of course is fulled
var NoCapacityLeftErr = errors.New("this course's capacity is filled")

// NoRemainingActionsErr means that user cannot disenroll or change group because of lack of
// remaining actions
var NoRemainingActionsErr = errors.New("no more remaining actions left")

// PlayedYourselfErr happens when user tries to change group to its own registered group
var PlayedYourselfErr = errors.New("source and destination group ID is same")
