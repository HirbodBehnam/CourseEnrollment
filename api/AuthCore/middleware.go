package AuthCore

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strconv"
	"strings"
)

// JWTAuthMiddleware is a middleware which authenticates the JWT token of user
func (a *API) JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get header
		const headerName = "Authorization"
		const prefix = "Bearer "
		header := c.Request.Header.Get(headerName)
		if !strings.HasPrefix(header, prefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{reasonKey: "empty auth"})
			return
		}
		header = header[len(prefix):]
		// Parse the JWT
		token, err := jwt.ParseWithClaims(header, new(JWTToken), func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return a.jwtKey, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{reasonKey: "invalid auth"})
			return
		}
		claims, ok := token.Claims.(*JWTToken)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{reasonKey: "invalid auth"})
			return
		}
		// Set the data
		authData := AuthData{
			Department: claims.Department,
			IsStaff:    claims.IsStaff,
		}
		authData.User, err = strconv.ParseUint(claims.Subject, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{reasonKey: "invalid auth"})
			return
		}
		// Set in map
		c.Set(authInfoKey, authData)
	}
}

// StudentOnly will only allow students to access this endpoint.
// It must be called after JWTAuthMiddleware
func StudentOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.MustGet(authInfoKey).(AuthData).IsStaff {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{reasonKey: "students only!"})
			return
		}
	}
}

// StaffOnly will only allow staffs to access this endpoint.
// It must be called after JWTAuthMiddleware
func StaffOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !c.MustGet(authInfoKey).(AuthData).IsStaff {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{reasonKey: "staff only!"})
			return
		}
	}
}

// ParseEnrollmentBody will parse the body of a request into CourseEnrollmentRequest
func ParseEnrollmentBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request CourseEnrollmentRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{reasonKey: err.Error()})
			return
		}
		c.Set(requestKey, request)
	}
}
