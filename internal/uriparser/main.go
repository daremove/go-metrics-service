package uriparser

import "strings"

type URIParser struct {
	fields map[string]string
}

func New(URI string, path string) *URIParser {
	parser := URIParser{map[string]string{}}
	splitURI := strings.Split(URI, "/")

	for i, val := range strings.Split(path, "/") {
		if i >= len(splitURI) {
			parser.fields[val] = ""
		} else {
			parser.fields[val] = splitURI[i]
		}
	}

	return &parser
}

func (r *URIParser) GetPathValue(pathName string) (string, bool) {
	val := r.fields[pathName]

	return val, val != ""
}
