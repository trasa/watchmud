package world

import (
	"testing"
	"time"
	"uuid"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/loader"
	"github.com/trasa/watchmud/memstore"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/rules"
	"github.com/trasa/watchmud/slot"
	"github.com/trasa/watchmud/spaces"
	"github.com/trasa/watchmud/zonereset"
)

type worldTestSuite struct {
	suite.Suite
	w *World
	r *player.Recorder
	p *player.Player
	c *gameserver.TestConn
}

// sent returns the ith message in the recorder
func sent[T any](t *testing.T, r *player.Recorder, i int) T {
	t.Helper()
	require.Greater(t, len(r.Sent), i, "wanted Sent[%d], only %d sent", i, len(r.Sent))
	v, ok := r.Sent[i].(T)
	require.Truef(t, ok, "Sent[%d] is %T, want %T", i, r.Sent[i], *new(T))
	return v
}

// handlerParameter builds the handler input for a command.
func (s *worldTestSuite) handlerParameter(cmd command.Command) *gameserver.HandlerParameter {
	s.T().Helper()
	return gameserver.NewHandlerParameter(s.c, cmd)
}

func (s *worldTestSuite) SetupTest() {
	w, err := NewTestWorld()
	require.NoError(s.T(), err)
	s.w = w
	s.r = &player.Recorder{}
	s.p = player.NewTestPlayer(uuid.New(), "testdood", s.r)
	s.w.AddPlayer(s.p)
	s.c = gameserver.NewTestConn(s.p)
}

func NewTestWorld() (*World, error) {

	voidZone := spaces.NewZone("void", "void", zonereset.NEVER, time.Duration(0))
	voidRoom := spaces.NewRoom(voidZone, "void", "void", "void")
	voidZone.AddRoom(voidRoom)

	startZone := spaces.NewZone("start", "start", zonereset.NEVER, time.Duration(0))
	startRoom := spaces.NewRoom(startZone, "start", "start", "this is a test room.")
	startZone.AddRoom(startRoom)

	// stuff that's in the start room
	knife := object.NewDefinition(
		"knife",
		"knife",
		startZone.Id,
		object.Weapon, []string{},
		"knife",
		"A knife is on the ground.",
		slot.Wield,
	)
	knife.RoleWeights = map[string]int{"striker": 2}
	knifeInstance := object.NewInstance(uuid.New(), knife)
	if err := startRoom.Inventory.Add(knifeInstance); err != nil {
		return nil, err
	}
	helmet := object.NewDefinition(
		"helmet",
		"helmet",
		startZone.Id,
		object.Armor,
		[]string{"helm", "iron", "helmet"},
		"iron helmet",
		"an iron helmet is on the ground",
		slot.Head,
	)
	helmet.RoleWeights = map[string]int{"tank": 2}
	helmetInstance := object.NewInstance(uuid.New(), helmet)
	if err := startRoom.Inventory.Add(helmetInstance); err != nil {
		return nil, err
	}
	mob := mobile.NewDefinition(
		"targetDrone",
		"Target Drone",
		startZone.Id,
		[]string{"target", "drone"},
		"Target Drone",
		"Target Drone buzzes around.",
		25,
		mobile.WanderingDefinition{CanWander: false},
		10,
	)
	startZone.AddMobileDefinition(mob)
	if err := startRoom.AddMobile(mobile.NewInstance(mob)); err != nil {
		return nil, err
	}

	zones := []*spaces.Zone{
		voidZone,
		startZone,
	}

	settings := loader.Settings{
		VoidZone:  "void",
		VoidRoom:  "void",
		StartZone: "start",
		StartRoom: "start",
	}

	catalog, err := rules.NewTestCatalog()
	if err != nil {
		return nil, err
	}

	store := memstore.New()
	content := loader.NewContent(&settings, catalog, zones)

	return New(content, store)
}
