package CourseEnrollmentServer

import (
	"CourseEnrollment/pkg/course"
	"CourseEnrollment/pkg/dbbatch"
	"CourseEnrollment/pkg/proto"
	"context"
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
	err = std.EnrollCourse(api.Courses, course.CourseID(r.CourseId), course.GroupID(r.GroupId))
	if err != nil {
		err = status.Error(codes.FailedPrecondition, err.Error())
		return
	}
	// Send to message broker for database
	err = api.Broker.UpdateDatabase(dbbatch.Message{
		Type:     dbbatch.MessageActionTypeEnroll,
		StdID:    std.ID,
		CourseID: course.CourseID(r.CourseId),
		GroupID:  course.GroupID(r.GroupId),
	})
	if err != nil {
		// FUCK
		// TODO
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
	err = std.DisenrollCourse(api.Courses, course.CourseID(r.CourseId))
	if err != nil {
		err = status.Error(codes.FailedPrecondition, err.Error())
		return
	}
	// Send to message broker for database
	err = api.Broker.UpdateDatabase(dbbatch.Message{
		Type:     dbbatch.MessageActionTypeDisenroll,
		StdID:    std.ID,
		CourseID: course.CourseID(r.CourseId),
	})
	if err != nil {
		// FUCK
		// TODO
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
	err = std.ChangeGroup(api.Courses, course.CourseID(r.CourseId), course.GroupID(r.NewGroupId))
	if err != nil {
		err = status.Error(codes.FailedPrecondition, err.Error())
		return
	}
	// Send to message broker for database
	err = api.Broker.UpdateDatabase(dbbatch.Message{
		Type:     dbbatch.MessageActionTypeChangeGroup,
		StdID:    std.ID,
		CourseID: course.CourseID(r.CourseId),
		GroupID:  course.GroupID(r.NewGroupId),
	})
	if err != nil {
		// FUCK
		// TODO
	}
	// Done
	return resp, nil
}
