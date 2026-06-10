package parser

type StartLine struct {
	Method  string
	Path    string
	Version string
}

type HeaderData *map[string]string
