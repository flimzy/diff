package diff

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

// HTTPRequest compares the metadata and bodies of the two HTTP requests, and
// returns the difference.
func HTTPRequest(expected, actual *http.Request) *Result {
	expDump, err := dumpRequest(expected)
	if err != nil {
		return &Result{err: fmt.Sprintf("Failed to dump expected request: %s", err)}
	}
	actDump, err := dumpRequest(actual)
	if err != nil {
		return &Result{err: fmt.Sprintf("Failed to dump actual request: %s", err)}
	}
	return Text(string(expDump), string(actDump))
}

func dumpRequest(req *http.Request) ([]byte, error) {
	if req == nil {
		return nil, nil
	}
	return httputil.DumpRequest(req, true)
}

// HTTPResponse compares the metadata and bodies of the two HTTP responses, and
// returns the difference.
func HTTPResponse(expected, actual *http.Response) *Result {
	expDump, err := dumpResponse(expected)
	if err != nil {
		return &Result{err: fmt.Sprintf("Failed to dump expected response: %s", err)}
	}
	actDump, err := dumpResponse(actual)
	if err != nil {
		return &Result{err: fmt.Sprintf("Failed to dump actual response: %s", err)}
	}
	return Text(string(expDump), string(actDump))
}

func dumpResponse(res *http.Response) ([]byte, error) {
	if res == nil {
		return nil, nil
	}
	return httputil.DumpResponse(res, true)
}
