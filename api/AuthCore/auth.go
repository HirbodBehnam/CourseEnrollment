package AuthCore

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

// LoginUser must check the authentication credentials of user and login it
func (a *API) LoginUser(c *gin.Context) {
	// Get username and password
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Check them
	userOk, departmentID, err := a.Database.AuthUser(c.Request.Context(), request.User, request.Password, request.IsStaff)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.WithError(err).WithField("request", request).Error("cannot check user credentials")
		return
	}
	// Check auth info
	if !userOk {
		c.Status(http.StatusUnauthorized)
		return
	}
	// Create the JWT
	now := time.Now()
	token, err := jwt.NewWithClaims(signingMethod, JWTToken{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtIssuer,
			Subject:   strconv.FormatUint(request.User, 10),
			ExpiresAt: jwt.NewNumericDate(now.Add(jwtTTL)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Department: departmentID,
		IsStaff:    request.IsStaff,
	}).SignedString(a.jwtKey)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.WithError(err).Error("cannot sign the jwt")
		return
	}
	// Send back the result
	c.JSON(http.StatusOK, TokenResult{token})
}
