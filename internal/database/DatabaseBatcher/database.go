package DatabaseBatcher

import (
	"CourseEnrollment/pkg/course"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	db *pgxpool.Pool
}

// NewDatabase creates a database for accessing the database
// which is just used in loading the courses for server startup
func NewDatabase(db *pgxpool.Pool) Database {
	return Database{db}
}

// EnrollCourse will enroll a student in a course
func (db Database) EnrollCourse(stdID course.StudentID, courseID course.CourseID, groupID course.GroupID, reserved bool) error {
	_, err := db.db.Exec(context.Background(), "INSERT INTO enrolled_courses (course_id, group_id, student_id, reserved) VALUES ($1, $2, $3, $4)", courseID, groupID, stdID, reserved)
	return err
}

// DisenrollCourse will disenroll a student in a course
func (db Database) DisenrollCourse(stdID course.StudentID, courseID course.CourseID) error {
	_, err := db.db.Exec(context.Background(), "DELETE FROM enrolled_courses WHERE course_id=$1 AND student_id=$2", courseID, stdID)
	return err
}

// ChangeCourseGroup will change the group of a user in an enrolled course
func (db Database) ChangeCourseGroup(stdID course.StudentID, courseID course.CourseID, newGroupID course.GroupID, reserved bool) error {
	_, err := db.db.Exec(context.Background(), "UPDATE enrolled_courses SET group_id=$1, reserved=$2 WHERE course_id=$3 AND student_id=$4", newGroupID, reserved, courseID, stdID)
	return err
}

// Close will close the connection
func (db Database) Close() {
	db.db.Close()
}
