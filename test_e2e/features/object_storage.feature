Feature: Saving and querying objects in the database using the EAV Model
  Scenario: Objects can be created with CreateObject
    When I request "/v1/objects/profile" with method "POST" with json
      | key         | value       | type    |
      | displayName | Jean Dupont | string  |
      | yearOfBirth | 1997        | integer |
    Then status code is "201"
    And response field "type" is "profile"
    And response field "attrs.displayName" is "Jean Dupont"
    And response field "attrs.yearOfBirth" is "1997"

  Scenario: Created object can be queried individually
    Given a "profile" object exists with properties
      | key         | value       | type    |
      | displayName | Jean Dupont | string  |
      | yearOfBirth | 1997        | integer |
    When I query a "profile" with the object id
    Then response field "type" is "profile"
    And response field "attrs.displayName" is "Jean Dupont"
    And response field "attrs.yearOfBirth" is "1997"

  Scenario: Created objects can be queried together
    Given a "profile" object exists with properties
      | key         | value       | type    |
      | displayName | Jean Dupont | string  |
      | yearOfBirth | 1997        | integer |
    And a "profile" object exists with properties
      | key         | value         | type    |
      | displayName | Pierre Martin | string  |
      | yearOfBirth | 2001          | integer |
    When I query all "profile" objects
    Then there are "2" "profile" objects
    And there is a "profile" object with properties
      | key         | value       | type   |
      | displayName | Jean Dupont | string |
      | yearOfBirth | 1997        | float  |
    And there is a "profile" object with properties
      | key         | value         | type   |
      | displayName | Pierre Martin | string |
      | yearOfBirth | 2001          | float  |

  Scenario: Created objects can be queried with parameters
    Given a "profile" object exists with properties
      | key         | value       | type    |
      | displayName | Jean Dupont | string  |
      | yearOfBirth | 1997        | integer |
    And a "profile" object exists with properties
      | key         | value         | type    |
      | displayName | Pierre Martin | string  |
      | yearOfBirth | 2001          | integer |
    When I query all "profile" objects with parameters
      | key         | value         | type   |
      | yearOfBirth | 2001          | string |
    Then there are "1" "profile" objects
    And there is a "profile" object with properties
      | key         | value         | type   |
      | displayName | Pierre Martin | string |
      | yearOfBirth | 2001          | float  |