package stun

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/alvarogf97/fox/pkg/msg"
	"github.com/stretchr/testify/require"
)

func TestDefaultStunClientConstructor(t *testing.T) {
	assert := require.New(t)

	t.Run("test_new_default_stun_client_success", func(t *testing.T) {
		conn := &UDPStunConnMock{}
		addr, _ := net.ResolveUDPAddr("udp4", ":50000")
		options := DefaultClientStunOptions()

		client := NewDefaultStunClient(conn, addr, options)

		assert.Equal(addr, client.addr)
	})

}

func TestDefaultStunClientLog(t *testing.T) {
	assert := require.New(t)

	t.Run("test_log", func(t *testing.T) {
		message := "Log me"
		conn := &UDPStunConnMock{}
		addr, _ := net.ResolveUDPAddr("udp4", ":50000")
		options := NewClientStunOptions(true, 10)

		client := NewDefaultStunClient(conn, addr, options)

		var str bytes.Buffer
		log.SetOutput(&str)

		client.log(message)
		assert.Contains(str.String(), message)
	})

	t.Run("test_not_log", func(t *testing.T) {
		message := "Log me"
		conn := &UDPStunConnMock{}
		addr, _ := net.ResolveUDPAddr("udp4", ":50000")
		options := NewClientStunOptions(false, 10)

		client := NewDefaultStunClient(conn, addr, options)

		var str bytes.Buffer
		log.SetOutput(&str)

		client.log(message)
		assert.NotContains(str.String(), message)
	})

}

func TestDefaultStunClientReadChannelWithTimeout(t *testing.T) {
	assert := require.New(t)

	t.Run("test_read_channel_with_timeout_success", func(t *testing.T) {
		conn := &UDPStunConnMock{}
		addr, _ := net.ResolveUDPAddr("udp4", ":50000")
		options := NewClientStunOptions(false, 10)
		client := NewDefaultStunClient(conn, addr, options)

		msgResponse := msg.NewMsgResponse("bonks", false, "dog", "godzilla")
		ch := make(chan *msg.MsgResponse)
		go func() {
			ch <- &msgResponse
		}()

		result, err := client.readChannelWithTimeout(ch, 10)

		assert.NoError(err)
		assert.Equal(msgResponse.Action, result.Action)
		assert.Equal(msgResponse.HasError, result.HasError)
		assert.Equal(msgResponse.Peername, result.Peername)
		assert.Equal(msgResponse.Message, result.Message)
	})

	t.Run("test_read_channel_with_timeout_fail", func(t *testing.T) {
		conn := &UDPStunConnMock{}
		addr, _ := net.ResolveUDPAddr("udp4", ":50000")
		options := NewClientStunOptions(false, 10)
		client := NewDefaultStunClient(conn, addr, options)

		ch := make(chan *msg.MsgResponse)

		_, err := client.readChannelWithTimeout(ch, 1)
		assert.Error(err)
	})

}

func TestDefaultStunClientGetActionChannel(t *testing.T) {
	assert := require.New(t)

	t.Run("test_get_action_channel_registrations", func(t *testing.T) {
		conn := &UDPStunConnMock{}
		addr, _ := net.ResolveUDPAddr("udp4", ":50000")
		options := NewClientStunOptions(false, 10)
		client := NewDefaultStunClient(conn, addr, options)

		ch, err := client.getActionChannel(msg.STUN_ACTION_NEW)

		assert.NoError(err)
		assert.Equal(ch, client.registrations)
	})

	t.Run("test_get_action_channel_requests", func(t *testing.T) {
		conn := &UDPStunConnMock{}
		addr, _ := net.ResolveUDPAddr("udp4", ":50000")
		options := NewClientStunOptions(false, 10)
		client := NewDefaultStunClient(conn, addr, options)

		ch, err := client.getActionChannel(msg.STUN_ACTION_GET)

		assert.NoError(err)
		assert.Equal(ch, client.requests)
	})

	t.Run("test_get_action_channel_disconnections", func(t *testing.T) {
		conn := &UDPStunConnMock{}
		addr, _ := net.ResolveUDPAddr("udp4", ":50000")
		options := NewClientStunOptions(false, 10)
		client := NewDefaultStunClient(conn, addr, options)

		ch, err := client.getActionChannel(msg.STUN_ACTION_DISCONNECT)

		assert.NoError(err)
		assert.Equal(ch, client.disconnections)
	})

	t.Run("test_get_action_channel_fail_unknown_action", func(t *testing.T) {
		conn := &UDPStunConnMock{}
		addr, _ := net.ResolveUDPAddr("udp4", ":50000")
		options := NewClientStunOptions(false, 10)
		client := NewDefaultStunClient(conn, addr, options)

		_, err := client.getActionChannel("Unchannel action")

		assert.Error(err)
	})

}

func TestDefaultStunClientCollect(t *testing.T) {
	assert := require.New(t)

	t.Run("test_client_collect_message", func(t *testing.T) {
		msgResponse := msg.NewMsgResponse("bonks", false, "dog", "godzilla")
		conn := &UDPStunConnMock{readFromUDPMock: &ReadFromUDPMock{response: &msgResponse}}
		addr, _ := net.ResolveUDPAddr("udp4", ":50000")
		options := DefaultClientStunOptions()
		client := NewDefaultStunClient(conn, addr, options)

		go client.Collect()

		result, err := client.readChannelWithTimeout(client.peerMsgs, 1)

		assert.NoError(err)
		assert.Equal(msgResponse.Action, result.Action)
		assert.Equal(msgResponse.HasError, result.HasError)
		assert.Equal(msgResponse.Peername, result.Peername)
		assert.Equal(msgResponse.Message, result.Message)
	})

	t.Run("test_client_collect_action", func(t *testing.T) {
		msgResponse := msg.NewMsgResponse(msg.PEER_ACTION_GET, false, "dog", "godzilla")
		conn := &UDPStunConnMock{readFromUDPMock: &ReadFromUDPMock{response: &msgResponse}}
		addr, _ := net.ResolveUDPAddr("udp4", ":50000")
		options := DefaultClientStunOptions()
		client := NewDefaultStunClient(conn, addr, options)

		go client.Collect()

		result, err := client.readChannelWithTimeout(client.requests, 1)

		assert.NoError(err)
		assert.Equal(msgResponse.Action, result.Action)
		assert.Equal(msgResponse.HasError, result.HasError)
		assert.Equal(msgResponse.Peername, result.Peername)
		assert.Equal(msgResponse.Message, result.Message)
	})

	t.Run("test_client_collect_disconnect", func(t *testing.T) {
		msgResponse := msg.NewMsgResponse(msg.PEER_ACTION_DISCONNECT, false, "dog", "godzilla")
		conn := &UDPStunConnMock{readFromUDPMock: &ReadFromUDPMock{response: &msgResponse}}
		addr, _ := net.ResolveUDPAddr("udp4", ":50000")
		options := DefaultClientStunOptions()
		client := NewDefaultStunClient(conn, addr, options)

		go client.Collect()

		result, err := client.readChannelWithTimeout(client.disconnections, 10)

		assert.NoError(err)
		assert.Equal(msgResponse.Action, result.Action)
		assert.Equal(msgResponse.HasError, result.HasError)
		assert.Equal(msgResponse.Peername, result.Peername)
		assert.Equal(msgResponse.Message, result.Message)
		assert.False(client.isListening)
	})

	t.Run("test_client_collect_fail_read_from_udp", func(t *testing.T) {
		rerr := fmt.Errorf("Error")
		msgResponse := msg.NewMsgResponse(msg.PEER_ACTION_DISCONNECT, false, "dog", "godzilla")
		conn := &UDPStunConnMock{readFromUDPMock: &ReadFromUDPMock{response: &msgResponse, err: rerr}}
		addr, _ := net.ResolveUDPAddr("udp4", ":50000")
		options := DefaultClientStunOptions()
		client := NewDefaultStunClient(conn, addr, options)

		go client.Collect()

		_, err := client.readChannelWithTimeout(client.disconnections, 1)

		assert.Error(err)
	})

	t.Run("test_client_collect_fail_already_listening", func(t *testing.T) {
		conn := &UDPStunConnMock{}
		addr, _ := net.ResolveUDPAddr("udp4", ":50000")
		options := DefaultClientStunOptions()
		client := NewDefaultStunClient(conn, addr, options)
		client.isListening = true

		err := client.Collect()

		assert.Error(err)
	})

	t.Run("test_client_collect_fail_unmarshal", func(t *testing.T) {
		msgResponse := msg.NewMsgResponse(msg.PEER_ACTION_DISCONNECT, false, "dog", "godzilla")
		conn := &UDPStunConnMock{readFromUDPMock: &ReadFromUDPMock{response: &msgResponse}}
		addr, _ := net.ResolveUDPAddr("udp4", ":50000")
		options := DefaultClientStunOptions()
		client := NewDefaultStunClient(conn, addr, options)
		client.unmarshal = FailUnmarshal(fmt.Errorf("Error"))

		go client.Collect()

		_, err := client.readChannelWithTimeout(client.disconnections, 1)

		assert.Error(err)
	})

}

func TestDefaultStunClientRequest(t *testing.T) {
	assert := require.New(t)

	t.Run("test_request_success", func(t *testing.T) {
		msgResponse := msg.NewMsgResponse(msg.PEER_ACTION_GET, false, "dog", "godzilla")
		conn := &UDPStunConnMock{readFromUDPMock: &ReadFromUDPMock{response: &msgResponse}, writeToUDPMock: &WriteToUDPMock{}}
		addr, _ := net.ResolveUDPAddr("udp4", ":50000")
		options := DefaultClientStunOptions()
		client := NewDefaultStunClient(conn, addr, options)

		go client.Collect()

		result, err := client.Request("dog", msg.STUN_ACTION_GET, "", 1)

		assert.NoError(err)
		assert.Equal(msgResponse.Action, result.Action)
		assert.Equal(msgResponse.HasError, result.HasError)
		assert.Equal(msgResponse.Peername, result.Peername)
		assert.Equal(msgResponse.Message, result.Message)
	})

	t.Run("test_request_fail_no_action_channel", func(t *testing.T) {
		conn := &UDPStunConnMock{}
		addr, _ := net.ResolveUDPAddr("udp4", ":50000")
		options := DefaultClientStunOptions()
		client := NewDefaultStunClient(conn, addr, options)

		_, err := client.Request("dog", "fake action", "", 1)

		assert.Error(err)
	})

	t.Run("test_request_fail_marshal", func(t *testing.T) {
		conn := &UDPStunConnMock{}
		addr, _ := net.ResolveUDPAddr("udp4", ":50000")
		options := DefaultClientStunOptions()
		client := NewDefaultStunClient(conn, addr, options)
		rerr := fmt.Errorf("Error")
		client.marshal = FailMarshal(rerr)

		_, err := client.Request("dog", msg.STUN_ACTION_GET, "", 1)

		assert.Error(err, rerr.Error())
	})

	t.Run("test_request_fail_write_to_udp", func(t *testing.T) {
		rerr := fmt.Errorf("Error")
		conn := &UDPStunConnMock{writeToUDPMock: &WriteToUDPMock{err: rerr}}
		addr, _ := net.ResolveUDPAddr("udp4", ":50000")
		options := DefaultClientStunOptions()
		client := NewDefaultStunClient(conn, addr, options)

		_, err := client.Request("dog", msg.STUN_ACTION_GET, "", 1)

		assert.Error(err)
	})

	t.Run("test_request_fail_timeout", func(t *testing.T) {
		conn := &UDPStunConnMock{writeToUDPMock: &WriteToUDPMock{}}
		addr, _ := net.ResolveUDPAddr("udp4", ":50000")
		options := DefaultClientStunOptions()
		client := NewDefaultStunClient(conn, addr, options)

		_, err := client.Request("dog", msg.STUN_ACTION_GET, "", 1)

		assert.Error(err)
	})

	t.Run("test_request_fail_response_error", func(t *testing.T) {
		msgResponse := msg.NewMsgResponse(msg.PEER_ACTION_GET, true, "dog", "godzilla")
		conn := &UDPStunConnMock{readFromUDPMock: &ReadFromUDPMock{response: &msgResponse}, writeToUDPMock: &WriteToUDPMock{}}
		addr, _ := net.ResolveUDPAddr("udp4", ":50000")
		options := DefaultClientStunOptions()
		client := NewDefaultStunClient(conn, addr, options)

		go client.Collect()

		_, err := client.Request("dog", msg.STUN_ACTION_GET, "", 1)

		assert.Error(err)
	})

}

func TestDefaultStunClientListen(t *testing.T) {
	assert := require.New(t)

	t.Run("test_listen_success", func(t *testing.T) {
		conn := &UDPStunConnMock{}
		addr, _ := net.ResolveUDPAddr("udp4", ":50000")
		options := NewClientStunOptions(false, 10)
		client := NewDefaultStunClient(conn, addr, options)

		msgResponse := msg.NewMsgResponse("bonks", false, "dog", "godzilla")
		go func() {
			client.peerMsgs <- &msgResponse
		}()

		result := client.Listen()

		assert.Equal(msgResponse.Action, result.Action)
		assert.Equal(msgResponse.HasError, result.HasError)
		assert.Equal(msgResponse.Peername, result.Peername)
		assert.Equal(msgResponse.Message, result.Message)
	})

}
