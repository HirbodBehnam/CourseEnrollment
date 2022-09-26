package AuthCore

import (
	"CourseEnrollment/pkg/course"
	"github.com/golang-jwt/jwt/v4"
)

// LoginRequest is the login request which user sends to us
type LoginRequest struct {
	Password string `json:"password" binding:"required"`
	User     uint64 `json:"user" binding:"required"`
	IsStaff  bool   `json:"staff"`
}

// JWTToken is the token stored in user's browsers
type JWTToken struct {
	jwt.RegisteredClaims
	Department course.DepartmentID `json:"department"`
	IsStaff    bool                `json:"staff"`
}

// AuthData is the struct which is passed to endpoints and contains the
// JWT authorization info
type AuthData struct {
	User       uint64
	Department course.DepartmentID
	IsStaff    bool
}

// TokenResult contains a JWT token only
type TokenResult struct {
	Token string `json:"token" binding:"required"`
}

// CourseEnrollmentRequest is the data which must be sent to us when
// a student wants to do anything with their courses except disenrolling.
type CourseEnrollmentRequest struct {
	CourseID course.CourseID `json:"course_id" binding:"required"`
	// On enrollment this is the group which user wants to enroll in.
	// On change group this is the destination group ID.
	GroupID course.GroupID `json:"group_id" binding:"required"`
}
