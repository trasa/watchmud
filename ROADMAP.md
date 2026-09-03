# WatchMUD: strangler-fig migration to telnet

## Context

This started as a learn-Go project: gRPC/protobuf transport, a separate console client
(`trasa/watchmud-client`), and a Postgres database. The destination is different — a real
telnet server that tintin++ and other MUD clients can connect to, no custom client, and
some persistence layer that isn't the current SQL schema.

The current tree is further from working than it looks. `cmd/watchmud/main.go` builds the
world and runs the tick loop, but **nothing calls `rpc.NewServer`, `web.Start`, or
`db.Init`** — there is no listener and `watchdb` is a nil handle, so the first thing a
connection would do (login → `db.GetPlayerData`) fails. `make vet` has never been green
(copylocks, from passing protobuf structs by value), one `playergenerator` test fails, and
the race→lineage/species rework is half-landed with `PLACEHOLDER` values in
`server.handleCreatePlayer` and commented-out DB calls in `world/h_stat.go`.

The good news: **the strangler-fig seam already exists.** `client.Client` (Send / SetPlayer /
GetPlayer / Close) and `gameserver.Instance` (Receive / Logout) already abstract the
transport. `rpc/client.go` is one implementation. Telnet is a second one. Nothing in
`world/` needs to know which is talking.

Decisions taken up front:
- **Telnet is the first live listener.** `rpc/` and `web/` stay on disk, unwired, as
  reference; they get deleted once telnet covers their ground.
- **Persistence becomes an interface with an in-memory implementation.** The `db` package
  and Postgres go away now; the real replacement is chosen later, against a working server.
- **Protobuf stays as the internal vocabulary for now.** Telnet parses text into
  `message.XRequest` and renders `message.XResponse` back out. `world/` is untouched.
  Re-evaluated in Phase 5.
- **The player model is fixed before telnet is written**, so the telnet connection is
  written once against the final shape.

Ordering principle throughout: never be more than one small step away from a server you can
start and connect to.

---

## Phase 0 — A baseline you can trust

Nothing else is measurable until `go test ./...` means something.

- `playergenerator/generator_test.go:35` fails (expects 16, got 15) — almost certainly
  fallout from the lineage/species rework in `e09628b`. Either fix the generator or pin the
  test to current intent; don't leave it red as background noise. **DONE September 2 2026**
- `make fmt-check` silently always passes: it calls `gofmt -files`, which is not a real
  flag, and the target swallows the error. Change to `gofmt -l .`. **FIXED September 2 2026**
- **Leave `make vet` red.** Every failure is `copylocks` from protobuf structs passed by
  value (`p.Send(message.LookResponse{...})`, `resp := x.(message.LookResponse)`). These
  disappear on their own if Phase 5 removes protobuf. Chasing them now is wasted work
  against code that may not survive. **NO ACTION REQUIRED September 2 2026**

Exit: `make test` green, and you know what red means again.
**Completed** September 2 2026**
All tests pass, `make fmt-check` succeeds, `make vet` still fails as expected.

---

## Phase 1 — Collapse `player.Player` to a struct

**The problem.** `player.Player` (`player/player.go`) is a 30-method interface with exactly
two implementations: `server.ClientPlayer` and `player.TestPlayer`. One of those is a test
double. The other lives in the wrong package. This is a C#/Java habit — "program to an
interface" — where Go wants a concrete type.

**Why it ended up that way.** `client` imports `player` (its `GetPlayer`/`SetPlayer` return
`player.Player`). A concrete player that holds a `client.Client` would make `player` import
`client` — an import cycle. The interface was the escape hatch. Verified: `player/` imports
only `combat` and `object` today.

**The Go-idiomatic fix.** Define the *narrow* interface at the point of need, inside
`player`, and depend on that instead:

```go
package player

// Sender is anything that can deliver a message to this player's connection.
type Sender interface {
    Send(msg interface{}) error
}

type Player struct {
    Id        int64
    Name      string
    out       Sender          // was: client.Client, via ClientPlayer
    inventory *Inventory
    slots     *object.Slots
    // curHealth, maxHealth, dirty, location, abilities ...
}

func New(name string, out Sender, ...) *Player
func (p *Player) Send(msg interface{}) error { return p.out.Send(msg) }
```

`client.Client` satisfies `player.Sender` structurally — no import, no cycle. The future
telnet connection satisfies it too, which is the whole point of doing this before Phase 4.

Work:
- New concrete `player.Player` struct in `player/player.go`, absorbing the fields and all 30
  method bodies from `server/clientplayer.go`.
- **Delete `server/clientplayer.go` and `player/testplayer.go`.** Replace the test double
  with a tiny `player.Recorder` implementing `Sender` and capturing sent messages; tests
  become `p := player.New("testdood", rec)` and `rec.Sent(0).(message.LookResponse)`,
  replacing today's `p.GetSentResponse(0)`. Same for `client.TestClient`, which mostly
  duplicates it.
- Mechanical sweep `player.Player` → `*player.Player` in `spaces/` (`room.go`,
  `roominventory.go`), `world/` (`playerroommap.go`, most `h_*.go`), `gameserver/handlerparameter.go`,
  `client/client.go`.
- `player/players.go`: `List` maps keyed on `*Player`. Its `sync.RWMutex` is vestigial —
  world state is single-goroutine by design (see the comment in `Room.CreateRoomDescription`)
  and stays that way after telnet, since connection goroutines only push to a channel. Leave
  the mutex or drop it, but don't let it imply the world is concurrent.
- **Keep `combat.Combatant`.** That interface is real polymorphism — `*player.Player` and
  `*mobile.Instance` both fight, and `combat/melee.go` genuinely must not care which. Don't
  collapse it along with the others.

Note after this lands: `client.Client` is now nearly redundant — `player.Sender` plus a
`Close()`. Worth deleting in Phase 4 rather than pre-emptively here.

Exit: `make test` green, `player.Player` is a struct, one fewer package in the cycle.

---

## Phase 2 — Persistence behind an interface

Files: new `player/store.go`, new `memstore/`, delete `db/`.

```go
package player

type Record struct {
    Id                          int64
    Name                        string
    CurHealth, MaxHealth        int64
    LineageId, ClassId          string   // string ids, not the int32s (see Phase 6)
    LastZoneId, LastRoomId      string
    Abilities                   Abilities
    Slots                       []SlotRecord
    Inventory                   []InventoryRecord
}

type Store interface {
    Load(name string) (*Record, bool, error)
    Create(r *Record) (*Record, error)
    Save(r *Record) error
}
```

`Record` mirrors what `db.PlayerData` + `db.PlayerInventoryData` + `db.SlotDataList` carry
today (see `db/sql/ddl.sql` for the full field list) — but **use string lineage/class ids
now**, even though nothing consumes them until Phase 6. Writing `int32` race ids into a new
store just to migrate them again is wasted motion.

Rewire the four call sites: `server.handleLogin`, `server.handleCreatePlayer`,
`world.HandleIncomingMessage` (the save-after-every-handler), and `world/h_logout.go`. Pass
the `Store` into `server.New` and `world.New` rather than reaching for a package global —
`db.watchdb` being a package-level var is part of why none of this is testable today.

Then delete `db/` and drop `sqlx`, `lib/pq`, `go-sql-driver/mysql`, `mattn/go-sqlite3`, and
the `golang.org/x/crypto/ssh` tunnel from `go.mod`; strip `DB` and `SSH` from
`serverconfig.Config` and `app.local.yaml`.

`memstore` is a `map[string]*player.Record` behind the interface. Players don't survive
restart. That is fine and it unblocks login today.

One thing to notice while you're here: `HandleIncomingMessage` saves the player after
*every* message. Free against a map, absurd against a real store. Leave it; the interface
means the eventual implementation can batch or debounce without touching `world/`.

Exit: no SQL in the tree, login path runs without a database.

---

## Phase 3 — Make the loop answer typed input promptly

`server/gameserver.go`: `incomingBuffer` is unbuffered and is only drained *inside*
`heartbeat`, which runs on the 1s `mudtime.PulseInterval` ticker. A telnet user would wait
up to a full second between pressing enter and seeing output. That is the difference between
a MUD that feels alive and one that feels broken, and you want it fixed before your first
telnet session forms an impression.

Restructure `Run` to select on both:

```go
for {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case msg := <-gs.incomingBuffer:
        gs.dispatch(msg)          // handle immediately
    case <-ticker.C:
        pulse++
        gs.heartbeat(pulse, delta) // zone/mobile/violence pulses only
    }
}
```

Move the message-draining loop out of `heartbeat`, leaving it purely for pulse work. This
preserves the single-goroutine invariant exactly — commands are still processed one at a
time, never concurrently with a pulse. Give `incomingBuffer` a modest buffer (say 64) so a
burst of input doesn't block connection goroutines.

Exit: command latency is bounded by processing, not by the tick.

---

## Phase 4 — The telnet listener

New package `telnet/`. Read `rpc/client.go` first — it is the template for what you're
writing, and `readPump`/`writePump` map almost directly.

**Connection** — `telnet.conn` implements `player.Sender` and holds a `net.Conn`:
- `Listen(addr string, gs gameserver.Instance)`: accept loop, goroutine per connection.
- Read side: `bufio.Scanner`/`Reader` over lines, strip `\r\n`, then parse → `message.XRequest`
  → `message.NewGameMessage` → `gs.Receive(gameserver.NewHandlerParameter(c, gm))`.
- Write side: keep the `sendQueue chan` + write-pump shape from `rpc/client.go`; it already
  handles the "world goroutine must never block on a slow socket" problem.
- On read error/EOF: `gs.Logout(c, cause)`, same as `rpc`.
- At minimum, strip `IAC` (0xFF) sequences from input now so a real client's opening
  negotiation doesn't get parsed as a command. Full negotiation is Phase 7.

**Login state machine** — before the connection is in the world, it isn't yet a player. A
small per-connection state (`awaitingName` → `awaitingConfirmNew` → `inGame`) that emits
`message.LoginRequest` / `message.CreatePlayerRequest` into the same pipe. Everything
downstream is unchanged.

**Command table** — new, and the piece with no existing analogue in this repo (it lived in
the client). Verb, aliases, and minimum abbreviation, mapping to a request constructor:

```go
{"north", []string{"n"}, func(args string) any { return message.MoveRequest{Direction: direction.North} }},
{"look",  []string{"l"}, ...},
{"get",   nil,           func(args string) any { return message.GetRequest{Target: args} }},
```

**Do not parse targets here.** `world.parseTarget` (`world/target_parser.go`) already handles
`all`, `all.knife`, `2.knife`, `20 coins`, and it is exercised by
`world/target_parser_test.go`. Pass the raw argument string through in the request and let
the server parse it, exactly as the gRPC client does today. Duplicating that grammar in the
transport is how the two drift apart.

**Renderer** — the substantial new work: one type switch over every `message.*Response` and
`message.*Notification`, producing text. `spaces.Room.CreateRoomDescription` already returns
a structured `message.RoomDescription` (name, description, exits, players, objects, mobs) —
the renderer turns that into the classic room block. This type switch is the single
chokepoint that Phase 5 would rewrite, which is why it is worth keeping it in one file and
free of game logic.

Wire it in `cmd/watchmud/main.go` alongside the existing setup, on a new `telnetPort` in
`app.local.yaml`. Delete `client/` (the interface serving the dead console client) and
`client.TestClient` once nothing references them.

Exit: `telnet localhost 4000`, log in, `look`, `north`, `get knife`, `inventory`, `who`, `quit`.

---

## Phase 5 — Decide protobuf's fate

Now, with the renderer written, judge it on evidence rather than taste. Protobuf earns its
place only if you want a second non-text transport (a web client, GMCP structured data). If
telnet is the only consumer, `message.GameMessage`'s oneof wrapper and
`message.DecodeTypeName` string-keyed dispatch are pure overhead — and every `copylocks`
error from Phase 0 is a symptom of value-copying types that were never meant to be copied.

If it goes, the path is mechanical and the blast radius is known: 39 files in `world/`, 6 in
`spaces/`, 4 in `object/`. Plain structs in a new `event` (out) and `command` (in) package;
`world.handlerMap` keys on a real type or enum instead of a decoded string; the renderer's
type switch is retargeted; `github.com/trasa/watchmud-message` and all of protobuf/gRPC leave
`go.mod`. `rpc/` and `web/` get deleted here if they haven't already.

---

## Phase 6 — Finish lineage/species, retire the int32 ids

The `e09628b` rework left `rules.Catalog` (`rules/catalog.go`) loaded from
`content/rules/{species,classes}.json` but connected to nothing. Meanwhile
`player.Player` still exposes `GetRaceId() int32` / `GetClassId() int32`, pointing at SQL
rows deleted in Phase 2, and `server.handleCreatePlayer` constructs a `rules.Lineage` full of
`"PLACEHOLDER"` strings.

- Replace those two methods with `LineageId`/`ClassId` strings resolved via
  `Catalog.Lineages[id]` / `Catalog.Classes[id]` — the `Record` from Phase 2 already stores
  them that way.
- Fix `world/h_stat.go`: the commented-out `db.GetSingleRaceData` block becomes a catalog
  lookup, filling the `Race`/`Class` fields currently hardcoded to `""`.
- Fix `handleCreatePlayer` to take a real lineage and class, and give `playergenerator` the
  catalog so ability generation reflects actual bonuses (likely the root of the Phase 0 test
  failure).
- Character creation over telnet — choose species → lineage → class — is a natural extension
  of the Phase 4 login state machine.

---

## Phase 7 — Real MUD-client protocol support

Additive once the byte loop exists:
- TELNET option negotiation (`IAC WILL/WONT/DO/DONT`), `ECHO` off for password entry.
- ANSI color, with a per-player toggle.
- `NAWS` (window size) for wrapping; wrap output to the client's width.
- Then the MUD-specific layer as it earns its keep: `MSSP`, `GMCP`, `MCCP`, `MXP`.
- Test against tintin++ specifically, since that's the target.

---

## Known problems deliberately *not* scheduled

Named so they don't get rediscovered as surprises:

- **Dual location bookkeeping.** A player's room is recorded in *both* `world.PlayerRoomMap`
  and the `Room`'s own `playerList` (plus `p.Location()`), and `World.movePlayer` must update
  all of them in step. Same pattern for `spaces.MobileRoomMap`. Any missed update silently
  desyncs the world. Worth collapsing to a single source of truth — after telnet works.
- **`web/`** serves a static page for a client that no longer exists. Delete with `rpc/`.
- **`world/settings.go`** is a single `VERBOSE_LOGGING` const, and logging is split between
  zerolog and stdlib `log` depending on file age. Worth one consolidating pass eventually.

---

## Verification

Per phase:
- **0–3:** `make test` green after each. Phase 3 additionally: server starts and idles
  without pegging a core.
- **4:** manual `telnet localhost 4000` — create a player, `look`, move between rooms in
  `content/world/wrathrock`, `get`/`drop`/`inventory`/`wear`, `kill` a mob from
  `content/world/sample`, `who`, `quit`, reconnect. Two simultaneous connections to confirm
  `say`/`tell` notifications reach the other session. Watch that mob wandering (10s pulse)
  and zone reset (3min lifetime in `zone_manifest.json`) still fire while a client is idle.
- **5:** `make vet` should be green for the first time — that's the signal the protobuf
  value-copying is actually gone, not just hidden.
- **7:** connect with tintin++ rather than raw telnet; verify negotiation, wrapping, color.

Regression net worth adding early (cheap, and it makes every later phase safer): a table test
that feeds command strings through the Phase 4 parser into a `newTestWorld()` and asserts on
rendered output text. That covers parser, dispatch, and renderer in one pass without a socket,
and it survives the protobuf decision in Phase 5 unchanged.
