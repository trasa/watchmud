# WatchMUD 

### Really Simple Text-Based MUD engine

This is a straightforward text MUD, written in Go and using gRPC instead
of telnet. 

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

This has more [history](codereview.md) about this project.

## Building the Server

Install Dependencies: you'll need Go 1.14+ to build this.
You'll need to install go stringer:

    $ go get golang.org/x/tools/cmd/stringer
    
and then make sure that the directory this is installed to (probably something
like ~/go/bin) is part of your path.
    
To compile and test the server:

    $ make
    
This will create a `watchmud` executable in the project/bin directory.

### Building the Database

WatchMud uses Postgres to hold some information about the users, game state
and other interesting stuff like that. You'll need to have an instance of
Postgres running for watchmud.

Installing postgres:

    $ brew install postgresql

Creating the user 'watchmud' and the schema:

See watchmud/db/sql/ddl.sql for table definitions and static data creation.

    
### Creating the Log Directory

By default, watchmud won't have permissions to create the directory where
it wants to write its log files. Unless you're overriding this directory with
the --logFile argument, you'll need to create the directory and give the 
watchmud user permission to write there:

    $ sudo mkdir -p /var/log/watchmud
    $ sudo chown myusername /var/log/watchmud
    $ sudo chmod 775 /var/log/watchmud

(Or, something similar to that.)

### Running the Server

Example Settings are shown in [worldfiles/server.yaml](worldfiles/server.yaml)

    $ ./watchmud -serverconfig ./worldfiles/myserverconfig.yaml
    
Ctrl-C to terminate the server. 

## Building and running the Client

See [the watchmud-client project](https://github.com/trasa/watchmud-client)
for more details, but the basics are the same:

    $ make
    
constructs the `watchmud-client` executable, and

    $ ./watchmud-client --player=YourNameGoesHere
    
Starts it up, with a login attempt to localhost for username "YourNameGoesHere".
(Settings, passwords, and other sorts of essentials also being on the "TODO" list.)

