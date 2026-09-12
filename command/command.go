// Package command is the game's inbound vocabulary: plain structs describing
// what a connection is asking the world to do.
//
// Like event, this is a leaf package -- it imports direction, slot and the
// standard library, and nothing else.
package command

// Command is anything a connection can ask the world to do.
//
// Verb is the canonical name of the command, not necessarily what the player
// typed ("l" and "look" are both command.Look, whose verb is "look"). It
// exists so a failing handler doesn't have to repeat its own name: see
// gameserver.HandlerParameter.Fail, which uses it to pick failure text.
type Command interface {
	Verb() string
}
