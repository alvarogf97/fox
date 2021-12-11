package p2p

import (
	"encoding/json"
	"net"

	"github.com/alvarogf97/fox/pkg/msg"
)

// Interface for P2P connection
type P2PConn interface {
	WriteToUDP(b []byte, addr *net.UDPAddr) (int, error)
}

// P2P message writer
type P2PWriter struct {
	name    string
	conn    P2PConn
	paddr   *net.UDPAddr
	marshal func(v interface{}) ([]byte, error)
}

// Writes the MsgRequest into the P2P connection
func (writer P2PWriter) Write(action string, message string) (int, error) {
	messageRequest := msg.NewMsgRequest(
		action,
		writer.name,
		message,
	)
	request, err := writer.marshal(messageRequest)
	if err != nil {
		return 0, err
	}
	return writer.conn.WriteToUDP(request, writer.paddr)
}

// Creates a new P2P writer
func NewP2PWriter(name string, conn P2PConn, paddr *net.UDPAddr) *P2PWriter {
	return &P2PWriter{
		name:    name,
		conn:    conn,
		paddr:   paddr,
		marshal: json.Marshal,
	}
}
