package tcp

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
	"tcp_to_http/internal/http"
	"time"
)

func handleConnection(conn net.Conn, timeoutInt int) error {
	log.Printf("New Connection: %v", conn.RemoteAddr())
	defer conn.Close()

	var builder strings.Builder
	timeout := time.Duration(timeoutInt) * time.Second

	conn.SetDeadline(time.Now().Add(timeout))

	reader := bufio.NewReader(conn)
	requestConfig := http.CreateNewHttpRequest()
	for {
		message, err := reader.ReadBytes('\n')

		if strings.TrimSpace(string(message)) == "" {
			break
		}

		if err != nil {
			break
		}

		builder.Write(message)
		requestConfig.ParseStringPerLinePopulateStruct(string(message))
	}

	bodyBuffer := make([]byte, requestConfig.ContentLength)
	_, err := io.ReadFull(reader, bodyBuffer)
	if err != nil {
		log.Printf("Failed to read the body of the request for conn: %v, err: %v", conn.RemoteAddr(), err)
	}
	log.Printf("BODY VALUE COLLECTED: %v", string(bodyBuffer))
	connectionError := conn.Close()
	log.Printf("Closing connection: %v", conn.RemoteAddr())
	if connectionError != nil {
		log.Printf("Failed to close connection: %v", conn.RemoteAddr())
	}
	return nil
}
