package AuthCore

import (
	"CourseEnrollment/pkg/course"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"
)

// createJWTToken will create a JWT token for user authorization
func createJWTToken(key []byte, userID uint64, department course.DepartmentID, isStaff bool) (string, error) {
	now := time.Now()
	return jwt.NewWithClaims(signingMethod, JWTToken{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtIssuer,
			Subject:   strconv.FormatUint(userID, 10),
			ExpiresAt: jwt.NewNumericDate(now.Add(jwtTTL)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Department: department,
		IsStaff:    isStaff,
	}).SignedString(key)
}
