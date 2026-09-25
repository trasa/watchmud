package server

import "github.com/watchmud/watchmud/command"

// loginChecked is a command.Command that contains the callback from comparing a password to the stored password hash,
// not exported and put over here so clients can't ever forge one.
type loginChecked struct {
	Name string
	Ok   bool
}

func (loginChecked) Verb() string { return "loginChecked" }

// createHashed is a command.Command that creates a hashed password from a plaintext password when
// creating a new player and record. Private for the same reason as loginChecked.
type createHashed struct {
	Name         string
	HashPassword command.Secret
	// Lineage is a rules.Lineage id, and the only choice creation makes. It
	// is cosmetic: there is no class to pick beside it, because what a
	// character is good at comes from the gear they put on.
	Lineage string
	// TODO other fields ... see command/commands.go
}

func (createHashed) Verb() string { return "createHashed" }
