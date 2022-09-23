package AuthCore

import (
	db "CourseEnrollment/internal/database/AuthCore"
	"crypto/rand"
)

// API contains the data needed to operate the auth core endpoints
type API struct {
	// The database to authorize users
	Database db.Database
	// The key to sign stuff with it
	jwtKey []byte
}

// GenerateJWTKey generates a random JWT key
func (a *API) GenerateJWTKey() {
	const keyLen = 32
	key := make([]byte, keyLen)
	_, _ = rand.Read(key)
	a.jwtKey = key
}
