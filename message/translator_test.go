package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTranslate_NoOp(t *testing.T) {
	req, err := TranslateLineToRequest("")

	assert.Nil(t, err, "should not error")
	assert.Equal(t, "no_op", req.GetMessageType(), "expected message_type no_op: %s", req)
}

func TestTranslate_Tell_shortest_message(t *testing.T) {
	req, err := TranslateLineToRequest("tell bob hello")

	assert.Nil(t, err, "no error")
	assert.Equal(t, "tell", req.GetMessageType(), "message_type tell: %s", req)

	tellReq := req.(TellRequest)
	assert.Equal(t, "bob", tellReq.ReceiverPlayerName, "wrong rec name: %s", req)
	assert.Equal(t, "hello", tellReq.Value, "wrong value: %s", req)
}

func TestTranslate_Tell(t *testing.T) {
	req, err := TranslateLineToRequest("tell bob hello there")

	assert.Nil(t, err, "no error")
	assert.Equal(t, "tell", req.GetMessageType(), "message_type tell: %s", req)

	tellReq := req.(TellRequest)
	assert.Equal(t, "bob", tellReq.ReceiverPlayerName, "wrong rec name: %s", req)
	assert.Equal(t, "hello there", tellReq.Value, "wrong value: %s", req)
}

func TestTranslate_Tell_shortcut(t *testing.T) {
	req, err := TranslateLineToRequest("t bob hello there")

	assert.Nil(t, err, "no error")
	assert.Equal(t, "tell", req.GetMessageType(), "message_type tell: %s", req)

	tellReq := req.(TellRequest)
	assert.Equal(t, "bob", tellReq.ReceiverPlayerName, "wrong rec name: %s", req)
	assert.Equal(t, "hello there", tellReq.Value, "wrong value: %s", req)
}

func TestTranslate_Tell_BadRequest(t *testing.T) {
	req, err := TranslateLineToRequest("tell bob")

	assert.NotNil(t, err, "expected error")
	assert.Nil(t, req, "req should be nil")
}

func TestTranslate_UnknownCommand(t *testing.T) {
	cmd := "asdjhfaksjdfhjk"
	req, err := TranslateLineToRequest(cmd)
	assert.Nil(t, req, "req should be nil")
	assert.NotNil(t, err, "should have error")
	assert.Equal(t, "Unknown request: "+cmd, err.Error(), "err text")
}

func TestTranslate_TellAll(t *testing.T) {
	req, err := TranslateLineToRequest("tellall hello")
	assert.Nil(t, err, "no error expected")
	assert.Equal(t, "tell_all", req.GetMessageType(), "expected tell_all")

	tellAllReq := req.(TellAllRequest)
	assert.Equal(t, "hello", tellAllReq.Value)
}

func TestTranslate_TellAll_shortcut(t *testing.T) {
	req, err := TranslateLineToRequest("ta hello")
	assert.Nil(t, err, "no error expected")
	assert.Equal(t, "tell_all", req.GetMessageType(), "expected tell_all")

	tellAllReq := req.(TellAllRequest)
	assert.Equal(t, "hello", tellAllReq.Value)
}

func TestTranslate_TellAll_LongMessage(t *testing.T) {
	req, err := TranslateLineToRequest("ta a b c d e f g h i j k l m n o p  q r s t    u v")
	assert.Nil(t, err, "no error")
	tellAllReq := req.(TellAllRequest)
	assert.Equal(t, "a b c d e f g h i j k l m n o p q r s t u v", tellAllReq.Value)
}
