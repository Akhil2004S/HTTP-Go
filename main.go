package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
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
	var isStartLineParsed bool
	var headerMap *map[string]string
	var headerDone bool
	var messageParsed bool
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
			response, statusCode := handleStartLine(data)
			if statusCode != 200 {
				conn.Write([]byte(response))
				conn.Close()
				return
			}
		} else if !headerDone {
			if data == "\r\n" {
				headerDone = true

			} else {
				headerMap = handleHeader(data, &header)
			}
			fmt.Println(isStartLineParsed, headerDone, messageParsed)
		} else if !messageParsed {
			fmt.Println("Hi")
			val, ok := (*headerMap)["Content-Length"]
			fmt.Println("Hi")
			if ok {
				fmt.Println("The length of the content is:", val)
				contentLength, err := strconv.Atoi(val)
				if err != nil {
					log.Fatal("Invalid content length", err)
				}
				msgBody, err := dataReader.Peek(contentLength)
				numer, err := io.ReadFull(conn, []byte(data))
				if err != nil {
					fmt.Println(err, "IO READ FULL")
				}
				fmt.Println(numer)
				if err != nil {
					log.Fatal("God save me")
				}
				fmt.Println("The alleged message body:", string(msgBody))
				messageParsed = true
			}
		}
	}
}

// bool is whether the line is parsed. String is the response and int is the status code
func handleStartLine(data string) (string, int) {
	allowedMethods := map[string]bool{
		"GET":    true,
		"POST":   true,
		"PUT":    true,
		"PATCH":  true,
		"DELETE": true,
	}

	allowedPaths := map[string]bool{
		"/home":     true,
		"/test":     true,
		"/willWork": true,
	}
	var parsedString []string
	var responseString string
	var code int

	parsedString = strings.Split(data, " ")

	if len(parsedString) != 3 {
		responseString = "HTTP/1.1 400 Bad Request\r\n\r\n"
		code = 400
		return responseString, code
	}

	if _, ok := allowedMethods[parsedString[0]]; !ok {
		log.Println("Invalid Method:", parsedString[2])
		responseString = "HTTP/1.1 501 Unsupported \r\n\r\n"
		code = 501
	} else if _, ok := allowedPaths[parsedString[1]]; !ok {
		log.Println("Wrong PATH:", parsedString[1])
		responseString = "HTTP/1.1 404 Not Found\r\n\r\n"
		code = 404
	} else if strings.TrimSuffix(parsedString[2], "\r\n") != "HTTP/1.1" {
		log.Println("Wrong VERSION:", parsedString[2])
		responseString = "HTTP/1.1 505 Unsupported Version\r\n\r\n"
		code = 505
	} else {
		responseString = "HTTP/1.1 200 OK\r\n\r\n"
		code = 200
	}
	return responseString, code
}

func handleHeader(data string, header *map[string]string) *map[string]string {
	var parsedHeader []string
	parsedHeader = strings.Split(data, ":")

	for data != "\r\n" {
		if len(parsedHeader) == 1 {
			break
		} else {
			(*header)[parsedHeader[0]] = strings.TrimSuffix(parsedHeader[1], "\r\n")
			break
		}
	}
	return header
}
