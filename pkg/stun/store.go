package stun

import (
	"fmt"
	"sync"
)

type PeerInfo struct {
	Peername string
	Addr     string
}

// Interface for structs that allow to store
// Peers information about their addresses
type PeerConnectionStore interface {
	SavePeerRemoteAddr(peer string, addr string) error
	DeletePeerRemoteAddr(peer string) error
	GetPeerRemoteAddr(peer string) (string, error)
	GetConnectedPeers() ([]PeerInfo, error)
}

// Peer connection store in memory.
// Please avoid to use it for large
// networks. It allow read/write
// async operations thanks to the RWMutex
// implementation
type memoryPeerConnectionStore struct {
	sync.RWMutex
	peers map[string]string
}

// Saves the given addr for the given peer
func (store *memoryPeerConnectionStore) SavePeerRemoteAddr(peer string, addr string) error {
	// checks that the given peer does not exist in the network
	_, err := store.GetPeerRemoteAddr(peer)
	if err == nil {
		return fmt.Errorf("peer `%s` already registered", peer)
	}

	store.Lock()
	store.peers[peer] = addr
	store.Unlock()

	return nil
}

// Removes the given peer addr if the peer exists
func (store *memoryPeerConnectionStore) DeletePeerRemoteAddr(peer string) error {
	// checks that the given peer exists in the network
	_, err := store.GetPeerRemoteAddr(peer)
	if err != nil {
		return fmt.Errorf("peer `%s` does not exist", peer)
	}

	store.Lock()
	delete(store.peers, peer)
	store.Unlock()
	return nil
}

// Retrieves the peer addr by giving his name
func (store *memoryPeerConnectionStore) GetPeerRemoteAddr(peer string) (string, error) {
	var err error
	store.RLock()
	addr, exists := store.peers[peer]
	if !exists {
		err = fmt.Errorf("addr of peer %s not found", peer)
	}

	store.RUnlock()
	return addr, err
}

// Get connected peers to the stun server
func (store *memoryPeerConnectionStore) GetConnectedPeers() ([]PeerInfo, error) {
	peernames := []PeerInfo{}
	store.RLock()
	for peername, addr := range store.peers {
		peernames = append(peernames, PeerInfo{peername, addr})
	}
	store.RUnlock()
	return peernames, nil
}

// Creates a new memory peer connection store
func NewMemoryPeerConnectionStore() *memoryPeerConnectionStore {
	return &memoryPeerConnectionStore{
		peers: map[string]string{},
	}
}
