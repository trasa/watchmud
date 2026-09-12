package telnet

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/rules"
)

type conn struct {
	gs         gameserver.Instance
	cat        *rules.Catalog // read-only; the creation menu is built from it
	netConn    net.Conn
	scanner    *bufio.Scanner
	sendQueue  chan any // NOT *message.GameMessage - see below.
	quit       chan struct{}
	closeOnce  sync.Once
	authResult chan bool // buffered 1; carries LoginResponse/CreatePlayerResponse success

	mu     sync.Mutex // guards the player
	player *player.Player
}

func newConn(c net.Conn, gs gameserver.Instance, cat *rules.Catalog) *conn {
	const sendQueueSize = 256
	return &conn{
		gs:         gs,
		cat:        cat,
		netConn:    c,
		sendQueue:  make(chan any, sendQueueSize),
		quit:       make(chan struct{}),
		authResult: make(chan bool, 1),
	}
}

// Listen takes the catalog as well as the game server because character
// creation is a conversation held on this side of the seam, before there is a
// player to hand a command to, and a menu of lineages is presentation. The
// renderer already depends on rules for the same reason.
func Listen(ctx context.Context, addr string, gs gameserver.Instance, cat *rules.Catalog) error {
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
		c := newConn(nc, gs, cat)
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
	// the login conversation, below, is waiting on these four.
	switch msg.(type) {
	case event.LoggedIn, event.PlayerCreated:
		c.signalAuth(true)
		return nil
	case event.LoginFailed, event.CreateFailed:
		c.signalAuth(false)
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
		c.emit(command.Login{Name: name})
		if c.awaitAuth() {
			return true
		}
		// PLAYER_LOGIN_FAILED: no such player
		yn, ok := c.prompt(fmt.Sprintf("No one by the name of %s. Create them? (yn) ", name))
		if !ok {
			return false
		}
		if strings.HasPrefix(strings.ToLower(yn), "y") {
			lineage, ok := c.chooseLineage()
			if !ok {
				return false
			}
			c.emit(command.CreatePlayer{Name: name, Lineage: lineage})
			if c.awaitAuth() {
				return true
			}
			c.Send("Something went wrong creating that character.\r\n")
		}

	}
}

// chooseLineage is the whole of character creation. There is no class step
// after it and no ability step: a lineage decides nothing but how the
// character is described, and what they are good at is decided later, by what
// they pick up and put on.
//
// The player answers with a number or with any unambiguous part of a lineage
// name, and an unrecognized answer re-asks rather than failing the creation.
func (c *conn) chooseLineage() (string, bool) {
	if c.cat == nil {
		return "", true // no catalog: the server picks the default
	}
	choices := lineageChoices(c.cat)
	if len(choices) == 0 {
		return "", true
	}

	c.Send(lineageMenu(c.cat))
	for {
		answer, ok := c.prompt("Which lineage? ")
		if !ok {
			return "", false
		}
		if id, found := matchLineage(choices, answer); found {
			return id, true
		}
		c.Send("That isn't one of them. Type a number, or the name.\r\n")
	}
}

// lineageChoices flattens the species tree into the numbered list the menu
// shows, in content order so the numbers are stable between sessions.
func lineageChoices(cat *rules.Catalog) []*rules.Lineage {
	var out []*rules.Lineage
	for _, s := range cat.SpeciesList() {
		out = append(out, s.Lineages...)
	}
	return out
}

// lineageMenu groups the numbered choices under their species, since the
// species is the only thing the grouping is still for.
func lineageMenu(cat *rules.Catalog) string {
	var b strings.Builder
	b.WriteString("\nChoose a lineage. It decides nothing but how you look;\n")
	b.WriteString("what you're good at comes from what you carry.\n")
	n := 0
	for _, s := range cat.SpeciesList() {
		b.WriteString("\n" + s.Name + "\n")
		for _, l := range s.Lineages {
			n++
			if l.Description != "" {
				fmt.Fprintf(&b, "  %2d) %-20s %s\n", n, l.Name, l.Description)
			} else {
				fmt.Fprintf(&b, "  %2d) %s\n", n, l.Name)
			}
		}
	}
	b.WriteString("\n")
	return b.String()
}

// matchLineage accepts the menu number, the lineage id, or a case-insensitive
// prefix of the name -- but only when exactly one lineage matches it, so
// "h" doesn't silently pick Hill Dwarf over High Elf.
func matchLineage(choices []*rules.Lineage, answer string) (string, bool) {
	answer = strings.TrimSpace(answer)
	if n, err := strconv.Atoi(answer); err == nil {
		if n >= 1 && n <= len(choices) {
			return choices[n-1].Id, true
		}
		return "", false
	}

	lower := strings.ToLower(answer)
	var match *rules.Lineage
	for _, l := range choices {
		if strings.EqualFold(l.Id, answer) || strings.EqualFold(l.Name, answer) {
			return l.Id, true // an exact hit beats any number of prefixes
		}
		if strings.HasPrefix(strings.ToLower(l.Name), lower) {
			if match != nil {
				return "", false // ambiguous
			}
			match = l
		}
	}
	if match == nil {
		return "", false
	}
	return match.Id, true
}

// emit hands a command to the game server.
// It is the only path from this connection into the world.
func (c *conn) emit(cmd command.Command) {
	c.gs.Receive(gameserver.NewHandlerParameter(c, cmd))
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
		name = c.Player().Name()
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
		cmd, err := parseCommand(strings.Fields(line))
		if err != nil {
			c.Send(err.Error() + "\r\n")
			continue
		}
		c.emit(cmd)
		if _, quitting := cmd.(command.Logout); quitting {
			c.Send("Goodbye.\r\n")
			return
		}
	}
}
