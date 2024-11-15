package main

import (
	"errors"
	"fmt"
	"strings"
)

// define type aliases for the verb, path, and version
type RequestLine struct {
	Verb    string
	Path    string
	Version string
}

type RequestParser interface {
	// parse the HTTP verb, path, and version from the request
	parseRequestLine(request string) (RequestLine, error)

	// parse any parameters included in the request
	parseRequestParameters(request string) (map[string]string, error)

	// parse the headers from the request
	parseHeaders(request string) (map[string]string, error)

	// parse the body from the request
	parseBody(request string) (string, error)
}

type RequestParserImpl struct{}

func (rp *RequestParserImpl) parseRequestLine(request string) (*RequestLine, error) {
	requestLine := strings.Split(request, "\r\n")[0]
	requestParts := strings.Split(requestLine, " ")

	if len(requestParts) != 3 {
		return &RequestLine{}, errors.New("invalid request")
	}

	// validate the verb
	verb := requestParts[0]
	if verb != "GET" && verb != "POST" && verb != "PUT" && verb != "DELETE" {
		return nil, errors.New(
			"invalid http verb",
		)
	}

	// todo: validate the path
	path := requestParts[1]

	// validate the version
	version := requestParts[2]
	if version != "HTTP/1.1" {
		return nil, errors.New("invalid http version")
	}

	return &RequestLine{
		Verb:    verb,
		Path:    path,
		Version: version,
	}, nil
}

func (rp *RequestParserImpl) parseRequestParameters(request string) (map[string]string, error) {
	requestLine, err := rp.parseRequestLine(request)
	if err != nil {
		fmt.Println("Error parsing request line")
		return nil, errors.New("invalid request")
	}
	path := requestLine.Path

	// parse the path for all request parameters
	// @todo this. parse a string like "/path?param1=value1&param2=value2"
	parts := strings.Split(path, "?")
	if len(parts) > 2 {
		return nil, errors.New("invalid input format")
	}

	if len(parts) == 1 {
		// no query string
		return make(map[string]string), nil
	}

	queryString := parts[1]
	// split the query string into key-value pairs
	// and store into a map
	pairs := strings.Split(queryString, "&")
	parameters := make(map[string]string)

	for _, pair := range pairs {
		// should look like key=value
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			return nil, errors.New("invalid query string format")
		}
		parameters[kv[0]] = kv[1]
	}

	return parameters, nil
}

func (rp *RequestParserImpl) parseHeaders(request string) (map[string]string, error) {
	headers := make(map[string]string)
	parts := strings.Split(request, "\r\n")
	if len(parts) < 2 {
		// no headers
		return headers, nil
	}
	for _, line := range parts[1:] {
		if line == "" {
			break
		}
		header := strings.Split(line, ": ")
		if len(header) != 2 {
			return nil, errors.New("invalid header format")
		}
		headers[header[0]] = header[1]
	}
	return headers, nil
}

func (rp *RequestParserImpl) parseBody(request string) (string, error) {
	lines := strings.Split(request, "\r\n")

	// find the empty line that separates the headers from the body
	for i, line := range lines {
		if line == "" {
			body := strings.Join(lines[i+1:], "\r\n")
			return body, nil
		}
	}
	// no empty line was found hence no headers and the body starts immediately
	if len(lines) > 1 {
		body := strings.Join(lines[1:], "\r\n")
		return body, nil
	}

	return "", errors.New("no body found")
}
