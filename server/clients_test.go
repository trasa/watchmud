package server

import "testing"

func TestClients_Add(t *testing.T) {
	clients := newClients()
	client := &TestClient{}

	clients.add(client)
	if _, ok := clients.clients[client]; !ok {
		t.Error("Client added to Clients but not found by key")
	}
}

func TestClients_Remove(t *testing.T) {
	clients := newClients()
	client := &TestClient{}

	clients.add(client)
	clients.remove(client)

	if _, ok := clients.clients[client]; ok {
		t.Error("Client removed from Clients but it is still there")
	}
}

func TestClients_Iter(t *testing.T) {
	clients := newClients()
	client := &TestClient{}

	clients.add(client)
	count := 0
	clients.iter(func(c Client) {
		count++
	})
	if count != 1 {
		t.Errorf("found %d clients expected 1", count)
	}

}
