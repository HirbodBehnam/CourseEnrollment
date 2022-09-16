package CourseEnrollmentServer

import (
	"CourseEnrollment/pkg/course"
	"CourseEnrollment/pkg/util"
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"time"
)

type Database struct {
	db *pgx.Conn
}

// NewDatabase creates a database for accessing the database
// which is just used in loading the courses for server startup
func NewDatabase(db *pgx.Conn) Database {
	return Database{db}
}

// GetDepartments will get the list of departments from database.
func (db *Database) GetDepartments() (course.Departments, error) {
	rows, err := db.db.Query(context.Background(), "SELECT id, name FROM departments")
	if err != nil {
		return nil, errors.Wrap(err, "cannot query departments")
	}
	defer rows.Close()
	// Reach each row
	result := make(course.Departments)
	for rows.Next() {
		var departmentID course.DepartmentID
		var departmentName string
		err = rows.Scan(&departmentID, &departmentName)
		if err != nil {
			return nil, errors.Wrap(err, "cannot scan")
		}
		result[departmentID] = departmentName
	}
	return result, nil
}

// GetCourses will get the list of courses from database.
// It also updates the courses registered list.
func (db *Database) GetCourses() (*course.Courses, error) {
	rows, err := db.db.Query(context.Background(), "SELECT course_id, group_id, for_department, units, capacity, reserve_capacity, exam_time, class_time, sex_lock FROM courses")
	if err != nil {
		return nil, errors.Wrap(err, "cannot query courses")
	}
	defer rows.Close()
	// Get results
	result := make(map[course.CourseID][]*course.Course)
	for rows.Next() {
		// Create the course
		currentCourse := new(course.Course)
		var examTime time.Time
		err = rows.Scan(&currentCourse.ID, &currentCourse.GroupID, &currentCourse.Department, &currentCourse.Units, &currentCourse.Capacity, &currentCourse.ReserveCapacity, &currentCourse,
			&examTime, &currentCourse.ClassHeldTime, &currentCourse.SexLock)
		if err != nil {
			return nil, errors.Wrap(err, "cannot scan course")
		}
		// Update some missing info based on scanned ones
		currentCourse.ReserveQueue = util.NewQueue[course.StudentID]()
		currentCourse.RegisteredStudents = make(map[course.StudentID]struct{}, currentCourse.Capacity)
		currentCourse.ExamTime.Store(examTime.Unix())
		// Get the courses registered list
		err = db.updateCourseRegistered(currentCourse)
		if err != nil {
			return nil, errors.Wrap(err, "cannot set registered users")
		}
		// Insert it into map
		result[currentCourse.ID] = append(result[currentCourse.ID], currentCourse)
	}
	return course.NewCourses(result), nil
}

// updateCourseRegistered updates the registered users in the
func (db *Database) updateCourseRegistered(c *course.Course) error {
	rows, err := db.db.Query(context.Background(), "SELECT student_id FROM enrolled_courses WHERE course_id=? AND group_id=? ORDER BY id")
	if err != nil {
		return errors.Wrapf(err, "cannot query course %d-%d", c.ID, c.GroupID)
	}
	defer rows.Close()
	// Get the students
	for rows.Next() {
		var stdID course.StudentID
		err = rows.Scan(&stdID)
		if err != nil {
			return errors.Wrap(err, "cannot scan row")
		}
		// Add to course
		if len(c.RegisteredStudents) == c.Capacity {
			// They will be queued in order
			c.ReserveQueue.Enqueue(stdID)
		} else {
			c.RegisteredStudents[stdID] = struct{}{}
		}
	}
	return nil
}

// GetStudents will get all students in the database as a map.
func (db *Database) GetStudents() (map[course.StudentID]*course.Student, error) {
	// Get all students
	rows, err := db.db.Query(context.Background(), "SELECT id, enrollment_start_time, max_units, remaining_actions, gender FROM students")
	if err != nil {
		return nil, errors.Wrap(err, "cannot query students")
	}
	defer rows.Close()
	// Get them
	result := make(map[course.StudentID]*course.Student)
	for rows.Next() {
		student := new(course.Student)
		var enrollmentStartTime time.Time
		err = rows.Scan(&student.ID, &enrollmentStartTime, &student.MaxUnits, &student.RemainingActions, &student.StudentSex)
		if err != nil {
			return nil, errors.Wrap(err, "cannot scan row")
		}
		student.RegisteredCourses, student.RegisteredUnits, err = db.getEnrolledCoursesOfStudent(student.ID)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot get student registered courses for %d", student.ID)
		}
		// Add to map
		result[student.ID] = student
	}
	return nil, nil
}

// getEnrolledCoursesOfStudent gets the list of enrolled (reserved and registered) of a user
// and the number of courses
func (db *Database) getEnrolledCoursesOfStudent(stdID course.StudentID) (map[course.CourseID]course.GroupID, uint8, error) {
	rows, err := db.db.Query(context.Background(), "SELECT enrolled_courses.course_id, enrolled_courses.group_id, c.units FROM enrolled_courses JOIN courses c on c.course_id = enrolled_courses.course_id and c.group_id = enrolled_courses.group_id WHERE student_id=?", stdID)
	if err != nil {
		return nil, 0, errors.Wrap(err, "cannot query")
	}
	defer rows.Close()
	// Get them
	var totalUnits uint8
	courses := make(map[course.CourseID]course.GroupID)
	for rows.Next() {
		var courseID course.CourseID
		var groupID course.GroupID
		var units uint8
		err = rows.Scan(&courseID, &groupID, &units)
		if err != nil {
			return nil, 0, errors.Wrap(err, "cannot scan")
		}
		// Apply
		totalUnits += units
		courses[courseID] = groupID
	}
	return courses, totalUnits, nil
}
