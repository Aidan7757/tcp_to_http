package http

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type QueryResponseBody struct {
	QueriedValues map[string]any `json:"query_result"` // keys -> values
	Result        []bool         `json:"result"`
}

type SetResponseBody struct {
	Key    string `json:"key"`
	Value  any    `json:"value"`
	Result bool   `json:"result"`
}

type CreateTableResponseBody struct {
	TableName string                    `json:"table_name"`
	Columns   map[string]map[string]any `json:"columns"` // column -> key -> value
	Result    bool                      `json:"result"`
}

type ErrorResponseBody struct {
	Error string `json:"error"`
}

type HttpResponse struct {
	StatusCode    int
	StatusString  string
	Date          string
	ContentType   string
	ContentLength int
	Body          any
}

func (response *HttpResponse) serializeHttpResponseIntoString() []byte {

	formatString := "HTTP/1.1 %v %v\r\n" +
		"Date: %v\r\n" +
		"Content-Type: %v; charset=UTF-8\r\n" +
		"Content-Length: %v\r\n" +
		"\r\n" +
		"%v"

	if response.Body == nil {
		response.Body = ErrorResponseBody{Error: "Error occurred."}
	}

	jsonBytes, err := json.Marshal(response.Body)

	if err != nil {
		log.Println("Failed to convert response body to JSON.")
		return nil
	}

	response.ContentType = "application/json"
	response.ContentLength = len(jsonBytes)

	formatWithPlacements := fmt.Sprintf(formatString, response.StatusCode, response.StatusString,
		response.Date, response.ContentType, response.ContentLength, string(jsonBytes))

	bytes := []byte(formatWithPlacements)
	return bytes
}

func CreateNewHttpResponse() *HttpResponse {
	response := HttpResponse{}
	response.Date = time.Now().UTC().Format(time.RFC1123)
	return &response
}
