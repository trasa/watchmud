package telnet

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	message "github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/player"
)

type conn struct {
	gs         gameserver.Instance
	netConn    net.Conn
	scanner    *bufio.Scanner
	sendQueue  chan any // NOT *message.GameMessage - see below.
	quit       chan struct{}
	closeOnce  sync.Once
	authResult chan bool // buffered 1; carries LoginResponse/CreatePlayerResponse success

	mu     sync.Mutex // guards the player
	player *player.Player
}

func newConn(c net.Conn, gs gameserver.Instance) *conn {
	const sendQueueSize = 256
	return &conn{
		gs:         gs,
		netConn:    c,
		sendQueue:  make(chan any, sendQueueSize),
		quit:       make(chan struct{}),
		authResult: make(chan bool, 1),
	}
}

func Listen(ctx context.Context, addr string, gs gameserver.Instance) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("telnet listen on %s: %w", addr, err)
	}
	log.Info().Msgf("telnet listening on %s", addr)

	go func() {
		<-ctx.Done()
		ln.Close() // unblocks the Accept, below
	}()

	for {
		nc, err := ln.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil // shutting down, expected
			}
			return fmt.Errorf("telnet accept: %w", err)
		}
		log.Info().Msgf("telnet connection from %s", nc.RemoteAddr())
		c := newConn(nc, gs)
		go c.writePump()
		go c.readPump()
		c.Send("Welcome to WatchMUD.\r\n")
	}
}

func (c *conn) Player() *player.Player {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.player
}

func (c *conn) SetPlayer(p *player.Player) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.player = p
}

func (c *conn) Send(msg any) error {
	switch m := msg.(type) {
	case message.LoginResponse:
		c.signalAuth(m.Success)
	case message.CreatePlayerResponse:
		c.signalAuth(m.Success)
	}
	select {
	case c.sendQueue <- msg:
		return nil
	default:
		c.Close()
		return errors.New("send queue full")
	}
}

func (c *conn) signalAuth(ok bool) {
	select {
	case c.authResult <- ok:
	default: // don't block if the authResult channel is full
	}
}

func (c *conn) awaitAuth() bool {
	select {
	case ok := <-c.authResult:
		return ok
	case <-c.quit:
		return false
	}
}

func (c *conn) login() bool {
	for {
		name, ok := c.prompt("By what name do you wish to be known? ")
		if !ok {
			return false // disconnected
		}
		if err := c.emit(message.LoginRequest{PlayerName: name}); err != nil {
			return false
		}
		if c.awaitAuth() {
			return true
		}
		// PLAYER_LOGIN_FAILED: no such player
		yn, ok := c.prompt(fmt.Sprintf("No one by the name of %s. Create them? (yn) ", name))
		if !ok {
			return false
		}
		if strings.HasPrefix(strings.ToLower(yn), "y") {
			if err := c.emit(message.CreatePlayerRequest{PlayerName: name}); err != nil {
				return false
			}
			if c.awaitAuth() {
				return true
			}
			c.Send("Something went wrong creating that character.\r\n")
		}

	}
}

// emit wraps a request and hands it to the game server.
// It is the only path from this connection into the world.
func (c *conn) emit(req any) error {
	gm, err := message.NewGameMessage(req)
	if err != nil {
		log.Error().Err(err).Msgf("telnet %s: cannot wrap %T", c.netConn.RemoteAddr(), req)
		return err
	}
	c.gs.Receive(gameserver.NewHandlerParameter(c, gm))
	return nil
}

// prompt writes text with no trailing newline, then waits
// for a reply. A bare Enter re-issues the prompt rather
// than returning an empty string.
func (c *conn) prompt(text string) (string, bool) {
	for {
		if err := c.Send(text); err != nil {
			return "", false // queue full; conn is being torn down
		}
		line, ok := c.readLine()
		if !ok {
			return "", false
		}
		if line != "" {
			return line, true
		}
	}
}

// readLine returns the next line from the client. ok is false once the
// connection is finished: EOF, read error, Close.
func (c *conn) readLine() (string, bool) {
	if !c.scanner.Scan() {
		return "", false
	}
	return strings.TrimSpace(c.scanner.Text()), true
}

func (c *conn) Close() {
	c.closeOnce.Do(func() {
		close(c.quit)
		// doesn't close netConn, since the writePump may still
		// have leftover data to send to the client as we're saying
		// goodbye.  we "quit", it drains the buffer, writes, and then
		// netConn.Close().
	})
}

const writeTimeout = 10 * time.Second

func (c *conn) writePump() {
	defer c.netConn.Close() // this unblocks a parked readPump
	for {
		select {
		case msg := <-c.sendQueue:
			if err := c.write(msg); err != nil {
				log.Warn().Err(err).Msgf("telnet %s write error", c.netConn.RemoteAddr())
				c.Close()
				return
			}
		case <-c.quit:
			// drain what's already queued so the goodbye actually lands
			for {
				select {
				case msg := <-c.sendQueue:
					if err := c.write(msg); err != nil {
						return
					}
				default:
					return
				}
			}
		}
	}
}

func (c *conn) write(msg any) error {
	// STEP C: this type switch becomes the renderer — one case per
	// message.XResponse / message.XNotification. Keep it in its own file
	// and free of game logic; it's the chokepoint Phase 5 would retarget.
	var text string
	switch m := msg.(type) {
	case string:
		text = m
	default:
		text = fmt.Sprintf("%v\r\n", m)
	}
	if err := c.netConn.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
		return err
	}
	_, err := io.WriteString(c.netConn, text)
	return err
}

func (c *conn) readPump() {
	defer c.Close()
	c.scanner = bufio.NewScanner(&iacFilter{src: bufio.NewReader(c.netConn)})
	if c.login() {
		c.commandLoop()
	}
	cause := "client disconnected"
	if err := c.scanner.Err(); err != nil {
		cause = fmt.Sprintf("read error: %v", err)
	}
	log.Info().Msgf("telnet %s: %s", c.netConn.RemoteAddr(), cause)
	c.gs.Logout(c, cause)
}

func (c *conn) commandLoop() {

}
