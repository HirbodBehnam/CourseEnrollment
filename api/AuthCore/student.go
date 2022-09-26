package AuthCore

import (
	pb "CourseEnrollment/pkg/proto"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// EnrollStudent will enroll the student in a course
func (a *API) EnrollStudent(c *gin.Context) {
	std := c.MustGet(authInfoKey).(AuthData)
	// Parse the course data
	var request CourseEnrollmentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{reasonKey: err.Error()})
		return
	}
	// Send data to enrollment core
	_, err := a.CoreClient.StudentEnroll(c.Request.Context(), &pb.StudentEnrollRequest{
		StudentId: std.User,
		CourseId:  int32(request.CourseID),
		GroupId:   uint32(request.GroupID),
	})
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.WithError(err).Error("cannot enroll student")
		return
	}
	// Done
	c.Status(http.StatusNoContent)
}
