package AuthCore

import "github.com/gin-gonic/gin"

// ForceEnroll will forcibly enroll a student in a course.
// It adds capacity to course if needed. The only possible way for this endpoint
// to fail is that user has this course already.
func (a *API) ForceEnroll(c *gin.Context) {
	// TODO
}

// ForceDisenroll will simply disenroll a student from course.
// It fails if student is not enrolled in the course.
func (a *API) ForceDisenroll(c *gin.Context) {
	// TODO
}

// CoursesOfStudent gets the list of courses which user is enrolled in
// or is in queue.
func (a *API) CoursesOfStudent(c *gin.Context) {
	// TODO
}

// StudentsOfCourse gets all the students enrolled in a course.
func (a *API) StudentsOfCourse(c *gin.Context) {
	// TODO
}
