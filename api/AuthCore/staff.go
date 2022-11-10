package AuthCore

import (
	"CourseEnrollment/pkg/proto"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// ForceEnroll will forcibly enroll a student in a course.
// It adds capacity to course if needed. The only possible way for this endpoint
// to fail is that user has this course already.
func (a *API) ForceEnroll(c *gin.Context) {
	var request StaffCourseEnrollmentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{reasonKey: err.Error()})
		return
	}
	_, err := a.CoreClient.ForceEnroll(c.Request.Context(), &proto.StudentEnrollRequest{
		StudentId: uint64(request.StudentID),
		CourseId:  int32(request.CourseID),
		GroupId:   uint32(request.GroupID),
	})
	handleEnrollmentRPCError(c, err)
}

// ForceDisenroll will simply disenroll a student from course.
// It fails if student is not enrolled in the course.
func (a *API) ForceDisenroll(c *gin.Context) {
	var request StaffCourseEnrollmentRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{reasonKey: err.Error()})
		return
	}
	_, err := a.CoreClient.ForceDisenroll(c.Request.Context(), &proto.StudentDisenrollRequest{
		StudentId: uint64(request.StudentID),
		CourseId:  int32(request.CourseID),
	})
	handleEnrollmentRPCError(c, err)
}

// CoursesOfStudent gets the list of courses which user is enrolled in
// or is in queue.
func (a *API) CoursesOfStudent(c *gin.Context) {
	// Get student ID
	stdID, err := strconv.ParseUint(c.Query("std_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{reasonKey: "invalid std_id"})
		return
	}
	// Get students from core
	courses, err := a.CoreClient.GetStudentEnrolledCourses(c.Request.Context(), &proto.GetStudentCoursesRequest{StudentId: stdID})
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.WithError(err).WithField("user id", stdID).Error("cannot get enrolled courses")
		return
	}
	c.JSON(http.StatusOK, courses)
}

// StudentsOfCourse gets all the students enrolled in a course.
func (a *API) StudentsOfCourse(c *gin.Context) {
	// Parse request
	var request CourseEnrollmentRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{reasonKey: err.Error()})
		return
	}
	// Get students
	result, err := a.CoreClient.GetStudentsInCourse(c.Request.Context(), &proto.StudentsOfCourseRequest{
		CourseId: int32(request.CourseID),
		GroupId:  uint32(request.GroupID),
	})
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.WithError(err).WithField("data", request).Error("cannot get enrolled students in course")
		return
	}
	// Send result
	c.JSON(http.StatusOK, result)
}
