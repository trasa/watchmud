package main

import (
	"github.com/emgee/go-xmpp/src/xmpp"
	"testing"
"github.com/trasa/jabmud/commands"
)

func TestErrorPresence(t *testing.T) {
	presence := xmpp.Presence{
		To:   "jabmud.localhost/someguy",
		From: "what@somewhere",
	}
	response := newErrorPresence(&presence)
	str := commands.Serialize(response)
	expected := "<presence type=\"error\" to=\"what@somewhere\" from=\"jabmud.localhost/someguy\"><x xmlns=\"http://jabber.org/protocol/muc\"></x><error type=\"cancel\"><conflict xmlns=\"urn:ietf:params:xml:ns:xmpp-stanzas\"></conflict></error></presence>"
	if str != expected {
		t.Errorf("serialize didn't get expected string\nexp=%s\nact=%s", expected, str)
	}
}

func TestSuccessPresence(t *testing.T) {
	presence := xmpp.Presence{
		To:   "jabmud.localhost/someguy",
		From: "what@somewhere",
	}
	response := newSuccessPresence(&presence)
	str := commands.Serialize(response)
	expected := "<presence to=\"what@somewhere\" from=\"jabmud.localhost/someguy\"><x xmlns=\"http://jabber.org/protocol/muc\"><item affiliation=\"member\" role=\"participant\"></item><status code=\"110\"></status></x></presence>"
	if str != expected {
		t.Errorf("serialize didn't get expected string\nexp=%s\nact=%s", expected, str)
	}
}
