package stun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/alvarogf97/fox/pkg/msg"
	"github.com/stretchr/testify/require"
)

func TestStunConstructor(t *testing.T) {
	assert := require.New(t)

	t.Run("test_new_stun_success", func(t *testing.T) {
		saddr := ":50000"
		store := NewMemoryPeerConnectionStore()
		options := DefaultStunOptions()

		stun, err := NewStun(saddr, store, options)
		stun.Close()

		assert.NoError(err)
		assert.Equal(saddr, stun.saddr)
	})

	t.Run("test_new_stun_fail_resolve_saddr", func(t *testing.T) {
		saddr := "fake"
		store := NewMemoryPeerConnectionStore()
		options := DefaultStunOptions()

		_, err := NewStun(saddr, store, options)

		assert.Error(err)
	})

	t.Run("test_new_stun_fail_listen_udp", func(t *testing.T) {
		saddr := ":50000"
		store := NewMemoryPeerConnectionStore()
		options := DefaultStunOptions()

		addr, _ := net.ResolveUDPAddr("udp4", saddr)
		conn, _ := net.ListenUDP("udp", addr)
		defer conn.Close()

		_, err := NewStun(saddr, store, options)

		assert.Error(err)
	})

}

func TestStunLog(t *testing.T) {
	assert := require.New(t)

	t.Run("test_log", func(t *testing.T) {
		message := "Hello"
		saddr := ":50000"
		store := NewMemoryPeerConnectionStore()
		options := NewStunOptions(true)

		stun, _ := NewStun(saddr, store, options)
		stun.Close()

		var str bytes.Buffer
		log.SetOutput(&str)

		stun.log(message)
		assert.Contains(str.String(), message)
	})

	t.Run("test_not_log", func(t *testing.T) {
		message := "Hello"
		saddr := ":50000"
		store := NewMemoryPeerConnectionStore()
		options := NewStunOptions(false)

		stun, _ := NewStun(saddr, store, options)
		stun.Close()

		var str bytes.Buffer
		log.SetOutput(&str)

		stun.log(message)
		assert.NotContains(str.String(), message)
	})

}

func TestStunSendResponse(t *testing.T) {
	assert := require.New(t)

	t.Run("test_send_response_success", func(t *testing.T) {
		saddr := ":50000"
		addr, _ := net.ResolveUDPAddr("udp4", saddr)
		store := NewMemoryPeerConnectionStore()
		options := NewStunOptions(true)
		send := 10
		conn := &UDPStunConnMock{writeToUDPMock: &WriteToUDPMock{send: send}}

		stun, _ := NewStun(saddr, store, options)
		stun.Close()

		stun.conn = conn

		rsend, err := stun.sendResponse("fake", false, "dog", "bonks", addr)

		assert.NoError(err)
		assert.Equal(send, rsend)
	})

	t.Run("test_send_response_fail_marshal", func(t *testing.T) {
		rerr := fmt.Errorf("error")
		saddr := ":50000"
		addr, _ := net.ResolveUDPAddr("udp4", saddr)
		store := NewMemoryPeerConnectionStore()
		options := NewStunOptions(true)
		send := 10
		conn := &UDPStunConnMock{writeToUDPMock: &WriteToUDPMock{send: send}}

		stun, _ := NewStun(saddr, store, options)
		stun.Close()

		stun.conn = conn
		stun.marshal = FailMarshal(rerr)

		_, err := stun.sendResponse("fake", false, "dog", "bonks", addr)

		assert.Error(err, rerr.Error())
	})

}

func TestStunHandleDisconnectRequest(t *testing.T) {
	assert := require.New(t)

	t.Run("test_disconnect_success", func(t *testing.T) {
		saddr := ":50000"
		addr, _ := net.ResolveUDPAddr("udp4", saddr)
		request := msg.NewMsgRequest(msg.STUN_ACTION_DISCONNECT, "dog", "bonks")
		store := &MockPeerConnectionStore{deletePeerRemoteAddrMock: &DeletePeerRemoteAddrMock{}}
		options := NewStunOptions(true)
		send := 10
		conn := &UDPStunConnMock{writeToUDPMock: &WriteToUDPMock{send: send}}

		stun, _ := NewStun(saddr, store, options)
		stun.Close()
		stun.conn = conn

		err := stun.handleDisconnectRequest(request, addr)

		assert.NoError(err)
		assert.Equal(request.Peername, store.deletePeerRemoteAddrMock.peer)
		assert.Equal(send, conn.writeToUDPMock.send)
	})

	t.Run("test_disconnect_fail_delete_peer_remote_addr", func(t *testing.T) {
		rerr := fmt.Errorf("Error")
		saddr := ":50000"
		addr, _ := net.ResolveUDPAddr("udp4", saddr)
		request := msg.NewMsgRequest(msg.STUN_ACTION_DISCONNECT, "dog", "bonks")
		store := &MockPeerConnectionStore{deletePeerRemoteAddrMock: &DeletePeerRemoteAddrMock{err: rerr}}
		options := NewStunOptions(true)
		send := 10
		conn := &UDPStunConnMock{writeToUDPMock: &WriteToUDPMock{send: send}}

		stun, _ := NewStun(saddr, store, options)
		stun.Close()
		stun.conn = conn

		err := stun.handleDisconnectRequest(request, addr)

		assert.Error(err, rerr.Error())
	})

	t.Run("test_disconnect_fail_response", func(t *testing.T) {
		rerr := fmt.Errorf("Error")
		saddr := ":50000"
		addr, _ := net.ResolveUDPAddr("udp4", saddr)
		request := msg.NewMsgRequest(msg.STUN_ACTION_DISCONNECT, "dog", "bonks")
		store := &MockPeerConnectionStore{deletePeerRemoteAddrMock: &DeletePeerRemoteAddrMock{}}
		options := NewStunOptions(true)
		conn := &UDPStunConnMock{writeToUDPMock: &WriteToUDPMock{err: rerr}}

		stun, _ := NewStun(saddr, store, options)
		stun.Close()
		stun.conn = conn

		err := stun.handleDisconnectRequest(request, addr)

		assert.Error(err, rerr.Error())
	})
}

func TestStunHandleNewRequest(t *testing.T) {
	assert := require.New(t)

	t.Run("test_new_request_success", func(t *testing.T) {
		saddr := ":50000"
		addr, _ := net.ResolveUDPAddr("udp4", saddr)
		remoteAddr := fmt.Sprintf("%s:%d", addr.IP, addr.Port)
		request := msg.NewMsgRequest(msg.STUN_ACTION_NEW, "dog", "bonks")
		store := &MockPeerConnectionStore{savePeerRemoteAddrMock: &SavePeerRemoteAddrMock{}}
		options := NewStunOptions(true)
		send := 10
		conn := &UDPStunConnMock{writeToUDPMock: &WriteToUDPMock{send: send}}

		stun, _ := NewStun(saddr, store, options)
		stun.Close()
		stun.conn = conn

		err := stun.handleNewRequest(request, addr)

		assert.NoError(err)
		assert.Equal(request.Peername, store.savePeerRemoteAddrMock.peer)
		assert.Equal(remoteAddr, store.savePeerRemoteAddrMock.addr)
		assert.Equal(send, conn.writeToUDPMock.send)
	})

	t.Run("test_new_request_fail_save_peer_remote_addr", func(t *testing.T) {
		rerr := fmt.Errorf("Error")
		saddr := ":50000"
		addr, _ := net.ResolveUDPAddr("udp4", saddr)
		request := msg.NewMsgRequest(msg.STUN_ACTION_NEW, "dog", "bonks")
		store := &MockPeerConnectionStore{
			savePeerRemoteAddrMock: &SavePeerRemoteAddrMock{err: rerr},
			getPeerRemoteAddrMock:  &GetPeerRemoteAddrMock{addr: "fake"},
		}
		options := NewStunOptions(true)
		send := 10
		conn := &UDPStunConnMock{writeToUDPMock: &WriteToUDPMock{send: send}}

		stun, _ := NewStun(saddr, store, options)
		stun.Close()
		stun.conn = conn

		err := stun.handleNewRequest(request, addr)

		assert.Equal(request.Peername, store.getPeerRemoteAddrMock.peer)
		assert.Equal(send, conn.writeToUDPMock.send)
		assert.Error(err, rerr.Error())
	})

	t.Run("test_new_request_fail_response", func(t *testing.T) {
		rerr := fmt.Errorf("Error")
		saddr := ":50000"
		addr, _ := net.ResolveUDPAddr("udp4", saddr)
		request := msg.NewMsgRequest(msg.STUN_ACTION_NEW, "dog", "bonks")
		store := &MockPeerConnectionStore{savePeerRemoteAddrMock: &SavePeerRemoteAddrMock{}}
		options := NewStunOptions(true)
		conn := &UDPStunConnMock{writeToUDPMock: &WriteToUDPMock{err: rerr}}

		stun, _ := NewStun(saddr, store, options)
		stun.Close()
		stun.conn = conn

		err := stun.handleNewRequest(request, addr)

		assert.Error(err, rerr.Error())
	})
}

func TestStunHandleGetRequest(t *testing.T) {
	assert := require.New(t)

	t.Run("test_get_request_success", func(t *testing.T) {
		saddr := ":50000"
		addr, _ := net.ResolveUDPAddr("udp4", saddr)
		request := msg.NewMsgRequest(msg.STUN_ACTION_GET, "dog", "bonks")
		store := &MockPeerConnectionStore{getPeerRemoteAddrMock: &GetPeerRemoteAddrMock{}}
		options := NewStunOptions(true)
		send := 10
		conn := &UDPStunConnMock{writeToUDPMock: &WriteToUDPMock{send: send}}

		stun, _ := NewStun(saddr, store, options)
		stun.Close()
		stun.conn = conn

		err := stun.handleGetRequest(request, addr)

		assert.NoError(err)
		assert.Equal(request.Message, store.getPeerRemoteAddrMock.peer)
		assert.Equal(send, conn.writeToUDPMock.send)
	})

	t.Run("test_get_request_fail_get_peer_remote_addr", func(t *testing.T) {
		rerr := fmt.Errorf("Error")
		assert := require.New(t)
		saddr := ":50000"
		addr, _ := net.ResolveUDPAddr("udp4", saddr)
		request := msg.NewMsgRequest(msg.STUN_ACTION_GET, "dog", "bonks")
		store := &MockPeerConnectionStore{getPeerRemoteAddrMock: &GetPeerRemoteAddrMock{err: rerr}}
		options := NewStunOptions(true)
		send := 10
		conn := &UDPStunConnMock{writeToUDPMock: &WriteToUDPMock{send: send}}

		stun, _ := NewStun(saddr, store, options)
		stun.Close()
		stun.conn = conn

		err := stun.handleGetRequest(request, addr)

		assert.Equal(request.Message, store.getPeerRemoteAddrMock.peer)
		assert.Equal(send, conn.writeToUDPMock.send)
		assert.Error(err, rerr.Error())
	})

	t.Run("test_get_request_fail_response", func(t *testing.T) {
		rerr := fmt.Errorf("Error")
		saddr := ":50000"
		addr, _ := net.ResolveUDPAddr("udp4", saddr)
		request := msg.NewMsgRequest(msg.STUN_ACTION_GET, "dog", "bonks")
		store := &MockPeerConnectionStore{getPeerRemoteAddrMock: &GetPeerRemoteAddrMock{}}
		options := NewStunOptions(true)
		conn := &UDPStunConnMock{writeToUDPMock: &WriteToUDPMock{err: rerr}}

		stun, _ := NewStun(saddr, store, options)
		stun.Close()
		stun.conn = conn

		err := stun.handleGetRequest(request, addr)

		assert.Error(err, rerr.Error())
	})

}

func TestStunHandle(t *testing.T) {
	assert := require.New(t)

	t.Run("test_handle_action_new", func(t *testing.T) {
		saddr := ":50000"
		addr, _ := net.ResolveUDPAddr("udp4", saddr)
		store := NewMemoryPeerConnectionStore()
		options := NewStunOptions(true)
		send := 10
		conn := &UDPStunConnMock{writeToUDPMock: &WriteToUDPMock{send: send}}

		stun, _ := NewStun(saddr, store, options)
		stun.Close()
		stun.conn = conn

		request := msg.NewMsgRequest(msg.STUN_ACTION_NEW, "dog", "bonks")
		brequest, _ := json.Marshal(&request)

		action, _ := stun.handle(brequest, addr)

		assert.Equal(msg.STUN_ACTION_NEW, action)
	})

	t.Run("test_handle_action_get", func(t *testing.T) {
		saddr := ":50000"
		addr, _ := net.ResolveUDPAddr("udp4", saddr)
		store := NewMemoryPeerConnectionStore()
		options := NewStunOptions(true)
		send := 10
		conn := &UDPStunConnMock{writeToUDPMock: &WriteToUDPMock{send: send}}

		stun, _ := NewStun(saddr, store, options)
		stun.Close()
		stun.conn = conn

		request := msg.NewMsgRequest(msg.STUN_ACTION_GET, "dog", "bonks")
		brequest, _ := json.Marshal(&request)

		action, _ := stun.handle(brequest, addr)

		assert.Equal(msg.STUN_ACTION_GET, action)
	})

	t.Run("test_handle_action_disconnect", func(t *testing.T) {
		saddr := ":50000"
		addr, _ := net.ResolveUDPAddr("udp4", saddr)
		store := NewMemoryPeerConnectionStore()
		options := NewStunOptions(true)
		send := 10
		conn := &UDPStunConnMock{writeToUDPMock: &WriteToUDPMock{send: send}}

		stun, _ := NewStun(saddr, store, options)
		stun.Close()
		stun.conn = conn

		request := msg.NewMsgRequest(msg.STUN_ACTION_DISCONNECT, "dog", "bonks")
		brequest, _ := json.Marshal(&request)

		action, _ := stun.handle(brequest, addr)

		assert.Equal(msg.STUN_ACTION_DISCONNECT, action)
	})

	t.Run("test_handle_action_fail_unknown", func(t *testing.T) {
		saddr := ":50000"
		addr, _ := net.ResolveUDPAddr("udp4", saddr)
		store := NewMemoryPeerConnectionStore()
		options := NewStunOptions(true)
		send := 10
		conn := &UDPStunConnMock{writeToUDPMock: &WriteToUDPMock{send: send}}

		stun, _ := NewStun(saddr, store, options)
		stun.Close()
		stun.conn = conn

		request := msg.NewMsgRequest("fake action", "dog", "bonks")
		brequest, _ := json.Marshal(&request)

		_, err := stun.handle(brequest, addr)

		assert.Error(err)
	})
}

func TestStunResponse(t *testing.T) {
	assert := require.New(t)

	t.Run("test_response", func(t *testing.T) {
		saddr := ":50000"
		addr, _ := net.ResolveUDPAddr("udp4", saddr)
		store := NewMemoryPeerConnectionStore()
		options := NewStunOptions(true)
		send := 10
		conn := &UDPStunConnMock{writeToUDPMock: &WriteToUDPMock{send: send}}

		stun, _ := NewStun(saddr, store, options)
		stun.Close()

		stun.conn = conn

		rsend, err := stun.Response("fake", "dog", "bonks", addr)

		assert.NoError(err)
		assert.Equal(send, rsend)
	})

}

func TestStunError(t *testing.T) {
	assert := require.New(t)

	t.Run("test_error", func(t *testing.T) {
		saddr := ":50000"
		addr, _ := net.ResolveUDPAddr("udp4", saddr)
		store := NewMemoryPeerConnectionStore()
		options := NewStunOptions(true)
		send := 10
		conn := &UDPStunConnMock{writeToUDPMock: &WriteToUDPMock{send: send}}

		stun, _ := NewStun(saddr, store, options)
		stun.Close()

		stun.conn = conn

		rsend, err := stun.Error("fake", "dog", "bonks", addr)

		assert.NoError(err)
		assert.Equal(send, rsend)
	})

}

func TestStunReadFromUDP(t *testing.T) {
	assert := require.New(t)

	t.Run("test_read_from_udp", func(t *testing.T) {
		saddr := ":50000"
		store := NewMemoryPeerConnectionStore()
		options := NewStunOptions(true)
		result := []byte("dog")
		conn := &UDPStunConnMock{readFromUDPMock: &ReadFromUDPMock{bb: result}}

		stun, _ := NewStun(saddr, store, options)
		stun.Close()

		stun.conn = conn

		buff := make([]byte, 2048)
		n, _, err := stun.ReadFromUDP(buff)

		assert.NoError(err)
		assert.Equal(result, buff[:n])
	})

}

func TestStunClose(t *testing.T) {
	assert := require.New(t)

	t.Run("test_stun_close", func(t *testing.T) {
		saddr := ":50000"
		store := NewMemoryPeerConnectionStore()
		options := NewStunOptions(true)
		conn := &UDPStunConnMock{closeMock: &CloseMock{}}

		stun, _ := NewStun(saddr, store, options)
		stun.Close()

		stun.conn = conn
		stun.Close()

		assert.True(conn.closeMock.hasBeenCalled)
	})

}

func TestStunGetConnectedPeers(t *testing.T) {
	assert := require.New(t)

	t.Run("test_get_connected_peers", func(t *testing.T) {
		expected := []PeerInfo{{"Fake", "FakeAddr"}}
		saddr := ":50000"
		store := &MockPeerConnectionStore{getConnectedPeersMock: &GetConnectedPeersMock{info: expected}}
		options := NewStunOptions(true)

		stun, _ := NewStun(saddr, store, options)
		stun.Close()

		peers, err := stun.GetConnectedPeers()

		assert.NoError(err)
		assert.Equal(expected, peers)
	})

}

func TestStunServe(t *testing.T) {
	assert := require.New(t)

	t.Run("test_serve", func(t *testing.T) {
		saddr := ":50000"
		store := NewMemoryPeerConnectionStore()
		options := NewStunOptions(true)
		request, _ := json.Marshal(msg.NewMsgRequest("bonks", "dog", "godzilla"))
		conn := &UDPStunConnMock{readFromUDPMock: &ReadFromUDPMock{bb: request}, writeToUDPMock: &WriteToUDPMock{}}

		var str bytes.Buffer
		log.SetOutput(&str)

		stun, _ := NewStun(saddr, store, options)
		stun.Close()
		stun.conn = conn

		go stun.Serve()

		time.Sleep(1 * time.Second)
		assert.Contains(str.String(), "Server is ready")
	})

	t.Run("test_serve_fail_read_from_udp", func(t *testing.T) {
		saddr := ":50000"
		rerr := fmt.Errorf("Error")
		store := NewMemoryPeerConnectionStore()
		options := NewStunOptions(true)
		conn := &UDPStunConnMock{readFromUDPMock: &ReadFromUDPMock{err: rerr}}

		var str bytes.Buffer
		log.SetOutput(&str)

		stun, _ := NewStun(saddr, store, options)
		stun.Close()
		stun.conn = conn

		go stun.Serve()

		time.Sleep(1 * time.Second)
		assert.Contains(str.String(), rerr.Error())
	})

}
