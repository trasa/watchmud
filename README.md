# WatchMUD

### Really Simple Text-Based MUD engine

A text MUD server written in Go, speaking telnet. Connect with any MUD client --
tintin++, mudlet, or plain `telnet` -- create a character, and wander around.

    $ make run
    $ telnet localhost 4000

### Where this is going

The project is mid-migration, and [ROADMAP.md](ROADMAP.md) is the authoritative
description of what's done and what isn't. The short version:

- **Telnet is the transport.** The old gRPC listener, the static web page, and the
  separate console client are deleted.
- **Persistence is an interface with an in-memory implementation.** Characters do not
  survive a restart yet. Postgres is gone; its replacement gets chosen later, against a
  server that actually runs.
- **Protobuf is still the internal vocabulary**, but it is on notice. Nothing serializes
  these messages any more -- the telnet layer reads their fields directly -- so the next
  decision is whether the `GameMessage` wrapper still earns its place.
- **Combat exists but nothing drives it.** The interfaces and the melee math are there and
  the violence pulse runs; what's missing is the gameplay layer that starts and sustains a
  fight. It's on the list.

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

`make vet` fails, and has for a long time. Nearly every finding is `copylocks` from
passing protobuf structs by value, which is the house style throughout; those go away on
their own if protobuf does. `make test` is the gate that means something.

## Running the Server

Configuration lives in [app.local.yaml](app.local.yaml) -- ports, content path, and log
destination.

    $ make run                                   # uses ./app.local.yaml
    $ ./bin/watchmud -config /path/to/other.yaml
    $ ./bin/watchmud -content /path/to/content   # override just the content path

Ctrl-C to terminate the server.

The world is loaded from `content/`: `content/rules/` holds species and classes,
`content/world/` holds the zones (`wrathrock`, `sample`, `void`) plus the settings and
zone manifest that say which of them load.

## Playing

    $ telnet localhost 4000

You'll be asked for a name; if there's no character by that name, you'll be offered the
chance to create one. From there the usual verbs work -- `look`, `north` (or just `n`),
`get`, `drop`, `inventory`, `wear`, `wield`, `say`, `tell`, `who`, `stat`, `exits`, and
`quit`. Most have the abbreviations you'd expect.

Characters live in memory only, so they vanish when the server stops.
