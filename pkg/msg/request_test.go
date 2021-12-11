package msg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgRequest(t *testing.T) {
	assert := require.New(t)

	t.Run("test_constructor_msg_request", func(t *testing.T) {
		action := STUN_ACTION_GET
		peername := "Dog"
		message := "Guau"

		request := NewMsgRequest(action, peername, message)

		assert.Equal(action, request.Action)
		assert.Equal(peername, request.Peername)
		assert.Equal(message, request.Message)
	})
}
