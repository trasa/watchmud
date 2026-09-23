package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegenAmount(t *testing.T) {
	assert.Equal(t, 5, RegenAmount(100))
	assert.Equal(t, 1, RegenAmount(25), "5% of 25 is 1.25, rounded down")
	assert.Equal(t, 1, RegenAmount(2), "at least a point")
}
