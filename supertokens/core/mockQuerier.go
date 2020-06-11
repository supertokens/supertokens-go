package core

import (
	"flag"
	"net/http"
)

// MockedHTTPClient mocked http client
type MockedHTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var idMap = map[string]MockedHTTPClient{}

// AddMockedHTTPHandler used during testing
func AddMockedHTTPHandler(requestID string, handler MockedHTTPClient) {
	if flag.Lookup("test.v") != nil {
		idMap[requestID] = handler
	}
}

// GetMockedHTTPClient is used during testing
func GetMockedHTTPClient(requestID string) MockedHTTPClient {
	if flag.Lookup("test.v") == nil {
		return nil
	}
	value := idMap[requestID]
	if value == nil {
		return nil
	}
	return value
}

// ResetHTTPMocking sets idMap to an empty map
func ResetHTTPMocking() {
	idMap = map[string]MockedHTTPClient{}
}
