package main

import (
	api "CourseEnrollment/api/AuthCore"
	pg "CourseEnrollment/internal/database"
	db "CourseEnrollment/internal/database/AuthCore"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create the API
	endpointApi := new(api.API)
	endpointApi.GenerateJWTKey()
	endpointApi.Database = setupDatabase()
	defer endpointApi.Database.Close()
	// Setup endpoints
	r := gin.Default()
	// Login and logout
	r.POST("/login", endpointApi.LoginUser)
	// Listen
	srv := &http.Server{
		Addr:    os.Getenv("LISTEN_ADDRESS"),
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	_ = srv.Shutdown(context.Background())
}

func setupDatabase() db.Database {
	// Check DB url
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalln("please set DATABASE_URL environment variable")
	}
	// Get the database url and connect to it
	database, err := pg.NewPostgresDatabase(dbURL)
	if err != nil {
		log.Fatalf("cannot connect to database: %s\n", err)
	}
	return db.NewDatabase(database)
}
