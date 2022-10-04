package AuthCore

import (
	"CourseEnrollment/pkg/proto"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// EnrolledCoursesOfStudent will return the enrolled courses of a student.
func (a *API) EnrolledCoursesOfStudent(c *gin.Context) {
	std := c.MustGet(authInfoKey).(AuthData)
	courses, err := a.CoreClient.GetStudentEnrolledCourses(c.Request.Context(), &proto.GetStudentCoursesRequest{StudentId: std.User})
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.WithError(err).WithField("user id", std.User).Error("cannot get enrolled courses")
		return
	}
	c.JSON(http.StatusOK, courses)
}

// CoursesOfDepartment will get all the courses in a department
func (a *API) CoursesOfDepartment(c *gin.Context) {
	// Get department
	department, err := strconv.ParseUint(c.Query("department"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{reasonKey: "cannot parse department id"})
		return
	}
	// Do the request
	courses, err := a.CoreClient.GetCoursesOfDepartment(c.Request.Context(), &proto.GetDepartmentCoursesRequest{DepartmentId: uint32(department)})
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.WithError(err).WithField("department", department).Error("cannot get department courses")
		return
	}
	c.JSON(http.StatusOK, courses)
}
