package CourseEnrollmentServer

import (
	"CourseEnrollment/pkg/course"
	"CourseEnrollment/pkg/proto"
)

// API is the server API which is used in course enrollment server
type API struct {
	proto.UnimplementedCourseEnrollmentServerServiceServer
	// Broker must handle the queries and batch them.
	Broker course.Batcher
	// List of all students
	Students map[course.StudentID]*course.Student
	// List of all courses
	Courses *course.Courses
}
