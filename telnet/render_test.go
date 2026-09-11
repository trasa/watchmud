package telnet

import (
	"strings"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	message "github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud-message/slot"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/world"
)

type commandCase struct {
	name      string
	setup     func(*world.World, *player.Player, *player.Player)
	input     string
	want      string
	wantOther string
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

			gm, err := message.TranslateLineToMessage(message.Tokenize(tc.input))
			require.NoError(t, err)
			require.NoError(t, w.HandleIncomingMessage(gameserver.NewHandlerParameter(c, gm)))

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
	return object.NewInstance(uuid.New(), object.NewDefinition(
		"knife", "knife", "start", object.Weapon, []string{},
		"knife", "A knife is on the ground.", slot.Wield))
}
