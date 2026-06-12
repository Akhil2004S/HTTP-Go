package types

import "httProtocol/parser"

type Request struct {
	StartLine parser.StartLine
	Headers   parser.HeaderData
	Message   string
}

type Response struct {
	StatusCode int
	Message    string
}

type AllowedPaths map[parser.StartLine]func(Request) Response
