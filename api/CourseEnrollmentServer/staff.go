package CourseEnrollmentServer

import (
	"CourseEnrollment/pkg/course"
	"CourseEnrollment/pkg/proto"
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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

// ForceEnroll will forcibly enroll a student in a course, increasing the capacity if needed
func (api *API) ForceEnroll(ctx context.Context, r *proto.StudentEnrollRequest) (*emptypb.Empty, error) {
	// Get student
	std, ok := api.Students[course.StudentID(r.StudentId)]
	if !ok {
		return nil, status.Error(codes.NotFound, "student_id")
	}
	// Enroll
	err := std.ForceEnrollCourse(ctx, api.Courses, course.CourseID(r.CourseId), course.GroupID(r.GroupId), api.Broker)
	if err != nil {
		if batchError, ok := err.(course.BatchError); ok {
			err = status.Error(codes.Internal, "")
			log.WithError(batchError).Error("cannot batch data")
		} else {
			err = status.Error(codes.FailedPrecondition, err.Error())
		}
		return nil, err
	}
	// Done
	return new(emptypb.Empty), nil
}
