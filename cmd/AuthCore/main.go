package main

import (
	api "CourseEnrollment/api/AuthCore"
	pg "CourseEnrollment/internal/database"
	db "CourseEnrollment/internal/database/AuthCore"
	pb "CourseEnrollment/pkg/proto"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
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
	// Setup the gRPC client
	var coreConnCloser func()
	endpointApi.CoreClient, coreConnCloser = setupGRPCClient()
	defer coreConnCloser()
	// Setup endpoints
	r := gin.New()
	r.Use(gin.Recovery())
	// Login and token refresh
	r.POST("/login", endpointApi.LoginUser)
	r.POST("/refresh", endpointApi.JWTAuthMiddleware(), endpointApi.RefreshJWTToken)
	// Student endpoints
	studentRouter := r.Group("/student", endpointApi.JWTAuthMiddleware(), api.StudentOnly())
	studentRouter.PUT("/course", api.ParseEnrollmentBody(), endpointApi.EnrollStudent)
	studentRouter.PATCH("/course", api.ParseEnrollmentBody(), endpointApi.ChangeGroupOfStudent)
	studentRouter.DELETE("/course", endpointApi.DisenrollStudent)
	studentRouter.GET("/course", endpointApi.EnrolledCoursesOfStudent)
	studentRouter.GET("/courses", endpointApi.CoursesOfDepartment)
	// Admin endpoints
	staffRouter := r.Group("/staff", endpointApi.JWTAuthMiddleware(), api.StaffOnly())
	staffRouter.PUT("/force-std", endpointApi.ForceEnroll)
	staffRouter.DELETE("/force-std", endpointApi.ForceDisenroll)
	staffRouter.GET("/student-courses", endpointApi.CoursesOfStudent)
	staffRouter.GET("/course-students", endpointApi.StudentsOfCourse)
	staffRouter.PATCH("/capacity", endpointApi.UpdateCourseCapacity)
	// Listen
	srv := &http.Server{
		Handler: r,
	}
	go func() {
		if err := srv.Serve(getListener()); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
		log.Fatal("please set DATABASE_URL environment variable")
	}
	// Get the database url and connect to it
	database, err := pg.NewPostgresDatabase(dbURL)
	if err != nil {
		log.Fatalf("cannot connect to database: %s\n", err)
	}
	return db.NewDatabase(database)
}

// setupGRPCClient will set up the grpc client for core.
// The function returned is the closer function which closes the
func setupGRPCClient() (pb.CourseEnrollmentServerServiceClient, func()) {
	// Get address
	coreAddress := os.Getenv("CORE_ADDRESS")
	if coreAddress == "" {
		log.Fatal("please set CORE_ADDRESS environment variable")
	}
	// Connect
	conn, err := grpc.Dial(coreAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return pb.NewCourseEnrollmentServerServiceClient(conn), func() {
		_ = conn.Close()
	}
}

// getListener will start a listener based on environment variables
func getListener() net.Listener {
	// Get protocol
	protocol := "tcp"
	if envProtocol := os.Getenv("LISTEN_PROTOCOL"); envProtocol != "" {
		protocol = envProtocol
	}
	// Listen
	listener, err := net.Listen(protocol, os.Getenv("LISTEN_ADDRESS"))
	if err != nil {
		log.Fatalf("cannot listen: %s", err)
	}
	return listener
}
