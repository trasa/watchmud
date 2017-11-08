package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTranslate_NoOp(t *testing.T) {
	msg, err := TranslateLineToMessage("")

	assert.NoError(t, err)
	assert.IsType(t, &GameMessage_Ping{}, msg.GetInner())
	assert.NotNil(t, msg.GetPing())
}

func TestTranslate_Tell_shortest_message(t *testing.T) {
	msg, err := TranslateLineToMessage("tell bob hello")

	assert.NoError(t, err)

	tellReq := msg.GetTellRequest()
	assert.Equal(t, "bob", tellReq.ReceiverPlayerName, "wrong rec name: %s", msg)
	assert.Equal(t, "hello", tellReq.Value, "wrong value: %s", msg)
}

func TestTranslate_Tell(t *testing.T) {
	msg, err := TranslateLineToMessage("tell bob hello there")

	assert.NoError(t, err)

	tellReq := msg.GetTellRequest()
	assert.Equal(t, "bob", tellReq.ReceiverPlayerName, "wrong rec name: %s", msg)
	assert.Equal(t, "hello there", tellReq.Value, "wrong value: %s", msg)
}

func TestTranslate_Tell_shortcut(t *testing.T) {
	msg, err := TranslateLineToMessage("t bob hello there")

	assert.NoError(t, err)

	tellReq := msg.GetTellRequest()
	assert.Equal(t, "bob", tellReq.ReceiverPlayerName, "wrong rec name: %s", msg)
	assert.Equal(t, "hello there", tellReq.Value, "wrong value: %s", msg)
}

func TestTranslate_Tell_BadRequest(t *testing.T) {
	msg, err := TranslateLineToMessage("tell bob")

	assert.Error(t, err)
	assert.Nil(t, msg, "should be nil")
}

func TestTranslate_UnknownCommand(t *testing.T) {
	cmd := "asdjhfaksjdfhjk"
	msg, err := TranslateLineToMessage(cmd)
	assert.Nil(t, msg, "req should be nil")
	assert.Error(t, err)
	assert.Equal(t, "Unknown request: "+cmd, err.Error(), "err text")
}

func TestTranslate_TellAll(t *testing.T) {
	msg, err := TranslateLineToMessage("tellall hello")

	assert.NoError(t, err)
	tellAllReq := msg.GetTellAllRequest()
	assert.Equal(t, "hello", tellAllReq.Value)
}

func TestTranslate_TellAll_shortcut(t *testing.T) {
	msg, err := TranslateLineToMessage("ta hello")

	assert.NoError(t, err)
	tellAllReq := msg.GetTellAllRequest()
	assert.Equal(t, "hello", tellAllReq.Value)
}

func TestTranslate_TellAll_LongMessage(t *testing.T) {
	msg, err := TranslateLineToMessage("ta a b c d e f g h i j k l m n o p  q r s t    u v")

	assert.NoError(t, err)
	tellAllReq := msg.GetTellAllRequest()
	assert.Equal(t, "a b c d e f g h i j k l m n o p q r s t u v", tellAllReq.Value)
}

func TestTranslate_Get_NoTarget(t *testing.T) {
	msg, err := TranslateLineToMessage("get")

	assert.NoError(t, err)
	getreq := msg.GetGetRequest()
	assert.Equal(t, 0, len(getreq.Targets))
}

func TestTranslate_Get_OneTarget(t *testing.T) {
	msg, err := TranslateLineToMessage("get foo")

	assert.NoError(t, err)
	getreq := msg.GetGetRequest()
	assert.Equal(t, 1, len(getreq.Targets))
	assert.Equal(t, "foo", getreq.Targets[0])
}

func TestTranslate_Get_TwoTargets(t *testing.T) {
	msg, err := TranslateLineToMessage("get foo bar")

	assert.NoError(t, err)
	getreq := msg.GetGetRequest()
	assert.Equal(t, 2, len(getreq.Targets))
	assert.Equal(t, "foo", getreq.Targets[0])
	assert.Equal(t, "bar", getreq.Targets[1])
}
