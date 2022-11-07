package course

import (
	"CourseEnrollment/pkg/proto"
	"CourseEnrollment/pkg/util"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCourseEnrollStudent(t *testing.T) {
	t.Run("general test", func(t *testing.T) {
		assertion := assert.New(t)
		course := Course{
			Capacity:           10,
			RegisteredStudents: make(map[StudentID]struct{}),
			ReserveCapacity:    4,
			ReserveQueue:       util.NewQueue[StudentID](),
		}
		// Add to registered courses
		for i := 0; i < course.Capacity; i++ {
			assertion.True(course.EnrollStudent(context.Background(), StudentID(i), noOpBatcher{}))
		}
		// Check
		assertion.Len(course.RegisteredStudents, course.Capacity)
		for i := 0; i < course.Capacity; i++ {
			_, exists := course.RegisteredStudents[StudentID(i)]
			assertion.True(exists)
		}
		assertion.Equal(0, course.ReserveQueue.Len())
		// Add to reserve queue
		for i := 0; i < course.ReserveCapacity; i++ {
			assertion.True(course.EnrollStudent(context.Background(), StudentID(course.ReserveCapacity+i), noOpBatcher{}))
		}
		assertion.Equal(course.ReserveCapacity, course.ReserveQueue.Len())
		for i := 0; i < course.ReserveCapacity; i++ {
			index := course.ReserveQueue.Exists(StudentID(course.ReserveCapacity + i))
			assertion.Equal(i, index)
		}
		// Add over capacity
		for i := 0; i < 100; i++ {
			assertion.False(course.EnrollStudent(context.Background(), StudentID(course.ReserveCapacity+course.Capacity+i), noOpBatcher{}))
		}
	})
	// Batcher
	t.Run("test batcher", func(t *testing.T) {
		assertion := assert.New(t)
		course := Course{
			ID:                 CourseID(2),
			GroupID:            GroupID(3),
			Department:         DepartmentID(4),
			Capacity:           2,
			RegisteredStudents: make(map[StudentID]struct{}),
			ReserveCapacity:    2,
			ReserveQueue:       util.NewQueue[StudentID](),
		}
		batcher := new(inMemoryBatcher)
		var expected []struct {
			data *proto.CourseDatabaseBatchMessage
			dep  DepartmentID
		}
		for i := 0; i < 4; i++ {
			assertion.True(course.EnrollStudent(context.Background(), StudentID(i), batcher))
			expected = append(expected, struct {
				data *proto.CourseDatabaseBatchMessage
				dep  DepartmentID
			}{
				data: &proto.CourseDatabaseBatchMessage{
					Action: &proto.CourseDatabaseBatchMessage_Enroll{
						Enroll: &proto.CourseDatabaseBatchEnrollMessage{
							StudentId: uint64(i),
							CourseId:  2,
							GroupId:   3,
							Reserved:  i >= 2,
						},
					},
				},
				dep: DepartmentID(4),
			})
		}
		for i := 0; i < 4; i++ {
			assertion.False(course.EnrollStudent(context.Background(), StudentID(i+4), batcher))
		}
		assertion.Equal(expected, batcher.messages)
	})
	t.Run("batcher error", func(t *testing.T) {
		assertion := assert.New(t)
		course := Course{
			ID:                 CourseID(2),
			GroupID:            GroupID(3),
			Department:         DepartmentID(4),
			Capacity:           2,
			RegisteredStudents: make(map[StudentID]struct{}),
			ReserveCapacity:    2,
			ReserveQueue:       util.NewQueue[StudentID](),
		}
		innerError := errors.New("test")
		batcher := errorBatcher{innerError}
		ok, err := course.EnrollStudent(context.Background(), StudentID(1), batcher)
		assertion.False(ok)
		assertion.ErrorIs(err, BatchError{innerError})
		assertion.Len(course.RegisteredStudents, 0)
	})
	t.Run("nil batcher", func(t *testing.T) {
		assert.PanicsWithValue(t, "nil batcher", func() {
			_, _ = new(Course).EnrollStudent(context.Background(), StudentID(1), nil)
		})
	})
}

func TestCourseUnrollStudent(t *testing.T) {
	t.Run("general test", func(t *testing.T) {
		assertion := assert.New(t)
		course := Course{
			ID:                 CourseID(1),
			GroupID:            GroupID(1),
			Capacity:           10,
			RegisteredStudents: make(map[StudentID]struct{}),
			ReserveCapacity:    4,
			ReserveQueue:       util.NewQueue[StudentID](),
		}
		// At first, we just run unroll on empty course
		assertion.PanicsWithValue("user 1 has lesson 1-1 in their registered courses but lesson map does not have this user", func() {
			_ = course.DisenrollStudent(context.Background(), StudentID(1), noOpBatcher{})
		})
		// Then, we add some users to registered user. We don't go to reserved capacity
		for i := 0; i < 5; i++ {
			assertion.True(course.EnrollStudent(context.Background(), StudentID(i), noOpBatcher{}))
		}
		// Then we unroll the first student
		assertion.NotPanics(func() {
			_ = course.DisenrollStudent(context.Background(), StudentID(0), noOpBatcher{})
		})
		assertion.Len(course.RegisteredStudents, 4) // zero must be removed
		for i := 1; i < 5; i++ {
			_, exists := course.RegisteredStudents[StudentID(i)]
			assertion.True(exists)
		}
		// Then we add them again
		course.RegisteredStudents = make(map[StudentID]struct{})
		for i := 0; i < course.Capacity; i++ {
			assertion.True(course.EnrollStudent(context.Background(), StudentID(i), noOpBatcher{}))
		}
		for i := 0; i < course.ReserveCapacity; i++ {
			assertion.True(course.EnrollStudent(context.Background(), StudentID(course.Capacity+i), noOpBatcher{}))
		}
		assertion.NotPanics(func() {
			_ = course.DisenrollStudent(context.Background(), StudentID(0), noOpBatcher{})
		})
		// Check it
		for i := 1; i < course.Capacity+1; i++ {
			_, exists := course.RegisteredStudents[StudentID(i)]
			assertion.True(exists)
		}
		// Check reserved
		{
			expectedReserved := make([]StudentID, 0, course.ReserveCapacity)
			for i := 1; i < course.ReserveCapacity; i++ {
				expectedReserved = append(expectedReserved, StudentID(course.Capacity+i))
			}
			assertion.Equal(expectedReserved, course.ReserveQueue.CopyAsArray())
		}
	})
	t.Run("batcher test", func(t *testing.T) {
		assertion := assert.New(t)
		course := Course{
			ID:         CourseID(1),
			GroupID:    GroupID(1),
			Department: DepartmentID(4),
			Capacity:   2,
			RegisteredStudents: map[StudentID]struct{}{
				StudentID(1): {},
				StudentID(2): {},
			},
			ReserveCapacity: 2,
			ReserveQueue:    util.NewQueue[StudentID](),
		}
		course.ReserveQueue.Enqueue(StudentID(3))
		course.ReserveQueue.Enqueue(StudentID(4))
		// Check if data is queued correctly in message broker
		broker := new(inMemoryBatcher)
		assertion.NoError(course.DisenrollStudent(context.Background(), StudentID(2), broker))
		assertion.NoError(course.DisenrollStudent(context.Background(), StudentID(1), broker))
		assertion.Equal([]struct {
			data *proto.CourseDatabaseBatchMessage
			dep  DepartmentID
		}{
			{
				data: &proto.CourseDatabaseBatchMessage{
					Action: &proto.CourseDatabaseBatchMessage_Disenroll{
						Disenroll: &proto.CourseDatabaseBatchDisenrollMessage{
							StudentId: uint64(2),
							CourseId:  1,
						},
					},
				},
				dep: DepartmentID(4),
			},
			{
				data: &proto.CourseDatabaseBatchMessage{
					Action: &proto.CourseDatabaseBatchMessage_Disenroll{
						Disenroll: &proto.CourseDatabaseBatchDisenrollMessage{
							StudentId: uint64(1),
							CourseId:  1,
						},
					},
				},
				dep: DepartmentID(4),
			}}, broker.messages)
		assertion.Equal(map[StudentID]struct{}{
			StudentID(3): {},
			StudentID(4): {},
		}, course.RegisteredStudents)
	})
	t.Run("error batcher", func(t *testing.T) {
		course := Course{
			ID:         CourseID(1),
			GroupID:    GroupID(1),
			Department: DepartmentID(4),
			Capacity:   2,
			RegisteredStudents: map[StudentID]struct{}{
				StudentID(1): {},
				StudentID(2): {},
			},
			ReserveCapacity: 2,
			ReserveQueue:    util.NewQueue[StudentID](),
		}
		course.ReserveQueue.Enqueue(StudentID(3))
		course.ReserveQueue.Enqueue(StudentID(4))
		// Check the data if broker returns error
		innerError := errors.New("my error")
		assert.ErrorIs(t, course.DisenrollStudent(context.Background(), StudentID(1), errorBatcher{innerError}), BatchError{innerError})
		assert.Equal(t, map[StudentID]struct{}{
			StudentID(1): {},
			StudentID(2): {},
		}, course.RegisteredStudents)
	})
	t.Run("nil batcher", func(t *testing.T) {
		assert.PanicsWithValue(t, "nil batcher", func() {
			_ = new(Course).DisenrollStudent(context.Background(), StudentID(1), nil)
		})
	})
}

func TestCourseChangeGroupOfStudent(t *testing.T) {
	t.Run("different course", func(t *testing.T) {
		course1 := Course{ID: CourseID(1)}
		course2 := Course{ID: CourseID(2)}
		assert.PanicsWithValue(t, "different courses provided", func() {
			_, _ = course1.ChangeGroupOfStudent(context.Background(), StudentID(1), &course2, noOpBatcher{})
		})
		assert.PanicsWithValue(t, "different courses provided", func() {
			_, _ = course2.ChangeGroupOfStudent(context.Background(), StudentID(1), &course1, noOpBatcher{})
		})
	})
	t.Run("same group", func(t *testing.T) {
		course := Course{ID: CourseID(1), GroupID: GroupID(1)}
		assert.PanicsWithValue(t, "same group provided", func() {
			_, _ = course.ChangeGroupOfStudent(context.Background(), StudentID(1), &course, noOpBatcher{})
		})
		assert.PanicsWithValue(t, "same group provided", func() {
			_, _ = course.ChangeGroupOfStudent(context.Background(), StudentID(1), &course, noOpBatcher{})
		})
	})
	t.Run("invalid student", func(t *testing.T) {
		course1 := Course{
			ID:                 CourseID(1),
			GroupID:            GroupID(1),
			RegisteredStudents: map[StudentID]struct{}{},
			ReserveQueue:       util.NewQueue[StudentID](),
		}
		course2 := Course{ID: CourseID(1), GroupID: GroupID(2)}
		assert.PanicsWithValue(t, "student 1 does not exists in course 1-1", func() {
			_, _ = course1.ChangeGroupOfStudent(context.Background(), StudentID(1), &course2, noOpBatcher{})
		})
	})
	t.Run("capacity reached", func(t *testing.T) {
		course1 := Course{
			ID:                 CourseID(1),
			GroupID:            GroupID(1),
			RegisteredStudents: map[StudentID]struct{}{StudentID(1): {}},
		}
		course2 := Course{
			ID:                 CourseID(1),
			GroupID:            GroupID(2),
			Capacity:           0,
			RegisteredStudents: map[StudentID]struct{}{},
			ReserveCapacity:    0,
			ReserveQueue:       util.NewQueue[StudentID](),
		}
		res, err := course1.ChangeGroupOfStudent(context.Background(), StudentID(1), &course2, noOpBatcher{})
		assert.NoError(t, err)
		assert.False(t, res)
	})
	t.Run("normal transfer into registered", func(t *testing.T) {
		course1 := Course{
			ID:                 CourseID(1),
			GroupID:            GroupID(1),
			Capacity:           1,
			RegisteredStudents: map[StudentID]struct{}{StudentID(1): {}},
			ReserveCapacity:    1,
			ReserveQueue:       util.NewQueue[StudentID](),
		}
		course2 := Course{
			ID:                 CourseID(1),
			GroupID:            GroupID(2),
			Capacity:           1,
			RegisteredStudents: map[StudentID]struct{}{},
			ReserveCapacity:    0,
			ReserveQueue:       util.NewQueue[StudentID](),
		}
		res, err := course1.ChangeGroupOfStudent(context.Background(), StudentID(1), &course2, noOpBatcher{})
		assert.NoError(t, err)
		assert.True(t, res)
		assert.Len(t, course1.RegisteredStudents, 0)
		assert.Equal(t, map[StudentID]struct{}{StudentID(1): {}}, course2.RegisteredStudents)
	})
	t.Run("normal transfer into queue", func(t *testing.T) {
		course1 := Course{
			ID:                 CourseID(1),
			GroupID:            GroupID(1),
			Capacity:           1,
			RegisteredStudents: map[StudentID]struct{}{StudentID(1): {}},
			ReserveCapacity:    1,
			ReserveQueue:       util.NewQueue[StudentID](),
		}
		course1.ReserveQueue.Enqueue(StudentID(2))
		course2 := Course{
			ID:                 CourseID(1),
			GroupID:            GroupID(2),
			Capacity:           1,
			RegisteredStudents: map[StudentID]struct{}{StudentID(3): {}},
			ReserveCapacity:    1,
			ReserveQueue:       util.NewQueue[StudentID](),
		}
		res, err := course1.ChangeGroupOfStudent(context.Background(), StudentID(1), &course2, noOpBatcher{})
		assert.NoError(t, err)
		assert.True(t, res)
		assert.Equal(t, map[StudentID]struct{}{StudentID(2): {}}, course1.RegisteredStudents)
		assert.Equal(t, 0, course1.ReserveQueue.Len())
		assert.Equal(t, map[StudentID]struct{}{StudentID(3): {}}, course2.RegisteredStudents)
		assert.Equal(t, []StudentID{1}, course2.ReserveQueue.CopyAsArray())
	})
	t.Run("batch test", func(t *testing.T) {
		course1 := Course{
			Department:         DepartmentID(4),
			ID:                 CourseID(1),
			GroupID:            GroupID(1),
			Capacity:           1,
			RegisteredStudents: map[StudentID]struct{}{StudentID(1): {}},
			ReserveCapacity:    1,
			ReserveQueue:       util.NewQueue[StudentID](),
		}
		course1.ReserveQueue.Enqueue(StudentID(2))
		course2 := Course{
			Department:         DepartmentID(4),
			ID:                 CourseID(1),
			GroupID:            GroupID(2),
			Capacity:           1,
			RegisteredStudents: map[StudentID]struct{}{StudentID(3): {}},
			ReserveCapacity:    1,
			ReserveQueue:       util.NewQueue[StudentID](),
		}
		broker := new(inMemoryBatcher)
		// No error test
		res, err := course1.ChangeGroupOfStudent(context.Background(), StudentID(1), &course2, broker)
		assert.NoError(t, err)
		assert.True(t, res)
		// Capacity full test
		res, err = course1.ChangeGroupOfStudent(context.Background(), StudentID(2), &course2, broker)
		assert.NoError(t, err)
		assert.False(t, res)
		// Test
		assert.Equal(t, []struct {
			data *proto.CourseDatabaseBatchMessage
			dep  DepartmentID
		}{
			{
				data: &proto.CourseDatabaseBatchMessage{
					Action: &proto.CourseDatabaseBatchMessage_ChangeGroup{
						ChangeGroup: &proto.CourseDatabaseBatchChangeGroupMessage{
							StudentId: uint64(1),
							CourseId:  1,
							GroupId:   2,
						},
					},
				},
				dep: DepartmentID(4),
			},
		}, broker.messages)
	})
}

func TestCourseForceEnrollStudent(t *testing.T) {
	t.Run("general test", func(t *testing.T) {
		assertion := assert.New(t)
		course := Course{
			Capacity:           2,
			RegisteredStudents: make(map[StudentID]struct{}),
			ReserveCapacity:    2,
			ReserveQueue:       util.NewQueue[StudentID](),
		}
		// Add to registered courses
		for i := 0; i < course.Capacity; i++ {
			assertion.NoError(course.ForceEnroll(context.Background(), StudentID(i), noOpBatcher{}))
		}
		assertion.Equal(map[StudentID]struct{}{
			StudentID(0): {},
			StudentID(1): {},
		}, course.RegisteredStudents)
		assertion.Equal(0, course.ReserveQueue.Len())
		// Add to main queue again!
		for i := 0; i < course.ReserveCapacity; i++ {
			assertion.NoError(course.ForceEnroll(context.Background(), StudentID(course.ReserveCapacity+i), noOpBatcher{}))
		}
		assertion.Equal(map[StudentID]struct{}{
			StudentID(0): {},
			StudentID(1): {},
			StudentID(2): {},
			StudentID(3): {},
		}, course.RegisteredStudents)
		assertion.Equal(0, course.ReserveQueue.Len())
		assertion.Equal(4, course.Capacity)
	})
	t.Run("batcher test", func(t *testing.T) {
		assertion := assert.New(t)
		course := Course{
			ID:                 CourseID(1),
			GroupID:            GroupID(1),
			Department:         DepartmentID(0),
			Capacity:           2,
			RegisteredStudents: make(map[StudentID]struct{}),
			ReserveCapacity:    2,
			ReserveQueue:       util.NewQueue[StudentID](),
		}
		batcher := new(inMemoryBatcher)
		// Add to registered courses
		for i := 0; i < course.Capacity; i++ {
			assertion.NoError(course.ForceEnroll(context.Background(), StudentID(i), batcher))
		}
		// Add to main queue again!
		for i := 0; i < course.ReserveCapacity; i++ {
			assertion.NoError(course.ForceEnroll(context.Background(), StudentID(course.ReserveCapacity+i), batcher))
		}
		// Check batcher
		assertion.Equal([]struct {
			data *proto.CourseDatabaseBatchMessage
			dep  DepartmentID
		}{
			{
				data: &proto.CourseDatabaseBatchMessage{
					Action: &proto.CourseDatabaseBatchMessage_Enroll{
						Enroll: &proto.CourseDatabaseBatchEnrollMessage{
							StudentId: 0,
							CourseId:  1,
							GroupId:   1,
							Reserved:  false,
						},
					},
				},
				dep: DepartmentID(0),
			},
			{
				data: &proto.CourseDatabaseBatchMessage{
					Action: &proto.CourseDatabaseBatchMessage_Enroll{
						Enroll: &proto.CourseDatabaseBatchEnrollMessage{
							StudentId: 1,
							CourseId:  1,
							GroupId:   1,
							Reserved:  false,
						},
					},
				},
				dep: DepartmentID(0),
			},
			{
				data: &proto.CourseDatabaseBatchMessage{
					Action: &proto.CourseDatabaseBatchMessage_UpdateCapacity{
						UpdateCapacity: &proto.CourseDatabaseBatchUpdateCapacity{
							CourseId:    1,
							GroupId:     1,
							NewCapacity: 3,
						},
					},
				},
				dep: DepartmentID(0),
			},
			{
				data: &proto.CourseDatabaseBatchMessage{
					Action: &proto.CourseDatabaseBatchMessage_Enroll{
						Enroll: &proto.CourseDatabaseBatchEnrollMessage{
							StudentId: 2,
							CourseId:  1,
							GroupId:   1,
							Reserved:  false,
						},
					},
				},
				dep: DepartmentID(0),
			},
			{
				data: &proto.CourseDatabaseBatchMessage{
					Action: &proto.CourseDatabaseBatchMessage_UpdateCapacity{
						UpdateCapacity: &proto.CourseDatabaseBatchUpdateCapacity{
							CourseId:    1,
							GroupId:     1,
							NewCapacity: 4,
						},
					},
				},
				dep: DepartmentID(0),
			},
			{
				data: &proto.CourseDatabaseBatchMessage{
					Action: &proto.CourseDatabaseBatchMessage_Enroll{
						Enroll: &proto.CourseDatabaseBatchEnrollMessage{
							StudentId: 3,
							CourseId:  1,
							GroupId:   1,
							Reserved:  false,
						},
					},
				},
				dep: DepartmentID(0),
			},
		}, batcher.messages)
	})
	t.Run("error batcher test", func(t *testing.T) {
		assertion := assert.New(t)
		course := Course{
			ID:                 CourseID(1),
			GroupID:            GroupID(1),
			Department:         DepartmentID(0),
			Capacity:           2,
			RegisteredStudents: make(map[StudentID]struct{}),
			ReserveCapacity:    2,
			ReserveQueue:       util.NewQueue[StudentID](),
		}
		innerError := errors.New("error")
		batcher := errorBatcher{err: innerError}
		assertion.ErrorIs(course.ForceEnroll(context.Background(), StudentID(0), batcher), BatchError{innerError})
		assertion.Empty(course.RegisteredStudents)
	})
	t.Run("nil batcher test", func(t *testing.T) {
		course := Course{
			Capacity:           2,
			RegisteredStudents: make(map[StudentID]struct{}),
			ReserveCapacity:    2,
			ReserveQueue:       util.NewQueue[StudentID](),
		}
		assert.PanicsWithValue(t, "nil batcher", func() {
			_ = course.ForceEnroll(context.Background(), StudentID(0), nil)
		})
	})
}
