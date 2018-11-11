package diff

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestCheckDir(t *testing.T) {
	tests := []struct {
		name     string
		dir      func(*testing.T) string
		expected map[string]string
		err      string
	}{
		{
			name: "empty dir",
			dir: func(t *testing.T) string {
				d, err := ioutil.TempDir("", "empty dir")
				if err != nil {
					t.Fatal(err)
				}
				return d
			},
			expected: map[string]string{},
		},
		{
			name: "files only",
			dir: func(t *testing.T) string {
				d, err := ioutil.TempDir("", "empty dir")
				if err != nil {
					t.Fatal(err)
				}
				if e := ioutil.WriteFile(d+"/foo", []byte("foo"), 0777); e != nil {
					t.Fatal(e)
				}
				return d
			},
			expected: map[string]string{
				"foo": "acbd18db4cc2f85cedef654fccc4a4d8",
			},
		},
		{
			name: "recursive",
			dir: func(t *testing.T) string {
				d, err := ioutil.TempDir("", "empty dir")
				if err != nil {
					t.Fatal(err)
				}
				if e := ioutil.WriteFile(d+"/foo", []byte("foo"), 0777); e != nil {
					t.Fatal(e)
				}
				if e := os.Mkdir(d+"/bar", 0777); e != nil {
					t.Fatal(e)
				}
				if e := ioutil.WriteFile(d+"/bar/baz", []byte("baz"), 0777); e != nil {
					t.Fatal(e)
				}
				return d
			},
			expected: map[string]string{
				"bar/":    "<dir>",
				"foo":     "acbd18db4cc2f85cedef654fccc4a4d8",
				"bar/baz": "73feffa4b7f6bb68e44cf984c85f6e88",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir := test.dir(t)
			defer func() {
				_ = os.RemoveAll(dir)
			}()
			result, err := checkDir(dir)
			var msg string
			if err != nil {
				msg = err.Error()
			}
			if msg != test.err {
				t.Errorf("Unexpected error: %s\n", msg)
			}
			if d := Interface(test.expected, result); d != nil {
				t.Error(d)
			}
		})
	}
}
