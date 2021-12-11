package p2p

import (
	"net"

	"github.com/alvarogf97/fox/pkg/msg"
)

// P2P connection mock
type P2PConnMock struct {
	writeToUDPMock P2PWriteToUDPMock
}

func (p2pConnMock *P2PConnMock) WriteToUDP(b []byte, addr *net.UDPAddr) (int, error) {
	p2pConnMock.writeToUDPMock.b = b
	p2pConnMock.writeToUDPMock.addr = addr
	return p2pConnMock.writeToUDPMock.send, p2pConnMock.writeToUDPMock.err
}

type P2PWriteToUDPMock struct {
	b    []byte
	addr *net.UDPAddr

	send int
	err  error
}

// Fake json marhsaler
func FailMarshal(err error) func(v interface{}) ([]byte, error) {
	return func(v interface{}) ([]byte, error) {
		return []byte{}, err
	}
}

// Stun mock client
type MockStunClient struct {
	collectMock CollectMock
	requestMock RequestMock
	listenMock  ListenMock
}

func (client *MockStunClient) Collect() error {
	return client.collectMock.err
}

func (client *MockStunClient) Request(peername string, action string, message string, timeout int) (*msg.MsgResponse, error) {
	client.requestMock.peername = peername
	client.requestMock.action = action
	client.requestMock.message = message
	client.requestMock.timeout = timeout
	return client.requestMock.response, client.requestMock.err
}

func (client MockStunClient) Listen() *msg.MsgResponse {
	return client.listenMock.response
}

type CollectMock struct {
	err error
}

type RequestMock struct {
	peername string
	action   string
	message  string
	timeout  int

	response *msg.MsgResponse
	err      error
}

type ListenMock struct {
	response *msg.MsgResponse
}
