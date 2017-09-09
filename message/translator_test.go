package message

import (
	"github.com/stretchr/testify/assert"
	"log"
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

func TestTranslateToResponse_LoginResponse(t *testing.T) {
	s := []byte("{\"Response\":{\"msg_type\":\"login_response\",\"success\":true,\"result_code\":\"OK\"},\"Player\":{\"Name\":\"somedood\"}}")
	resp, err := TranslateToResponse(s)
	log.Println("resp ", resp)
	log.Println("err", err)
	if err != nil {
		t.Fatal("error", err)
	}
	lr := resp.(*LoginResponse)
	assert.True(t, lr.IsSuccessful())
	assert.Equal(t, "OK", lr.GetResultCode())
	assert.Equal(t, "somedood", lr.Player.Name)

}

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

func TestTranslateToResponse_LookResponse(t *testing.T) {
	s := []byte("{\"Response\":{\"msg_type\":\"look\",\"success\":true,\"result_code\":\"OK\"},\"RoomDescription\":{\"name\":\"Central Portal\",\"description\":\"desc\",\"exits\":\"ns\",\"players\":[\"player1\",\"player2\"]}}")
	resp, err := TranslateToResponse(s)

	assert.Nil(t, err)

	lr := resp.(*LookResponse)

	assert.True(t, lr.IsSuccessful())
	assert.Equal(t, "OK", lr.GetResultCode())
	assert.Equal(t, "ns", lr.RoomDescription.Exits)
	assert.Equal(t, "Central Portal", lr.RoomDescription.Name)
	assert.Equal(t, "desc", lr.RoomDescription.Description)
	assert.Equal(t, []string{"player1", "player2"}, lr.RoomDescription.Players)
}

func TestTranslateToResponse_ExitsResponse(t *testing.T) {
	s := []byte("{\"Response\":{\"msg_type\":\"exits\",\"success\":true,\"result_code\":\"OK\"},\"exitinfo\":{\"n\":\"North Room\"}}")
	resp, err := TranslateToResponse(s)
	assert.Nil(t, err)
	r := resp.(*ExitsResponse)

	assert.True(t, r.IsSuccessful())
	assert.Equal(t, "OK", r.GetResultCode())
	assert.Equal(t, "North Room", r.ExitInfo["n"])
}

func TestDecodeResponse(t *testing.T) {
	rawMap := make(map[string]interface{})
	rawMap["Value"] = "hello"
	orig := &SayResponse{}
	response := decodeResponse("say", orig, rawMap)

	assert.Equal(t, "hello", orig.Value)
	// response.ResponseBase stuff won't be set yet - only the
	// values within the top level SayResponse

	//assert.Equal(t, "OK", orig.GetResultCode())
	//assert.Equal(t, "say", orig.GetMessageType())
	//assert.Equal(t, true, orig.IsSuccessful())

	assert.Equal(t, "hello", response.(*SayResponse).Value)

}

func TestFillResponseBase(t *testing.T) {
	responseMap := map[string]interface{}{
		"result_code": "FINE",
		"success":     true,
		"msg_type":    "say",
	}

	rawMap := make(map[string]interface{})
	rawMap["Value"] = "hello"
	orig := &SayResponse{}
	response := decodeResponse("say", orig, rawMap)

	fillResponseBase(response, responseMap)

	assert.Equal(t, "hello", orig.Value)

	assert.Equal(t, "hello", response.(*SayResponse).Value)
	assert.Equal(t, "FINE", response.GetResultCode())
	assert.Equal(t, "say", response.GetMessageType())
	assert.Equal(t, true, response.IsSuccessful())
}
