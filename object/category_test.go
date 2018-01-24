package object

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCategorySet_Add(t *testing.T) {
	cs := make(CategorySet)
	cs.Add(Food)
	cs.Add(Weapon)

	assert.True(t, cs[Food])
	assert.True(t, cs[Weapon])
	assert.False(t, cs[Armor])
}

func TestCategorySet_ToList(t *testing.T) {
	cs := make(CategorySet)
	cs.Add(Food)
	cs.Add(Weapon)

	list := cs.ToInt32List()
	assert.Contains(t, list, int32(Food))
	assert.Contains(t, list, int32(Weapon))
	assert.NotContains(t, list, int32(Armor))
}

func TestCategoriesToString(t *testing.T) {
	cats := []Category{
		Food, Weapon, Armor,
	}
	str := CategoriesToString(cats)
	assert.Equal(t, "Food, Weapon, Armor", str)

	// empty slice
	assert.Equal(t, "", CategoriesToString([]Category{}))
}

func TestCategorySet_Contains(t *testing.T) {
	cs := make(CategorySet)
	cs.Add(Weapon)

	assert.True(t, cs.Contains(Weapon))
	assert.False(t, cs.Contains(Food))
}

func TestStringToCategory(t *testing.T) {
	runTest := func(c Category) {
		result, err := StringToCategory(c.String())
		assert.NoError(t, err)
		assert.Equal(t, c, result)

		result, err = StringToCategory(strings.ToUpper(c.String()))
		assert.NoError(t, err)
		assert.Equal(t, c, result)
	}

	runTest(None)
	runTest(Weapon)
	runTest(Wand)
	runTest(Staff)
	runTest(Treasure)
	runTest(Armor)
	runTest(Food)
	runTest(Other)

	c, err := StringToCategory("")
	assert.Equal(t, c, None)
	assert.NoError(t, err)

	_, err = StringToCategory("asdflkjas")
	assert.Error(t, err)
}
