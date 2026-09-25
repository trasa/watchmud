package telnet

import (
	"bufio"
	"context"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// dial connects and returns the first line the server says.
func dial(t *testing.T, addr string) (net.Conn, string) {
	t.Helper()
	nc, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	t.Cleanup(func() { _ = nc.Close() })
	require.NoError(t, nc.SetReadDeadline(time.Now().Add(2*time.Second)))
	line, err := bufio.NewReader(nc).ReadString('\n')
	require.NoError(t, err)
	return nc, strings.TrimSpace(line)
}

// One address gets a few connections -- a household, a second window -- and
// no more, so one client can't eat every slot.
func TestServe_capsConnectionsPerAddress(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go func() { _ = serve(ctx, ln, &fakeServer{passwords: map[string]string{}}, nil, 2) }()
	addr := ln.Addr().String()

	first, greeting := dial(t, addr)
	assert.Equal(t, "Welcome to WatchMUD.", greeting)
	_, greeting = dial(t, addr)
	assert.Equal(t, "Welcome to WatchMUD.", greeting)

	refused, greeting := dial(t, addr)
	assert.Equal(t, "Too many connections from your address. Try again later.", greeting)
	_, err = bufio.NewReader(refused).ReadString('\n')
	assert.Error(t, err, "and then hung up on")

	// a slot comes back when a connection ends
	require.NoError(t, first.Close())
	require.Eventually(t, func() bool {
		nc, err := net.Dial("tcp", addr)
		if err != nil {
			return false
		}
		defer nc.Close()
		_ = nc.SetReadDeadline(time.Now().Add(time.Second))
		line, _ := bufio.NewReader(nc).ReadString('\n')
		return strings.HasPrefix(line, "Welcome")
	}, 2*time.Second, 20*time.Millisecond)
}
