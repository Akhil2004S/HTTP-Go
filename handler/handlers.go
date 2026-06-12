package handler

import (
	"fmt"
	"httProtocol/network"
	"httProtocol/parser"
	"httProtocol/types"
	"net"
)

func HandleError(conn net.Conn, errCode int) {
	switch errCode {
	case 400:
		fmt.Println("400")
		response := "HTTP/1.1 400 Bad Request\r\n\r\n"
		conn.Write([]byte(response))
		conn.Close()
	case 404:
		response := "HTTP/1.1 404 Not Found\r\n\r\n"
		conn.Write([]byte(response))
		conn.Close()
	case 501:
		response := "HTTP/1.1 501 Not Implemented. Invalid version.\r\n\r\n"
		conn.Write([]byte(response))
		conn.Close()
	}
}

func HandleRequest(conn net.Conn, startLineData *parser.StartLine, header *parser.HeaderData, body string, routes types.AllowedPaths) {
	fmt.Println("In handle request. Method:", startLineData.Method)
	request := types.Request{Method: startLineData.Method, Route: startLineData.Path}
	response := routes[request]()
	fmt.Println(response)
	network.WriteToConn(conn, response)
}

// func handleGet(path string) (string, error) {
// 	var response string
// 	switch path {
// 	case "/home":
// 		body := "body is home"

// 		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length:%d\r\nContent-Type:text/plain\r\n\r\n%s\r\n", len(body), body)
// 	}
// 	return response, nil
// }
