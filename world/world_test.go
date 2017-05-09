package world

import (
	"log"
	"testing"
)

func TestWorldHasStartRoom(t *testing.T) {
	log.Printf("Start Room: %v", WorldInstance.StartRoom)
	if WorldInstance.StartRoom == nil {
		t.Error("StartRoom should not be nil")
	}
}
