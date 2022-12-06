package utils_test

import (
	"testing"

	"github.com/ditrit/badaas/services/eavservice/utils"
)

type jsonTest struct {
	input  []string
	output string
}

func TestBuildJsonFromString(t *testing.T) {
	testData := []jsonTest{
		{
			input:  []string{"\"voila\":1"},
			output: "{\"voila\":1}",
		},
		{
			input: []string{
				"\"voila\":1",
				"\"45\":\"231115421 351635 6351sd6351 354qsd35 4qs\"",
			},
			output: "{\"voila\":1,\"45\":\"231115421 351635 6351sd6351 354qsd35 4qs\"}",
		},
	}
	for _, td := range testData {
		result := utils.BuildJsonFromStrings(td.input)
		if result != td.output {
			t.Errorf("Expected %s, got %s", td.output, result)
		}
	}

}

func TestBuildJsonListFromString(t *testing.T) {
	testData := []jsonTest{
		{
			input:  []string{"{\"voila\":1}"},
			output: "[{\"voila\":1}]",
		},
		{
			input: []string{
				"{\"voila\":1}",
				"{\"45\":\"231115421 351635 6351sd6351 354qsd35 4qs\"}",
			},
			output: "[{\"voila\":1},{\"45\":\"231115421 351635 6351sd6351 354qsd35 4qs\"}]",
		},
	}
	for _, td := range testData {
		result := utils.BuildJsonListFromStrings(td.input)
		if result != td.output {
			t.Errorf("Expected %s, got %s", td.output, result)
		}
	}

}
