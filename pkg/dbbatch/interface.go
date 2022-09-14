package dbbatch

// Interface must be used to send a database query to our message broker
type Interface interface {
	// UpdateDatabase must be called when we want to queue a message
	UpdateDatabase(Message) error
}
