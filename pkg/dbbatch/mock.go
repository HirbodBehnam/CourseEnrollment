package dbbatch

import "sync"

// Mock is a mock for Interface.
// It only stores the
type Mock struct {
	Messages []Message
	mu       sync.Mutex
}

func (m *Mock) UpdateDatabase(msg Message) error {
	m.mu.Lock()
	m.Messages = append(m.Messages, msg)
	m.mu.Unlock()
	return nil
}
