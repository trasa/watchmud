package player

import (
	"errors"
	"strings"
)

const (
	MinNameLength = 3
	MaxNameLength = 16
)

var (
	ErrNameLength  = errors.New("a name is 3 to 16 letters")
	ErrNameLetters = errors.New("a name is letters and nothing else")
)

// CanonicalName is the one form a name is stored and shown in: "bOB" is Bob.
// It fails for anything that can't be a name at all. Whether the name is
// free, or reserved by the world, is someone else's question.
//
// Letters means a-z. Anything wider and "Bob" has look-alikes a player could
// register to be mistaken for him, and a name nobody can type can't be told.
func CanonicalName(raw string) (string, error) {
	name := strings.TrimSpace(raw)
	if len(name) < MinNameLength || len(name) > MaxNameLength {
		return "", ErrNameLength
	}
	for _, r := range name {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return "", ErrNameLetters
		}
	}
	return strings.ToUpper(name[:1]) + strings.ToLower(name[1:]), nil
}

// NameKey is what names are compared by, so "bob" typed at a prompt or in a
// tell finds Bob.
func NameKey(name string) string {
	return strings.ToLower(name)
}
