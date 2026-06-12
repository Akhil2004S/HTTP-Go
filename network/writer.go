package network

import (
	"fmt"
	"httProtocol/types"
	"net"
)

func WriteToConn(conn net.Conn, returnResponse types.Response) {
	var response string
	message := returnResponse.Message
	statusCode := returnResponse.StatusCode
	contentLength := len(returnResponse.Message)
	if message != "" {
		response = fmt.Sprintf("HTTP/1.1 %d\r\nContent-Length:%d\r\nContent-Type:text/plain\r\n\r\n%s\r\n", statusCode, contentLength, message)
	} else {
		response = fmt.Sprintf("HTTP/1.1 %d\r\n", statusCode)
	}

	conn.Write([]byte(response))
	conn.Close()
}
