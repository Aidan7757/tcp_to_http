package main

import (
	"log"
	"tcp_to_http/internal/http"
)

func queryCallBackFunc(body any, httpResponse *http.HttpResponse) (int, error) {
	fakeQueryValues := make(map[string]any)
	fakeQueryValues["test_key"] = "test_value"
	fakeResultValues := []bool{true}
	responseBody := http.QueryResponseBody{}

	responseBody.QueriedValues = fakeQueryValues
	responseBody.Result = fakeResultValues

	httpResponse.Body = responseBody
	httpResponse.StatusCode = 200
	log.Printf("Response Body HTTP: %+v", httpResponse)
	return 200, nil
}

func main() {

	queryRoute := http.Route{RouterFunc: queryCallBackFunc, RouteType: http.QueryRoute}

	address := "localhost:8080"
	router := http.CreateNewRouter(address)
	router.RegisterNewRoute("/query", http.QUERY, queryRoute)

	router.CreateAndRunListener("tcp", 20)
}
