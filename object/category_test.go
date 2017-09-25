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

	list := cs.ToList()
	assert.Contains(t, list, FOOD)
	assert.Contains(t, list, WEAPON)
	assert.NotContains(t, list, ARMOR)
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
