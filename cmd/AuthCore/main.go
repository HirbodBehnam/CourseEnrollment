package main

import (
	api "CourseEnrollment/api/AuthCore"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create the API
	endpointApi := new(api.API)
	endpointApi.GenerateJWTKey()

	// Setup endpoints
	r := gin.Default()
	// Login and logout
	r.POST("/login", endpointApi.LoginUser)

}
