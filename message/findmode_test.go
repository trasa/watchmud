package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindMode(t *testing.T) {
	assert.Equal(t, FindAll, determineFindMode("all"))
	assert.Equal(t, FindAllDot, determineFindMode("all.thing"))
	assert.Equal(t, FindIndividual, determineFindMode("thing"))
	assert.Equal(t, FindIndividual, determineFindMode("1.thing"))
}

func TestParseFindMode_Individual(t *testing.T) {
	findMode, dotPart, target := parseFindToken("foo")
	assert.Equal(t, FindIndividual, findMode)
	assert.Equal(t, "", dotPart)
	assert.Equal(t, "foo", target)
}

func TestParseFindMode_Individual_with_dot(t *testing.T) {
	findMode, dotPart, target := parseFindToken("1.foo")
	assert.Equal(t, FindIndividual, findMode)
	assert.Equal(t, "1", dotPart)
	assert.Equal(t, "foo", target)
}

func TestParseFindMode_All(t *testing.T) {
	findMode, dotPart, target := parseFindToken("all")
	assert.Equal(t, FindAll, findMode)
	assert.Equal(t, "", dotPart)
	assert.Equal(t, "all", target)
}

func TestParseFindMode_All_with_dot(t *testing.T) {
	findMode, dotPart, target := parseFindToken("all.foo")
	assert.Equal(t, FindAllDot, findMode)
	assert.Equal(t, "all", dotPart)
	assert.Equal(t, "foo", target)
}
