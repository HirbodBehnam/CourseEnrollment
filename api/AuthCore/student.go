package AuthCore

import (
	pb "CourseEnrollment/pkg/proto"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
)

// EnrollStudent will enroll the student in a course
func (a *API) EnrollStudent(c *gin.Context) {
	std := c.MustGet(authInfoKey).(AuthData)
	request := c.MustGet(requestKey).(CourseEnrollmentRequest)
	// Send data to enrollment core
	_, err := a.CoreClient.StudentEnroll(c.Request.Context(), &pb.StudentEnrollRequest{
		StudentId: std.User,
		CourseId:  int32(request.CourseID),
		GroupId:   uint32(request.GroupID),
	})
	handleEnrollmentRPCError(c, err)
}

// ChangeGroupOfStudent will change the group of a student in one of their enrolled courses.
func (a *API) ChangeGroupOfStudent(c *gin.Context) {
	std := c.MustGet(authInfoKey).(AuthData)
	request := c.MustGet(requestKey).(CourseEnrollmentRequest)
	// Send data to enrollment core
	_, err := a.CoreClient.StudentChangeGroup(c.Request.Context(), &pb.StudentChangeGroupRequest{
		StudentId:  std.User,
		CourseId:   int32(request.CourseID),
		NewGroupId: uint32(request.GroupID),
	})
	handleEnrollmentRPCError(c, err)
}

// DisenrollStudent will disenroll the student from a course
func (a *API) DisenrollStudent(c *gin.Context) {
	std := c.MustGet(authInfoKey).(AuthData)
	// Get the course ID from query
	courseID, err := strconv.ParseInt(c.Query("course_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{reasonKey: "cannot parse course_id: " + err.Error()})
		return
	}
	// Send data to enrollment core
	_, err = a.CoreClient.StudentDisenroll(c.Request.Context(), &pb.StudentDisenrollRequest{
		StudentId: std.User,
		CourseId:  int32(courseID),
	})
	handleEnrollmentRPCError(c, err)
}

// handleEnrollmentRPCError will handle the error returned from a gRPC request which corresponds to
// an action which a student does.
func handleEnrollmentRPCError(c *gin.Context, err error) {
	if err != nil {
		if statusError, ok := status.FromError(err); ok {
			switch statusError.Code() {
			case codes.NotFound, codes.FailedPrecondition:
				// Some errors like "not your enrollment time or etc."
				// We directly send the error message
				c.JSON(http.StatusBadRequest, gin.H{reasonKey: statusError.Message()})
				return
			}
		}
		c.Status(http.StatusInternalServerError)
		log.WithError(err).Error("cannot enroll student")
		return
	}
	// Done
	c.Status(http.StatusNoContent)
}
