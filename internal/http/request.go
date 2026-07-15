package http

import (
	"log"
	"strconv"
	"strings"
)

type HttpMethod string

const (
	SET    HttpMethod = "SET"
	DELETE HttpMethod = "DELETE"
	QUERY  HttpMethod = "QUERY"
	POST   HttpMethod = "POST"
)

var VALID_HTTP_METHODS = []HttpMethod{
	POST,
	DELETE,
	SET,
	QUERY,
}

type HttpRequest struct {
	Connection        string
	Method            HttpMethod
	Path              string
	ContentLength     int
	Host              string
	AcceptEncoding    []string
	Accept            string
	ContentType       string
	Body              string
	AdditionalHeaders map[string]string
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

func CreateNewHttpRequest() *HttpRequest {
	HttpRequest := HttpRequest{}
	HttpRequest.AdditionalHeaders = make(map[string]string)
	return &HttpRequest
}

func (request *HttpRequest) checkHttpMethodLineAndPopulate(line string) bool {
	for _, method := range VALID_HTTP_METHODS {
		if strings.HasPrefix(line, string(method)) {
			suffix, _ := strings.CutPrefix(line, string(method))
			suffix = strings.TrimSpace(suffix)
			firstSlash := strings.Index(suffix, "/")
			firstSpace := strings.Index(suffix, " ")
			path := suffix[firstSlash : firstSpace+1]
			request.Path = strings.TrimSpace(path)
			request.Method = method
			return true
		}
	}
	return false
}

func (request *HttpRequest) parseHttpFieldsAndPopulateRequestConfig(prefix string, value string) {
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

func (request *HttpRequest) ParseStringPerLinePopulateStruct(line string) {
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
