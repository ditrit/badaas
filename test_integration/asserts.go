package integrationtests

import (
	"log"
	"sort"

	"github.com/ditrit/badaas/persistence/models"
	"github.com/stretchr/testify/suite"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func EqualList[T any](ts *suite.Suite, expectedList, actualList []T) {
	expectedLen := len(expectedList)
	equalLen := ts.Len(actualList, expectedLen)

	if equalLen {
		for i := 0; i < expectedLen; i++ {
			j := 0
			for ; j < expectedLen; j++ {
				if is.DeepEqual(
					actualList[j],
					expectedList[i],
				)().Success() {
					break
				}
			}
			if j == expectedLen {
				// TODO mejorar esto
				ts.Fail("Lists not equal", "element %v not in list %v", expectedList[i], actualList)
				for _, element := range actualList {
					log.Println(element)
				}
			}
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
