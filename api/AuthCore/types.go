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

// TokenResult contains a JWT token only
type TokenResult struct {
	Token string `json:"token" binding:"required"`
}
