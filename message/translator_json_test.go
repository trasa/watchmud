package message

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestTranslateToJson(t *testing.T) {

	lr := LookResponse{
		Response: NewSuccessfulResponse("look"),
		RoomDescription: RoomDescription{
			Name:        "foo",
			Description: "desc",
			Players:     []string{"player1", "player2"},
			Exits:       "ns",
		},
	}
	str, err := TranslateToJson(lr)
	assert.Nil(t, err)
	log.Println(str)
}
