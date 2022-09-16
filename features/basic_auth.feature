Feature: Login as superadmin using the basic authentification

Scenario: Server should let us login and access a protected endpoint
  When I sign-in as "superadmin@badaas.test" with password "1234"
  Then I expect status code is "200"
