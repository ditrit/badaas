package integrationtests

import (
	"reflect"

	"github.com/stretchr/testify/suite"
	is "gotest.tools/assert/cmp"
)

func EqualList(ts *suite.Suite, expected, actual any) {
	v := reflect.ValueOf(expected)
	v2 := reflect.ValueOf(actual)

	ts.Len(actual, v.Len())

	for i := 0; i < v.Len(); i++ {
		j := 0
		for ; j < v.Len(); j++ {
			if is.DeepEqual(v2.Index(j).Interface(), v.Index(i).Interface())().Success() {
				break
			}
		}
		if j == v.Len() {
			ts.Fail("element %v not in list %v", v.Index(i).Interface(), actual)
		}
	}
}
