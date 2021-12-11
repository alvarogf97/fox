package stun

import (
	"encoding/json"
	"net"

	"github.com/alvarogf97/fox/pkg/msg"
)

// UDP stun connection mock
type UDPStunConnMock struct {
	readFromUDPMock *ReadFromUDPMock
	writeToUDPMock  *WriteToUDPMock
	closeMock       *CloseMock
}

func (conn *UDPStunConnMock) WriteToUDP(b []byte, addr *net.UDPAddr) (int, error) {
	conn.writeToUDPMock.b = b
	conn.writeToUDPMock.addr = addr
	return conn.writeToUDPMock.send, conn.writeToUDPMock.err
}

func (conn *UDPStunConnMock) ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error) {
	if conn.readFromUDPMock.response != nil {
		buff, _ := json.Marshal(conn.readFromUDPMock.response)
		conn.readFromUDPMock.b = b

		copy(b, buff)
		return len(buff), nil, conn.readFromUDPMock.err
	}

	if len(conn.readFromUDPMock.bb) != 0 {
		copy(b, conn.readFromUDPMock.bb)
		return len(conn.readFromUDPMock.bb), nil, conn.readFromUDPMock.err
	}

	return 0, nil, conn.readFromUDPMock.err
}

func (conn *UDPStunConnMock) Close() error {
	conn.closeMock.hasBeenCalled = true
	return conn.closeMock.err
}

type ReadFromUDPMock struct {
	b []byte

	bb       []byte
	response *msg.MsgResponse
	err      error
}

type WriteToUDPMock struct {
	b    []byte
	addr *net.UDPAddr

	send int
	err  error
}

type CloseMock struct {
	hasBeenCalled bool
	err           error
}

// Fake json marhsaler
func FailMarshal(err error) func(v interface{}) ([]byte, error) {
	return func(v interface{}) ([]byte, error) {
		return []byte{}, err
	}
}

func FailUnmarshal(err error) func(data []byte, v interface{}) error {
	return func(data []byte, v interface{}) error {
		return err
	}
}

// MockStore
type MockPeerConnectionStore struct {
	savePeerRemoteAddrMock   *SavePeerRemoteAddrMock
	deletePeerRemoteAddrMock *DeletePeerRemoteAddrMock
	getPeerRemoteAddrMock    *GetPeerRemoteAddrMock
	getConnectedPeersMock    *GetConnectedPeersMock
}

func (store *MockPeerConnectionStore) SavePeerRemoteAddr(peer string, addr string) error {
	store.savePeerRemoteAddrMock.peer = peer
	store.savePeerRemoteAddrMock.addr = addr
	return store.savePeerRemoteAddrMock.err
}

func (store *MockPeerConnectionStore) DeletePeerRemoteAddr(peer string) error {
	store.deletePeerRemoteAddrMock.peer = peer
	return store.deletePeerRemoteAddrMock.err
}
func (store *MockPeerConnectionStore) GetPeerRemoteAddr(peer string) (string, error) {
	store.getPeerRemoteAddrMock.peer = peer
	return store.getPeerRemoteAddrMock.addr, store.getPeerRemoteAddrMock.err
}
func (store *MockPeerConnectionStore) GetConnectedPeers() ([]PeerInfo, error) {
	return store.getConnectedPeersMock.info, store.getConnectedPeersMock.err
}

type SavePeerRemoteAddrMock struct {
	peer string
	addr string

	err error
}

type DeletePeerRemoteAddrMock struct {
	peer string

	err error
}

type GetPeerRemoteAddrMock struct {
	peer string

	addr string
	err  error
}

type GetConnectedPeersMock struct {
	info []PeerInfo
	err  error
}
