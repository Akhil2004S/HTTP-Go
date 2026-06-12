package types

type Request struct {
	Method string
	Route  string
}

type Response struct {
	StatusCode int
	Message    string
}

type AllowedPaths map[Request]func() Response
