package main

import (
	"tcp_to_http/internal/http"
)

// Fake query function - will be replaced with the DB calls once it is implemented
func queryCallBackFunc(body any, httpResponse *http.HttpResponse) (int, error) {
	fakeQueryValues := make(map[string]any)
	fakeQueryValues["test_key"] = "test_value"
	fakeResultValues := []bool{true}
	responseBody := http.QueryResponseBody{}

	responseBody.QueriedValues = fakeQueryValues
	responseBody.Result = fakeResultValues

	httpResponse.Body = responseBody
	httpResponse.StatusCode = 200
	return 200, nil
}

func main() {

	queryRoute := http.Route{RouterFunc: queryCallBackFunc, RouteType: http.QueryRoute}

	address := "localhost:8080"
	router := http.CreateNewRouter(address)
	router.RegisterNewRoute("/query", http.QUERY, queryRoute)

	router.CreateAndRunListener("tcp", 20)
}
