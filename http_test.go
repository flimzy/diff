package diff

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHTTPRequest(t *testing.T) {
	tests := []struct {
		name             string
		expected, actual interface{}
		result           string
	}{
		{
			name:     "two GET requests",
			expected: httptest.NewRequest("GET", "/", nil),
			actual:   httptest.NewRequest("GET", "/", nil),
		},
		{
			name:     "different URLs requests",
			expected: httptest.NewRequest("GET", "/", nil),
			actual:   httptest.NewRequest("GET", "/foo.html", nil),
			result:   "--- expected\n+++ actual\n@@ -1,3 +1,3 @@\n-GET / HTTP/1.1\r\n+GET /foo.html HTTP/1.1\r\n Host: example.com\r\n \r\n",
		},
		{
			name:     "nil request",
			expected: nil,
			actual:   httptest.NewRequest("GET", "/foo.html", nil),
			result:   "--- expected\n+++ actual\n@@ -1 +1,3 @@\n-\n+GET /foo.html HTTP/1.1\r\n+Host: example.com\r\n+\r\n",
		},
		{
			name: "string",
			expected: `GET / HTTP/1.1
Host: localhost:6005
User-Agent: curl/7.52.1
Accept: */*

`,
			actual: &http.Request{
				Method:     http.MethodGet,
				ProtoMajor: 1,
				ProtoMinor: 1,
				URL:        &url.URL{Host: "localhost:6005"},
				Header: http.Header{
					"Accept":     []string{"*/*"},
					"User-Agent": []string{"curl/7.52.1"},
				},
			},
		},
		{
			name: "byte slice",
			expected: []byte(`GET / HTTP/1.1
Host: localhost:6005
User-Agent: curl/7.52.1
Accept: */*

`),
			actual: &http.Request{
				Method:     http.MethodGet,
				ProtoMajor: 1,
				ProtoMinor: 1,
				URL:        &url.URL{Host: "localhost:6005"},
				Header: http.Header{
					"Accept":     []string{"*/*"},
					"User-Agent": []string{"curl/7.52.1"},
				},
			},
		},
		{
			name:     "file",
			expected: &File{Path: "testdata/request.raw"},
			actual: &http.Request{
				Method:     http.MethodGet,
				ProtoMajor: 1,
				ProtoMinor: 1,
				URL:        &url.URL{Host: "localhost:6005"},
				Header: http.Header{
					"Accept":     []string{"*/*"},
					"User-Agent": []string{"curl/7.52.1"},
				},
			},
		},
		{
			name:     "nil",
			expected: nil,
			actual:   func() interface{} { return nil }(),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := HTTPRequest(test.expected, test.actual)
			var resultText string
			if result != nil {
				resultText = result.String()
			}
			if d := Text(test.result, resultText); d != nil {
				t.Error(d)
			}
		})
	}
}

func TestHTTPResponse(t *testing.T) {
	tests := []struct {
		name             string
		expected, actual interface{}
		result           string
	}{
		{
			name:     "two empty responses",
			expected: &http.Response{},
			actual:   &http.Response{},
		},
		{
			name: "Different headers",
			expected: &http.Response{
				Header: http.Header{"Foo": []string{"bar"}},
			},
			actual: &http.Response{
				Header: http.Header{"Foo": []string{"qux"}},
			},
			result: "--- expected\n+++ actual\n@@ -1,4 +1,4 @@\n HTTP/0.0 000 status code 0\r\n-Foo: bar\r\n+Foo: qux\r\n Content-Length: 0\r\n \r\n",
		},
		{
			name:     "nil response body",
			expected: nil,
			actual: &http.Response{
				Header: http.Header{"Foo": []string{"qux"}},
			},
			result: "--- expected\n+++ actual\n@@ -1 +1,4 @@\n-\n+HTTP/0.0 000 status code 0\r\n+Foo: qux\r\n+Content-Length: 0\r\n+\r\n",
		},
		{
			name: "two nil bodies",
			expected: &http.Response{
				StatusCode: 400,
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     http.Header{"Foo": []string{"qux"}},
			},
			actual: &http.Response{
				StatusCode: 400,
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     http.Header{"Foo": []string{"qux"}},
			},
		},
		{
			name:     "read from file",
			expected: &File{Path: "testdata/response.raw"},
			actual: &http.Response{
				StatusCode:    http.StatusOK,
				ProtoMajor:    1,
				ProtoMinor:    1,
				ContentLength: 11,
				Header: http.Header{
					"Content-Type": []string{"application/json"},
					"Date":         []string{"Tue, 22 Jan 2019 18:44:09 GMT"},
				},
				Body: ioutil.NopCloser(strings.NewReader(`{"ok":true}`)),
			},
		},
		{
			name: "string",
			expected: `HTTP/1.1 200 OK
Content-Length: 11
Content-Type: application/json
Date: Tue, 22 Jan 2019 18:44:09 GMT

{"ok":true}
`,
			actual: &http.Response{
				StatusCode:    http.StatusOK,
				ProtoMajor:    1,
				ProtoMinor:    1,
				ContentLength: 11,
				Header: http.Header{
					"Content-Type": []string{"application/json"},
					"Date":         []string{"Tue, 22 Jan 2019 18:44:09 GMT"},
				},
				Body: ioutil.NopCloser(strings.NewReader(`{"ok":true}`)),
			},
		},
		{
			name: "byte slice",
			expected: []byte(`HTTP/1.1 200 OK
Content-Length: 11
Content-Type: application/json
Date: Tue, 22 Jan 2019 18:44:09 GMT

{"ok":true}
`),
			actual: &http.Response{
				StatusCode:    http.StatusOK,
				ProtoMajor:    1,
				ProtoMinor:    1,
				ContentLength: 11,
				Header: http.Header{
					"Content-Type": []string{"application/json"},
					"Date":         []string{"Tue, 22 Jan 2019 18:44:09 GMT"},
				},
				Body: ioutil.NopCloser(strings.NewReader(`{"ok":true}`)),
			},
		},
		{
			name:     "unknown input type",
			expected: int(123),
			result:   "Failed to dump expected response: Unable to convert int to *http.Response",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := HTTPResponse(test.expected, test.actual)
			var resultText string
			if result != nil {
				resultText = result.String()
			}
			if d := Text(test.result, resultText); d != nil {
				t.Error(d)
			}
		})
	}
}
