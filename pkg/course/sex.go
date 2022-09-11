package course

// Sex represents the sex of a student
type Sex uint8

const (
	SexMale Sex = iota + 1
	SexFemale
)

// SexLock Some courses can be picked by certain gender. This property is stored as SexLock.
type SexLock uint8

const (
	// SexLockUnlocked means that both genders can pick this course
	SexLockUnlocked SexLock = iota
	SexLockMaleOnly
	SexLockFemaleOnly
)

// sexLockCompatible checks if a locked course can be picked up by a user which "sex" Sex
func sexLockCompatible(lock SexLock, sex Sex) bool {
	if lock == SexLockUnlocked {
		return true
	} else {
		return uint8(sex) == uint8(lock)
	}
}
