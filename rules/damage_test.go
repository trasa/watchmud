package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDamageRoll(t *testing.T) {
	for _, good := range []string{"1d6", "2d4+1", "1d8-1", "10d10"} {
		d, err := ParseDamageRoll(good)
		assert.NoError(t, err)
		assert.Equal(t, DamageRoll(good), d)
	}
	for _, bad := range []string{"", "d6", "1d", "0d6", "1d0", "banana 1d6", "1d6 fire", "1d6+", "1D6"} {
		_, err := ParseDamageRoll(bad)
		assert.Error(t, err, bad)
	}
}
