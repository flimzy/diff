// Package diff provides some convenience functions for comparing text in various
// forms. It's primary use case is in automated testing.
package diff

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"strings"

	"github.com/davecgh/go-spew/spew"
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

func isJSON(i interface{}) (bool, []byte) {
	if r, ok := i.(io.Reader); ok {
		buf := &bytes.Buffer{}
		buf.ReadFrom(r)
		return true, buf.Bytes()
	}
	switch t := i.(type) {
	case []byte:
		return true, t
	case json.RawMessage:
		return true, t
	}
	return false, nil
}

func marshal(i interface{}) []byte {
	if isJ, buf := isJSON(i); isJ {
		var x interface{}
		_ = json.Unmarshal(buf, &x)
		i = x
	}
	j, _ := json.MarshalIndent(i, "", "    ")
	return j
}

// AsJSON marshals two objects as JSON, then compares the output. Marshaling
// errors are ignored. If an input object is an io.Reader, it is treated as
// a JSON stream. If it is a []byte or json.RawMessage, it is treated as raw
// JSON. Any raw JSON source is unmarshaled then remarshaled with indentation
// for normalization and comparison.
func AsJSON(expected, actual interface{}) (diff string) {
	expectedJSON := marshal(expected)
	actualJSON := marshal(actual)
	var e, a interface{}
	_ = json.Unmarshal(expectedJSON, &e)
	_ = json.Unmarshal(actualJSON, &a)
	if reflect.DeepEqual(e, a) {
		return ""
	}
	return Text(string(expectedJSON)+"\n", string(actualJSON)+"\n")
}

// JSON unmarshals two JSON strings, then calls AsJSON on them.
func JSON(expected, actual []byte) (diff string) {
	var expectedInterface, actualInterface interface{}
	_ = json.Unmarshal(expected, &expectedInterface)
	_ = json.Unmarshal(actual, &actualInterface)
	return AsJSON(expectedInterface, actualInterface)
}

// Interface compares two objects with reflect.DeepEqual, and if they differ,
// it returns a diff of the spew.Dump() outputs
func Interface(expected, actual interface{}) (diff string) {
	if reflect.DeepEqual(expected, actual) {
		return ""
	}
	scs := spew.ConfigState{
		Indent:                  "  ",
		DisableMethods:          true,
		SortKeys:                true,
		DisablePointerAddresses: true,
		DisableCapacities:       true,
	}
	expString := scs.Sdump(expected)
	actString := scs.Sdump(actual)
	return Text(expString, actString)
}
