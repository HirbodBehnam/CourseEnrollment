package main

import (
	"context"
	"encoding/hex"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func setupDatabase(databaseUrl string, studentCount uint64) map[uint64]*student {
	// Create the map and fill it
	log.Info("creating users")
	students := make(map[uint64]*student, studentCount)
	for i := uint64(1); i <= studentCount; i++ {
		students[i] = &student{id: i, password: randomPassword()}
	}
	// Connect to database
	log.Info("connecting to database")
	conn, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		log.WithError(err).Fatal("cannot connect to database")
	}
	defer conn.Close(context.Background())
	// Clear tables and start a transaction
	_, err = conn.Exec(context.Background(), "TRUNCATE TABLE students, enrolled_courses")
	if err != nil {
		log.WithError(err).Fatal("cannot truncate tables")
	}
	tx, err := conn.Begin(context.Background())
	if err != nil {
		log.WithError(err).Fatal("cannot start transaction")
	}
	// Add users
	log.Info("registering users")
	for userID, data := range students {
		addUser(tx, userID, data.password)
	}
	// Commit the transaction
	log.Info("applying into database")
	err = tx.Commit(context.Background())
	if err != nil {
		log.WithError(err).Fatal("cannot commit transaction")
	}
	// Load courses
	log.Info("getting courses")
	courses := loadCourses(conn)
	// Assign students to courses
	log.Info("creating random courses")
	for _, std := range students {
		for i := 0; i < 7; i++ {
			std.toPickCourses = append(std.toPickCourses, courses[rng.Intn(len(courses))])
		}
	}
	// Return created users
	return students
}

func addUser(conn pgx.Tx, userID uint64, password string) {
	// Bcrypt the password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	// Insert it into database
	_, err := conn.Exec(context.Background(), "INSERT INTO students (id, password, enrollment_start_time, max_units, remaining_actions, department_id, entry_year, gender) VALUES ($1,$2,now(),255,12,40,99,'male')",
		userID, string(hashedPassword))
	if err != nil {
		log.WithError(err).WithField("userID", userID).Fatal("cannot create user")
	}
}

func loadCourses(conn *pgx.Conn) []course {
	// Do the query
	courses := make([]course, 0)
	rows, err := conn.Query(context.Background(), "SELECT course_id, group_id FROM courses")
	if err != nil {
		log.WithError(err).Fatal("cannot query courses")
	}
	defer rows.Close()
	// Read results
	for rows.Next() {
		var course course
		err = rows.Scan(&course.courseID, &course.groupID)
		if err != nil {
			log.WithError(err).Fatal("cannot scan row")
		}
		courses = append(courses, course)
	}
	return courses
}

func randomPassword() string {
	b := make([]byte, 4)
	rng.Read(b)
	return hex.EncodeToString(b)
}
