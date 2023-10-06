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

// This file contains endpoints which are exposed to student.
// Only three actions are supported which are: Enroll, Disenroll and Change group.

// StudentEnroll must be called with PUT to enroll a student.
func (api *API) StudentEnroll(ctx context.Context, r *proto.StudentEnrollRequest) (*emptypb.Empty, error) {
	// Get student
	std, ok := api.Students[course.StudentID(r.StudentId)]
	if !ok {
		return nil, status.Error(codes.NotFound, "student_id")
	}
	// Enroll
	err := std.EnrollCourse(ctx, api.Courses, course.CourseID(r.CourseId), course.GroupID(r.GroupId), api.Broker)
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

// StudentDisenroll must be called with DELETE to disenroll a student.
func (api *API) StudentDisenroll(ctx context.Context, r *proto.StudentDisenrollRequest) (*emptypb.Empty, error) {
	// Get student
	std, ok := api.Students[course.StudentID(r.StudentId)]
	if !ok {
		return nil, status.Error(codes.NotFound, "student_id")
	}
	// Disenroll
	err := std.DisenrollCourse(ctx, api.Courses, course.CourseID(r.CourseId), api.Broker)
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
	return new(emptypb.Empty), err
}

// StudentChangeGroup must be called with PATCH to change group of a student.
func (api *API) StudentChangeGroup(ctx context.Context, r *proto.StudentChangeGroupRequest) (*emptypb.Empty, error) {
	// Get student
	std, ok := api.Students[course.StudentID(r.StudentId)]
	if !ok {
		return nil, status.Error(codes.NotFound, "student_id")
	}
	// Change group
	err := std.ChangeGroup(ctx, api.Courses, course.CourseID(r.CourseId), course.GroupID(r.NewGroupId), api.Broker)
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
