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
		name           string
		updateMode     bool
		expected       *File
		actual         string
		diff           *Result
		finalExpected  string
		expectedResult string
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
			diff:          &Result{diff: "xxx"},
			actual:        "testing",
			finalExpected: "testing",
		},
		{
			name:           "Create failure",
			updateMode:     true,
			expected:       &File{Path: dir + "/foo/create.txt"},
			diff:           &Result{diff: "xxx"},
			actual:         "testing more",
			expectedResult: "Update failed: open " + dir + "/foo/create.txt: no such file or directory",
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
				diff:          &Result{diff: "xxx"},
				actual:        "testing update",
				finalExpected: "testing update",
			}
		}(),
		{
			name:          "No difference",
			updateMode:    true,
			expected:      &File{Path: "testdata/test.txt"},
			diff:          nil,
			actual:        "bogus text that should not be written",
			finalExpected: "Test Content",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := update(test.updateMode, test.expected, test.actual, test.diff)
			if d := Interface(test.expectedResult, result.String()); d != nil {
				t.Errorf("Unexpected result:\n%s\n", d)
			}
			if result != nil {
				return
			}
			if d := Text(test.finalExpected, &File{Path: test.expected.Path}); d != nil {
				t.Errorf("Unexpected file contents:\n%s\n", d)
			}
		})
	}
}
