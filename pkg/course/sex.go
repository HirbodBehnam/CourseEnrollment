package course

import "errors"

// Sex represents the sex of a student
type Sex uint8

const (
	SexMale Sex = iota + 1
	SexFemale
)

// Scan will scan the sex enum from database
func (s *Sex) Scan(value interface{}) error {
	data, ok := value.([]byte)
	if !ok {
		return databaseInvalidTypeErr
	}
	switch string(data) {
	case "male":
		*s = SexMale
	case "female":
		*s = SexFemale
	default:
		return errors.New("unexpected value: " + string(data))
	}
	return nil
}

// SexLock Some courses can be picked by certain gender. This property is stored as SexLock.
type SexLock uint8

const (
	// SexLockUnlocked means that both genders can pick this course
	SexLockUnlocked SexLock = iota
	SexLockMaleOnly
	SexLockFemaleOnly
)

// Scan will scan the SexLock from database. A null value in database means no lock.
func (s *SexLock) Scan(value interface{}) error {
	// Without lock check
	if value == nil {
		*s = SexLockUnlocked
		return nil
	}
	// Check normal stuff
	var sex Sex
	err := sex.Scan(value)
	if err != nil {
		return err
	}
	// One to one mapping
	*s = SexLock(sex)
	return nil
}

// sexLockCompatible checks if a locked course can be picked up by a user which "sex" Sex
func sexLockCompatible(lock SexLock, sex Sex) bool {
	if lock == SexLockUnlocked {
		return true
	} else {
		return uint8(sex) == uint8(lock)
	}
}
