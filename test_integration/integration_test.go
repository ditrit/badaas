//go:build integration
// +build integration

package integration_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestAll(t *testing.T) {
	integrationTestSuite := IntegrationTestSuite{}

	suite.Run(t, &EAVServiceIntTestSuite{
		IntegrationTestSuite: integrationTestSuite,
	})
}
