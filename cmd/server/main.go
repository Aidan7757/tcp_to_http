package main

import (
	"tcp_to_http/internal/http"
)

func main() {
	address := "localhost:8080"
	http.CreateNewRouter(address)
}
