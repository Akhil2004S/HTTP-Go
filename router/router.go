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

func Route(conn net.Conn, requestLine parser.StartLine, headerData parser.HeaderData, body string) bool {
	fmt.Println("In router")
	_, validMethod := allowedMethods[requestLine.Method]
	_, validPath := allowedRoutes[requestLine]

	if requestLine.Version != "HTTP/1.1" {
		// err := errors.New("Invalid path. Version ain't available")
		errCode := 501
		handler.HandleError(conn, errCode)
	}

	if validMethod && validPath {
		// Call the appropriate handler
		request := types.Request{
			StartLine: requestLine,
			Headers:   headerData,
			Message:   body,
		}
		handler.HandleRequest(conn, request, allowedRoutes)
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

func RegisterRoute(method string, route string, routerFunction func(types.Request) types.Response) {
	request := parser.StartLine{
		Method:  method,
		Path:    route,
		Version: "HTTP/1.1",
	}
	allowedRoutes[request] = routerFunction
}
