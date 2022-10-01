package main

import (
	pg "CourseEnrollment/internal/database"
	db "CourseEnrollment/internal/database/DatabaseBatcher"
	"CourseEnrollment/internal/shared"
	"CourseEnrollment/pkg/broker"
	"CourseEnrollment/pkg/course"
	"CourseEnrollment/pkg/proto"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

const consumerName = "course-enrollment-database-batcher"

func main() {
	database := setupDatabase()
	defer database.Close()
	mqBroker := setupMessageBroker()
	defer mqBroker.Close()
	// Listen to changes
	data, err := mqBroker.Consume(consumerName)
	if err != nil {
		log.WithError(err).Fatalf("cannot consumer queue")
	}
	// Setup the signal trap
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down...")
		err := mqBroker.CancelConsumer(consumerName)
		if err != nil {
			log.WithError(err).Warn("cannot cancel consumer")
		}
	}()
	// Loop over them
	for query := range data {
		processQuery(database, query)
	}
	// Done
	log.Info("clean shutdown")
}

// processQuery will apply the query in database
func processQuery(database db.Database, query *proto.CourseDatabaseBatchMessage) {
	var err error
	switch query.Action {
	case proto.CourseDatabaseBatchAction_Enroll:
		err = database.EnrollCourse(course.StudentID(query.StudentId), course.CourseID(query.CourseId), course.GroupID(query.GroupId))
	case proto.CourseDatabaseBatchAction_Disenroll:
		err = database.DisenrollCourse(course.StudentID(query.StudentId), course.CourseID(query.CourseId))
	case proto.CourseDatabaseBatchAction_ChangeGroup:
		err = database.ChangeCourseGroup(course.StudentID(query.StudentId), course.CourseID(query.CourseId), course.GroupID(query.GroupId))
	default:
		log.WithField("query", query).Error("invalid action")
		return
	}
	if err != nil {
		log.WithField("query", query).WithError(err).Error("cannot apply action")
	} else {
		log.WithField("query", query).Debug("applied")
	}
}

// setupDatabase will set up the database which is used to apply the courses
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

// setupMessageBroker creates and connects to our message broker
func setupMessageBroker() broker.RabbitMQBroker {
	address := os.Getenv("RABBITMQ_ADDRESS")
	if address == "" {
		log.Fatal("please set RABBITMQ_ADDRESS environment variable")
	}
	mq, err := broker.NewRabbitMQBroker(address, shared.CourseEnrollmentServerDatabaseQueueName)
	if err != nil {
		log.Fatalf("cannot instantiate the RabbitMQ client: %s", err)
	}
	return mq
}
