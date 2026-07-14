package http

import (
	"fmt"
	"log"
	"time"
)

type HttpResponse struct {
	StatusCode    int
	StatusString  string
	Method        HttpMethod
	Date          string
	ContentType   string
	ContentLength int
	Body          string
}

func (response *HttpResponse) serializeHttpResponseIntoString() []byte {
	formatString := "HTTP/1.1 %v %v\r\n" +
		"Date: %v\r\n" +
		"Content-Type: %v; charset=UTF-8\r\n" +
		"Content-Length: %v\r\n" +
		"\r\n" +
		"%v"

	formatWithPlacements := fmt.Sprintf(formatString, response.StatusCode, response.StatusString,
		response.Date, response.ContentType, response.ContentLength, response.Body)
	log.Printf("Format with placements: %v", formatWithPlacements)
	return []byte(formatWithPlacements)
}

func CreateNewHttpResponse() *HttpResponse {
	response := HttpResponse{}
	response.Date = time.Now().UTC().Format(time.RFC1123)
	response.ContentLength = 0
	response.Body = ""

	return &response
}
