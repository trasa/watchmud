package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeTypeName(t *testing.T) {
	payload := DropResponse{}
	ptr := &payload

	assert.Equal(t, "DropResponse", DecodeTypeName(payload))
	assert.Equal(t, "DropResponse", DecodeTypeName(ptr))
}
