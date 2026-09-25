// Package writebehind is a player.Store that takes saves off of the world
// goroutine. Save queues a record and returns; one background goroutine
// writes them to the store it wraps around.
package writebehind

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"uuid"

	"github.com/rs/zerolog/log"
	"github.com/watchmud/watchmud/player"
)

type Store struct {
	inner   player.Store
	mu      sync.Mutex
	pending map[uuid.UUID]*player.Record // newest unwritten record per player; guarded by mu
	written map[uuid.UUID]*player.Record // last record that reached inner; writer goroutine only

	wake chan struct{} // buffered 1: "there is something to write"
	quit chan struct{} // closed by Close
	done chan struct{} // closed by the writer once it has finished
}

// batchSaver is an inner store that can write many records in one go.
//
// failed lists the index in recs of every record that was not written, and
// err says why; everything not in failed was written. A batch that failed
// outright lists every index. One bad record must not count against the rest
// of the batch, or it gets rewritten on every flush forever.
type batchSaver interface {
	SaveAll(recs []*player.Record) (failed []int, err error)
}

func New(inner player.Store) *Store {
	s := &Store{
		inner:   inner,
		pending: make(map[uuid.UUID]*player.Record),
		written: make(map[uuid.UUID]*player.Record),
		wake:    make(chan struct{}, 1), // doorbell pattern
		quit:    make(chan struct{}),
		done:    make(chan struct{}),
	}
	go s.run()
	return s
}

func (s *Store) run() {
	defer close(s.done)
	for {
		select {
		case <-s.wake:
			s.flush()
		case <-s.quit:
			s.flush()
			return
		}
	}
}

func (s *Store) flush() {
	// copy queue under the lock
	s.mu.Lock()
	batch := make([]*player.Record, 0, len(s.pending))
	for _, r := range s.pending {
		batch = append(batch, r)
	}
	s.mu.Unlock()
	// now do the slow part without the lock:
	// copy the changed entries over to a new slice
	changed := make([]*player.Record, 0, len(batch))
	for _, r := range batch {
		if prev, ok := s.written[r.Id]; ok && reflect.DeepEqual(prev, r) {
			s.forget(r) // nothing changed since last write
			continue
		}
		changed = append(changed, r)
	}
	if bs, ok := s.inner.(batchSaver); ok {
		failed, err := bs.SaveAll(changed)
		unwritten := make(map[int]bool, len(failed))
		for _, i := range failed {
			unwritten[i] = true
		}
		if err != nil {
			log.Error().Err(err).Int("failed", len(failed)).Int("records", len(changed)).Msg("writebehind: batch save failed")
		}
		for i, r := range changed {
			if unwritten[i] {
				continue // stays pending; next wake tries again
			}
			s.written[r.Id] = r
			s.forget(r)
		}
		return
	}
	// save one at a time
	for _, r := range changed {
		if err := s.inner.Save(r); err != nil {
			// left in pending; next wake tries again
			log.Error().Err(err).Str("player", r.Name).Msg("writebehind: save failed")
			continue
		}
		s.written[r.Id] = r
		s.forget(r)
	}
}

// Forget takes r out of the queue, unless a newer record replaced it
// while r was being written, in which case that one still needs writing.
func (s *Store) forget(r *player.Record) {
	s.mu.Lock()
	if s.pending[r.Id] == r {
		delete(s.pending, r.Id)
	}
	s.mu.Unlock()
}

func (s *Store) Save(r *player.Record) error {
	if r == nil {
		return nil
	}
	s.mu.Lock()
	s.pending[r.Id] = r
	s.mu.Unlock()
	// if doorbell is quiet, send goes into the buffer.
	// if it's already rung, default branch skips the send instead of blocking.
	select {
	case s.wake <- struct{}{}:
	default: // bell rung; writer will see this too
	}
	return nil
}

// Load answers from the queue before the inner store: a record still
// waiting to be written is newer than anything in the inner store.
func (s *Store) Load(name string) (*player.Record, bool, error) {
	s.mu.Lock()
	for _, r := range s.pending {
		if r.Name == name {
			s.mu.Unlock()
			return r, true, nil
		}
	}
	s.mu.Unlock()
	return s.inner.Load(name)
}

func (s *Store) Close(ctx context.Context) error {
	close(s.quit)
	select {
	case <-s.done:
	case <-ctx.Done():
		return fmt.Errorf("writebehind: gave up waiting for the last writes: %w", ctx.Err())
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if n := len(s.pending); n > 0 {
		return fmt.Errorf("writebehind: %d records were never written", n)
	}
	return nil
}
