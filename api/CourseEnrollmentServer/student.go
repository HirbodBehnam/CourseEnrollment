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

// This file contains endpoints which are exposed to student.
// Only three actions are supported which are: Enroll, Disenroll and Change group.

// StudentEnroll must be called with PUT to enroll a student.
func (api *API) StudentEnroll(_ context.Context, r *proto.StudentEnrollRequest) (resp *emptypb.Empty, err error) {
	resp = new(emptypb.Empty)
	// Get student
	std, ok := api.Students[course.StudentID(r.StudentId)]
	if !ok {
		err = status.Error(codes.NotFound, "student_id")
		return
	}
	// Enroll
	err = std.EnrollCourse(api.Courses, course.CourseID(r.CourseId), course.GroupID(r.GroupId), api.Broker)
	if err != nil {
		if batchError, ok := err.(course.BatchError); ok {
			err = status.Error(codes.Internal, "")
			log.WithError(batchError).Error("cannot batch data")
		} else {
			err = status.Error(codes.FailedPrecondition, err.Error())
		}
		return
	}
	// Done
	return resp, nil
}

// StudentDisenroll must be called with DELETE to disenroll a student.
func (api *API) StudentDisenroll(_ context.Context, r *proto.StudentDisenrollRequest) (resp *emptypb.Empty, err error) {
	resp = new(emptypb.Empty)
	// Get student
	std, ok := api.Students[course.StudentID(r.StudentId)]
	if !ok {
		err = status.Error(codes.NotFound, "student_id")
		return
	}
	// Disenroll
	err = std.DisenrollCourse(api.Courses, course.CourseID(r.CourseId), api.Broker)
	if err != nil {
		if batchError, ok := err.(course.BatchError); ok {
			err = status.Error(codes.Internal, "")
			log.WithError(batchError).Error("cannot batch data")
		} else {
			err = status.Error(codes.FailedPrecondition, err.Error())
		}
	}
	// Done
	return resp, err
}

// StudentChangeGroup must be called with PATCH to change group of a student.
func (api *API) StudentChangeGroup(_ context.Context, r *proto.StudentChangeGroupRequest) (resp *emptypb.Empty, err error) {
	resp = new(emptypb.Empty)
	// Get student
	std, ok := api.Students[course.StudentID(r.StudentId)]
	if !ok {
		err = status.Error(codes.NotFound, "student_id")
		return
	}
	// Change group
	err = std.ChangeGroup(api.Courses, course.CourseID(r.CourseId), course.GroupID(r.NewGroupId), api.Broker)
	if err != nil {
		if batchError, ok := err.(course.BatchError); ok {
			err = status.Error(codes.Internal, "")
			log.WithError(batchError).Error("cannot batch data")
		} else {
			err = status.Error(codes.FailedPrecondition, err.Error())
		}
	}
	// Done
	return resp, nil
}
