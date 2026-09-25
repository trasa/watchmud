package mongostore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// The uri is logged at startup, and it carries the database password.
func TestRedactURI(t *testing.T) {
	cases := map[string]string{
		"mongodb://watchmud:hunter22@mongo:27017/watchmud?authSource=watchmud": "mongodb://watchmud:xxxxx@mongo:27017/watchmud?authSource=watchmud",
		"mongodb://localhost:27018":             "mongodb://localhost:27018",
		"mongodb://user:pw@a:1,b:2/db":          "mongodb://user:xxxxx@a:1,b:2/db",
		"mongodb+srv://user:pw@cluster.example": "mongodb+srv://user:xxxxx@cluster.example",
	}
	for uri, want := range cases {
		assert.Equal(t, want, RedactURI(uri), uri)
		assert.NotContains(t, RedactURI(uri), "hunter22")
	}
	// not a uri at all: say nothing rather than risk printing it
	assert.Equal(t, "(unparseable uri)", RedactURI("mongodb://user:p%zzw@host"))
}
