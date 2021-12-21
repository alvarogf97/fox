package stun

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/alvarogf97/fox/pkg/msg"
)

// Interface for udp stun server connection
type UDPStunConn interface {
	ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error)
	WriteToUDP(b []byte, addr *net.UDPAddr) (int, error)
	Close() error
}

// Interface for Stun client
type StunClient interface {
	Collect() error
	Request(peername string, action string, message string, timeout int) (*msg.MsgResponse, error)
	Listen() *msg.MsgResponse
}

// Default stun client that handles stun
// server comunication in the easiest possible way.
type DefaultStunClient struct {
	conn           UDPStunConn
	addr           *net.UDPAddr
	requests       chan *msg.MsgResponse
	registrations  chan *msg.MsgResponse
	disconnections chan *msg.MsgResponse
	peerMsgs       chan *msg.MsgResponse
	isListening    bool
	options        ClientStunOptions

	// marshaller
	marshal   func(v interface{}) ([]byte, error)
	unmarshal func(data []byte, v interface{}) error
}

// Logs the given messages only if logging
// is enabled for the stun
func (client DefaultStunClient) log(v ...interface{}) {
	if client.options.logging {
		log.Println(v[:])
	}
}

// Reads message from the given channel
// without blocking forever due to the
// given timeout in seconds
func (client DefaultStunClient) readChannelWithTimeout(ch chan *msg.MsgResponse, timeout int) (*msg.MsgResponse, error) {
	select {
	case response := <-ch:
		return response, nil
	case <-time.After(time.Duration(timeout) * time.Second):
		return nil, fmt.Errorf("timeout")
	}
}

// Returns the channel that handles the given action,
// If there's no channel that handles it, and error
// will be returned
func (client DefaultStunClient) getActionChannel(action string) (chan *msg.MsgResponse, error) {
	switch action {
	case msg.STUN_ACTION_NEW, msg.PEER_ACTION_NEW:
		return client.registrations, nil
	case msg.STUN_ACTION_GET, msg.PEER_ACTION_GET:
		return client.requests, nil
	case msg.STUN_ACTION_DISCONNECT, msg.PEER_ACTION_DISCONNECT:
		return client.disconnections, nil
	default:
		return nil, fmt.Errorf("unrecognized Stun action `%s`", action)
	}
}

// Starts a goroutine that handles stun server
// responses. This goroutine will saves responses
// into the channels to whom them belongs. This
// method cannot be invoked twice.
func (client *DefaultStunClient) Collect() error {
	if client.isListening {
		return fmt.Errorf("client is already listening connections")
	}

	go func() {
		buff := make([]byte, msg.MAX_MESSAGE_SIZE)
		for {
			var response msg.MsgResponse

			// Wait until there's something in the socket
			// that needs to be read
			bytesRead, _, err := client.conn.ReadFromUDP(buff)
			if err != nil {
				client.log("Get response from server failed ", err)
				continue
			}

			err = client.unmarshal(buff[:bytesRead], &response)
			if err != nil {
				client.log("Unmarshal server response failed ", err)
				continue
			}

			// Saves the response into the channel to whom it belongs
			channel, err := client.getActionChannel(response.Action)
			if err != nil {
				channel = client.peerMsgs
			}

			channel <- &response

			// exit goroutine if disconnect from the P2P network
			if response.Action == msg.PEER_ACTION_DISCONNECT && !response.HasError {
				client.isListening = false
				break
			}
		}
	}()
	client.isListening = true
	return nil
}

// Request stun server with the given paramenters
func (client DefaultStunClient) Request(peername string, action string, message string, timeout int) (*msg.MsgResponse, error) {
	// get the channel that handles the request
	channel, err := client.getActionChannel(action)
	if err != nil {
		return nil, err
	}

	// serializes the request
	request := msg.NewMsgRequest(
		action,
		peername,
		message,
	)

	payload, err := client.marshal(request)
	if err != nil {
		return nil, fmt.Errorf("cannot serialize the request %s", err)
	}

	// send request to the stun server
	if _, err := client.conn.WriteToUDP(payload, client.addr); err != nil {
		return nil, fmt.Errorf("write to UDP failed: %s", err)
	}

	// read the response from the channel
	response, err := client.readChannelWithTimeout(channel, timeout)
	if err != nil {
		return nil, err
	}

	if response.HasError {
		return nil, fmt.Errorf(response.Message)
	}
	return response, nil
}

// Listen for incoming P2P messages.
// This method should be used inside goroutine
// or infinite loop and handle the returned
// messages as wanted
func (client DefaultStunClient) Listen() *msg.MsgResponse {
	return <-client.peerMsgs
}

// Creates a new Stun client
func NewDefaultStunClient(conn UDPStunConn, addr *net.UDPAddr, options ClientStunOptions) *DefaultStunClient {
	return &DefaultStunClient{
		peerMsgs:       make(chan *msg.MsgResponse, options.maxMsgInQueue),
		conn:           conn,
		addr:           addr,
		requests:       make(chan *msg.MsgResponse),
		registrations:  make(chan *msg.MsgResponse),
		disconnections: make(chan *msg.MsgResponse),
		options:        options,
		marshal:        json.Marshal,
		unmarshal:      json.Unmarshal,
	}
}
