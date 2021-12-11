package p2p

import (
	"encoding/json"
	"fmt"
	"net"
	"testing"

	"github.com/alvarogf97/fox/pkg/msg"
	"github.com/stretchr/testify/require"
)

func TestP2PWriter(t *testing.T) {
	assert := require.New(t)

	t.Run("test_new_p2p_writer", func(t *testing.T) {
		name := "fakeP2PWriter"
		port := "50000"
		addr, _ := net.ResolveUDPAddr("udp4", port)
		p2pConn := P2PConnMock{P2PWriteToUDPMock{}}

		writer := NewP2PWriter(name, &p2pConn, addr)

		assert.Equal(name, writer.name)
	})

	t.Run("test_write_success", func(t *testing.T) {
		action := "FakeAction"
		message := "FakeMessage"
		name := "fakeP2PWriter"
		port := "50000"
		addr, _ := net.ResolveUDPAddr("udp4", port)

		returnSend := 5
		expectedBytes, _ := json.Marshal(msg.NewMsgRequest(
			action,
			name,
			message,
		))
		p2pConn := P2PConnMock{P2PWriteToUDPMock{send: returnSend}}

		writer := NewP2PWriter(name, &p2pConn, addr)

		send, err := writer.Write(action, message)
		assert.NoError(err)
		assert.Equal(returnSend, send)
		assert.Equal(expectedBytes, p2pConn.writeToUDPMock.b)
		assert.Equal(addr, p2pConn.writeToUDPMock.addr)
	})

	t.Run("test_write_fail_marshal", func(t *testing.T) {
		expectedError := fmt.Errorf("Fail")
		action := "fake"
		message := "FakeMessage"
		name := "fakeP2PWriter"
		port := "50000"
		addr, _ := net.ResolveUDPAddr("udp4", port)
		p2pConn := P2PConnMock{P2PWriteToUDPMock{}}

		writer := NewP2PWriter(name, &p2pConn, addr)
		writer.marshal = FailMarshal(expectedError)

		_, err := writer.Write(action, message)
		assert.Error(err, expectedError.Error())
	})
}
