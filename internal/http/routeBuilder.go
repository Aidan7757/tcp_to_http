package http

import (
	"fmt"
)

type RouteBuilder struct {
	router *Router
	path   string
	method string
	route  Route
}

func (rb *RouteBuilder) Path(name string) *RouteBuilder {
	currentIndex := len(rb.route.Params)
	rb.route.Params = append(rb.route.Params, ParamMetadata{
		Source: SourcePath,
		Name:   name,
		Type:   rb.route.Args[currentIndex],
	})
	return rb
}

func (rb *RouteBuilder) Query() *RouteBuilder {
	currentIndex := len(rb.route.Params)
	rb.route.Params = append(rb.route.Params, ParamMetadata{
		Source: SourceQuery,
		Type:   rb.route.Args[currentIndex],
	})
	return rb
}

func (rb *RouteBuilder) Body() *RouteBuilder {
	currentIndex := len(rb.route.Params)
	rb.route.Params = append(rb.route.Params, ParamMetadata{
		Source: SourceBody,
		Type:   rb.route.Args[currentIndex],
	})
	return rb
}

func (rb *RouteBuilder) Header() *RouteBuilder {
	currentIndex := len(rb.route.Params)
	rb.route.Params = append(rb.route.Params, ParamMetadata{
		Source: SourceHeader,
		Type:   rb.route.Args[currentIndex],
	})
	return rb
}

func (rb *RouteBuilder) Register() {
	rb.router.mu.Lock()
	defer rb.router.mu.Unlock()

	if len(rb.route.Params) != len(rb.route.Args) {
		panic(fmt.Sprintf(
			"Router error for %s %s: handler expects %d args, but builder configured %d",
			rb.method, rb.path, len(rb.route.Args), len(rb.route.Params),
		))
	}

	if rb.router.routerFunctionMaps[rb.path] == nil {
		rb.router.routerFunctionMaps[rb.path] = make(map[string]Route)
	}

	rb.router.routerFunctionMaps[rb.path][rb.method] = rb.route
}
