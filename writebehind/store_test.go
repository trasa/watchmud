package writebehind

import (
	"context"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trasa/watchmud/memstore"
	"github.com/trasa/watchmud/player"
)

// gatedStore holds every write until the test lets it through, so a test can
// catch the writer mid-save. Writes are counted by the writer goroutine and
// read only after Close, which is what makes that safe.
type gatedStore struct {
	*memstore.Store
	release chan struct{}
	saves   int
}

func newGated() *gatedStore {
	return &gatedStore{
		Store:   memstore.New(),
		release: make(chan struct{}),
	}
}

func (g *gatedStore) Save(r *player.Record) error {
	<-g.release
	g.saves++
	return g.Store.Save(r)
}

func TestSaveReachesTheInnerStore(t *testing.T) {
	inner := memstore.New()
	s := New(inner)
	require.NoError(t, s.Save(&player.Record{Id: uuid.New(), Name: "dood", CurHealth: 5}))
	require.NoError(t, s.Close(context.Background()))

	got, found, err := inner.Load("dood")
	require.NoError(t, err)
	require.True(t, found)
	assert.Equal(t, 5, got.CurHealth)
}

// quit and log straight back in: the record is still on its way to the
// inner store, and Load has to find it anyway.
func TestLoadSeesARecordStillBeingWritten(t *testing.T) {
	inner := newGated()
	s := New(inner)
	require.NoError(t, s.Save(&player.Record{Id: uuid.New(), Name: "dood", CurHealth: 5}))

	got, found, err := s.Load("dood") // parked in inner.Save
	require.NoError(t, err)
	require.True(t, found)
	assert.Equal(t, 5, got.CurHealth)

	close(inner.release)
	require.NoError(t, s.Close(context.Background()))
}

// Saves that pile up behind a slow write become one write of the newest.
func TestSavesBehindASlowWriteCollapse(t *testing.T) {
	inner := newGated()
	s := New(inner)
	id := uuid.New()

	require.NoError(t, s.Save(&player.Record{Id: id, Name: "dood", CurHealth: 1}))
	for hp := 2; hp <= 10; hp++ {
		require.NoError(t, s.Save(&player.Record{Id: id, Name: "dood", CurHealth: hp}))
	}
	close(inner.release)
	require.NoError(t, s.Close(context.Background()))

	got, _, _ := inner.Load("dood")
	assert.Equal(t, 10, got.CurHealth, "newest one wins")
	assert.LessOrEqual(t, inner.saves, 2, "not ten writes")
}
