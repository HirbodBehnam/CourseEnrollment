package course

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSexLockCompatible(t *testing.T) {
	assert.True(t, sexLockCompatible(SexLockUnlocked, SexFemale))
	assert.True(t, sexLockCompatible(SexLockUnlocked, SexMale))
	assert.True(t, sexLockCompatible(SexLockMaleOnly, SexMale))
	assert.False(t, sexLockCompatible(SexLockMaleOnly, SexFemale))
	assert.True(t, sexLockCompatible(SexLockFemaleOnly, SexFemale))
	assert.False(t, sexLockCompatible(SexLockFemaleOnly, SexMale))
}
