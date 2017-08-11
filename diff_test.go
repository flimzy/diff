package diff

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestResultString(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var r *Result
		expected := ""
		if result := r.String(); result != expected {
			t.Errorf("Unexpected result: %s", result)
		}
	})
	t.Run("diff", func(t *testing.T) {
		expected := "foo"
		r := &Result{diff: expected}
		if result := r.String(); result != expected {
			t.Errorf("Unexpected result: %s", result)
		}
	})
}

func TestSliceDiff(t *testing.T) {
	tests := []struct {
		name             string
		expected, actual []string
		result           string
	}{
		{
			name:     "equal",
			expected: []string{"foo"},
			actual:   []string{"foo"},
		},
		{
			name:     "different",
			expected: []string{"foo"},
			actual:   []string{"bar"},
			result:   "--- expected\n+++ actual\n@@ -1 +1 @@\n-foo+bar",
		},
	}
	for _, test := range tests {
		result := sliceDiff(test.expected, test.actual)
		var resultText string
		if result != nil {
			resultText = result.String()
		}
		if resultText != test.result {
			t.Errorf("Unexpected result:\n%s\n", resultText)
		}
	}
}

func TestTextSlices(t *testing.T) {
	tests := []struct {
		name             string
		expected, actual []string
		result           string
	}{
		{
			name:     "equal",
			expected: []string{"foo", "bar"},
			actual:   []string{"foo", "bar"},
		},
		{
			name:     "different",
			expected: []string{"foo", "bar"},
			actual:   []string{"bar", "bar"},
			result:   "--- expected\n+++ actual\n@@ -1,2 +1,2 @@\n-foo\n bar\n+bar\n",
		},
	}
	for _, test := range tests {
		result := TextSlices(test.expected, test.actual)
		var resultText string
		if result != nil {
			resultText = result.String()
		}
		if resultText != test.result {
			t.Errorf("Unexpected result:\n%s\n", resultText)
		}
	}
}

func TestText(t *testing.T) {
	tests := []struct {
		name             string
		expected, actual string
		result           string
	}{
		{
			name:     "equal",
			expected: "foo\nbar\n",
			actual:   "foo\nbar\n",
		},
		{
			name:     "different",
			expected: "foo\nbar",
			actual:   "bar\nbar",
			result:   "--- expected\n+++ actual\n@@ -1,2 +1,2 @@\n-foo\n bar\n+bar\n",
		},
	}
	for _, test := range tests {
		result := Text(test.expected, test.actual)
		var resultText string
		if result != nil {
			resultText = result.String()
		}
		if resultText != test.result {
			t.Errorf("Unexpected result:\n%s\n", resultText)
		}
	}
}

func TestIsJSON(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		isJSON bool
		result string
	}{
		{
			name:   "io.Reader",
			input:  strings.NewReader("foo"),
			isJSON: true,
			result: "foo",
		},
		{
			name:   "byte slice",
			input:  []byte("foo"),
			isJSON: true,
			result: "foo",
		},
		{
			name:   "json.RawMessage",
			input:  json.RawMessage("foo"),
			isJSON: true,
			result: "foo",
		},
		{
			name:   "string",
			input:  "foo",
			isJSON: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isJSON, result := isJSON(test.input)
			if isJSON != test.isJSON {
				t.Errorf("Unexpected result: %t", isJSON)
			}
			if string(result) != test.result {
				t.Errorf("Unexpected result: %s", string(result))
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "byte slice",
			input:    []byte(`"foo"`),
			expected: `"foo"`,
		},
		{
			name:     "string",
			input:    "foo",
			expected: `"foo"`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := marshal(test.input)
			if string(result) != test.expected {
				t.Errorf("Unexpected result: %s", string(result))
			}
		})
	}
}

func TestAsJSON(t *testing.T) {
	tests := []struct {
		name             string
		expected, actual interface{}
		result           string
	}{
		{
			name:     "equal",
			expected: []string{"foo", "bar"},
			actual:   []string{"foo", "bar"},
		},
		{
			name:     "different",
			expected: []string{"foo", "bar"},
			actual:   []string{"bar", "bar"},
			result: `--- expected
+++ actual
@@ -1,4 +1,4 @@
 [
-    "foo",
+    "bar",
     "bar"
 ]
`,
		},
	}
	for _, test := range tests {
		result := AsJSON(test.expected, test.actual)
		var resultText string
		if result != nil {
			resultText = result.String()
		}
		if resultText != test.result {
			t.Errorf("Unexpected result:\n%s\n", resultText)
		}
	}
}

func TestJSON(t *testing.T) {
	tests := []struct {
		name             string
		expected, actual string
		result           string
	}{
		{
			name:     "equal",
			expected: `["foo","bar"]`,
			actual:   `["foo","bar"]`,
		},
		{
			name:     "different",
			expected: `["foo","bar"]`,
			actual:   `["bar","bar"]`,
			result: `--- expected
+++ actual
@@ -1,4 +1,4 @@
 [
-    "foo",
+    "bar",
     "bar"
 ]
`,
		},
	}
	for _, test := range tests {
		result := JSON([]byte(test.expected), []byte(test.actual))
		var resultText string
		if result != nil {
			resultText = result.String()
		}
		if resultText != test.result {
			t.Errorf("Unexpected result:\n%s\n", resultText)
		}
	}
}

func TestInterface(t *testing.T) {
	tests := []struct {
		name             string
		expected, actual interface{}
		result           string
	}{
		{
			name:     "equal",
			expected: []string{"foo", "bar"},
			actual:   []string{"foo", "bar"},
		},
		{
			name:     "different",
			expected: []string{"foo", "bar"},
			actual:   []string{"bar", "bar"},
			result: `--- expected
+++ actual
@@ -1,4 +1,4 @@
 ([]string) (len=2) {
-  (string) (len=3) "foo",
+  (string) (len=3) "bar",
   (string) (len=3) "bar"
 }
`,
		},
	}
	for _, test := range tests {
		result := Interface(test.expected, test.actual)
		var resultText string
		if result != nil {
			resultText = result.String()
		}
		if resultText != test.result {
			t.Errorf("Unexpected result:\n%s\n", resultText)
		}
	}
}
