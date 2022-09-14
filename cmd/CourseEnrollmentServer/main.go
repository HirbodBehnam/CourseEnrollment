package main

import (
	api "CourseEnrollment/api/CourseEnrollmentServer"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// TODO: connect to database and get data
	apiData := new(api.API)
	// Create the middleware
	router := gin.Default()
	studentRoutes := router.Group("/student/:stdID/:courseID")
	studentRoutes.PUT("", apiData.StudentEnroll)
	studentRoutes.DELETE("", apiData.StudentDisenroll)
	studentRoutes.PATCH("", apiData.StudentChangeGroup)
	// Listen
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Cannot listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Graceful shutdown initiated...")
	_ = srv.Shutdown(context.Background())
}
