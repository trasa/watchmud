package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeResponse(t *testing.T) {
	rawMap := make(map[string]interface{})
	rawMap["Value"] = "hello"
	orig := &SayResponse{}
	response := decodeResponse(orig, rawMap)

	assert.Equal(t, "hello", orig.Value)
	// response.ResponseBase stuff won't be set yet - only the
	// values within the top level SayResponse
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
	response := decodeResponse(orig, rawMap)

	fillResponseBase(response, responseMap)

	assert.Equal(t, "hello", orig.Value)
	assert.Equal(t, "hello", response.(*SayResponse).Value)

	assert.Equal(t, "FINE", response.GetResultCode())
	assert.Equal(t, "say", response.GetMessageType())
	assert.Equal(t, true, response.IsSuccessful())
}
