package CourseEnrollmentServer

import (
	"CourseEnrollment/pkg/course"
	"database/sql"
)

type Database struct {
	db *sql.DB
}

// NewDatabase creates a database for accessing the database
// which is just used in loading the courses for server startup
func NewDatabase(db *sql.DB) Database {
	return Database{db}
}

// GetStudents will get all students in the database as a map.
func (db *Database) GetStudents() (map[course.StudentID]*course.Student, error) {
	// TODO
	return nil, nil
}

func (db *Database) GetCourses() (*course.Courses, error) {
	// TODO
	return nil, nil
}
