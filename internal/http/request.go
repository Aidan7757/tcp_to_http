package http

import (
	"log"
	"strconv"
	"strings"
)

type HttpMethod string

const (
	GET     HttpMethod = "GET"
	SET     HttpMethod = "SET"
	POST    HttpMethod = "POST"
	DELETE  HttpMethod = "DELETE"
	PUT     HttpMethod = "PUT"
	PATCH   HttpMethod = "PATCH"
	HEAD    HttpMethod = "HEAD"
	OPTIONS HttpMethod = "OPTIONS"
)

type httpRequest struct {
	Connection        string
	Method            HttpMethod
	ContentLength     int
	Host              string
	AcceptEncoding    []string
	Accept            string
	ContentType       string
	Body              string
	AdditionalHeaders map[string]string
}

var VALID_HTTP_METHODS = []HttpMethod{
	GET,
	POST,
	PUT,
	DELETE,
	SET,
	PATCH,
	HEAD,
	OPTIONS,
}

// POST / HTTP/1.1
// Test: Hello
// Content-Type: text/plain
// User-Agent: PostmanRuntime/7.54.0
// Accept: */*
// Cache-Control: no-cache
// Postman-Token: 9220833d-6f7b-4fc6-b898-876fc2ce4530
// Host: localhost:8080
// Accept-Encoding: gzip, deflate, br
// Connection: keep-alive
// Content-Length: 10

func CreateNewHttpRequest() *httpRequest {
	httpRequest := httpRequest{}
	httpRequest.AdditionalHeaders = make(map[string]string)
	return &httpRequest
}

func (request *httpRequest) checkHttpMethodLineAndPopulate(line string) bool {
	for _, method := range VALID_HTTP_METHODS {
		if strings.HasPrefix(line, string(method)) {
			request.Method = method
			return true
		}
	}
	return false
}

func (request *httpRequest) parseHttpFieldsAndPopulateRequestConfig(prefix string, value string) {
	switch prefix {
	case "Accept":
		request.Accept = value
	case "Connection":
		request.Connection = value
	case "Content-Length":
		num, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("Unable to parse content length header for int value: %v", value)
			return
		}
		request.ContentLength = num
	case "Accept-Encoding":
		// split by , over strings, split then loop over
		encodings := strings.Split(value, ",")
		for _, value := range encodings {
			request.AcceptEncoding = append(request.AcceptEncoding, value)
		}
	case "Host":
		request.Host = value
	case "Content-Type":
		request.ContentType = value
	default:
		request.AdditionalHeaders[prefix] = value
	}

}

func (request *httpRequest) ParseStringPerLinePopulateStruct(line string) {
	line = strings.TrimSpace(line)
	methodResult := request.checkHttpMethodLineAndPopulate(line)

	if methodResult {
		return
	}

	before, after, result := strings.Cut(line, ":")
	before, after = strings.TrimSpace(before), strings.TrimSpace(after)
	request.parseHttpFieldsAndPopulateRequestConfig(before, after)
	if !result {
		return
	}
}
