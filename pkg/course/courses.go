package course

import (
	"CourseEnrollment/pkg/proto"
	"CourseEnrollment/pkg/util"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

// CourseID is the type of the course's ID
type CourseID int32

// GroupID is the type of the group's ID used in Course
type GroupID uint8

// Courses held a list of all courses
type Courses struct {
	courses map[CourseID][]*Course
	mu      sync.RWMutex
}

// Course represents a single course
type Course struct {
	// The course ID
	ID CourseID
	// The group ID
	GroupID GroupID
	// What is the department of this course?
	Department DepartmentID
	// Who lectures this course?
	Lecturer string
	// Number of units
	Units uint8
	// The total capacity of this course
	Capacity int
	// List of students which have registered in this course
	RegisteredStudents map[StudentID]struct{}
	// Total number of students which can be in reserved queue
	ReserveCapacity int
	// The queue of reserved students
	ReserveQueue util.Queue[StudentID]
	// When is the exam of this course? In unix epoch (seconds)
	// If this value is zero, it means that there is no exam for this course
	ExamTime atomic.Int64
	// The time and days which class is held on
	ClassHeldTime ClassTime
	// Does this class has a sex lock?
	SexLock SexLock
	// The mutex to work with this course
	mu sync.RWMutex
}

// NewCourses creates Courses from its map
func NewCourses(courses map[CourseID][]*Course) *Courses {
	return &Courses{courses: courses}
}

// GetCourse will get the course based on group ID and course ID. If the course does not exist,
// it returns nil
func (c *Courses) GetCourse(courseID CourseID, groupID GroupID) *Course {
	var result *Course
	c.mu.RLock()
	for _, course := range c.courses[courseID] {
		if course.GroupID == groupID {
			result = course
			break
		}
	}
	c.mu.RUnlock()
	return result
}

// GetDepartmentCoursesProto gets all courses in a department
func (c *Courses) GetDepartmentCoursesProto(id DepartmentID) *proto.DepartmentCourses {
	c.mu.RLock()
	result := new(proto.DepartmentCourses)
	for _, courseWithSameGroups := range c.courses {
		if courseWithSameGroups[0].Department == id {
			for _, course := range courseWithSameGroups {
				result.Courses = append(result.Courses, course.ToProtoCourse())
			}
		}
	}
	c.mu.RUnlock()
	return result
}

// EnrollStudent will enroll the student in this course.
//
// NOTE: This method does not check for class conflicts or etc. It just adds a student to a course
// if possible. At first, it tries to add student to Course.RegisteredStudents and if it's not possible,
// (due to limitations), it ties to add they in Course.ReserveQueue. If that's also not possible,
// it returns false.
// It also does not check if the student is enrolled in this course or not
func (c *Course) EnrollStudent(ctx context.Context, studentID StudentID, batcher Batcher) (bool, error) {
	if batcher == nil {
		panic("nil batcher")
	}
	// Lock the course and unlock it when we are leaving
	c.mu.Lock()
	defer c.mu.Unlock()
	// Could not enroll the course
	return c.threadUnsafeEnrollStudent(ctx, studentID, batcher)
}

// threadUnsafeEnrollStudent does EnrollStudent but without locking the course
func (c *Course) threadUnsafeEnrollStudent(ctx context.Context, studentID StudentID, batcher Batcher) (bool, error) {
	// We check the space of this course and return early
	if !c.threadUnsafeCanBeEnrolled() {
		return false, nil
	}
	if batcher != nil {
		// We queue the message in batcher
		err := batcher.ProcessDatabaseQuery(ctx,
			c.Department,
			&proto.CourseDatabaseBatchMessage{
				Action: &proto.CourseDatabaseBatchMessage_Enroll{
					Enroll: &proto.CourseDatabaseBatchEnrollMessage{
						StudentId: uint64(studentID),
						CourseId:  int32(c.ID),
						GroupId:   uint32(c.GroupID),
						Reserved:  len(c.RegisteredStudents) == c.Capacity,
					},
				},
			})
		if err != nil {
			return false, BatchError{err}
		}
	}

	// At first check the registered count
	if len(c.RegisteredStudents) < c.Capacity {
		c.RegisteredStudents[studentID] = struct{}{}
		return true, nil
	}
	// Next check the reserve queue
	if c.ReserveQueue.Len() < c.ReserveCapacity {
		c.ReserveQueue.Enqueue(studentID)
		return true, nil
	}
	// Should never happen because we checked before
	panic("could not enroll in course because it was full and now the message is in the fucking broker")
}

// DisenrollStudent will remove the student from course.
//
// Will panic if the student is not enrolled in course.
func (c *Course) DisenrollStudent(ctx context.Context, studentID StudentID, batcher Batcher) error {
	if batcher == nil {
		panic("nil batcher")
	}
	// Lock the course and unlock it when we are leaving
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.threadUnsafeDisenrollStudent(ctx, studentID, batcher)
}

// threadUnsafeDisenrollStudent is basically DisenrollStudent but without locking
// the course
func (c *Course) threadUnsafeDisenrollStudent(ctx context.Context, studentID StudentID, batcher Batcher) error {
	if batcher != nil {
		// Put data in batcher
		err := batcher.ProcessDatabaseQuery(
			ctx,
			c.Department,
			&proto.CourseDatabaseBatchMessage{
				Action: &proto.CourseDatabaseBatchMessage_Disenroll{
					Disenroll: &proto.CourseDatabaseBatchDisenrollMessage{
						StudentId: uint64(studentID),
						CourseId:  int32(c.ID),
					},
				},
			})
		if err != nil {
			return BatchError{err}
		}
	}
	// Check registered list
	if _, registered := c.RegisteredStudents[studentID]; registered {
		delete(c.RegisteredStudents, studentID)
		// Now put first person from reserve queue into registered users
		// (if exists)
		if c.ReserveQueue.Len() != 0 {
			c.RegisteredStudents[c.ReserveQueue.Dequeue()] = struct{}{}
		}
		// Done
		return nil
	}
	// Otherwise remove from queue
	if !c.ReserveQueue.Remove(studentID) {
		panic(fmt.Sprintf("user %d has lesson %d-%d in their registered courses but lesson map does not have this user", studentID, c.ID, c.GroupID))
	}
	return nil
}

// ChangeGroupOfStudent tries to change the group of a student between two courses
func (c *Course) ChangeGroupOfStudent(ctx context.Context, studentID StudentID, other *Course, batcher Batcher) (bool, error) {
	if batcher == nil {
		panic("nil batcher")
	}
	// Check courses
	if c.ID != other.ID {
		panic("different courses provided")
	}
	if c.GroupID == other.GroupID {
		panic("same group provided")
	}
	// For locking, we at first lock the smaller group ID
	// to avoid deadlocks
	if c.GroupID < other.GroupID {
		c.mu.Lock()
		other.mu.Lock()
	} else {
		other.mu.Lock()
		c.mu.Lock()
	}
	defer c.mu.Unlock()
	defer other.mu.Unlock()
	// Check if user exists in this course
	if _, existsInRegistered := c.RegisteredStudents[studentID]; !existsInRegistered {
		if c.ReserveQueue.Exists(studentID) == -1 {
			panic(fmt.Sprintf("student %d does not exists in course %d-%d", studentID, c.ID, c.GroupID))
		}
	}
	// Check the capacity
	if !other.threadUnsafeCanBeEnrolled() {
		return false, nil
	}
	// Send data in batcher
	err := batcher.ProcessDatabaseQuery(
		ctx,
		c.Department,
		&proto.CourseDatabaseBatchMessage{
			Action: &proto.CourseDatabaseBatchMessage_ChangeGroup{
				ChangeGroup: &proto.CourseDatabaseBatchChangeGroupMessage{
					StudentId: uint64(studentID),
					CourseId:  int32(c.ID),
					GroupId:   uint32(other.GroupID),
				},
			},
		})
	if err != nil {
		return false, BatchError{err}
	}
	// Now try to add it to other course
	if ok, _ := other.threadUnsafeEnrollStudent(ctx, studentID, nil); !ok {
		panic("could not change group due to capacity and a is message in broker")
	}
	// Now remove the user from this course
	_ = c.threadUnsafeDisenrollStudent(ctx, studentID, nil) // no error because no batcher
	return true, nil
}

// ForceEnroll will forcibly enroll a student in this course.
// If this course had free space it will register the user directly in it.
// Otherwise, it will add capacity to course and then add the student
func (c *Course) ForceEnroll(ctx context.Context, studentID StudentID, batcher Batcher) error {
	if batcher == nil {
		panic("nil batcher")
	}
	// Lock the course and unlock it when we are leaving
	c.mu.Lock()
	defer c.mu.Unlock()
	// Check the capacity
	if len(c.RegisteredStudents) == c.Capacity {
		// Update capacity
		err := batcher.ProcessDatabaseQuery(
			ctx,
			c.Department,
			&proto.CourseDatabaseBatchMessage{
				Action: &proto.CourseDatabaseBatchMessage_UpdateCapacity{
					UpdateCapacity: &proto.CourseDatabaseBatchUpdateCapacity{
						CourseId:    int32(c.ID),
						GroupId:     uint32(c.GroupID),
						NewCapacity: int32(c.Capacity + 1),
					},
				},
			})
		if err != nil {
			return BatchError{err}
		}
		c.Capacity++
	}
	// Register user
	err := batcher.ProcessDatabaseQuery(ctx,
		c.Department,
		&proto.CourseDatabaseBatchMessage{
			Action: &proto.CourseDatabaseBatchMessage_Enroll{
				Enroll: &proto.CourseDatabaseBatchEnrollMessage{
					StudentId: uint64(studentID),
					CourseId:  int32(c.ID),
					GroupId:   uint32(c.GroupID),
					Reserved:  false,
				},
			},
		})
	if err != nil {
		return BatchError{err}
	}
	// Add user
	c.RegisteredStudents[studentID] = struct{}{}
	return nil
}

// threadUnsafeCanBeEnrolled checks if one student can enroll in a course.
// It doesn't lock anything, so it's thread unsafe.
func (c *Course) threadUnsafeCanBeEnrolled() bool {
	return len(c.RegisteredStudents) < c.Capacity || c.ReserveQueue.Len() < c.ReserveCapacity
}

// GetStudentQueuePosition gets the position of a user in queue.
// It returns false as the second argument if user is not registered at all in this course.
// For the first argument, it returns 0 if user is in the registered users. Otherwise, it returns
// the index of user in queue (1-indexed)
//
// This method is not thread safe
func (c *Course) getStudentQueuePosition(id StudentID) (uint, bool) {
	// Check normal registered users
	if _, exists := c.RegisteredStudents[id]; exists {
		return 0, true
	}
	// Check normal queue
	index := c.ReserveQueue.Exists(id)
	if index == -1 {
		return 0, false
	}
	return uint(index + 1), true
}

// ToProtoCourse converts current course to proto.CourseData
func (c *Course) ToProtoCourse() *proto.CourseData {
	c.mu.RLock()
	result := c.threadUnsafeToProtoCourse()
	c.mu.RUnlock()
	return result
}

// threadUnsafeToProtoCourse does not lock the mutex in course and
// gets the protobuf representation of this Course
func (c *Course) threadUnsafeToProtoCourse() *proto.CourseData {
	result := &proto.CourseData{
		CourseId:        int32(c.ID),
		GroupId:         uint32(c.GroupID),
		Units:           uint32(c.Units),
		Capacity:        int32(c.Capacity),
		RegisteredCount: uint32(len(c.RegisteredStudents)),
		ExamTime:        c.ExamTime.Load(),
		Lecturer:        c.Lecturer,
	}
	// Get the class time
	days, start, end := c.ClassHeldTime.Get()
	result.ClassTime = make([]*proto.ClassTime, len(days))
	for i := range days {
		result.ClassTime[i] = &proto.ClassTime{
			Day:         proto.Weekday(days[i]),
			StartMinute: uint32(start.t),
			EndMinute:   uint32(end.t),
		}
	}
	return result
}

// ToStudentCourseDataProto gets the data which needs to be sent when user wants to see their
// enrolled courses.
//
// Passing a student ID which is not enrolled in this course causes this method to panic
func (c *Course) ToStudentCourseDataProto(std StudentID) *proto.StudentCourseData {
	c.mu.RLock()
	position, ok := c.getStudentQueuePosition(std)
	if !ok {
		c.mu.RUnlock()
		panic(fmt.Sprintf("requested student course data of a %d which is not registered in course %d-%d", std, c.ID, c.GroupID))
	}
	result := &proto.StudentCourseData{
		Course:               c.threadUnsafeToProtoCourse(),
		ReserveQueuePosition: uint32(position),
	}
	c.mu.RUnlock()
	return result
}

// ToStudentsOfCourseResponseProto gets all students enrolled in this course
// including the ones in reserve queue.
func (c *Course) ToStudentsOfCourseResponseProto() *proto.StudentsOfCourseResponse {
	c.mu.RLock()
	result := &proto.StudentsOfCourseResponse{
		RegisteredStudents:    make([]uint64, 0, len(c.RegisteredStudents)),
		ReservedQueueStudents: make([]uint64, c.ReserveQueue.Len()),
	}
	// Add main users
	for std := range c.RegisteredStudents {
		result.RegisteredStudents = append(result.RegisteredStudents, uint64(std))
	}
	// Add reserve queue
	queue := c.ReserveQueue.CopyAsArray()
	for i, std := range queue {
		result.ReservedQueueStudents[i] = uint64(std)
	}
	c.mu.RUnlock()
	return result
}

// UpdateCapacity will update the courses
func (c *Course) UpdateCapacity(ctx context.Context, newCapacity int, batcher Batcher) error {
	if batcher == nil {
		panic("nil batcher")
	}
	// Lock to update the course
	c.mu.Lock()
	defer c.mu.Unlock()
	// Check if new capacity is less than registered amount
	if len(c.RegisteredStudents) > newCapacity {
		return LowerCapacityThanRegistered
	}
	// Check if new capacity is old capacity
	if c.Capacity == newCapacity {
		return nil // do nothing
	}
	// Get the users which are going to be moved from reserve queue to main registered users.
	// We must add to registered users for Min(capacity difference, reserve queue len) times.
	// We also get the max with 0 to avoid panics when we are reducing the capacity.
	reservedMovedUsers := c.ReserveQueue.CopyAsArray()[:util.Max(util.Min(newCapacity-c.Capacity, c.ReserveQueue.Len()), 0)]
	reservedMovedUsersUint := make([]uint64, len(reservedMovedUsers))
	for i, v := range reservedMovedUsers {
		reservedMovedUsersUint[i] = uint64(v)
	}
	// Batch it
	err := batcher.ProcessDatabaseQuery(
		ctx,
		c.Department,
		&proto.CourseDatabaseBatchMessage{
			Action: &proto.CourseDatabaseBatchMessage_UpdateCapacity{
				UpdateCapacity: &proto.CourseDatabaseBatchUpdateCapacity{
					CourseId:      int32(c.ID),
					GroupId:       uint32(c.GroupID),
					NewCapacity:   int32(newCapacity),
					MovedStudents: reservedMovedUsersUint,
				},
			},
		})
	if err != nil {
		return BatchError{err}
	}
	// Remove from queue and add to main registered users.
	for i := len(reservedMovedUsers); i > 0; i-- {
		c.RegisteredStudents[c.ReserveQueue.Dequeue()] = struct{}{}
	}
	// Update the capacity
	c.Capacity = newCapacity
	// Done
	return nil
}
