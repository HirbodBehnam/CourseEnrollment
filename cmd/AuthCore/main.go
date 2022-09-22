package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	// Login and logout
	r.POST("/login")
}
