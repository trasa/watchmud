package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trasa/watchmud/command"
)

func TestSecretsRedacted(t *testing.T) {
	s := command.Secret("foo").String()
	assert.Equal(t, "[redacted]", s)
}

func TestGoSecretRedacted(t *testing.T) {
	s := command.Secret("foo").GoString()
	assert.Equal(t, "[redacted]", s)
}
