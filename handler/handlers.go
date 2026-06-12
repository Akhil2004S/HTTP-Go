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

func HandleRequest(conn net.Conn, requestData types.Request, routes types.AllowedPaths) {
	fmt.Println("In handle request. Method:", requestData.StartLine.Method)
	request := parser.StartLine{Method: requestData.StartLine.Method, Path: requestData.StartLine.Path, Version: requestData.StartLine.Version}
	response := routes[request](requestData)
	fmt.Println(response)
	network.WriteToConn(conn, response)
}
