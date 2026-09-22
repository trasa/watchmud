package server

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/loader"
	"github.com/trasa/watchmud/memstore"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/rules"
	"github.com/trasa/watchmud/testdice"
	"github.com/trasa/watchmud/world"
)

// testConn is a connection with nobody on it yet, which is what creation
// starts from. It records rather than forwarding to the player, since the
// player it is about to be given sends through this same conn.
type testConn struct {
	p    *player.Player
	sent []any
}

func (c *testConn) Send(msg any)               { c.sent = append(c.sent, msg) }
func (c *testConn) Player() *player.Player     { return c.p }
func (c *testConn) SetPlayer(p *player.Player) { c.p = p }
func (c *testConn) Close()                     {}

func newTestGameServer(t *testing.T) (*GameServer, *memstore.Store) {
	t.Helper()

	content, err := loader.LoadContent(os.DirFS("../testcontent"))
	require.NoError(t, err)

	store := memstore.New()
	w, err := world.New(content, store, testdice.New())
	require.NoError(t, err)

	return New(w, content.Catalog, store), store
}

// A brand new character arrives dressed, and the record written for them says
// so -- so the gear is still there on their next login.
func TestCreatePlayer_startingGear(t *testing.T) {
	gs, store := newTestGameServer(t)
	c := &testConn{}

	err := gs.dispatch(gameserver.NewHandlerParameter(c, command.CreatePlayer{
		Name:    "newbie",
		Lineage: "human",
	}))
	require.NoError(t, err)

	p := c.Player()
	require.NotNil(t, p)
	// created, then shown where they are -- in that order, so the telnet login
	// conversation is over before the description arrives
	require.Len(t, c.sent, 2)
	assert.IsType(t, event.PlayerCreated{}, c.sent[0])
	assert.IsType(t, event.RoomDescription{}, c.sent[1])

	// testcontent's kit: a knife, a helmet, and a rope that is only carried.
	assert.Equal(t, 3, p.Inventory().Len())

	knife := p.Equipment().At(rules.SlotWield)
	require.NotNil(t, knife)
	assert.Equal(t, "knife", knife.Definition.Name)
	require.NotNil(t, p.Equipment().At(rules.SlotHead))

	rec, found, err := store.Load("newbie")
	require.NoError(t, err)
	require.True(t, found)
	assert.Len(t, rec.Inventory, 3)
	assert.Len(t, rec.Equipment, 2)

	// and it comes back the same way it went in
	restored, err := player.FromRecord(rec, &player.Recorder{}, gs.catalog, gs.world)
	require.NoError(t, err)
	assert.Equal(t, 3, restored.Inventory().Len())
	require.NotNil(t, restored.Equipment().At(rules.SlotWield))
	require.NotNil(t, restored.Equipment().At(rules.SlotHead))
}

// Once a character is in the world, the prompt reaches them -- which is what
// puts the first "> " on the screen after login.
func TestPrompt_reachesPlayersInTheWorld(t *testing.T) {
	gs, _ := newTestGameServer(t)
	c := &testConn{}
	require.NoError(t, gs.dispatch(gameserver.NewHandlerParameter(c, command.CreatePlayer{Name: "newbie"})))

	c.Player().TakeMeleeDamage(30)
	gs.prompt()

	assert.Equal(t, event.Prompt{CurrentHealth: 70, MaxHealth: 100}, c.sent[len(c.sent)-1])
}

// A failed login never joined the world, so the login conversation isn't
// interrupted by a prompt it didn't ask for.
func TestPrompt_skipsAFailedLogin(t *testing.T) {
	gs, _ := newTestGameServer(t)
	c := &testConn{}
	require.NoError(t, gs.dispatch(gameserver.NewHandlerParameter(c, command.Login{Name: "nobody"})))

	gs.prompt()

	require.Len(t, c.sent, 1)
	assert.IsType(t, event.LoginFailed{}, c.sent[0])
}

// A returning player comes back where they left off.
func TestLogin_returnsToTheLastRoom(t *testing.T) {
	gs, store := newTestGameServer(t)
	c := &testConn{}
	require.NoError(t, gs.dispatch(gameserver.NewHandlerParameter(c, command.CreatePlayer{Name: "wanderer"})))
	rec, _, err := store.Load("wanderer")
	require.NoError(t, err)
	gs.world.RemovePlayer(c.Player())

	rec.LastZoneId, rec.LastRoomId = "wrathrock", "market_square"
	require.NoError(t, store.Save(rec))

	back := &testConn{}
	require.NoError(t, gs.dispatch(gameserver.NewHandlerParameter(back, command.Login{Name: "wanderer"})))

	require.NotNil(t, back.Player())
	require.Len(t, back.sent, 2)
	assert.IsType(t, event.LoggedIn{}, back.sent[0])
	assert.Equal(t, "Market Square", back.sent[1].(event.RoomDescription).Name, "shown where they came back to")

	// a look is saved like any other command, and the save says where they are
	require.NoError(t, gs.dispatch(gameserver.NewHandlerParameter(back, command.Look{})))
	rec, _, err = store.Load("wanderer")
	require.NoError(t, err)
	assert.Equal(t, "market_square", rec.LastRoomId)
}
