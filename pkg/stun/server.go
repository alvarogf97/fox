package stun

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/alvarogf97/fox/pkg/msg"
)

type Stun struct {
	saddr   string
	conn    UDPStunConn
	store   PeerConnectionStore
	options StunOptions

	// marshaller
	marshal   func(v interface{}) ([]byte, error)
	unmarshal func(data []byte, v interface{}) error
}

// Logs the given messages only if logging
// is enabled for the stun
func (stun Stun) log(v ...interface{}) {
	if stun.options.logging {
		log.Println(v[:])
	}
}

// Response the peer request
// with the given information
func (stun Stun) sendResponse(action string, hasError bool, peername string, message string, addr *net.UDPAddr) (int, error) {
	response := msg.NewMsgResponse(
		action,
		hasError,
		peername,
		message,
	)
	serialized, err := stun.marshal(response)
	if err != nil {
		return 0, err
	}
	return stun.conn.WriteToUDP(serialized, addr)
}

// Handles peer disconnect request
func (stun Stun) handleDisconnectRequest(request msg.MsgRequest, addr *net.UDPAddr) error {
	if err := stun.store.DeletePeerRemoteAddr(request.Peername); err != nil {
		stun.Error(msg.PEER_ACTION_DISCONNECT, request.Peername, err.Error(), addr)
		return err
	}

	// Returns the peer that he has been disconnected
	// successfully
	if _, err := stun.Response(msg.PEER_ACTION_DISCONNECT, request.Peername, "", addr); err != nil {
		ferr := fmt.Sprintf("Error sending response with action %s to %s : %s", msg.PEER_ACTION_GET, addr, err)
		stun.Error(msg.PEER_ACTION_DISCONNECT, request.Peername, ferr, addr)
		return err
	}

	return nil
}

// Handles peer registration request by saving
// the incoming address and peername into the
// stun store
func (stun Stun) handleNewRequest(request msg.MsgRequest, addr *net.UDPAddr) error {
	remoteAddr := fmt.Sprintf("%s:%d", addr.IP, addr.Port)
	if err := stun.store.SavePeerRemoteAddr(request.Peername, remoteAddr); err != nil {
		// Checks the addressed due to the peer could be down
		// and tries to recconect, so if the saved address and
		// the incoming one are the same return as normal
		savedAddr, _ := stun.store.GetPeerRemoteAddr(request.Peername)
		if remoteAddr != savedAddr {
			stun.Error(msg.PEER_ACTION_NEW, request.Peername, err.Error(), addr)
			return err
		}
	}

	if _, err := stun.Response(msg.PEER_ACTION_NEW, request.Peername, remoteAddr, addr); err != nil {
		ferr := fmt.Sprintf("Error sending response with action %s to %s : %s", msg.PEER_ACTION_NEW, addr, err)
		stun.Error(msg.PEER_ACTION_NEW, request.Peername, ferr, addr)
		return err
	}

	return nil
}

// Handles peer connect request by send to the
// requested peer the addr of the peer he wants
// to establish a connection. This method will
// check the peername exists in the network and
// is online
func (stun Stun) handleGetRequest(request msg.MsgRequest, addr *net.UDPAddr) error {
	peername := request.Message

	// checks the requested peer is registered in the network
	peerAddr, err := stun.store.GetPeerRemoteAddr(peername)
	if err != nil {
		stun.Error(msg.PEER_ACTION_GET, request.Peername, err.Error(), addr)
		return err
	}

	// Returns the requested address to peer
	if _, err := stun.Response(msg.PEER_ACTION_GET, request.Peername, peerAddr, addr); err != nil {
		ferr := fmt.Sprintf("Error sending response with action %s to %s : %s", msg.PEER_ACTION_GET, addr, err)
		stun.Error(msg.PEER_ACTION_GET, request.Peername, ferr, addr)
		return err
	}

	return nil
}

// Handles incomming request data and
// returns the handled action name if
// it can be handled, otherwise returns
// an error
func (stun Stun) handle(data []byte, addr *net.UDPAddr) (string, error) {
	// marshal data into request struct
	var request msg.MsgRequest
	err := stun.unmarshal(data, &request)
	if err != nil {
		stun.log("Cannot unmarhsal request: ", err)
		return "", err
	}

	// handle request action
	switch request.Action {
	case msg.STUN_ACTION_NEW:
		err := stun.handleNewRequest(request, addr)
		return msg.STUN_ACTION_NEW, err
	case msg.STUN_ACTION_GET:
		err := stun.handleGetRequest(request, addr)
		return msg.STUN_ACTION_GET, err
	case msg.STUN_ACTION_DISCONNECT:
		err := stun.handleDisconnectRequest(request, addr)
		return msg.STUN_ACTION_DISCONNECT, err
	default:
		message := fmt.Sprintf("unknown action `%s`", request.Action)
		stun.Error(request.Action, request.Peername, message, addr)
		return "", fmt.Errorf(message)
	}
}

// shortcut for `sendResponse` that not
// sends the error flag
func (stun Stun) Response(action string, peername string, message string, addr *net.UDPAddr) (int, error) {
	return stun.sendResponse(action, false, peername, message, addr)
}

// shortcut for `sendResponse` that
// sends the error flag
func (stun Stun) Error(action string, peername string, message string, addr *net.UDPAddr) (int, error) {
	return stun.sendResponse(action, true, peername, message, addr)
}

// Reads data from the udp connection.
// this method exposes the stun connection
// for those ones that requires handle
// the request in their own implementation.
// Feel free to handle the data whathever
// want :D
func (stun Stun) ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error) {
	return stun.conn.ReadFromUDP(b)
}

// Closes stun connection
func (stun Stun) Close() error {
	return stun.conn.Close()
}

// Get saved connected peers into the
// stun store
func (stun Stun) GetConnectedPeers() ([]PeerInfo, error) {
	return stun.store.GetConnectedPeers()
}

// Starts server infinite loop
func (stun Stun) Serve() {
	var buf [2048]byte
	defer stun.Close()
	stun.log("Server is ready to accept UDP connections in ", stun.saddr)

	for {
		n, addr, err := stun.ReadFromUDP(buf[0:])
		if err != nil {
			stun.log(err)
			continue
		}

		go stun.handle(buf[:n], addr)
	}
}

// Creates a new Stun server
func NewStun(saddr string, store PeerConnectionStore, options StunOptions) (*Stun, error) {
	addr, err := net.ResolveUDPAddr("udp4", saddr)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	return &Stun{
		saddr:     saddr,
		conn:      conn,
		store:     store,
		options:   options,
		marshal:   json.Marshal,
		unmarshal: json.Unmarshal,
	}, nil
}
