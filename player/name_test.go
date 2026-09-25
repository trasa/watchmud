package player

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
)

func TestCanonicalName(t *testing.T) {
	cases := []struct {
		raw, want string
		problem   error
	}{
		{"bob", "Bob", nil},
		{"BOB", "Bob", nil},
		{"bOb", "Bob", nil},
		{"  bob  ", "Bob", nil}, // readLine trims, but don't count on it
		{"abc", "Abc", nil},
		{"abcdefghijklmnop", "Abcdefghijklmnop", nil}, // 16
		{"ab", "", ErrNameLength},
		{"", "", ErrNameLength},
		{"abcdefghijklmnopq", "", ErrNameLength}, // 17
		{"bob2", "", ErrNameLetters},
		{"bo b", "", ErrNameLetters},
		{"bob-o", "", ErrNameLetters},
		{"böb", "", ErrNameLetters}, // letters, but not ones everyone can type
		{"Ωmega", "", ErrNameLetters},
	}
	for _, tc := range cases {
		t.Run(tc.raw, func(t *testing.T) {
			got, err := CanonicalName(tc.raw)
			assert.ErrorIs(t, err, tc.problem)
			if tc.problem == nil {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

// Whatever case they're typed in, names find the same player.
func TestList_findsByNameInAnyCase(t *testing.T) {
	l := NewList()
	bob := NewTestPlayer(uuid.New(), "Bob", nil)
	l.Add(bob)

	for _, typed := range []string{"Bob", "bob", "BOB"} {
		assert.Same(t, bob, l.FindByName(typed), typed)
	}
	assert.Nil(t, l.FindByName("bobby"))
}
