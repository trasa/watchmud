What is this?

A long time ago, I played and wrote code for "Multi User Dungeons",
aka MUDs. It was great fun and hazardous to my GPA. The text-based
multiplayer game is still one of my favorite environments when I
want to learn a new language or technology stack. There's plenty
of things going on, it's much more complex than "Hello World",
and it's not too boring.

I wanted to learn Go and GRPC - I was considering using them for
a project at work, and wanted something beyond the small
examples that I had seen. I wanted to see what would happen if 
you tried to build something "big" in Go. Go's structure for
organizing code is intentionally light, and it also enforces
at compile-time strong opinions about how code should be
organized. 

What I found is that I was bringing habits from C# and Java
to Go with me. Some of those habits caused more problems than
others. Some problems have been refactored out in this codebase,
but others haven't been.

This is a work-in-progress and probably always will be. Tech gets
swapped in and out, better ways of doing things become possible.
(Go Generics anybody?) Some of this code hasn't been touched for
years.

Some places to start looking at code:

[main.go](main.go)
The ever popular `func main`. Sets things up and starts them off.

[server/gameserver.go](server/gameserver.go)
This code runs the main game server loop which keeps track of time,
triggers events, and receives messages from clients.

[world/handlers.go](world/handlers.go) and [world/world.go](world/world.go)
Handlers for incoming GRPC messages and how to find them.
The "state of the world" and where stuff is inside the world.


