<p  align="center">

<img  width="300"  height="300"  src="https://i.imgur.com/XAGh4os.png">

</p>

  

# Fox in a nutshell



Fox is a lightweight framework that allows you to build a P2P connection that works without forwarding ports into the clients firewall. It works by using the `UDP hole punching` whose workflow is

  

<p  align="center">

<img  width="561"  height="281"  src="https://camo.githubusercontent.com/2bcc9106bf8b9f021e536549084e9b48e461400d/687474703a2f2f692e696d6775722e636f6d2f645a4e456870772e706e67">

</p>

  

Fox provides you both interfaces, one for build the Stun (Rendezvous) server and  another one for build Peers which whom estabilish P2P connections.

  

# How it works

  

There are many predefined actions for both, stun server and peers which indicates the actor how to handle them.

  

**Registration workflow**:

  

- Peer sends `STUN_ACTION_NEW` action to the Stun server

- Stun server responses Peer with `PEER_ACTION_NEW`

- Peer is now registered in the network and ready to read and connect to other peers

  

**Estabilish connection workflow**:

  

- Peer send `STUN_ACTION_GET` action to the Stun server with the peer he wants connect to

- Stun server responses Peer with `PEER_ACTION_GET`

- Peer receives the address of the requested peer and connection is estabilished

  

**Disconnect workflow**:

  

- Peer send `STUN_ACTION_DISCONNECT` action to the Stun server

- Stun server responses Peer with `PEER_ACTION_DISCONNECT`

- Peer is no longer registered in the network and stops listen for incomming connections

  

# How to use it

  

## Stun server

  

The following example shows you how to build a basic stun server

  

```go
package main

import (
	"fmt"
	"github.com/alvarogf97/fox/pkg/stun"
)

const (
	DEFAULT_ADDR = "0.0.0.0:60001"
)

func main() {
	server, err := stun.NewStun(
		DEFAULT_ADDR,
		stun.NewMemoryPeerConnectionStore(),
		stun.DefaultStunOptions(),
	)

	if err != nil {
		fmt.Println(err)
	}

	// starts server
	server.Serve()

}
```

  

If you need a more customizable stun server you may want to handle the connection in your own way:

  

```go
package main

import (
	"fmt"
	"log"
	"net"
	"github.com/alvarogf97/fox/pkg/stun"
)

const (
	DEFAULT_ADDR = "0.0.0.0:60001"
)

func handle(data []byte, addr *net.UDPAddr) {
	// handles stun server requests here
}

func main() {
	var buf [2048]byte
	server, err := stun.NewStun(
		DEFAULT_ADDR,
		stun.NewMemoryPeerConnectionStore(), stun.DefaultStunOptions(),
	)

	if err != nil {
		fmt.Println(err)

	}

	for {
		n, addr, err := server.ReadFromUDP(buf[0:])
		if err != nil {
			log.Println(err)
			continue
		}

		go handle(buf[:n], addr)

	}

}
```

  
Please take a look of the  **Stun server** example in the [examples](examples/) folder

  

## Peer

Peer is the main actor in the P2P network, it's easy to create a new peer and starts writing to other ones by using **fox**. The following code shows an example of how to build a basic Peer:
  


  

```go
package main

import (
	"fmt"
	"log"
	"github.com/alvarogf97/fox/pkg/p2p"
)

const (
	DEFAULT_ADDR      = "0.0.0.0:50001"
	DEFAULT_PEER      = "MyPeerName"
	DEFAULT_CONN_PEER = "OtherPeer"
	DEFAULT_STUN_ADDR = "localhost:60001"
)

// Listen P2P messages
func listen(peer *p2p.Peer) {
	for {
		message, err := peer.Listen()
		if err != nil {
			log.Fatal(err)
		}
		// Handle messages here ...
	}

}

func main() {

	// creates the peer
	peer, err := p2p.NewPeer(
		DEFAULT_PEER,
		DEFAULT_STUN_ADDR,
		DEFAULT_ADDR,
		p2p.DefaultPeerOptions(),
	)
	defer peer.Close()

	if err != nil {
		log.Fatal(err)
	}

	// register peer and initialize it into the
	// network
	fmt.Println("Connecting peer to network... ")

	err = peer.Init()
	if err != nil {
		log.Fatal(err)

	}

	fmt.Println("Connected :D")

	// starts listening messages
	go listen(peer)

	// Connect to a peer

	writer, err := peer.Connect(DEFAULT_CONN_PEER)
	if err != nil {
		fmt.Println(err)
	}

	// write some message
	if _, err := writer.Write("Chat", "Hello Friend"); err != nil {
		log.Fatal(err)

	}

}
```

Please take a look of the peer examples in the [examples](examples/) folder.
