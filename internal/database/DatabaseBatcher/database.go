package DatabaseBatcher

import (
	"CourseEnrollment/pkg/course"
	"context"
	"github.com/go-faster/errors"
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

// UpdateCapacity will update the capacity of a course
func (db Database) UpdateCapacity(courseID course.CourseID, groupID course.GroupID, newCapacity int32, movedStudents []uint64) error {
	// Start a transaction
	tx, err := db.db.Begin(context.Background())
	if err != nil {
		return errors.Wrap(err, "cannot start transaction")
	}
	defer tx.Rollback(context.Background())
	// Put people from reserve into main class capacity if needed
	if len(movedStudents) != 0 {
		_, err = tx.Exec(context.Background(), "UPDATE enrolled_courses SET reserved=FALSE WHERE id = ANY($1)", movedStudents)
		if err != nil {
			return errors.Wrap(err, "cannot update reserved status")
		}
	}
	// Update the course capacity
	_, err = tx.Exec(context.Background(), "UPDATE courses SET capacity=$1 WHERE course_id=$2 AND group_id=$3", newCapacity, courseID, groupID)
	if err != nil {
		return errors.Wrap(err, "cannot update course capacity")
	}
	// Done
	err = tx.Commit(context.Background())
	if err != nil {
		return errors.Wrap(err, "cannot commit")
	}
	return nil
}

// Close will close the connection
func (db Database) Close() {
	db.db.Close()
}
