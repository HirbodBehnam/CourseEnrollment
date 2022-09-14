package CourseEnrollmentServer

import (
	globalapi "CourseEnrollment/api"
	"CourseEnrollment/pkg/course"
	"CourseEnrollment/pkg/dbbatch"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// This file contains endpoints which are exposed to student.
// Only three actions are supported which are: Enroll, Disenroll and Change group.
// My preferred way to distinguish between these actions is to use different HTTP methods:
// PUT, DELETE, PATCH

// These errors must be sent if the data provided in request is not valid
const invalidStudentIDErr = "invalid student ID"
const invalidCourseIDErr = "invalid course ID"
const invalidGroupIDErr = "invalid group ID"

// StudentEnroll must be called with PUT to enroll a student.
func (api *API) StudentEnroll(c *gin.Context) {
	// Parse values
	std, courseID, ok := api.parseCourseAndStudentID(c)
	if !ok {
		return
	}
	// Get group ID
	groupID, err := strconv.ParseUint(c.Query("group"), 10, 8)
	if err != nil {
		c.JSON(http.StatusBadRequest, globalapi.Error{Message: invalidGroupIDErr})
		return
	}
	// Enroll
	err = std.EnrollCourse(api.Courses, courseID, course.GroupID(groupID))
	if err != nil {
		c.JSON(http.StatusNotAcceptable, globalapi.Error{Message: err.Error()})
		return
	}
	// Send to message broker for database
	err = api.Broker.UpdateDatabase(dbbatch.Message{
		Type:     dbbatch.MessageActionTypeEnroll,
		StdID:    std.ID,
		CourseID: courseID,
		GroupID:  course.GroupID(groupID),
	})
	if err != nil {
		// FUCK
		// TODO
	}
	// Done
	c.Status(http.StatusNoContent)
}

// StudentDisenroll must be called with DELETE to disenroll a student.
func (api *API) StudentDisenroll(c *gin.Context) {
	// Parse values
	std, courseID, ok := api.parseCourseAndStudentID(c)
	if !ok {
		return
	}
	// Disenroll
	err := std.DisenrollCourse(api.Courses, courseID)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, globalapi.Error{Message: err.Error()})
		return
	}
	// Send to message broker for database
	err = api.Broker.UpdateDatabase(dbbatch.Message{
		Type:     dbbatch.MessageActionTypeDisenroll,
		StdID:    std.ID,
		CourseID: courseID,
	})
	if err != nil {
		// FUCK
		// TODO
	}
	// Done
	c.Status(http.StatusNoContent)
}

// StudentChangeGroup must be called with PATCH to change group of a student.
func (api *API) StudentChangeGroup(c *gin.Context) {
	// Parse values
	std, courseID, ok := api.parseCourseAndStudentID(c)
	if !ok {
		return
	}
	// Get group ID
	groupID, err := strconv.ParseUint(c.Query("group"), 10, 8)
	if err != nil {
		c.JSON(http.StatusBadRequest, globalapi.Error{Message: invalidGroupIDErr})
		return
	}
	// Enroll
	err = std.ChangeGroup(api.Courses, courseID, course.GroupID(groupID))
	if err != nil {
		c.JSON(http.StatusNotAcceptable, globalapi.Error{Message: err.Error()})
		return
	}
	// Send to message broker for database
	err = api.Broker.UpdateDatabase(dbbatch.Message{
		Type:     dbbatch.MessageActionTypeChangeGroup,
		StdID:    std.ID,
		CourseID: courseID,
		GroupID:  course.GroupID(groupID),
	})
	if err != nil {
		// FUCK
		// TODO
	}
}

// parseCourseAndStudentID will parse course ID and student ID from the url.
// It returns false if it couldn't parse them otherwise true.
func (api *API) parseCourseAndStudentID(c *gin.Context) (*course.Student, course.CourseID, bool) {
	// Parse the student ID
	stdID, err := strconv.ParseUint(c.Param("stdID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, globalapi.Error{Message: invalidStudentIDErr})
		return nil, 0, false
	}
	// Parse the course ID
	courseID, err := strconv.ParseInt(c.Param("stdID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, globalapi.Error{Message: invalidCourseIDErr})
		return nil, 0, false
	}
	// Get student
	std, ok := api.Students[course.StudentID(stdID)]
	if !ok {
		c.JSON(http.StatusBadRequest, globalapi.Error{Message: invalidStudentIDErr})
		return nil, 0, false
	}
	// Done
	return std, course.CourseID(courseID), true
}
