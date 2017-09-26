# WatchMUD 

### Really Simple Text-Based MUD engine

I started out writing something in Go and using eJabber as a backbone.
This was fun, but XMPP / Jabber has an awful lot of overhead. So this
is the same idea, but using web sockets instead and trying to keep
things simple...ish.

## Building and Running the Server

Makefile will compile and test the server:

    $ make
    
This will create a `watchmud` executable in the project directory.

### Running the Server

Settings and clever command line switches are still TODO, so for now

    $ ./watchmud
    
Ctrl-C to terminate the server. There's also no console or anything yet,
... also TODO.    

## Building and running the Client

See [the watchmud-client project](https://github.com/trasa/watchmud-client)
for more details, but the basics are the same:

    $ make
    
constructs the `watchmud-client` executable, and

    $ ./watchmud-client --player=YourNameGoesHere
    
Starts it up, with a login attempt to localhost for username "YourNameGoesHere".
(Settings, passwords, and other sorts of essentials also being on the "TODO" list.)

