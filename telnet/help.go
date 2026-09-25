package telnet

import (
	"fmt"
	"strings"
)

// help is answered by the connection, like a parse error, rather than sent to
// the world: it is a list of what to type, and only parse.go knows that. Each
// entry names the verbs it shows so help_test.go can hold it to the parser --
// a verb help offers and the parser refuses is a new player's first
// "Unknown request". Builder commands stay out; to a player they don't exist.

type helpEntry struct {
	usage string
	what  string
	verbs []string
}

type helpSection struct {
	title   string
	entries []helpEntry
}

var helpSections = []helpSection{
	{"Moving around", []helpEntry{
		{"n s e w u d", "walk that way", []string{"n", "s", "e", "w", "u", "d"}},
		{"exits", "where you can go from here", []string{"exits"}},
		{"recall", "back to the start", []string{"recall"}},
	}},
	{"Looking", []helpEntry{
		{"look, look <thing>", "the room, or something in it (l)", []string{"look", "l"}},
		{"look in <corpse>", "what it's holding", []string{"look"}},
		{"consider <mob>", "how a fight with it would go (con)", []string{"consider", "con"}},
	}},
	{"Things", []helpEntry{
		{"get <item>", "pick it up; get all, get 2.knife", []string{"get"}},
		{"get <item> from <corpse>", "loot", []string{"get"}},
		{"drop <item>", "put it down", []string{"drop"}},
		{"wear <item>, wield <item>", "put it on, take up a weapon", []string{"wear", "wield"}},
		{"remove <item>", "take it off", []string{"remove"}},
		{"inventory", "what you're carrying (i)", []string{"inventory", "i"}},
		{"equipment", "what you're wearing (eq)", []string{"equipment", "eq"}},
	}},
	{"Fighting", []helpEntry{
		{"kill <mob>", "start a fight", []string{"kill"}},
		{"flee", "get out of one", []string{"flee"}},
	}},
	{"You", []helpEntry{
		{"stat", "health and power", []string{"stat"}},
		{"role", "what your gear makes you, and why", []string{"role"}},
	}},
	{"Talking", []helpEntry{
		{"say <words>", "to the room (')", []string{"say", "'"}},
		{"tell <who> <words>", "to one player, anywhere", []string{"tell"}},
		{"shout <words>", "to everyone playing", []string{"shout"}},
		{"who", "who's playing", []string{"who"}},
		{"quit", "save and leave", []string{"quit"}},
	}},
}

// helpText is built once: the sections never change while the server runs.
var helpText = func() string {
	var b strings.Builder
	for i, s := range helpSections {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(s.title + "\n")
		for _, e := range s.entries {
			fmt.Fprintf(&b, "  %-26s %s\n", e.usage, e.what)
		}
	}
	return b.String()
}()

// isHelp is every way of asking.
func isHelp(line string) bool {
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "help", "?", "commands":
		return true
	}
	return false
}
