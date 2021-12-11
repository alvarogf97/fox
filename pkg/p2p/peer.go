package p2p

import (
	"fmt"
	"net"

	"github.com/alvarogf97/fox/pkg/msg"
	"github.com/alvarogf97/fox/pkg/stun"
)

// Peer is a node into a P2P connection,
// it can send messages to other peers and
// listen the incoming ones
type Peer struct {
	name        string
	options     PeerOptions
	initialized bool
	conn        *net.UDPConn
	client      stun.StunClient
}

// Register the current peer into the stun
// server and starts listening for incoming
// messages
func (peer *Peer) Init() error {
	// starts listening incoming messages by
	// using stun client
	go peer.client.Collect()

	// requests stun server in order to register
	// the current peer in the p2p network
	_, err := peer.client.Request(peer.name, msg.STUN_ACTION_NEW, "", peer.options.timeout)
	if err != nil {
		return err
	}
	peer.initialized = true
	return nil
}

// Connects to a peer by givin his name.
// If the peer does no exist in the P2P
// network an error will be raised
func (peer Peer) Connect(peername string) (*P2PWriter, error) {
	if !peer.initialized {
		return nil, fmt.Errorf("Peer needs to be initialized first")
	}
	// requests server about the given peername
	response, err := peer.client.Request(peer.name, msg.STUN_ACTION_GET, peername, peer.options.timeout)
	if err != nil {
		return nil, err
	}

	// tries to resolve given udp address
	paddr, err := net.ResolveUDPAddr("udp4", response.Message)
	if err != nil {
		return nil, fmt.Errorf("resolve peer address failed %s", err)
	}

	// return a P2P wirter through the one you can write
	// messages to the connected peer
	return NewP2PWriter(peer.name, peer.conn, paddr), nil
}

// Disconnects from the P2P network
// so initialized will be back to false
func (peer *Peer) Disconnect() error {
	if !peer.initialized {
		return fmt.Errorf("Peer needs to be initialized first")
	}

	_, err := peer.client.Request(peer.name, msg.STUN_ACTION_DISCONNECT, "", peer.options.timeout)
	return err
}

// Recover P2P messages from the stun server
// queue. This function shoudl be used by a
// goroutine in order to handle the incoming messages
func (peer Peer) Listen() (*msg.MsgResponse, error) {
	if !peer.initialized {
		return nil, fmt.Errorf("Peer needs to be initialized first")
	}
	return peer.client.Listen(), nil
}

// Closes peer connection
func (peer Peer) Close() error {
	return peer.conn.Close()
}

// Creates a new peer
func NewPeer(name string, stunaddr string, addr string, options PeerOptions) (*Peer, error) {
	saddr, err := net.ResolveUDPAddr("udp4", stunaddr)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve server address: %s", err)
	}

	laddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return nil, fmt.Errorf("cannot resolver local address: %s", err)
	}

	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return nil, fmt.Errorf("address already in use: %s", err)
	}

	client := stun.NewDefaultStunClient(conn, saddr, stun.NewClientStunOptions(true, options.maxMsgInQueue))

	return &Peer{
		name:        name,
		options:     options,
		initialized: false,
		conn:        conn,
		client:      client,
	}, nil
}
