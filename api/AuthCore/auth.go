package AuthCore

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
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
	token, err := createJWTToken(a.jwtKey, request.User, departmentID, request.IsStaff)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.WithError(err).Error("cannot sign the jwt")
		return
	}
	// Send back the result
	c.JSON(http.StatusOK, TokenResult{token})
}

// RefreshJWTToken refreshes the JWT token of a user
func (a *API) RefreshJWTToken(c *gin.Context) {
	// Get auth data
	auth := c.MustGet(authInfoKey).(AuthData)
	// Sign again
	token, err := createJWTToken(a.jwtKey, auth.User, auth.Department, auth.IsStaff)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.WithError(err).Error("cannot sign the jwt")
		return
	}
	// Send back the result
	c.JSON(http.StatusOK, TokenResult{token})
}
