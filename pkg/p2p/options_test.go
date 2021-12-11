package p2p

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPeerOptions(t *testing.T) {
	assert := require.New(t)

	t.Run("test_new_peer_options", func(t *testing.T) {
		maxMsgInQueue := 5
		timeout := 5

		options := NewPeerOptions(maxMsgInQueue, timeout)

		assert.Equal(maxMsgInQueue, options.maxMsgInQueue)
		assert.Equal(timeout, options.timeout)
	})

	t.Run("test_new_peer_default_options", func(t *testing.T) {
		options := DefaultPeerOptions()

		assert.Equal(DEFAULT_MAX_MSG_IN_QUEUE, options.maxMsgInQueue)
		assert.Equal(DEFAULT_SECONDS_TIMEOUT, options.timeout)
	})
}
