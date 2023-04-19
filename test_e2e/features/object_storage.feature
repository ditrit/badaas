Feature: Saving and querying objects in the database using the EAV Model
  Scenario: Objects can be created with CreateObject
    When I request "/v1/objects/profile" with method "POST" with json
      | key         | value       | type    |
      | displayName | Jean Dupont | string  |
      | yearOfBirth | 1997        | integer |
    Then I expect status code is "200"
    And I expect response field "type" is "profile"
    And I expect response field "attrs.displayName" is "Jean Dupont"
    And I expect response field "attrs.yearOfBirth" is "1997"

  # Scenario: Created objects can be obtained with GetAll
    # When I request "/login" with method "POST" with json
    # Then I expect status code is "200"
    # And I expect response field "username" is "admin"
    # And I expect response field "email" is "admin-no-reply@badaas.com"