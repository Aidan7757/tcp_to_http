package main

import (
	"tcp_to_http/internal/http"
)

func queryCallBackFunc(body *any) any {

	responseBody := http.QueryResponseBody{}
	return responseBody
}

func main() {

	queryRoute := http.Route{RouterFunc: queryCallBackFunc, RouteType: http.QueryRoute}

	address := "localhost:8080"
	router := http.CreateNewRouter(address)
	router.RegisterNewRoute("/query", http.QUERY, queryRoute)

	router.CreateAndRunListener("tcp", 20)
}
