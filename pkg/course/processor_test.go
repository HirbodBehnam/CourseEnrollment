package course

import (
	"CourseEnrollment/pkg/proto"
	"sync"
)

type noOpBatcher struct {
}

func (noOpBatcher) Process(DepartmentID, *proto.CourseDatabaseBatchMessage) error {
	return nil
}

type inMemoryBatcher struct {
	messages []struct {
		data *proto.CourseDatabaseBatchMessage
		dep  DepartmentID
	}
	mu sync.Mutex
}

func (b *inMemoryBatcher) Process(dep DepartmentID, data *proto.CourseDatabaseBatchMessage) error {
	b.mu.Lock()
	b.messages = append(b.messages, struct {
		data *proto.CourseDatabaseBatchMessage
		dep  DepartmentID
	}{data: data, dep: dep})
	b.mu.Unlock()
	return nil
}

type errorBatcher struct {
	err error
}

func (b errorBatcher) Process(DepartmentID, *proto.CourseDatabaseBatchMessage) error {
	return b.err
}
