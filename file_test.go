package diff

import (
	"io/ioutil"
	"testing"
)

func TestFile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
		err      string
	}{
		{
			name: "not found",
			path: "./not_found",
			err:  "open ./not_found: no such file or directory",
		},
		{
			name:     "found",
			path:     "testdata/test.txt",
			expected: "Test Content\n",
		},
	}
	for _, test := range tests {
		content, err := ioutil.ReadAll(&File{Path: test.path})
		var msg string
		if err != nil {
			msg = err.Error()
		}
		if msg != test.err {
			t.Errorf("Unexpected error: %s\n", msg)
		}
		if test.expected != string(content) {
			t.Errorf("Unexpected content: %s\n", string(content))
		}
	}
}
