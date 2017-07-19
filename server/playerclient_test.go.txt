package server

import (
	"testing"
	"log"
)

type TestPlayer struct {
	Name       string
	TestClient TestClient
}


func (p *TestPlayer) Send(msg interface{}) {
	p.TestClient.Send(msg)
}

type TestClient interface {
	Send(msg interface{})
}


type WebTestClient struct {
	conn *FancyWebConnection
}

func (c *WebTestClient) Send(msg interface{}) {
	c.conn.Send(msg)
}

type FancyWebConnection struct {
}

func (f *FancyWebConnection) Send(msg interface{}) {
	log.Println("send msg through FancyWebConnection")
}


func NewWebTestClient(c *FancyWebConnection) *WebTestClient {
	return &WebTestClient{
		conn: c,
	}
}

type NoOpTestClient struct {
	tosend []interface{}
}

func newNoOpTestClient() *NoOpTestClient {
	return &NoOpTestClient{
		//tosend: []interface{},
	}
}

func (c *NoOpTestClient) Send(msg interface{}) {
	log.Println("NoOp is sending")
	c.tosend = append(c.tosend, msg)
	log.Printf("len is %d", len(c.tosend))
}

func TestNoOp(t *testing.T) {
	c := &NoOpTestClient{}
	p := &TestPlayer{
		Name: "Bob",
		TestClient: c,
	}
	p.Send("hi")
	log.Printf("here len is %d", len(c.tosend))
	if len(c.tosend) != 1 {
		t.Error("expected len 1")
	}

}
