# WatchMUD: strangler-fig migration to telnet

## Status

**Phases 0-5 are complete** (September 12 2026). The server is a working telnet MUD:
`make run`, then `telnet localhost 4000`, create a character, and play. Protobuf is gone;
the vocabulary is `command/` in and `event/` out. Phase 6 -- finishing lineage/species --
is next.

The Context section below describes the tree as it was in September 2026, before any of
this landed. It is kept for its reasoning, not as a description of the present.

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
  reference; they get deleted once telnet covers their ground. **DONE** -- both deleted
  September 2026, along with `client/`.
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
All tests pass, `make fmt-check` succeeds, `make vet` still fails as expected.
**Completed** September 2 2026**

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
  method bodies from `server/clientplayer.go`. **DONE**
- **Delete `server/clientplayer.go` and `player/testplayer.go`.** Replace the test double
  with a tiny `player.Recorder` implementing `Sender` and capturing sent messages; tests
  become `p := player.New("testdood", rec)` and `rec.Sent(0).(message.LookResponse)`,
  replacing today's `p.GetSentResponse(0)`. Same for `client.TestClient`, which mostly
  duplicates it. **DONE**
- Mechanical sweep `player.Player` → `*player.Player` in `spaces/` (`room.go`,
  `roominventory.go`), `world/` (`playerroommap.go`, most `h_*.go`), `gameserver/handlerparameter.go`,
  `client/client.go`. **DONE**
- `player/players.go`: `List` maps keyed on `*Player`. Its `sync.RWMutex` is vestigial —
  world state is single-goroutine by design (see the comment in `Room.CreateRoomDescription`)
  and stays that way after telnet, since connection goroutines only push to a channel. Leave
  the mutex or drop it, but don't let it imply the world is concurrent. **DONE**
- **Keep `combat.Combatant`.** That interface is real polymorphism — `*player.Player` and
  `*mobile.Instance` both fight, and `combat/melee.go` genuinely must not care which. Don't
  collapse it along with the others. **DONE**

Note after this lands: `client.Client` is now nearly redundant — `player.Sender` plus a
`Close()`. Worth deleting in Phase 4 rather than pre-emptively here.

Exit: `make test` green, `player.Player` is a struct, one fewer package in the cycle.
**DONE September 4 2026**
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
store just to migrate them again is wasted motion. **DONE**

Rewire the four call sites: `server.handleLogin`, `server.handleCreatePlayer`,
`world.HandleIncomingMessage` (the save-after-every-handler), and `world/h_logout.go`. Pass
the `Store` into `server.New` and `world.New` rather than reaching for a package global —
`db.watchdb` being a package-level var is part of why none of this is testable today.
**DONE**

Then delete `db/` and drop `sqlx`, `lib/pq`, `go-sql-driver/mysql`, `mattn/go-sqlite3`, and
the `golang.org/x/crypto/ssh` tunnel from `go.mod`; strip `DB` and `SSH` from
`serverconfig.Config` and `app.local.yaml`. **DONE**

`memstore` is a `map[string]*player.Record` behind the interface. Players don't survive
restart. That is fine and it unblocks login today. **DONE**

One thing to notice while you're here: `HandleIncomingMessage` saves the player after
*every* message. Free against a map, absurd against a real store. Leave it; the interface
means the eventual implementation can batch or debounce without touching `world/`.

Exit: no SQL in the tree, login path runs without a database.
**DONE September 7 2026**
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
**DONE September 7 2026**
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
**DONE September 11 2026.** `telnet/` is `conn.go`, `protocol.go`, `render.go`,
`resultcode.go`, with tests for the renderer and the IAC filter.

What actually happened, where it differed from the plan above:

- **`client.Client` didn't just get deleted, it moved and got a setter.** It became
  `gameserver.Conn` (`player.Sender` + `Player`/`SetPlayer`/`Close`), next to `Instance`
  and `HandlerParameter` in the package that already is the transport seam. Two bugs fell
  out first, in their own commit: `gameserver.Instance.Logout` took `*client.Client` --
  pointer to an interface -- so `*server.GameServer` had never satisfied `Instance`; and
  nothing could attach a player to a connection, so `handleLogin`'s `msg.Player = p` was
  discarded and every command after login would have seen a nil player.
- **`player.Sender.Send` lost its `error` return.** 87 call sites, 5 used the error, and
  the only error any implementation could return meant "this connection is already dead
  and I already tore it down." Deleting it also deleted three `// TODO error handling`
  comments in `spaces/room.go` that were never going to be resolved.
- **The command table was never written.** `message.TranslateLineToMessage` already
  existed in `watchmud-message` (written for the dead console client) and covers 20 of the
  21 handlers, aliases included. `telnet` owns only `quit` -- which has no request type of
  its own and becomes a `LogoutRequest` -- plus a bounds guard for `drop`, because
  `translator.go` indexes `tokens[1]` unguarded and a bare `drop` panics the process.
- **The login state machine is straight-line code, not a state field.** A login is a
  sequential conversation, so it wants a stack, not a state machine: `login()` runs on the
  read goroutine and blocks on an `authResult` channel that `Send` signals when a
  `LoginResponse` or `CreatePlayerResponse` goes past. The connection mutex ended up
  guarding only `player`.
- **The renderer is two files, not one.** `render.go` is keyed on protobuf types and is
  what Phase 5 retargets; `resultcode.go` maps the `ResultCode` strings `world/` emits to
  player-facing text and survives Phase 5 untouched. `TARGET_NOT_FOUND` is used by six
  handlers with two different meanings ("not here" vs "not carrying it"), so the table
  takes a verb as well as a code.
- **"Don't parse targets in the transport" was half right.** True for `drop` and `equip`,
  which take a raw string and call `world.parseTarget`. False for `get`, which reads
  `FindMode`/`Index`/`Target` and expects the *client* to have parsed them. Two grammars
  live in the tree; the translator happens to satisfy both. Phase 5 resolves it.
- **`kill` was blocked by something unrelated.** `h_kill.go`'s `fightLedger.Fight` call was
  commented out because it did not compile: after Phase 1 collapsed `Player` to a struct,
  `*player.Player` implemented exactly *one* of `combat.Combatant`'s ten methods, and
  nothing caught it because nothing ever assigned one to the other. Fixed in this phase --
  see the combat notes below.
- **Two latent world bugs surfaced during manual testing.** `World.RemovePlayer` never
  removed the player from the `Room`'s own list, so quitting left a ghost; and
  `RoomInventory.GetAll` iterated a map, so a room's contents shuffled on every `look`.

Combat was repaired here rather than deferred, since `kill` is in the exit criteria:
`Combatant` split into `Attacker` and `Defender` (the roles in a single swing, which swap
every round) plus `Combatant` (the entity that persists across them, holding `Id`,
`Dead` and `TakeMeleeDamage`). `Type() CombatantType` and its enum are gone -- it existed
so callers could un-abstract, and `corpse.go` now does an honest type switch. `FightLedger`
is keyed on `uuid.UUID` instead of on the interface itself.

That made combat *compile and cohere*, not work. The violence pulse runs and the melee math
is correct, but nothing drives a fight -- no aggression, nothing that sustains or resolves
one in play. Treat combat as structurally sound and behaviourally absent.

---

## Phase 5 — Decide protobuf's fate

**Decided: it goes. DONE September 12 2026.** Telnet was the only consumer, so the
`GameMessage` oneof wrapper and `DecodeTypeName` string-keyed dispatch were pure overhead,
and every `copylocks` finding from Phase 0 was a symptom of value-copying types that were
never meant to be copied.

`command/` (in) and `event/` (out) are plain Go structs. `make check` — fmt-check, vet and
test — is green for the first time in the project's history, which is the signal the
value-copying is actually gone rather than hidden.

Where it differed from the sketch above:

- **Three shape decisions mattered more than the mechanical conversion**, and were taken
  before the structs were written rather than discovered afterwards:
  - **`Success bool` + `ResultCode string` did not survive.** A success event carries its
    payload and nothing else; failure is a single `event.Failed{Verb, Code}` with a typed
    `event.ResultCode`. `render.go` had ~25 cases each opening with
    `if !m.Success { return failureText("<verb>", m.ResultCode) }` where the verb was a
    per-case constant — a table pretending to be a switch. It is now one case. The verb
    comes from `command.Command.Verb()` via `HandlerParameter.Fail`, so a handler never
    repeats its own name, and the `"PARSE_ERROR_" + err.Error()` string-splicing (which
    could put a raw `strconv` message in front of a player) is gone.
  - **Response and Notification merged** wherever the difference was only audience.
    `renderViolence` already took `self` and did exactly this. Drop, get, say, tell and
    shout became one event each with an `Actor`, and look/move/recall — which had always
    rendered identically — became one `event.RoomDescription`. Roughly half as many types,
    and response/notification can no longer drift apart.
  - **Naming is imperative in, past tense out**: `command.Drop` / `event.Dropped`, not
    `DropRequest` / `DropResponse`, which stuttered inside packages already named for the
    direction of travel.
- **`world.handlerMap` became a type switch, not a typed key.** Handlers take their command
  as a second parameter, so `msg.Message.GetDropRequest()` is gone from every one of them.
  The cost is that a type switch has no exhaustiveness check — a missing case compiles fine
  — so `world_unknownMessage_test.go` now covers the default arm.
- **`direction/` and `slot/` moved into this repo** as top-level packages. They were always
  domain types rather than wire types: `object`, `player` and `loader` imported them
  directly, with no `message` involved. `message.FindMode` was deleted outright rather than
  moved — `get` now passes a raw target string like every other command, which closes the
  "two target grammars" problem the Phase 4 notes left open.
- **The parser came home too.** `message.TranslateLineToMessage` became `telnet/parse.go`,
  returning `(command.Command, error)`. `quit` folded in as `command.Logout`, and the
  bare-`drop` bounds guard disappeared entirely: an empty target is a normal `NO_TARGET`
  failure, so the input that used to panic the process is now just a command that fails.
- **Three bugs surfaced, two of them pre-existing.** `server.dispatch` dereferenced
  `msg.Message` unconditionally and panicked the process on the first converted command
  (found by a live telnet session, not by any test — nothing in `world/`'s unit tests
  reaches `server.dispatch`). Recall rendered as "alice leaves none!." because
  `movePlayerMagically` moves with `direction.None`; that had always been true and is now
  guarded in the renderer. And `rules/species.go`'s malformed `json:"id""` tag — the one
  real vet finding hiding in the copylocks noise — is fixed.
- **`go mod tidy` dropped testify to v1.2.2**, which has no `require.Greater`. The message
  module had been raising it through MVS all along. Bumped to v1.11.1.
- **Deleted along the way**: `handleDataRequest` (nothing had emitted a `DataRequest` since
  the console client died), `combat.CombatantType` (a type with no values), and from
  `go.mod`: `watchmud-message`, grpc, protobuf, genproto, and `gorilla/mux` — the last two
  already had zero Go references, left over from the deleted `web/` and `rpc/`.

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
  desyncs the world. Worth collapsing to a single source of truth. This is not theoretical:
  `World.RemovePlayer` was missing the `Room.RemovePlayer` half, so quitting left a ghost in
  the room until it was fixed during Phase 4. Now unblocked.
- **`spaces.Room` conflates definition and instance.** One struct holds both the static
  topology loaded from `content/` (`Id`, `Name`, `Description`, `Zone`, `directions`, `flags`)
  and the live contents that change every tick (`playerList`, `Inventory`, `mobs`). Because
  the loader is necessarily two-pass — exits are cyclic and cross-zone, so rooms are all
  constructed before `Content.connectRooms` wires them — `Room.Connect` has to be exported,
  and a handler can rewrite world topology at runtime by calling it. Hiding `Connect` would
  fix little: `Name`, `Description`, `Zone` and `Inventory` are exported fields anyway. The
  real fix is the split this codebase already applies everywhere else (see "Definition vs
  Instance" in CLAUDE.md): a `RoomDefinition` owned by the `Zone`, immutable once loaded and
  holding the exits, and a live `Room` pointing at it. Same refactor as the dual-bookkeeping
  item above, seen from the other side — do them together. Now unblocked.
- **Nothing drives combat.** `DoViolence` is wired to `PulseViolence` and the melee
  calculation is done, but there is no aggression behaviour and nothing that carries a
  fight through to a conclusion in normal play. The pieces are in place; the gameplay layer
  on top of them is not written.
- **The fight ledger leaks third-party attackers.** `Fight(A, B)` writes two entries,
  `A->B` and `B->A`. When B kills A, `becomeCorpse` ends A's and `violence.go` ends B's --
  but a third combatant C who was also attacking A keeps its entry forever. Every violence
  pulse thereafter fetches it, sees `Fightee.Dead()`, and `continue`s. C is never told the
  target died and the entry never goes away. `IsBeingFought`'s linear scan is the ledger
  admitting it has no index for "who is attacking X"; an `EndAllFightsWith(id)` is the fix.
- **`Fight` snapshots `ZoneId`/`RoomId`** at the moment it starts, so a fight that somehow
  outlives its room notifies the wrong one. Same family as the location bookkeeping above.
- **`server.handleLogin`** logs the error from `player.FromRecord` and then falls through
  and uses the player anyway.
- **`world/settings.go`** is a single `VERBOSE_LOGGING` const, and logging is split between
  zerolog and stdlib `log` depending on file age. Worth one consolidating pass eventually.

---

## Verification

Per phase:
- **0–3:** `make test` green after each. Phase 3 additionally: server starts and idles
  without pegging a core.
- **4: DONE.** All of the below was exercised by hand, plus two automated suites that
  need no socket: `telnet/render_test.go` drives command strings through `NewTestWorld`
  and asserts on rendered text (parser, dispatch and renderer in one pass), and
  `telnet/protocol_test.go` covers the IAC filter including subnegotiation payloads that
  contain `0x00` and doubled `0xFF`.
- **4 (original plan):** manual `telnet localhost 4000` — create a player, `look`, move between rooms in
  `content/world/wrathrock`, `get`/`drop`/`inventory`/`wear`, `kill` a mob from
  `content/world/sample`, `who`, `quit`, reconnect. Two simultaneous connections to confirm
  `say`/`tell` notifications reach the other session. Watch that mob wandering (10s pulse)
  and zone reset (3min lifetime in `zone_manifest.json`) still fire while a client is idle.
- **5: DONE.** `make check` (fmt-check + vet + test) is green — vet for the first time ever.
  Verified live as well, with two simultaneous telnet sessions, because the merged
  audience-aware events are exactly the thing a single-connection test cannot check: the
  actor must see "Dropped." while the bystander sees "testdood drops knife." Both
  `telnet/render_test.go` cases for that passed with their expected text **unchanged**
  through the whole conversion, which is the evidence no player-visible string moved.
- **7:** connect with tintin++ rather than raw telnet; verify negotiation, wrapping, color.

Regression net worth adding early (cheap, and it makes every later phase safer): a table test
that feeds command strings through the Phase 4 parser into a `newTestWorld()` and asserts on
rendered output text. That covers parser, dispatch, and renderer in one pass without a socket,
and it survives the protobuf decision in Phase 5 unchanged.
