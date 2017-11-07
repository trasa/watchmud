package spaces

import (
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/object"
	"testing"
)

func TestRoomInventory_Empty(t *testing.T) {
	ri := NewRoomInventory()

	inst, exists := ri.GetByName("nothere")
	assert.False(t, exists)
	assert.Nil(t, inst)

	all := ri.GetAll()
	assert.Equal(t, 0, len(all))

	nothing, exists := ri.GetByInstanceId(uuid.NewV4())
	assert.Nil(t, nothing)
	assert.False(t, exists)
}

func TestRoomInventory_AddOne(t *testing.T) {
	ri := NewRoomInventory()

	defn := object.NewDefinition("id", "name", "zoneid", object.OTHER, []string{}, "short desc", "on ground")
	inst := object.NewInstance(defn)

	ri.Add(inst)

	bydefn, exists := ri.GetByName("name")
	assert.True(t, exists)
	assert.Equal(t, inst, bydefn)

	retrieved, exists := ri.GetByInstanceId(inst.InstanceId)
	assert.True(t, exists)
	assert.Equal(t, inst, retrieved)

	all := ri.GetAll()
	assert.Equal(t, 1, len(all))
	assert.Equal(t, inst, all[0])
}

func TestRoomInventory_AddMany(t *testing.T) {
	ri := NewRoomInventory()

	defn := object.NewDefinition("id", "name", "zoneid", object.OTHER, []string{}, "short desc", "on ground")
	instOne := object.NewInstance(defn)
	instTwo := object.NewInstance(defn)

	ri.Add(instOne)
	ri.Add(instTwo)

	inst, exists := ri.GetByName("name")
	assert.True(t, exists)
	assert.NotNil(t, inst)

	all := ri.GetAll()
	assert.Equal(t, 2, len(all))

	retone, exists := ri.GetByInstanceId(instOne.InstanceId)
	assert.True(t, exists)
	assert.Equal(t, retone, instOne)

	rettwo, exists := ri.GetByInstanceId(instTwo.InstanceId)
	assert.True(t, exists)
	assert.Equal(t, rettwo, instTwo)

	nothing, exists := ri.GetByInstanceId(uuid.NewV4())
	assert.False(t, exists)
	assert.Nil(t, nothing)
}

func TestRoomInventory_Remove(t *testing.T) {
	ri := NewRoomInventory()

	defn := object.NewDefinition("id", "name", "zoneid", object.OTHER, []string{}, "short desc", "on ground")
	inst := object.NewInstance(defn)

	assert.NoError(t, ri.Add(inst))

	assert.NoError(t, ri.Remove(inst))

	assert.Equal(t, 0, len(ri.GetAll()))
	i, exists := ri.GetByName("name")
	assert.Nil(t, i)
	assert.False(t, exists)

	ret, exists := ri.GetByInstanceId(inst.InstanceId)
	assert.False(t, exists)
	assert.Nil(t, ret)
}

func TestRoomInventory_RemoveEmpty(t *testing.T) {
	ri := NewRoomInventory()

	defn := object.NewDefinition("id", "name", "zoneid", object.OTHER, []string{}, "short desc", "on ground")
	inst := object.NewInstance(defn)

	assert.Error(t, ri.Remove(inst))
}
