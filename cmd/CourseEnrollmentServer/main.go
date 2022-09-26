package main

import (
	api "CourseEnrollment/api/CourseEnrollmentServer"
	pg "CourseEnrollment/internal/database"
	database "CourseEnrollment/internal/database/CourseEnrollmentServer"
	"CourseEnrollment/internal/shared"
	"CourseEnrollment/pkg/broker"
	"CourseEnrollment/pkg/course"
	"CourseEnrollment/pkg/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Connect to database and get initial data
	apiData := new(api.API)
	_, apiData.Courses, apiData.Students = getInitialData()
	// Connect to message broker
	var closeBroker func()
	apiData.Broker, closeBroker = setupMessageBroker()
	defer closeBroker()
	// Create the listener
	listener, err := net.Listen("tcp", os.Getenv("LISTEN_ADDRESS"))
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterCourseEnrollmentServerServiceServer(grpcServer, apiData)
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Cannot listen: %s", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Graceful shutdown initiated...")
	grpcServer.GracefulStop()
}

func getInitialData() (course.Departments, *course.Courses, map[course.StudentID]*course.Student) {
	// Check DB url
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalln("please set DATABASE_URL environment variable")
	}
	// Get the database url and connect to it
	db, err := pg.NewPostgresDatabase(dbURL)
	if err != nil {
		log.Fatalf("cannot connect to database: %s", err)
	}
	defer db.Close()
	pgDB := database.NewDatabase(db)
	// Fetch data
	departments, err := pgDB.GetDepartments()
	if err != nil {
		log.Fatalf("cannot get departments: %s", err)
	}
	courses, err := pgDB.GetCourses()
	if err != nil {
		log.Fatalf("cannot get courses: %s", err)
	}
	students, err := pgDB.GetStudents()
	if err != nil {
		log.Fatalf("cannot get students: %s", err)
	}
	// Done
	return departments, courses, students
}

// setupMessageBroker creates and connects to our message broker.
// The first returned value is the broker itself.
// The second value is a closer function. IT MUST BE CALLED BEFORE APPLICATION EXITS.
func setupMessageBroker() (course.Batcher, func()) {
	address := os.Getenv("RABBITMQ_ADDRESS")
	if address == "" {
		log.Fatalln("please set RABBITMQ_ADDRESS environment variable")
	}
	mq, err := broker.NewRabbitMQBroker(address, shared.CourseEnrollmentServerDatabaseQueueName)
	if err != nil {
		log.Fatalf("cannot instantiate the RabbitMQ client: %s", err)
	}
	return mq, func() {
		_ = mq.Close()
	}
}
