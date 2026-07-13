package main

import (
	"tcp_to_http/internal/tcp"
)

func main() {
	method := "tcp"
	address := "localhost:8080"
	defaultTimeout := 5
	tcp.CreateAndRunListener(method, address, defaultTimeout)
}
