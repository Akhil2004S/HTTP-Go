package parser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

// bool is whether the line is parsed. String is the response and int is the status code
func ParseStartLine(data string) (*StartLine, error) {
	startLine := &StartLine{}

	var parsedString []string
	// var responseString string
	var method string
	var path string
	var version string

	parsedString = strings.Split(data, " ")

	if len(parsedString) != 3 {
		err := errors.New("Invalid start line")
		return startLine, err
	}
	method = parsedString[0]
	path = parsedString[1]
	version = strings.TrimSuffix(parsedString[2], "\r\n")

	if method == "" || path == "" || version == "" {
		err := errors.New("Invalid start line")
		return startLine, err
	}
	startLine.Method = method
	startLine.Path = path
	startLine.Version = version
	return startLine, nil
}

func ParseHeader(data string, header *map[string]string) *map[string]string {
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

func ParseBody(dataReader *bufio.Reader, headerMap *map[string]string) string {
	var data []byte
	val, ok := (*headerMap)["Content-Length"]
	val = strings.TrimPrefix(val, " ")
	if ok {
		contentLength, err := strconv.Atoi(val)
		if err != nil {
			log.Fatal("Invalid content length", err)
		}
		data = make([]byte, contentLength)
		// msgBody, err := dataReader.Peek(contentLength)
		// fmt.Println("The alleged message body:", string(msgBody))
		_, err = io.ReadFull(dataReader, []byte(data))
		if err != nil {
			fmt.Println(err, "IO READ FULL")
		}
		if err != nil {
			log.Fatal("God save me")
		}
	}
	return string(data)
}
