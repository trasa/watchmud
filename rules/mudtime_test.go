package rules

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMudTime_parseDurations(t *testing.T) {
	var m MudTime
	require.NoError(t, json.Unmarshal([]byte(
		`{"mobile":"10s", "violence": "1s", "zone": "1m", "playerSave": "1m"}`), &m))

	assert.Equal(t, 10*time.Second, m.Mobile)
	assert.Equal(t, time.Minute, m.PlayerSave)
}

func TestMudTime_missingKeyFails(t *testing.T) {
	var m MudTime
	assert.Error(t, json.Unmarshal([]byte(`{"mobile":"10s", "violence": "1s", "zone": "1m"}`), &m))
}
