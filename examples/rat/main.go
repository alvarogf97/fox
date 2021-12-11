package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alvarogf97/fox/pkg/p2p"
)

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

// handle command
func handle(command string, writer *p2p.P2PWriter) {
	// handle command execution here ...
	result := "Result! :D"
	writer.Write("COMMAND_EXEC", fmt.Sprintf("Command execution result is: %s", result))
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

	// Handle connections
	for {
		message, err := peer.Listen()
		if err != nil {
			log.Fatal(err)
		}

		command := message.Message
		writer, err := peer.Connect(message.Peername)
		if err != nil {
			log.Println("Requested peer ", message.Peername, " is now offline")
			continue
		}

		go handle(command, writer)
	}
}
