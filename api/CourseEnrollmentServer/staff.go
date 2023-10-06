package CourseEnrollmentServer

import (
	"CourseEnrollment/pkg/course"
	"CourseEnrollment/pkg/proto"
	"context"
	"github.com/go-faster/errors"
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
func (api *API) ForceEnroll(ctx context.Context, req *proto.StudentEnrollRequest) (*emptypb.Empty, error) {
	// Get student
	std, ok := api.Students[course.StudentID(req.StudentId)]
	if !ok {
		return nil, status.Error(codes.NotFound, "student_id")
	}
	// Enroll
	err := std.ForceEnrollCourse(ctx, api.Courses, course.CourseID(req.CourseId), course.GroupID(req.GroupId), api.Broker)
	if err != nil {
		var batchError course.BatchError
		if errors.As(err, &batchError) {
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

// ForceDisenroll will forcibly remove a user from a course.
// This means that the api won't check for registration time nor remaining actions.
// This call won't change the remaining actions of user.
func (api *API) ForceDisenroll(ctx context.Context, req *proto.StudentDisenrollRequest) (*emptypb.Empty, error) {
	// Get student
	std, ok := api.Students[course.StudentID(req.StudentId)]
	if !ok {
		return nil, status.Error(codes.NotFound, "student_id")
	}
	// Disenroll
	err := std.ForceDisenrollCourse(ctx, api.Courses, course.CourseID(req.CourseId), api.Broker)
	if err != nil {
		var batchError course.BatchError
		if errors.As(err, &batchError) {
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

// ChangeCapacity will update a course's capacity. It can fail if we try to shrink the capacity
// while users are registered in course.
func (api *API) ChangeCapacity(ctx context.Context, req *proto.ChangeCourseCapacityRequest) (*emptypb.Empty, error) {
	// Get the course
	c := api.Courses.GetCourse(course.CourseID(req.CourseId), course.GroupID(req.CourseId))
	if c == nil {
		return nil, status.Error(codes.NotFound, "course")
	}
	// Update capacity
	err := c.UpdateCapacity(ctx, int(req.NewCapacity), api.Broker)
	if err != nil {
		var batchError course.BatchError
		if errors.As(err, &batchError) {
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
