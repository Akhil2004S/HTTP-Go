package router

import (
	"fmt"
	"httProtocol/handler"
	"httProtocol/parser"
	"httProtocol/types"
	"net"
)

var allowedMethods = map[string]bool{
	"GET":    true,
	"POST":   true,
	"PUT":    true,
	"PATCH":  true,
	"DELETE": true,
}

var allowedRoutes = make(types.AllowedPaths)

func Route(conn net.Conn, startLineData *parser.StartLine, headerData *parser.HeaderData, body string) bool {
	fmt.Println("In router")
	request := types.Request{Method: startLineData.Method, Route: startLineData.Path}
	_, validMethod := allowedMethods[startLineData.Method]
	_, validPath := allowedRoutes[request]

	if startLineData.Version != "HTTP/1.1" {
		// err := errors.New("Invalid path. Version ain't available")
		errCode := 501
		handler.HandleError(conn, errCode)
	}

	if validMethod && validPath {
		// Call the appropriate handler
		handler.HandleRequest(conn, startLineData, headerData, body, allowedRoutes)
		return true
	} else if !validMethod {
		// Handle error where the handler is sent the error and the error code
		// err := errors.New("Invalid method. Cannot process")
		errCode := 400
		handler.HandleError(conn, errCode)
	} else if !validPath {
		// Handle error where the handler is sent the error and the error code
		// Sending error here cause it cannot be routed to any appropriate handler
		// err := errors.New("Invalid path. Path ain't here")
		errCode := 404
		handler.HandleError(conn, errCode)
	}
	return true
}

func RegisterRoute(method string, route string, routerFunction func() types.Response) {
	request := types.Request{
		Method: method,
		Route:  route,
	}
	allowedRoutes[request] = routerFunction
}
