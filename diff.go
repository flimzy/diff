// Package diff provides some convenience functions for comparing text in various
// forms. It's primary use case is in automated testing.
package diff

import (
	"strings"

	"github.com/pmezard/go-difflib/difflib"
)

// sliceDiff expects two slices of \n-terminated strings to compare.
func sliceDiff(expected, actual []string) string {
	udiff := difflib.UnifiedDiff{
		A:        expected,
		FromFile: "expected",
		B:        actual,
		ToFile:   "actual",
		Context:  2,
	}
	diff, err := difflib.GetUnifiedDiffString(udiff)
	if err != nil {
		panic(err)
	}
	return diff

}

// TextSlices compares two slices of text, treating each element as a line of
// text. Newlines are added to each element,if they are found to be missing.
func TextSlices(expected, actual []string) (diff string) {
	e := make([]string, len(expected))
	a := make([]string, len(actual))
	for i, str := range expected {
		e[i] = strings.TrimRight(str, "\n") + "\n"
	}
	for i, str := range actual {
		a[i] = strings.TrimRight(str, "\n") + "\n"
	}
	return sliceDiff(e, a)
}

// Text compares two strings, line-by-line, for differences. If the slices are
// identical, the return value will be the empty string.
func Text(expected, actual string) (diff string) {
	return sliceDiff(
		strings.SplitAfter(expected, "\n"),
		strings.SplitAfter(actual, "\n"),
	)
}
