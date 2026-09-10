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
		_ = ln.Close() // unblocks the Accept, below
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

// Send sends a message for the player.Sender interface
func (c *conn) Send(msg any) {
	_ = c.send(msg)
}

// Send handles sending messages to the telnet client.
// This can return an error if there's a legitimate sending error,
// which makes it distinctive from player.Sender.Send.
// Hooks into the login state signaling completion of login/create.
func (c *conn) send(msg any) error {
	switch m := msg.(type) {
	case message.LoginResponse:
		c.signalAuth(m.Success)
		return nil
	case message.CreatePlayerResponse:
		c.signalAuth(m.Success)
		return nil
	}
	select {
	case c.sendQueue <- msg:
		return nil
	default:
		log.Warn().Msgf("telnet %s: send queue full, closing", c.netConn.RemoteAddr())
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
		// log AND error, because this is a bad programming error (missing entry in table)
		// and we want that to stand out...
		log.Error().Err(err).Msgf("telnet %s: cannot wrap %T", c.netConn.RemoteAddr(), req)
		return err
	}
	c.dispatch(gm)
	return nil
}

func (c *conn) dispatch(gm *message.GameMessage) {
	c.gs.Receive(gameserver.NewHandlerParameter(c, gm))
}

// prompt writes text with no trailing newline, then waits
// for a reply. A bare Enter re-issues the prompt rather
// than returning an empty string.
func (c *conn) prompt(text string) (string, bool) {
	for {
		if err := c.send(text); err != nil {
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
	defer func() {
		// this unblocks a parked readPump
		_ = c.netConn.Close()
	}()

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
	name := ""
	if c.Player() != nil {
		name = c.Player().Name
	}
	text := render(msg, name)
	text = strings.ReplaceAll(text, "\r\n", "\n") // normalize
	text = strings.ReplaceAll(text, "\n", "\r\n") // replace with \r\n
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
	for {
		line, ok := c.readLine()
		if !ok {
			return
		}
		if line == "" {
			continue // bare Enter: ignore. Note this is the opposite of prompt(), which re-asks. Different context, different policy.
		}
		tokens := message.Tokenize(line)
		switch strings.ToLower(tokens[0]) {
		case "quit":
			_ = c.emit(message.LogoutRequest{Cause: "quit"})
			c.Send("Goodbye.\r\n")
			return

		case "drop":
			// TODO special case handling because the message defn doesn't protect against invalid drop messages
			// fix later.
			if len(tokens) < 2 {
				c.Send("Drop what?\r\n")
				continue
			}
		}
		gm, err := message.TranslateLineToMessage(tokens)
		if err != nil {
			c.Send(err.Error() + "\r\n")
			continue
		}
		c.dispatch(gm)
	}
}
