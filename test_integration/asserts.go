package integrationtests

import (
	"reflect"
	"sort"

	"github.com/ditrit/badaas/persistence/models"
	"github.com/stretchr/testify/suite"
	"gotest.tools/assert"
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

func EqualEntityList(ts *suite.Suite, expected, actual []*models.Entity) {
	ts.Len(actual, len(expected))

	sort.SliceStable(expected, func(i, j int) bool {
		return expected[i].ID.String() < expected[j].ID.String()
	})

	sort.SliceStable(actual, func(i, j int) bool {
		return actual[i].ID.String() < actual[j].ID.String()
	})

	for i := range actual {
		EqualEntity(ts, expected[i], actual[i])
	}
}

func EqualEntity(ts *suite.Suite, expected, actual *models.Entity) {
	assert.DeepEqual(ts.T(), expected, actual)
	EqualList(
		ts,
		expected.Fields,
		actual.Fields,
	)
}
