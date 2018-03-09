package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTranslate_NoOp(t *testing.T) {
	msg, err := TranslateLineToMessage(Tokenize(""))

	assert.NoError(t, err)
	assert.IsType(t, &GameMessage_Ping{}, msg.GetInner())
	assert.NotNil(t, msg.GetPing())
}

func TestTranslate_Tell_shortest_message(t *testing.T) {
	msg, err := TranslateLineToMessage(Tokenize("tell bob hello"))

	assert.NoError(t, err)

	tellReq := msg.GetTellRequest()
	assert.Equal(t, "bob", tellReq.ReceiverPlayerName, "wrong rec name: %s", msg)
	assert.Equal(t, "hello", tellReq.Value, "wrong value: %s", msg)
}

func TestTranslate_Tell(t *testing.T) {
	msg, err := TranslateLineToMessage(Tokenize("tell bob hello there"))

	assert.NoError(t, err)

	tellReq := msg.GetTellRequest()
	assert.Equal(t, "bob", tellReq.ReceiverPlayerName, "wrong rec name: %s", msg)
	assert.Equal(t, "hello there", tellReq.Value, "wrong value: %s", msg)
}

func TestTranslate_Tell_shortcut(t *testing.T) {
	msg, err := TranslateLineToMessage(Tokenize("t bob hello there"))

	assert.NoError(t, err)

	tellReq := msg.GetTellRequest()
	assert.Equal(t, "bob", tellReq.ReceiverPlayerName, "wrong rec name: %s", msg)
	assert.Equal(t, "hello there", tellReq.Value, "wrong value: %s", msg)
}

func TestTranslate_Tell_BadRequest(t *testing.T) {
	msg, err := TranslateLineToMessage(Tokenize("tell bob"))

	assert.Error(t, err)
	assert.Nil(t, msg, "should be nil")
}

func TestTranslate_UnknownCommand(t *testing.T) {
	cmd := "asdjhfaksjdfhjk"
	msg, err := TranslateLineToMessage(Tokenize(cmd))
	assert.Nil(t, msg, "req should be nil")
	assert.Error(t, err)
	assert.Equal(t, "Unknown request: "+cmd, err.Error(), "err text")
}

func TestTranslate_TellAll(t *testing.T) {
	msg, err := TranslateLineToMessage(Tokenize("tellall hello"))

	assert.NoError(t, err)
	tellAllReq := msg.GetTellAllRequest()
	assert.Equal(t, "hello", tellAllReq.Value)
}

func TestTranslate_TellAll_shortcut(t *testing.T) {
	msg, err := TranslateLineToMessage(Tokenize("ta hello"))

	assert.NoError(t, err)
	tellAllReq := msg.GetTellAllRequest()
	assert.Equal(t, "hello", tellAllReq.Value)
}

func TestTranslate_TellAll_LongMessage(t *testing.T) {
	msg, err := TranslateLineToMessage(Tokenize("ta a b c d e f g h i j k l m n o p  q r s t    u v"))

	assert.NoError(t, err)
	tellAllReq := msg.GetTellAllRequest()
	assert.Equal(t, "a b c d e f g h i j k l m n o p q r s t u v", tellAllReq.Value)
}

func TestTranslate_Get_NoTarget(t *testing.T) {
	_, err := TranslateLineToMessage(Tokenize("get"))

	assert.Error(t, err)
}

func TestTranslate_Get_OneTarget(t *testing.T) {
	msg, err := TranslateLineToMessage(Tokenize("get foo"))

	assert.NoError(t, err)
	getreq := msg.GetGetRequest()
	assert.Equal(t, "foo", getreq.Target)
	assert.Equal(t, int32(FindIndividual), getreq.FindMode)
	assert.Equal(t, "", getreq.Index)
}

func TestTranslate_Get_Dots(t *testing.T) {
	msg, err := TranslateLineToMessage(Tokenize("get 2.foo"))

	assert.NoError(t, err)
	getreq := msg.GetGetRequest()
	assert.Equal(t, "foo", getreq.Target)
	assert.Equal(t, int32(FindIndividual), getreq.FindMode)
	assert.Equal(t, "2", getreq.Index)
}

func TestTranslate_Get_All(t *testing.T) {
	msg, err := TranslateLineToMessage(Tokenize("get all"))

	assert.NoError(t, err)
	getreq := msg.GetGetRequest()
	assert.Equal(t, "all", getreq.Target)
	assert.Equal(t, int32(FindAll), getreq.FindMode)
	assert.Equal(t, "", getreq.Index)
}

func TestTranslate_Get_All_Dot(t *testing.T) {
	msg, err := TranslateLineToMessage(Tokenize("get all.foo"))

	assert.NoError(t, err)
	getreq := msg.GetGetRequest()
	assert.Equal(t, "foo", getreq.Target)
	assert.Equal(t, int32(FindAllDot), getreq.FindMode)
	assert.Equal(t, "all", getreq.Index)
}
