package main

import (
	"errors"
	"testing"
)

func TestParseRequestLine(t *testing.T) {

	testCases := []struct {
		verb    string
		path    string
		version string
		err     error
	}{
		{"GET", "/path/to/resource", "HTTP/1.1", nil},
		{"INVALID", "/path/to/resource", "HTTP/1.1", errors.New(
			"invalid http verb",
		)},
	}

	rp := RequestParserImpl{}

	for _, tc := range testCases {
		request := tc.verb + " " + tc.path + " " + tc.version + "\r\n"
		requestLine, err := rp.parseRequestLine(request)

		if tc.err != nil {
			if err.Error() != tc.err.Error() {
				t.Errorf("Expected error [%v], got [%v]", tc.err, err)
			}
			continue
		}

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if requestLine.Verb != tc.verb {
			t.Errorf("Expected %s, got %s", tc.verb, requestLine.Verb)
		}
		if requestLine.Path != tc.path {
			t.Errorf("Expected %s, got %s", tc.path, requestLine.Path)
		}
		if requestLine.Version != tc.version {
			t.Errorf("Expected %s, got %s", tc.version, requestLine.Version)
		}
	}
}

func TestParseRequestParameters(t *testing.T) {
	testCases := []struct {
		request        string
		expectedParams map[string]string
		ExpectedError  error
	}{
		{"GET /path/to/resource?param1=value1&param2=value2 HTTP/1.1\r\n", map[string]string{"param1": "value1", "param2": "value2"}, nil},
		{"GET /path/to/resource HTTP/1.1\r\n", map[string]string{}, nil},
		{"GET /path/to/resource?param1:value1 HTTP/1.1\r\n", nil, errors.New("invalid query string format")},
	}

	rp := RequestParserImpl{}

	for _, tc := range testCases {
		params, err := rp.parseRequestParameters(tc.request)

		if tc.ExpectedError != nil {
			if err.Error() != tc.ExpectedError.Error() {
				t.Errorf("Expected error [%v], got [%v]", tc.ExpectedError, err)
			}
			continue
		}

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(params) != len(tc.expectedParams) {
			t.Errorf("Expected %d parameters, got %d", len(tc.expectedParams), len(params))
		}
		for k, v := range tc.expectedParams {
			if params[k] != v {
				t.Errorf("Expected %s, got %s", v, params[k])
			}
		}
	}
}

func TestParseRequestHeaders(t *testing.T) {
	testCases := []struct {
		request         string
		expectedHeaders map[string]string
		ExpectedError   error
	}{
		{"GET /path/to/resource HTTP/1.1\r\nHeader1: value1\r\nHeader2: value2\r\n", map[string]string{"Header1": "value1", "Header2": "value2"}, nil},
		{"GET /path/to/resource HTTP/1.1", map[string]string{}, nil},
		{"GET /path/to/resource HTTP/1.1\r\nHeader1", map[string]string{}, errors.New("invalid header format")},
	}

	rp := RequestParserImpl{}
	for _, tc := range testCases {
		headers, err := rp.parseHeaders(tc.request)

		if tc.ExpectedError != nil {
			if err.Error() != tc.ExpectedError.Error() {
				t.Errorf("Expected error [%v], got [%v]", tc.ExpectedError, err)
			}
			continue
		}

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(headers) != len(tc.expectedHeaders) {
			t.Errorf("Expected %d headers, got %d", len(tc.expectedHeaders), len(headers))
		}
		for k, v := range tc.expectedHeaders {
			if headers[k] != v {
				t.Errorf("Expected %s, got %s", v, headers[k])
			}
		}
	}
}

func TestParseBody(t *testing.T) {
	testCases := []struct {
		request       string
		expectedBody  string
		ExpectedError error
	}{
		{"GET /path/to/resource HTTP/1.1\r\nHeader1: value1\r\nHeader2: value2\r\n\r\nBody", "Body", nil},
		{"GET /path/to/resource HTTP/1.1\r\nHeader1: value1\r\nHeader2: value2\r\n", "", nil},
	}

	rp := RequestParserImpl{}
	for _, tc := range testCases {
		body, err := rp.parseBody(tc.request)

		if tc.ExpectedError != nil {
			if err.Error() != tc.ExpectedError.Error() {
				t.Errorf("Expected error [%v], got [%v]", tc.ExpectedError, err)
			}
			continue
		}

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if body != tc.expectedBody {
			t.Errorf("Expected %s, got %s", tc.expectedBody, body)
		}
	}
}
