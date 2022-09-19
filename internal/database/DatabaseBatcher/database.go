package DatabaseBatcher

import (
	"CourseEnrollment/pkg/course"
	"context"
	"github.com/jackc/pgx/v4"
)

type Database struct {
	db *pgx.Conn
}

// EnrollCourse will enroll a student in a course
func (db Database) EnrollCourse(stdID course.StudentID, courseID course.CourseID, groupID course.GroupID) error {
	_, err := db.db.Exec(context.Background(), "INSERT INTO enrolled_courses (course_id, group_id, student_id) VALUES ($1, $2, $3)", courseID, groupID, stdID)
	return err
}

// DisenrollCourse will disenroll a student in a course
func (db Database) DisenrollCourse(stdID course.StudentID, courseID course.CourseID, groupID course.GroupID) error {
	_, err := db.db.Exec(context.Background(), "DELETE FROM enrolled_courses WHERE course_id=$1 AND group_id=$2 AND student_id=$3", courseID, groupID, stdID)
	return err
}

// ChangeCourseGroup will change the group of a user in an enrolled course
func (db Database) ChangeCourseGroup(stdID course.StudentID, courseID course.CourseID, newGroupID course.GroupID) error {
	_, err := db.db.Exec(context.Background(), "UPDATE enrolled_courses SET group_id=$1 WHERE course_id=$2 AND student_id=$3", newGroupID, courseID, stdID)
	return err
}
