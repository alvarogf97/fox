package main

import (
	"fmt"
	"strings"

	"github.com/alvarogf97/fox/pkg/stun"
	"github.com/gosuri/uilive"
)

const (
	DEFAULT_ADDR = "0.0.0.0:60001"
	SPACES       = 30
)

// Prints connected peers to the stun server
func PrintConnectedPeers(stun *stun.Stun) {
	writer := uilive.New()
	writer.Start()
	peers, _ := stun.GetConnectedPeers()
	previouseLenght := len(peers)

	for {
		peers, _ := stun.GetConnectedPeers()
		if previouseLenght != len(peers) {
			out := []string{fmt.Sprintf("Peer%sAddress", strings.Repeat(" ", SPACES-len("Peer"))), ""}
			for _, peerInfo := range peers {
				sp := SPACES - len(peerInfo.Peername)
				out = append(out, fmt.Sprintf("%s%s%s", peerInfo.Peername, strings.Repeat(" ", sp), peerInfo.Addr))
			}
			fmt.Fprintln(writer, strings.Join(out, "\n"))
			previouseLenght = len(peers)
		}
	}
}

func main() {
	server, err := stun.NewStun(DEFAULT_ADDR, stun.NewMemoryPeerConnectionStore(), stun.DefaultStunOptions())
	if err != nil {
		fmt.Println(err)
	}

	// Prints connected peers to the network
	go PrintConnectedPeers(server)

	// starts server
	server.Serve()
}
