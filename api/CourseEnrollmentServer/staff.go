package CourseEnrollmentServer

import (
	"CourseEnrollment/pkg/course"
	"CourseEnrollment/pkg/proto"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetStudentsInCourse gets all the students in a course
func (api *API) GetStudentsInCourse(_ context.Context, req *proto.StudentsOfCourseRequest) (*proto.StudentsOfCourseResponse, error) {
	// Get the course
	c := api.Courses.GetCourse(course.CourseID(req.CourseId), course.GroupID(req.GroupId))
	if c == nil {
		return nil, status.Error(codes.NotFound, "course")
	}
	// Done!
	return c.ToStudentsOfCourseResponseProto(), nil
}
