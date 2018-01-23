package object

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCategorySet_Add(t *testing.T) {
	cs := make(CategorySet)
	cs.Add(FOOD)
	cs.Add(WEAPON)

	assert.True(t, cs[FOOD])
	assert.True(t, cs[WEAPON])
	assert.False(t, cs[ARMOR])
}

func TestCategorySet_ToList(t *testing.T) {
	cs := make(CategorySet)
	cs.Add(FOOD)
	cs.Add(WEAPON)

	list := cs.ToInt32List()
	assert.Contains(t, list, int32(FOOD))
	assert.Contains(t, list, int32(WEAPON))
	assert.NotContains(t, list, int32(ARMOR))
}

func TestCategoriesToString(t *testing.T) {
	cats := []Category{
		FOOD, WEAPON, ARMOR,
	}
	str := CategoriesToString(cats)
	assert.Equal(t, "FOOD, WEAPON, ARMOR", str)

	// empty slice
	assert.Equal(t, "", CategoriesToString([]Category{}))
}

func TestCategorySet_Contains(t *testing.T) {
	cs := make(CategorySet)
	cs.Add(WEAPON)

	assert.True(t, cs.Contains(WEAPON))
	assert.False(t, cs.Contains(FOOD))
}

func TestStringToCategory(t *testing.T) {
	runTest := func(c Category) {
		result, err := StringToCategory(c.String())
		assert.NoError(t, err)
		assert.Equal(t, c, result)
	}

	runTest(NONE)
	runTest(WEAPON)
	runTest(WAND)
	runTest(STAFF)
	runTest(TREASURE)
	runTest(ARMOR)
	runTest(FOOD)
	runTest(OTHER)

	_, err := StringToCategory("")
	assert.Error(t, err)

	_, err = StringToCategory("asdflkjas")
	assert.Error(t, err)
}
