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


## Building the Server

To compile, test the server::

    $ make
    
This will create a `watchmud` executable in the project directory.

### Running the Server

Settings and clever command line switches are still TODO, so for now

    $ ./watchmud
    
Ctrl-C to terminate the server. 

## Building and running the Client

See [the watchmud-client project](https://github.com/trasa/watchmud-client)
for more details, but the basics are the same:

    $ make
    
constructs the `watchmud-client` executable, and

    $ ./watchmud-client --player=YourNameGoesHere
    
Starts it up, with a login attempt to localhost for username "YourNameGoesHere".
(Settings, passwords, and other sorts of essentials also being on the "TODO" list.)

