package telnet

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *session) waitClosed() {
	s.t.Helper()
	select {
	case <-s.closed:
	case <-time.After(2 * time.Second):
		s.t.Fatalf("still connected; transcript %q", s.transcript())
	}
}

// Sitting at the name prompt costs a connection slot and nothing else, so
// that wait is short.
func TestIdle_atLogin(t *testing.T) {
	gs := &fakeServer{passwords: map[string]string{}}
	s := startSessionWith(t, gs, func(c *conn) {
		c.loginIdle = 50 * time.Millisecond
		c.playIdle = time.Hour
	})

	s.waitClosed()
	assert.Contains(t, s.transcript(), "Idle too long. Goodbye.")
}

// In the game the wait is long, and ends the way quitting does: the server
// is told, so the character is saved and leaves the world.
func TestIdle_inGame(t *testing.T) {
	gs := &fakeServer{passwords: map[string]string{"Bob": "sekrit99"}}
	s := startSessionWith(t, gs, func(c *conn) {
		c.loginIdle = time.Hour
		c.playIdle = 50 * time.Millisecond
	})
	s.answer("known? ", "Bob")
	s.answer("Password: ", "sekrit99")

	s.waitClosed()
	assert.Contains(t, s.transcript(), "Idle too long. Goodbye.")
	loggedOut, cause := gs.logoutCause()
	require.True(t, loggedOut)
	assert.Equal(t, "idle", cause)
}

// Each line starts the clock again: a player typing now and then is not idle,
// however long the session.
func TestIdle_typingResetsTheClock(t *testing.T) {
	gs := &fakeServer{passwords: map[string]string{}}
	s := startSessionWith(t, gs, func(c *conn) {
		c.loginIdle = 150 * time.Millisecond
	})

	for range 4 { // 4 × 80ms: well past one timeout, never past it between lines
		time.Sleep(80 * time.Millisecond)
		s.answer("known? ", "xy") // refused, and asked again
	}
	select {
	case <-s.closed:
		t.Fatalf("hung up on a player who was typing: %q", s.transcript())
	default:
	}
}
