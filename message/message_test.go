package message

import (
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestSerDe(t *testing.T) {
	// create a drop request
	dropRequest := &GameMessage_DropRequest{
		DropRequest: &DropRequest{Target: "dropTarget"},
	}
	msg := &GameMessage{}
	msg.Inner = dropRequest

	// marshal the message containing the drop request
	out, err := proto.Marshal(msg)
	log.Printf("marshal to bytes: error: %s bytes: '%b'", err, out)
	assert.NoError(t, err)
	assert.NotNil(t, out)

	// extract the generic request from the bytes
	resultMessage := GameMessage{}
	err = proto.Unmarshal(out, &resultMessage)
	log.Printf("unmarshal to generic request: error: %s req: '%s'", err, resultMessage.Inner)
	assert.NoError(t, err)

	switch r := resultMessage.Inner.(type) {
	case *GameMessage_GetRequest:
		log.Printf("get %s", r.GetRequest)
		t.Fatal("shouldn't end up here")
	case *GameMessage_DropRequest:
		dr := resultMessage.GetDropRequest()
		log.Printf("drop this: '%s', '%s', '%s', ", r, r.DropRequest.Target, dr.Target)
		assert.Equal(t, "dropTarget", dr.Target)
	}
}

func TestBuildMessage(t *testing.T) {
	// create a drop request
	dropRequest := &DropRequest{Target: "dropTarget"}
	innerMessage := &GameMessage_DropRequest{
		DropRequest: dropRequest,
	}
	msg := &GameMessage{}
	msg.Inner = innerMessage

}
