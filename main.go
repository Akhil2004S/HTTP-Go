package main

import (
	"bufio"
	"fmt"
	"httProtocol/parser"
	"httProtocol/router"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}

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
	var startLineData *parser.StartLine
	var isStartLineParsed bool
	var headerMap parser.HeaderData
	var headerDone bool

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
				return
			}
			fmt.Println("The start line data in struct is:", startLineData)
		} else if !headerDone {
			if data == "\r\n" {
				headerDone = true
				messageBody := parser.ParseBody(dataReader, headerMap)
				fmt.Println("Getting the Start line data:", startLineData)
				fmt.Println("Getting the header:", headerMap)
				fmt.Println("The message body is:", messageBody)
				router.Route(startLineData, headerMap, messageBody)
			} else {
				headerMap = parser.ParseHeader(data, &header)
			}
		}
	}
}
