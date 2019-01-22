package diff

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestUpdate(t *testing.T) {
	dir, err := ioutil.TempDir("", "diff-file")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	type updateTest struct {
		name          string
		updateMode    bool
		expected      *File
		actual        string
		finalExpected string
		err           string
	}
	tests := []updateTest{
		{
			name:          "Non-update mode",
			updateMode:    false,
			expected:      &File{Path: "testdata/test.txt"},
			actual:        "oink",
			finalExpected: "Test Content\n",
		},
		{
			name:          "Create file",
			updateMode:    true,
			expected:      &File{Path: dir + "/create.txt"},
			actual:        "testing",
			finalExpected: "testing",
		},
		{
			name:       "Create failure",
			updateMode: true,
			expected:   &File{Path: dir + "/foo/create.txt"},
			actual:     "testing more",
			err:        "open " + dir + "/foo/create.txt: no such file or directory",
		},
		func() updateTest {
			file := dir + "/update.txt"
			f, err := os.Create(file)
			if err != nil {
				t.Fatal(err)
			}
			f.WriteString("something else")
			f.Close()
			return updateTest{
				name:          "Ovwerrite",
				updateMode:    true,
				expected:      &File{Path: file},
				actual:        "testing update",
				finalExpected: "testing update",
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := update(test.updateMode, test.expected, test.actual)
			var msg string
			if err != nil {
				msg = err.Error()
			}
			if msg != test.err {
				t.Errorf("Unexpected error: %s\n", msg)
			}
			if err != nil {
				return
			}
			if d := Text(test.finalExpected, &File{Path: test.expected.Path}); d != nil {
				t.Error(d)
			}
		})
	}
}
