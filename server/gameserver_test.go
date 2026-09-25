package server

import (
	"os"
	"testing"
	"time"

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
	"golang.org/x/crypto/bcrypt"
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

	gs := New(w, content.Catalog, store)
	gs.bcryptCost = bcrypt.MinCost
	return gs, store
}

// settle does what Run would: takes the callback a login or creation queued
// from its bcrypt goroutine and dispatches it.
func settle(t *testing.T, gs *GameServer) {
	t.Helper()
	select {
	case msg := <-gs.incomingBuffer:
		require.NoError(t, gs.dispatch(msg))
	case <-time.After(5 * time.Second):
		t.Fatal("nothing came back from the bcrypt goroutine")
	}
}

// create and login are the whole two-step conversation.
func create(t *testing.T, gs *GameServer, c gameserver.Conn, name, password string) {
	t.Helper()
	require.NoError(t, gs.dispatch(gameserver.NewHandlerParameter(c, command.CreatePlayer{
		Name:     name,
		Lineage:  "human",
		Password: command.Secret(password),
	})))
	settle(t, gs)
}

func login(t *testing.T, gs *GameServer, c gameserver.Conn, name, password string) {
	t.Helper()
	require.NoError(t, gs.dispatch(gameserver.NewHandlerParameter(c, command.Login{
		Name:     name,
		Password: command.Secret(password),
	})))
	settle(t, gs)
}

// A brand new character arrives dressed, and the record written for them says
// so -- so the gear is still there on their next login.
// create player is two steps now because of hashing the password.
func TestCreatePlayer_startingGear(t *testing.T) {
	gs, store := newTestGameServer(t)
	c := &testConn{}

	create(t, gs, c, "newbie", "sekrit")

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
	create(t, gs, c, "newbie", "sekrit")

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
	create(t, gs, c, "wanderer", "sekrit")
	rec, _, err := store.Load("wanderer")
	require.NoError(t, err)
	gs.world.RemovePlayer(c.Player())

	rec.LastZoneId, rec.LastRoomId = "wrathrock", "market_square"
	require.NoError(t, store.Save(rec))

	back := &testConn{}
	login(t, gs, back, "wanderer", "sekrit")

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

// The store no longer refuses a duplicate name itself (saves are queued), so
// creation has to.
func TestCreatePlayer_nameTaken(t *testing.T) {
	gs, _ := newTestGameServer(t)
	create(t, gs, &testConn{}, "bob", "sekrit")

	second := &testConn{}
	require.NoError(t, gs.dispatch(gameserver.NewHandlerParameter(second, command.CreatePlayer{Name: "bob", Password: "other"})))

	assert.Nil(t, second.Player(), "no second bob")
	require.Len(t, second.sent, 1)
	assert.Equal(t, event.CreateFailed{Reason: event.NameTaken}, second.sent[0])
}

// One character, one session: a second login is refused while the first is
// playing, and allowed once they have left.
func TestLogin_alreadyPlaying(t *testing.T) {
	gs, _ := newTestGameServer(t)
	first := &testConn{}
	create(t, gs, first, "bob", "sekrit")

	second := &testConn{}
	require.NoError(t, gs.dispatch(gameserver.NewHandlerParameter(second, command.Login{Name: "bob", Password: "sekrit"})))
	assert.Nil(t, second.Player())
	require.Len(t, second.sent, 1)
	assert.Equal(t, event.LoginFailed{Reason: event.AlreadyPlaying}, second.sent[0])

	require.NoError(t, gs.dispatch(gameserver.NewHandlerParameter(first, command.Logout{})))

	third := &testConn{}
	login(t, gs, third, "bob", "sekrit")
	assert.NotNil(t, third.Player(), "back in once the first session is gone")
}

// The wrong password is refused, with a reason: login() reads an empty one as
// success.
func TestLogin_wrongPassword(t *testing.T) {
	gs, _ := newTestGameServer(t)
	first := &testConn{}
	create(t, gs, first, "bob", "sekrit")
	require.NoError(t, gs.dispatch(gameserver.NewHandlerParameter(first, command.Logout{})))

	c := &testConn{}
	login(t, gs, c, "bob", "guess")

	assert.Nil(t, c.Player())
	require.Len(t, c.sent, 1)
	assert.Equal(t, event.LoginFailed{Reason: event.BadPassword}, c.sent[0])
}

// Two creations of one name can both pass the first check while their hashes
// are being made; the second to come back must lose.
func TestCreatePlayer_nameTakenWhileHashing(t *testing.T) {
	gs, _ := newTestGameServer(t)
	first, second := &testConn{}, &testConn{}
	for _, c := range []*testConn{first, second} {
		require.NoError(t, gs.dispatch(gameserver.NewHandlerParameter(c, command.CreatePlayer{Name: "bob", Password: "sekrit"})))
	}
	settle(t, gs)
	settle(t, gs)

	// either may have come back first
	winner, loser := first, second
	if first.Player() == nil {
		winner, loser = second, first
	}
	assert.NotNil(t, winner.Player())
	assert.Nil(t, loser.Player())
	assert.Equal(t, []any{event.CreateFailed{Reason: event.NameTaken}}, loser.sent)
}
