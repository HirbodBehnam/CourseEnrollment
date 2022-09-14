package dbbatch

import "CourseEnrollment/pkg/course"

// MessageActionType contains the action which must be done on database
type MessageActionType uint8

const (
	MessageActionTypeEnroll MessageActionType = iota
	MessageActionTypeDisenroll
	MessageActionTypeChangeGroup
)

// Message is a type of data which must be queued in our message broker
type Message struct {
	// What shall be done
	Type MessageActionType `json:"type"`
	// For what student we do this?
	StdID course.StudentID `json:"stdID"`
	// What is the course ID?
	CourseID course.CourseID `json:"courseID"`
	// The group ID. In MessageActionTypeDisenroll, this field is empty,
	// in MessageActionTypeChangeGroup this is the destination group ID.
	// And in MessageActionTypeEnroll it's obvious what is this...
	GroupID course.GroupID `json:"groupID"`
}
