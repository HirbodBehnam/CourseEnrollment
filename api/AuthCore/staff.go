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
	// TODO
}

// ForceDisenroll will simply disenroll a student from course.
// It fails if student is not enrolled in the course.
func (a *API) ForceDisenroll(c *gin.Context) {
	// TODO
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
	// TODO
}
