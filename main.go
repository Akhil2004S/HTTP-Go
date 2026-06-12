package main

import (
	"bufio"
	"fmt"
	"httProtocol/parser"
	"httProtocol/router"
	"httProtocol/types"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}

	router.RegisterRoute("GET", "/home", handleHome)
	router.RegisterRoute("POST", "/homer", handleHomePOST)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go read(conn)
	}
}

func read(conn net.Conn) {
	dataReader := bufio.NewReader(conn)
	var startLineData parser.StartLine
	var isStartLineParsed bool
	var headerMap parser.HeaderData
	var headerDone bool
	var messageBody string
	// var requestRouted bool

	header := make(map[string]string)
	for {
		data, err := dataReader.ReadString('\n')
		fmt.Println(data)
		if err != nil {
			log.Println("Client disconnected.")
			conn.Close()
			return
		}

		if !isStartLineParsed {
			isStartLineParsed = true
			startLineData, err = parser.ParseStartLine(data)
			if err != nil {
				fmt.Println(err)
				conn.Close()
				return
			}
		} else if !headerDone {
			if data == "\r\n" {
				headerDone = true
				messageBody = parser.ParseBody(dataReader, headerMap)
				router.Route(conn, startLineData, headerMap, messageBody)
			} else {
				headerMap = parser.ParseHeader(data, &header)
			}
		}
	}
}

func handleHome(request types.Request) types.Response {
	response := types.Response{
		StatusCode: 200,
		Message:    "You've called home",
	}
	return response
}

func handleHomePOST(request types.Request) types.Response {
	fmt.Println("Post can do something with the obtained request data")
	fmt.Println(request.Headers)
	fmt.Println(request.Message)
	response := types.Response{
		StatusCode: 200,
		Message:    "Did something",
	}
	return response
}
