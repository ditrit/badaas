package integration_test

import "github.com/ditrit/badaas/services/eavservice"

type EAVServiceIntTestSuite struct {
	IntegrationTestSuite
	eavService eavservice.EAVService
}

func NewEAVServiceIntTestSuite(
	ts *IntegrationTestSuite,
	eavService eavservice.EAVService,
) *EAVServiceIntTestSuite {
	return &EAVServiceIntTestSuite{
		IntegrationTestSuite: *ts,
		eavService:           eavService,
	}
}

func (ts *EAVServiceIntTestSuite) TestSomething() {
	ts.Assert().Equal(1, 1)
}
