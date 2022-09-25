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
			c.JSON(http.StatusUnauthorized, gin.H{"reason": "empty auth"})
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
			c.JSON(http.StatusUnauthorized, gin.H{"reason": "invalid auth"})
			return
		}
		claims, ok := token.Claims.(*JWTToken)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"reason": "invalid auth"})
			return
		}
		// Set the data
		authData := AuthData{
			Department: claims.Department,
			IsStaff:    claims.IsStaff,
		}
		authData.User, err = strconv.ParseUint(claims.Subject, 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"reason": "invalid auth"})
			return
		}
		// Set in map
		c.Set(authInfoKey, authData)
		// Continue
		c.Next()
	}
}
