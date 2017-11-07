package message

import (
	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/direction"
	"log"
	"testing"
)

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

func TestSerialize_ExitsResponse(t *testing.T) {
	exitresponse := ExitsResponse{
		Response: NewSuccessfulResponse("exits"),
		ExitInfo: []DirectionToRoomName{
			{Direction: direction.NORTH, RoomName: "North Room"},
		},
	}
	s, err := TranslateToJson(exitresponse)
	assert.NoError(t, err)
	log.Printf(s)
}

func TestTranslateToResponse_ExitsResponse(t *testing.T) {
	s := []byte("{\"Response\":{\"msg_type\":\"exits\",\"success\":true,\"result_code\":\"OK\"},\"ExitInfo\":[{\"Direction\":1,\"RoomName\":\"North Room\"}]}")
	resp, err := TranslateToResponse(s)
	assert.Nil(t, err)
	r := resp.(*ExitsResponse)

	assert.True(t, r.IsSuccessful())
	assert.Equal(t, "OK", r.GetResultCode())
	assert.Equal(t, "North Room", r.ExitInfo[0].RoomName)
}

func TestTranslateToResponse_SayResponse(t *testing.T) {
	s := []byte("{\"Response\":{\"msg_type\":\"say\",\"success\":true,\"result_code\":\"WHATEVER\"},\"value\":\"hello\"}")
	resp, err := TranslateToResponse(s)
	assert.Nil(t, err)
	r := resp.(*SayResponse)
	assert.True(t, r.IsSuccessful())
	assert.Equal(t, "WHATEVER", r.GetResultCode())
	assert.Equal(t, "hello", r.Value)
	log.Println(r)
}

func TestTranslateToResponse_NoOp(t *testing.T) {
	s := []byte("{\"Response\":{\"msg_type\":\"no_op\",\"success\":false,\"result_code\":\"UNKNOWN_MESSAGE_TYPE\"}}")
	resp, err := TranslateToResponse(s)
	r := resp.(*NoOpResponse)
	assert.Nil(t, err)
	assert.False(t, r.IsSuccessful())
	assert.Equal(t, "UNKNOWN_MESSAGE_TYPE", r.GetResultCode())
}
