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
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/player"
)

type conn struct {
	gs        gameserver.Instance
	netConn   net.Conn
	sendQueue chan any // NOT *message.GameMessage - see below.
	quit      chan struct{}
	closeOnce sync.Once

	mu     sync.Mutex // guards the two fields below, and only those
	player *player.Player
	//state  loginState
}

func newConn(c net.Conn, gs gameserver.Instance) *conn {
	const sendQueueSize = 256
	return &conn{
		gs:        gs,
		netConn:   c,
		sendQueue: make(chan any, sendQueueSize),
		quit:      make(chan struct{}),
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
	select {
	case c.sendQueue <- msg:
		return nil
	default:
		c.Close()
		return errors.New("send queue full")
	}
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

	scanner := bufio.NewScanner(&iacFilter{src: bufio.NewReader(c.netConn)})
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		// STEP B: login state machine, then command table ->
		// message.XRequest -> message.NewGameMessage ->
		// c.gs.Receive(gameserver.NewHandlerParameter(c, gm))
		if err := c.Send("you said: " + line + "\r\n"); err != nil {
			return
		}
	}
	cause := "client disconnected"
	if err := scanner.Err(); err != nil {
		cause = fmt.Sprintf("read error: %v", err)
	}
	log.Info().Msgf("telnet %s: %s", c.netConn.RemoteAddr(), cause)
	c.gs.Logout(c, cause)
}
