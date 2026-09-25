# WatchMUD

### Really Simple Text-Based MUD engine

A text MUD server written in Go, speaking telnet. Connect with any MUD client --
tintin++, mudlet, or plain `telnet` -- create a character, and wander around.

    $ telnet watchmud.com 4000      # or 4443 for TLS

(The public server at watchmud.com isn't up yet. To run your own: `make run`, then
`telnet localhost 4000`.)

### Where this is going

[ROADMAP.md](ROADMAP.md) is the authoritative description of what's done and what isn't.
The short version:

- **Telnet is the transport**, with a second port speaking TLS so passwords don't cross
  the internet in the clear. The old gRPC listener, web page, console client and protobuf
  are all gone.
- **Characters are kept in mongo**, one document each, with bcrypt-hashed passwords.
- **There are no classes.** What a character is good at -- Tank, Healer, Striker -- is
  read off the gear they're wearing, and changes when they change it.
- **Fighting, loot, corpses and gear that wears out** all work. Next is content, and
  bots to play against.

### History

I started out writing something in Java and using eJabber for communication.
This was fun, but XMPP has an awful lot of overhead and that turned into
a lot of work. Also, Go seemed like a fun language to learn. So the Java
code was scrapped for Go, and eventually the XMPP / eJabber implementation
was scrapped for JSON over Web Sockets. The original client was a single
web page app using JQuery, with the intention of replacing JQuery with
something better...

I found that I was having to write a great amount of code translating
JSON to Go structs and back, both on the server and in the client. So
I replaced the JQuery web page with a Go Client application,
[watchmud-client](https://github.com/trasa/watchmud-client).

But there was still too much serializing-deserializing code between
client and server and websocket. So I replaced that with gRPC.

What will I rewrite next??

The transport, as it turns out. gRPC meant every player needed my custom client, which
is a strange thing to ask of a MUD -- the genre has had a perfectly good wire protocol
since 1978, and people already own clients they like. So gRPC came out and telnet went
in, the console client was retired, and Postgres went with it. Writing a telnet server
also turned out to be the cheapest way to find out which of the remaining abstractions
were real: several of them, it turned out, had no working implementations at all.

This has more [history](codereview.md) about this project.

## Building

You'll need Go 1.27 or later. That's it -- `stringer` is declared as a tool
dependency in `go.mod`, so there's nothing to install separately.

    $ make            # build -> bin/watchmud
    $ make test
    $ make help       # list every target

`make check` -- gofmt, vet and the tests -- is the gate, and it's green.

## Running the Server

Configuration lives in [app.local.yaml](app.local.yaml) -- ports, content path, and log
destination. Running it for real players is [deploy/README.md](deploy/README.md): docker
compose, mongo with auth, backups, and TLS.

    $ make run                                   # uses ./app.local.yaml
    $ ./bin/watchmud -config /path/to/other.yaml
    $ ./bin/watchmud -content /path/to/content   # override just the content path

Ctrl-C to terminate the server.

### Persistence

Characters are stored in mongo, one document per character, in the `players` collection.
[docker-compose.yml](docker-compose.yml) has one for local development:

    $ make db-up      # start it (host port 27018)
    $ make db-shell   # poke at it: db.players.find()
    $ make db-down    # stop it, keeping the data
    $ make db-reset   # stop it and throw the data away

It listens on **27018** rather than the usual 27017, so it can't be confused with a mongo
you already have installed -- `mongo.uri` in app.local.yaml points at it.

Empty out `mongo.uri` and the server runs on the in-memory store instead, which is fine
for a throwaway session and loses everything on exit. A uri that is set and unreachable
fails startup on purpose: a server that comes up anyway looks healthy right until it has
silently discarded an evening of play.

The tests that need a real mongo skip themselves unless you point them at one:

    $ make db-up && make test-db

The world is loaded from `content/`: `content/rules/` holds lineages, roles and the rules tables,
`content/world/` holds the zones (`wrathrock`, `sample`, `void`) plus the settings and
zone manifest that say which of them load.

## Playing

    $ telnet localhost 4000

You'll be asked for a name, then its password; if there's no character by that name,
you'll be offered the chance to create one. Once you're in, `help` lists the commands.
Characters are saved, and come back where you left them.
