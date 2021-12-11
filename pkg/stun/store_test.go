package stun

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemoryPeerConnectionStoreSavePeerRemoteAddr(t *testing.T) {
	assert := require.New(t)

	t.Run("test_save_peer_remote_addr_success", func(t *testing.T) {
		peer := "dog"
		addr := "127.0.0.1:50000"
		store := NewMemoryPeerConnectionStore()

		err := store.SavePeerRemoteAddr(peer, addr)

		value, exists := store.peers[peer]

		assert.NoError(err)
		assert.True(exists)
		assert.Equal(addr, value)
	})

	t.Run("test_save_peer_remote_addr_fail_exists", func(t *testing.T) {
		peer := "dog"
		addr := "127.0.0.1:50000"
		store := NewMemoryPeerConnectionStore()
		store.peers[peer] = addr

		err := store.SavePeerRemoteAddr(peer, addr)

		assert.Error(err)
	})

}

func TestMemoryPeerConnectionStoreDeletePeerRemoteAddr(t *testing.T) {
	assert := require.New(t)

	t.Run("test_delete_peer_remote_addr_success", func(t *testing.T) {
		peer := "dog"
		addr := "127.0.0.1:50000"
		store := NewMemoryPeerConnectionStore()
		store.peers[peer] = addr

		err := store.DeletePeerRemoteAddr(peer)

		assert.NoError(err)
	})

	t.Run("test_delete_peer_remote_addr_fail_not_exists", func(t *testing.T) {
		peer := "dog"
		store := NewMemoryPeerConnectionStore()

		err := store.DeletePeerRemoteAddr(peer)

		assert.Error(err)
	})

}

func TestMemoryPeerConnectionStoreGetPeerRemoteAddr(t *testing.T) {
	assert := require.New(t)

	t.Run("test_get_peer_remote_addr_success", func(t *testing.T) {
		peer := "dog"
		addr := "127.0.0.1:50000"
		store := NewMemoryPeerConnectionStore()
		store.peers[peer] = addr

		result, err := store.GetPeerRemoteAddr(peer)

		assert.NoError(err)
		assert.Equal(addr, result)
	})

	t.Run("test_get_peer_remote_addr_fail_not_exists", func(t *testing.T) {
		peer := "dog"
		store := NewMemoryPeerConnectionStore()

		_, err := store.GetPeerRemoteAddr(peer)

		assert.Error(err)
	})

}

func TestMemoryPeerConnectionStoreGetConnectedPeers(t *testing.T) {
	assert := require.New(t)

	t.Run("test_get_connected_peers_success", func(t *testing.T) {
		peer := "dog"
		addr := "127.0.0.1:50000"
		expected := []PeerInfo{{peer, addr}}
		store := NewMemoryPeerConnectionStore()
		store.peers[peer] = addr

		result, err := store.GetConnectedPeers()

		assert.NoError(err)
		assert.Equal(expected, result)
	})

}
