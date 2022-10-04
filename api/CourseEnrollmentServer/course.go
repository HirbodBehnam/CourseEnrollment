package CourseEnrollmentServer

import (
	"CourseEnrollment/pkg/course"
	"CourseEnrollment/pkg/proto"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetStudentEnrolledCourses will get the enrolled courses of a student
func (api *API) GetStudentEnrolledCourses(_ context.Context, req *proto.GetStudentCoursesRequest) (*proto.StudentCourseDataArray, error) {
	// Get student
	std, exists := api.Students[course.StudentID(req.GetStudentId())]
	if !exists {
		return nil, status.Error(codes.NotFound, "student_id")
	}
	// Get enrolled courses
	return std.GetEnrolledCoursesProto(api.Courses), nil
}

// GetCoursesOfDepartment returns all the courses in a department
func (api *API) GetCoursesOfDepartment(_ context.Context, req *proto.GetDepartmentCoursesRequest) (*proto.DepartmentCourses, error) {
	return api.Courses.GetDepartmentCoursesProto(course.DepartmentID(req.GetDepartmentId())), nil
}
