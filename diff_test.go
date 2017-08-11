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
