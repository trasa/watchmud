package telnet

import (
	"strings"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/slot"
	"github.com/trasa/watchmud/world"
)

type commandCase struct {
	name      string
	setup     func(*world.World, *player.Player, *player.Player)
	input     string
	want      string
	wantOther string
	// wantParseError is set for input the parser rejects before the world
	// sees it; the connection shows the player this text directly.
	wantParseError string
}

var commandCases = []commandCase{
	{
		name:  "get nonexistent",
		input: "get nonexistent",
		want:  "You don't see that here.\n",
	},
	{
		name:  "look shows the room",
		input: "look",
		want:  startRoomBlock,
	},
	{
		name:  "move with no exit",
		input: "north",
		want:  "You can't go that way.\n",
	},
	{
		name:  "get something that isn't here",
		input: "get sword",
		want:  "You don't see that here.\n",
	},
	{
		name:      "get from the room",
		input:     "get knife",
		want:      "Taken.\n",
		wantOther: "testdood gets knife.\n",
	},
	{
		// the verb override in resultcode.go: TARGET_NOT_FOUND means
		// something different to drop than it does to get
		name:  "drop something you aren't carrying",
		input: "drop sword",
		want:  "You aren't carrying that.\n",
	},
	{
		name:      "drop what you have",
		setup:     func(_ *world.World, p *player.Player, o *player.Player) { p.Inventory().Add(testKnife()) },
		input:     "drop knife",
		want:      "Dropped.\n",
		wantOther: "testdood drops knife.\n",
	},
	{
		name:  "empty inventory",
		input: "inventory",
		want:  "You aren't carrying anything.\n",
	},
	{
		name:  "inventory with one item",
		setup: func(_ *world.World, p *player.Player, o *player.Player) { p.Inventory().Add(testKnife()) },
		input: "inv",
		want:  "You are carrying:\n\tknife\n",
	},
	{
		name:      "say reaches the room",
		input:     "say hello there",
		want:      "You say, \"hello there\".\n",
		wantOther: "testdood says, \"hello there\".\n",
	},
	{
		name:  "tell to someone who isn't playing",
		input: "tell nobody hi",
		want:  "No one by that name is playing.\n",
	},
	{
		name:      "tell reaches the target",
		input:     "tell otherdood hi",
		want:      "Ok.\n",
		wantOther: "testdood tells you, \"hi\".\n",
	},
	{
		name:      "tellall reaches everyone else",
		input:     "tellall listen up",
		want:      "Ok.\n",
		wantOther: "testdood shouts, \"listen up\".\n",
	},
	{
		name:  "tellall with nothing to say",
		input: "tellall",
		want:  "Say what?\n",
	},
	{
		name:  "exits",
		input: "exits",
		want:  "Exits:\nNone!\n",
	},
	{
		name:  "who lists everyone",
		input: "who",
		want:  "-- Who Is Here --\notherdood the Human - start - start\ntestdood the Human - start - start\n",
	},
	{
		name:  "nothing equipped",
		input: "equipment",
		want:  "Nothing equipped.\n",
	},
	{
		name:  "wear something you aren't carrying",
		input: "wear helmet",
		want:  "You aren't carrying that.\n",
	},
	{
		name:  "wield with no target",
		input: "wield",
		want:  "Wield what?\n",
	},
	{
		name:  "kill something that isn't here",
		input: "kill dragon",
		want:  "You don't see that here.\n",
	},
	{
		name:  "a verb nobody knows",
		input: "florb the thing",
		want:  "",
		// parseCommand rejects it before the world ever sees it; the
		// connection prints the parser's error itself.
		wantParseError: "Unknown request: florb",
	},
	{
		name:  "role with nothing equipped",
		input: "role",
		want:  roleBlockNoGear,
	},
	{
		// the phase in one test case: gear alone decides the role
		name: "wearing armor makes you a tank",
		setup: func(_ *world.World, p *player.Player, o *player.Player) {
			p.Slots().Set(slot.Head, testHelmet())
		},
		input: "role",
		want:  roleBlockTank,
	},
	{
		name: "wielding a knife makes you a striker",
		setup: func(_ *world.World, p *player.Player, o *player.Player) {
			p.Slots().Set(slot.Wield, testKnife())
		},
		input: "roles",
		want:  roleBlockStriker,
	},
	{
		name: "remove something you aren't using",
		setup: func(_ *world.World, p *player.Player, o *player.Player) {
			p.Inventory().Add(testHelmet()) // carrying it is not using it
		},
		input: "remove helmet",
		want:  "You aren't using that.\n",
	},
	{
		name:  "remove with no target",
		input: "remove",
		want:  "Remove what?\n",
	},
	{
		name: "taking the gear off takes the role with it",
		setup: func(_ *world.World, p *player.Player, o *player.Player) {
			p.Slots().Set(slot.Head, testHelmet())
		},
		input: "remove helmet",
		want:  "You stop using iron helmet.\n",
	},
	{
		name:  "stat shows lineage and role, not class",
		input: "stat",
		want:  statBlockNoGear,
	},
	{
		name: "stat reflects what is equipped right now",
		setup: func(_ *world.World, p *player.Player, o *player.Player) {
			p.Slots().Set(slot.Head, testHelmet())
		},
		input: "stat",
		want:  statBlockTank,
	},
	{
		// a role shows up in who beside the lineage, where a class used to
		name: "who shows the role",
		setup: func(_ *world.World, p *player.Player, o *player.Player) {
			p.Slots().Set(slot.Head, testHelmet())
		},
		input: "who",
		want:  "-- Who Is Here --\notherdood the Human - start - start\ntestdood the Human Tank - start - start\n",
	},
	{
		// recall moves with direction.None, which must not render as "none!"
		name:      "recall leaves in no direction",
		input:     "recall",
		want:      startRoomBlock,
		wantOther: "testdood leaves.\ntestdood enters.\n",
	},
}

func TestCommandRendering(t *testing.T) {
	for _, tc := range commandCases {
		t.Run(tc.name, func(t *testing.T) {
			w, err := world.NewTestWorld()
			require.NoError(t, err)

			rec := &player.Recorder{}
			p := player.NewTestPlayer(uuid.New(), "testdood", rec)
			w.AddPlayer(p)
			c := gameserver.NewTestConn(p)

			// second player
			otherRec := &player.Recorder{}
			o := player.NewTestPlayer(uuid.New(), "otherdood", otherRec)
			w.AddPlayer(o)
			_ = gameserver.NewTestConn(o)

			if tc.setup != nil {
				tc.setup(w, p, o)
			}

			// the same path the connection takes: parse, then dispatch.
			cmd, err := parseCommand(strings.Fields(tc.input))
			if tc.wantParseError != "" {
				require.EqualError(t, err, tc.wantParseError)
				return
			}
			require.NoError(t, err)
			require.NoError(t, w.HandleIncomingMessage(gameserver.NewHandlerParameter(c, cmd)))

			var got strings.Builder
			for _, m := range rec.Sent {
				got.WriteString(render(m, p.Name()))
			}
			assert.Equal(t, tc.want, got.String())

			var otherGot strings.Builder
			for _, m := range otherRec.Sent {
				otherGot.WriteString(render(m, o.Name()))
			}
			assert.Equal(t, tc.wantOther, otherGot.String())
		})
	}
}

// what NewTestWorld's start room looks like
const startRoomBlock = `start
 this is a test room.
[ Exits: None! ]
A knife is on the ground.
an iron helmet is on the ground
Target Drone buzzes around.
otherdood is here.
`

func testKnife() *object.Instance {
	d := object.NewDefinition(
		"knife", "knife", "start", object.Weapon, []string{},
		"knife", "A knife is on the ground.", slot.Wield)
	d.RoleWeights = map[string]int{"striker": 2}
	return object.NewInstance(uuid.New(), d)
}

func testHelmet() *object.Instance {
	d := object.NewDefinition(
		"helmet", "helmet", "start", object.Armor, []string{"helm"},
		"iron helmet", "an iron helmet is on the ground", slot.Head)
	d.RoleWeights = map[string]int{"tank": 2}
	return object.NewInstance(uuid.New(), d)
}

// The role listing always shows every role, including the ones with nothing
// behind them: "you are a Tank" with no standings is a verdict the player
// can't argue with or work out how to change.
const roleBlockNoGear = `You aren't wearing anything that argues for a role.
  Tank     0
  Healer   0
  Striker  0
Change what you're wearing to change your role.
`

const roleBlockTank = `You are fighting as a Tank.
  Tank     2  (iron helmet 2)
  Healer   0
  Striker  0
Change what you're wearing to change your role.
`

const roleBlockStriker = `You are fighting as a Striker.
  Tank     0
  Healer   0
  Striker  2  (knife 2)
Change what you're wearing to change your role.
`

// Tabs, so these are quoted rather than raw. There is no ability block: what
// a character can do is their gear, and the six scores were a number nothing
// read.
const statBlockNoGear = "Status:\n" +
	"Player:\ttestdood\n" +
	"Lineage:\tHuman\tRole: none\n" +
	"Health:\t100 of 100\n" +
	"Location:\t(start - start)\n\n"

const statBlockTank = "Status:\n" +
	"Player:\ttestdood\n" +
	"Lineage:\tHuman\tRole: Tank\n" +
	"Health:\t100 of 100\n" +
	"Location:\t(start - start)\n\n"
