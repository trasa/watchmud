package logging

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// No file is stdout only, which is what a container wants: docker keeps and
// rotates it.
func TestInitialize_noFile(t *testing.T) {
	closeLog, err := Initialize("", "info")
	require.NoError(t, err)
	closeLog()
}
