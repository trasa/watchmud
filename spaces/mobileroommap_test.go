package spaces

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/mobile"
	"testing"
)

func TestMobileRoomMap_GetAllMobiles(t *testing.T) {
	defn := mobile.NewDefinition("id", "name", "zone", []string{}, "shortdesc", "roomdesc", mobile.WanderingDefinition{})
	instOne := mobile.NewInstance(defn)
	instTwo := mobile.NewInstance(defn)
	r := NewTestRoom("test")

	mobileMap := NewMobileRoomMap()
	mobileMap.Add(instOne, r)
	mobileMap.Add(instTwo, r)

	result := mobileMap.GetAllMobiles()
	assert.Equal(t, 2, len(result))
	assert.True(t, instOne == result[0] || instOne == result[1])
	assert.True(t, instTwo == result[0] || instTwo == result[1])
}

func TestMobileRoomMap_GetAllMobiles_Empty(t *testing.T) {
	mobileMap := NewMobileRoomMap()

	result := mobileMap.GetAllMobiles()
	assert.Equal(t, 0, len(result))
}
