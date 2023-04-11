package integration_test

type EAVServiceIntTestSuite struct {
	IntegrationTestSuite
}

func (ts *EAVServiceIntTestSuite) TestSomething() {
	ts.Assert().Equal(1, 2)
}
