package http

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type RouteError struct {
	routeError   int
	errorMessage string
}

type Router struct {
	routerFunctionMaps map[string]map[string]func()
	mu                 sync.RWMutex
	address            string
}

func CreateNewRouter(newAddress string) *Router {
	router := Router{
		address: newAddress,
	}
	router.routerFunctionMaps = make(map[string]map[string]func())
	router.CreateAndRunListener("tcp", newAddress, 10)
	return &router
}

func (router *Router) verifyRouteExistence(request *HttpRequest) (error, int) {
	router.mu.RLock()
	if router.routerFunctionMaps[request.Path] == nil {
		return fmt.Errorf("Not Found"), 404
	}

	if router.routerFunctionMaps[request.Path][string(request.Method)] == nil {
		return fmt.Errorf("Method Not Allowed"), 405
	}

	return nil, -1
}

func (router *Router) RegisterNewRoute(path string, method HttpMethod, routerFunc func()) {
	router.mu.Lock()
	defer router.mu.Unlock()
	if router.routerFunctionMaps[path] == nil {
		router.routerFunctionMaps[path] = make(map[string]func())
	}

	router.routerFunctionMaps[path][string(method)] = routerFunc
}

func (router *Router) CreateAndRunListener(method string, address string, defaultTimeout int) error {
	listener, err := net.Listen(method, address)

	if err != nil {
		log.Fatalf("Error creating %v listener for %v address: %v\n", strings.ToUpper(method), address, err)
	}

	defer listener.Close()
	log.Printf("%v listener active on address: %v", strings.ToUpper(method), address)

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
	}
	return *httpResponse
}
