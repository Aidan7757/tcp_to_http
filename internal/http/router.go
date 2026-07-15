package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type CreateTableRouteBody struct {
	TableName *string                    `json:"table_name"`
	Columns   *map[string]map[string]any `json:"columns"` // column -> key -> value
}

type QueryRouteBody struct {
	TableName *string `json:"table_name"`
	Column    *string `json:"column_name"`
	Key       *any    `json:"key"`
}

type SetRouteBody struct {
	TableName *string `json:"table_name"`
	Column    *string `json:"column"`
	Key       *any    `json:"key"`
	Value     *any    `json:"value"`
}

type RouteType int

const (
	CreateTableRoute RouteType = iota
	QueryRoute
	SetRoute
)

type Route struct {
	RouterFunc func(body *any) any
	RouteType  RouteType
}

type Router struct {
	routerFunctionMaps map[string]map[string]Route
	mu                 sync.RWMutex
	address            string
}

func CreateNewRouter(newAddress string) *Router {
	router := Router{
		address: newAddress,
	}
	router.routerFunctionMaps = make(map[string]map[string]Route)
	return &router
}

func (router *Router) verifyRouteExistence(request *HttpRequest) (error, int) {
	router.mu.RLock()

	methodMap, result := router.routerFunctionMaps[request.Path]

	if !result {
		return fmt.Errorf("Not Found"), 404
	}

	_, result = methodMap[string(request.Method)]

	if !result {
		return fmt.Errorf("Method Not Allowed"), 405
	}

	return nil, -1
}

func (router *Router) RegisterNewRoute(path string, method HttpMethod, route Route) error {
	router.mu.Lock()
	defer router.mu.Unlock()
	if router.routerFunctionMaps[path] == nil {
		router.routerFunctionMaps[path] = make(map[string]Route)
	}

	router.routerFunctionMaps[path][string(method)] = route
	return nil
}

func (router *Router) CreateAndRunListener(method string, defaultTimeout int) error {
	listener, err := net.Listen(method, router.address)

	if err != nil {
		log.Fatalf("Error creating %v listener for %v address: %v\n", strings.ToUpper(method), router.address, err)
	}

	defer listener.Close()
	log.Printf("%v listener active on address: %v", strings.ToUpper(method), router.address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go router.HandleConnection(conn, defaultTimeout)
	}
}

func (router *Router) Serve(httpRequest *HttpRequest) HttpResponse {
	httpResponse := CreateNewHttpResponse()
	err, errorResponseCode := router.verifyRouteExistence(httpRequest)
	if err != nil {
		httpResponse.StatusCode = errorResponseCode
		httpResponse.StatusString = err.Error()
		return *httpResponse
	}
	route := router.routerFunctionMaps[httpRequest.Path][string(httpRequest.Method)]

	log.Printf("ROUTE ROUTE TYPE: %+v", route.RouteType)
	if route.RouteType == SetRoute {
		var createTableRouteBody CreateTableRouteBody
		err := json.Unmarshal([]byte(httpRequest.Body), &createTableRouteBody)
		if err != nil {
			httpResponse.Body = ErrorResponseBody{Error: "Invalid request body."}
			httpResponse.StatusCode = 400
			httpResponse.StatusString = "Bad Request"
			return *httpResponse
		}

		if createTableRouteBody.Columns == nil {
			httpResponse.Body = ErrorResponseBody{Error: "No columns field included."}
			httpResponse.StatusCode = 400
			httpResponse.StatusString = "Bad Request"
			return *httpResponse
		}

		if createTableRouteBody.TableName == nil {
			httpResponse.Body = ErrorResponseBody{Error: "No table_name field included."}
			httpResponse.StatusCode = 400
			httpResponse.StatusString = "Bad Request"
			return *httpResponse
		}
		log.Printf("JSON Request Body: %+v", createTableRouteBody)
	}
	if route.RouteType == SetRoute {
		var setRouteBody SetRouteBody
		err := json.Unmarshal([]byte(httpRequest.Body), &setRouteBody)
		if err != nil {
			httpResponse.Body = ErrorResponseBody{Error: "Invalid request body."}
			httpResponse.StatusCode = 400
			httpResponse.StatusString = "Bad Request"
			return *httpResponse
		}

		if setRouteBody.Column == nil {
			httpResponse.Body = ErrorResponseBody{Error: "No column field included."}
			httpResponse.StatusCode = 400
			httpResponse.StatusString = "Bad Request"
			return *httpResponse
		}

		if setRouteBody.Key == nil {
			httpResponse.Body = ErrorResponseBody{Error: "No key field included."}
			httpResponse.StatusCode = 400
			httpResponse.StatusString = "Bad Request"
			return *httpResponse
		}

		if setRouteBody.TableName == nil {
			httpResponse.Body = ErrorResponseBody{Error: "No table_name field included."}
			httpResponse.StatusCode = 400
			httpResponse.StatusString = "Bad Request"
			return *httpResponse
		}
		log.Printf("JSON Request Body: %+v", setRouteBody)
	}
	if route.RouteType == QueryRoute {
		var queryRouteBody QueryRouteBody
		err := json.Unmarshal([]byte(httpRequest.Body), &queryRouteBody)
		if err != nil {
			httpResponse.Body = ErrorResponseBody{Error: "Invalid request body."}
			httpResponse.StatusCode = 400
			httpResponse.StatusString = "Bad Request"
			return *httpResponse
		}

		if queryRouteBody.Column == nil {
			httpResponse.Body = ErrorResponseBody{Error: "No column field included."}
			httpResponse.StatusCode = 400
			httpResponse.StatusString = "Bad Request"
			return *httpResponse
		}

		if queryRouteBody.Key == nil {
			httpResponse.Body = ErrorResponseBody{Error: "No key field included."}
			httpResponse.StatusCode = 400
			httpResponse.StatusString = "Bad Request"
			return *httpResponse
		}

		if queryRouteBody.TableName == nil {
			httpResponse.Body = ErrorResponseBody{Error: "No table_name field included."}
			httpResponse.StatusCode = 400
			httpResponse.StatusString = "Bad Request"
			return *httpResponse
		}
		log.Printf("JSON Request Body: %+v", queryRouteBody)
	}

	// need to implement the working route execution rather than just failing to find
	return *httpResponse
}
