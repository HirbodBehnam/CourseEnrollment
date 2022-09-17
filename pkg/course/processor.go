package course

import "CourseEnrollment/pkg/proto"

// Batcher must send the request in a queue to be processed later
type Batcher interface {
	// Process must send it into queue
	Process(DepartmentID, *proto.CourseDatabaseBatchMessage) error
}

// BatchError is an error which Batcher.Process can return
type BatchError struct {
	err error
}

func (e BatchError) Error() string {
	return "cannot batch data: " + e.err.Error()
}

func (e BatchError) Unwrap() error {
	return e.err
}
