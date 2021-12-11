package msg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgResponse(t *testing.T) {
	assert := require.New(t)

	t.Run("test_constructor_msg_response", func(t *testing.T) {
		action := STUN_ACTION_GET
		hasError := false
		peername := "Dog"
		message := "Guau"

		response := NewMsgResponse(action, hasError, peername, message)

		assert.Equal(action, response.Action)
		assert.Equal(hasError, response.HasError)
		assert.Equal(peername, response.Peername)
		assert.Equal(message, response.Message)
	})
}
