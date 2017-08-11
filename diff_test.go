package diff

import "testing"

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
