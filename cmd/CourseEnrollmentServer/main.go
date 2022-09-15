package main

import (
	api "CourseEnrollment/api/CourseEnrollmentServer"
	"CourseEnrollment/pkg/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// TODO: connect to database and get data
	apiData := new(api.API)
	// Create the listener
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterCourseEnrollmentServerServiceServer(grpcServer, apiData)
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Cannot listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Graceful shutdown initiated...")
	grpcServer.GracefulStop()
}
