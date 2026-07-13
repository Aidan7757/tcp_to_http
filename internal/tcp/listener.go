package tcp

import (
	"log"
	"net"
	"strings"
)

func CreateAndRunListener(method string, address string, defaultTimeout int) error {
	listener, err := net.Listen(method, address)

	if err != nil {
		log.Fatalf("Error creating %v listener for %v address: %v\n", strings.ToUpper(method), address, err)
	}

	defer listener.Close()
	log.Printf("%v listener active on address: %v", strings.ToUpper(method), address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go handleConnection(conn, defaultTimeout)
	}
}
