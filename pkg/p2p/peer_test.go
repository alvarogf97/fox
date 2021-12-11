package p2p

import (
	"fmt"
	"net"
	"testing"

	"github.com/alvarogf97/fox/pkg/msg"
	"github.com/stretchr/testify/require"
)

func TestPeerConstructors(t *testing.T) {
	assert := require.New(t)

	t.Run("test_peer_constructor_success", func(t *testing.T) {
		name := "FakePeer"
		stunAddr := "127.0.0.1:60001"
		addr := ":50000"
		options := DefaultPeerOptions()

		peer, err := NewPeer(name, stunAddr, addr, options)
		peer.Close()

		assert.NoError(err)
		assert.Equal(name, peer.name)
	})

	t.Run("test_peer_constructor_fail_resolve_stunaddr", func(t *testing.T) {
		name := "FakePeer"
		stunAddr := "malformedaddr"
		addr := ":50000"
		options := DefaultPeerOptions()

		_, err := NewPeer(name, stunAddr, addr, options)
		assert.Error(err)
	})

	t.Run("test_peer_constructor_fail_resolve_laddr", func(t *testing.T) {
		name := "FakePeer"
		stunAddr := "127.0.0.1:60001"
		addr := "fakeport"
		options := DefaultPeerOptions()

		_, err := NewPeer(name, stunAddr, addr, options)
		assert.Error(err)
	})

	t.Run("test_peer_constructor_fail_listen_udp", func(t *testing.T) {
		name := "FakePeer"
		stunAddr := "127.0.0.1:60001"
		addr := ":50000"
		options := DefaultPeerOptions()

		laddr, _ := net.ResolveUDPAddr("udp4", addr)
		conn, _ := net.ListenUDP("udp", laddr)
		defer conn.Close()

		_, err := NewPeer(name, stunAddr, addr, options)
		assert.Error(err)
	})

}

func TestPeerInit(t *testing.T) {
	assert := require.New(t)

	t.Run("test_peer_init_success", func(t *testing.T) {
		name := "FakePeer"
		stunAddr := "127.0.0.1:60001"
		addr := ":50000"
		options := DefaultPeerOptions()
		client := &MockStunClient{collectMock: CollectMock{}, requestMock: RequestMock{}}

		peer, _ := NewPeer(name, stunAddr, addr, options)
		peer.client = client
		defer peer.Close()

		err := peer.Init()

		assert.NoError(err)
		assert.True(peer.initialized)
		assert.Equal(name, client.requestMock.peername)
		assert.Equal(msg.STUN_ACTION_NEW, client.requestMock.action)
		assert.Equal("", client.requestMock.message)
		assert.Equal(options.timeout, client.requestMock.timeout)
	})

	t.Run("test_peer_init_fail_register", func(t *testing.T) {
		expectedError := fmt.Errorf("Whops")
		name := "FakePeer"
		stunAddr := "127.0.0.1:60001"
		addr := ":50000"
		options := DefaultPeerOptions()
		client := MockStunClient{collectMock: CollectMock{}, requestMock: RequestMock{err: expectedError}}

		peer, _ := NewPeer(name, stunAddr, addr, options)
		peer.client = &client
		defer peer.Close()

		err := peer.Init()

		assert.Error(err, expectedError.Error())
	})

}

func TestPeerConnect(t *testing.T) {
	assert := require.New(t)

	t.Run("test_peer_connect_success", func(t *testing.T) {
		peername := "anotherPeer"
		name := "FakePeer"
		stunAddr := "127.0.0.1:60001"
		addr := ":50000"
		options := DefaultPeerOptions()
		expectedMsg := ":50000"
		response := msg.NewMsgResponse("", false, name, expectedMsg)
		client := &MockStunClient{requestMock: RequestMock{response: &response}}

		peer, _ := NewPeer(name, stunAddr, addr, options)
		peer.client = client
		peer.initialized = true
		defer peer.Close()

		writer, err := peer.Connect(peername)

		assert.NoError(err)
		assert.Equal(name, writer.name)
		assert.Equal(name, client.requestMock.peername)
		assert.Equal(msg.STUN_ACTION_GET, client.requestMock.action)
		assert.Equal(peername, client.requestMock.message)
		assert.Equal(options.timeout, client.requestMock.timeout)
	})

	t.Run("test_peer_connect_fail_not_initialized", func(t *testing.T) {
		peername := "anotherPeer"
		name := "FakePeer"
		stunAddr := "127.0.0.1:60001"
		addr := ":50000"
		options := DefaultPeerOptions()
		expectedMsg := ":50000"
		response := msg.NewMsgResponse("", false, name, expectedMsg)
		client := &MockStunClient{requestMock: RequestMock{response: &response}}

		peer, _ := NewPeer(name, stunAddr, addr, options)
		peer.client = client
		defer peer.Close()

		_, err := peer.Connect(peername)

		assert.Error(err)
	})

	t.Run("test_peer_connect_fail_get_peer_addr", func(t *testing.T) {
		peername := "anotherPeer"
		name := "FakePeer"
		stunAddr := "127.0.0.1:60001"
		addr := ":50000"
		options := DefaultPeerOptions()
		responseError := fmt.Errorf("Fail")
		client := &MockStunClient{requestMock: RequestMock{err: responseError}}

		peer, _ := NewPeer(name, stunAddr, addr, options)
		peer.client = client
		peer.initialized = true
		defer peer.Close()

		_, err := peer.Connect(peername)

		assert.Error(err)
	})

	t.Run("test_peer_connect_fail_resolve_peer_addr", func(t *testing.T) {
		peername := "anotherPeer"
		name := "FakePeer"
		stunAddr := "127.0.0.1:60001"
		addr := ":50000"
		options := DefaultPeerOptions()
		expectedMsg := "fakeaddr"
		response := msg.NewMsgResponse("", false, name, expectedMsg)
		client := &MockStunClient{requestMock: RequestMock{response: &response}}

		peer, _ := NewPeer(name, stunAddr, addr, options)
		peer.client = client
		peer.initialized = true
		defer peer.Close()

		_, err := peer.Connect(peername)

		assert.Error(err)
	})

}

func TestPeerDisconnect(t *testing.T) {
	assert := require.New(t)

	t.Run("test_peer_diconnect_success", func(t *testing.T) {
		name := "FakePeer"
		stunAddr := "127.0.0.1:60001"
		addr := ":50000"
		options := DefaultPeerOptions()
		expectedMsg := ":50000"
		response := msg.NewMsgResponse("", false, name, expectedMsg)
		client := &MockStunClient{requestMock: RequestMock{response: &response}}

		peer, _ := NewPeer(name, stunAddr, addr, options)
		peer.client = client
		peer.initialized = true
		defer peer.Close()

		err := peer.Disconnect()

		assert.NoError(err)
		assert.Equal(name, client.requestMock.peername)
		assert.Equal(msg.STUN_ACTION_DISCONNECT, client.requestMock.action)
		assert.Equal("", client.requestMock.message)
		assert.Equal(options.timeout, client.requestMock.timeout)
	})

	t.Run("test_peer_diconnect_fail_not_initialized", func(t *testing.T) {
		name := "FakePeer"
		stunAddr := "127.0.0.1:60001"
		addr := ":50000"
		options := DefaultPeerOptions()
		client := &MockStunClient{requestMock: RequestMock{}}

		peer, _ := NewPeer(name, stunAddr, addr, options)
		peer.client = client
		peer.initialized = false
		defer peer.Close()

		err := peer.Disconnect()

		assert.Error(err)
	})

}

func TestPeerListen(t *testing.T) {
	assert := require.New(t)

	t.Run("test_peer_listen_success", func(t *testing.T) {
		name := "FakePeer"
		stunAddr := "127.0.0.1:60001"
		addr := ":50000"
		options := DefaultPeerOptions()
		response := msg.NewMsgResponse("", false, name, "bonks")
		client := &MockStunClient{listenMock: ListenMock{response: &response}}

		peer, _ := NewPeer(name, stunAddr, addr, options)
		peer.client = client
		peer.initialized = true
		defer peer.Close()

		msg, err := peer.Listen()

		assert.NoError(err)
		assert.Equal(response.Action, msg.Action)
		assert.Equal(response.Peername, msg.Peername)
		assert.Equal(response.HasError, msg.HasError)
		assert.Equal(response.Message, msg.Message)
	})

	t.Run("test_peer_listen_fail_not_initialized", func(t *testing.T) {
		name := "FakePeer"
		stunAddr := "127.0.0.1:60001"
		addr := ":50000"
		options := DefaultPeerOptions()
		client := &MockStunClient{}

		peer, _ := NewPeer(name, stunAddr, addr, options)
		peer.client = client
		peer.initialized = false
		defer peer.Close()

		_, err := peer.Listen()

		assert.Error(err)
	})

}
