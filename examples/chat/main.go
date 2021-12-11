package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/alvarogf97/fox/pkg/p2p"
)

// Listen P2P messages
func listen(peer *p2p.Peer) {
	for {
		message, err := peer.Listen()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("\n", strings.TrimSpace(message.Peername), ": ", message.Message)
		fmt.Print("Input message: ")
	}
}

// Disconnect from stun server on ctrl^C
func onSigterm(peer *p2p.Peer, sigs chan os.Signal) {
	<-sigs
	if err := peer.Disconnect(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("\nDisconnected")
	defer peer.Close()
	os.Exit(0)
}

func main() {
	if len(os.Args) < 4 {
		log.Fatal("Usage: ", os.Args[0], " port peername stunaddr")
	}
	port := fmt.Sprintf(":%s", os.Args[1])
	peername := os.Args[2]
	stunaddr := os.Args[3]

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// creates the peer
	peer, err := p2p.NewPeer(peername, stunaddr, port, p2p.DefaultPeerOptions())
	if err != nil {
		log.Fatal(err)
	}
	go onSigterm(peer, sigs)

	// register peer and initialize it into the
	// network
	fmt.Println("Connecting peer to network... ")
	err = peer.Init()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected :D")

	// starts listening
	go listen(peer)

	// select a peer and starts chatting :D
	var writer *p2p.P2PWriter
	for writer == nil {
		fmt.Print("Enter peer name: ")
		peername := make([]byte, 2048)
		fmt.Scanln(&peername)

		writer, err = peer.Connect(string(peername))
		if err != nil {
			fmt.Println(err)
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	message := ""
	for {
		fmt.Print("Input message: ")
		if scanner.Scan() {
			message = scanner.Text()
		}
		if _, err := writer.Write("Chat", message); err != nil {
			log.Fatal(err)
		}
	}

}
