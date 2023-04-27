Feature: Saving and querying objects in the database using the EAV Model
  Scenario: Objects can be created with CreateObject
    When I request "/objects/profile" with method "POST" with json
      | key         | value       | type    |
      | displayName | Jean Dupont | string  |
      | yearOfBirth | 1997        | integer |
    Then status code is "201"
    And response field "type" is "profile"
    And response field "attrs.displayName" is "Jean Dupont"
    And response field "attrs.yearOfBirth" is "1997"

  Scenario: Created object can be queried individually
    Given a "profile" object exists with attributes
      | key         | value       | type    |
      | displayName | Jean Dupont | string  |
      | yearOfBirth | 1997        | integer |
    When I query a "profile" with the object id
    Then response field "type" is "profile"
    And response field "attrs.displayName" is "Jean Dupont"
    And response field "attrs.yearOfBirth" is "1997"

  Scenario: Created objects can be queried together
    Given a "profile" object exists with attributes
      | key         | value       | type    |
      | displayName | Jean Dupont | string  |
      | yearOfBirth | 1997        | integer |
    And a "profile" object exists with attributes
      | key         | value         | type    |
      | displayName | Pierre Martin | string  |
      | yearOfBirth | 2001          | integer |
    When I query all "profile" objects
    Then there are "2" "profile" objects
    And there is a "profile" object with attributes
      | key         | value       | type   |
      | displayName | Jean Dupont | string |
      | yearOfBirth | 1997        | float  |
    And there is a "profile" object with attributes
      | key         | value         | type   |
      | displayName | Pierre Martin | string |
      | yearOfBirth | 2001          | float  |

  Scenario: Created objects can be queried by a property
    Given a "profile" object exists with attributes
      | key         | value       | type    |
      | displayName | Jean Dupont | string  |
      | yearOfBirth | 1997        | integer |
    And a "profile" object exists with attributes
      | key         | value         | type    |
      | displayName | Pierre Martin | string  |
      | yearOfBirth | 2001          | integer |
    When I query all "profile" objects with conditions
      | key         | value         | type    |
      | yearOfBirth | 2001          | integer |
    Then there are "1" "profile" objects
    And there is a "profile" object with attributes
      | key         | value         | type   |
      | displayName | Pierre Martin | string |
      | yearOfBirth | 2001          | float  |

  Scenario: Created objects can be queried by multiple properties
    Given a "profile" object exists with attributes
      | key         | value       | type    |
      | displayName | Jean Dupont | string  |
      | yearOfBirth | 1997        | integer |
    And a "profile" object exists with attributes
      | key         | value         | type    |
      | displayName | Pierre Martin | string  |
      | yearOfBirth | 2001          | integer |
    And a "profile" object exists with attributes
      | key         | value           | type    |
      | displayName | Gabriel Bernard | string  |
      | yearOfBirth | 2001            | integer |
    When I query all "profile" objects with conditions
      | key         | value         | type    |
      | displayName | Pierre Martin | string  |
      | yearOfBirth | 2001          | integer |
    Then there are "1" "profile" objects
    And there is a "profile" object with attributes
      | key         | value         | type   |
      | displayName | Pierre Martin | string |
      | yearOfBirth | 2001          | float  |

  Scenario: Created objects can be queried doing joins
    Given a "user" object exists with attributes
      | key  | value | type    |
      | name | user1 | string  |
    And a "profile" object exists with property "userID" related to last object and properties
      | key         | value       | type    |
      | displayName | Jean Dupont | string  |
      | yearOfBirth | 1997        | integer |
    And a "user" object exists with attributes
      | key  | value | type    |
      | name | user2 | string  |
    And a "profile" object exists with property "userID" related to last object and properties
      | key         | value         | type    |
      | displayName | Pierre Martin | string  |
      | yearOfBirth | 2001          | integer |
    When I query all "profile" objects with conditions
      | key         | value             | type   |
      | displayName | Jean Dupont       | string |
      | userID      | {"name": "user1"} | json   |
    Then there are "1" "profile" objects
    And there is a "profile" object with attributes
      | key         | value       | type   |
      | displayName | Jean Dupont | string |
      | yearOfBirth | 1997        | float  |

  Scenario: Created objects can be deleted
    Given a "profile" object exists with attributes
      | key         | value       | type    |
      | displayName | Jean Dupont | string  |
      | yearOfBirth | 1997        | integer |
    When I delete a "profile" with the object id
    And I query all "profile" objects
    Then there are "0" "profile" objects

  Scenario: Created objects can be modified
    Given a "profile" object exists with attributes
      | key         | value       | type    |
      | displayName | Jean Dupont | string  |
      | yearOfBirth | 1997        | integer |
    When I modify a "profile" with attributes
      | key         | value       | type    |
      | yearOfBirth | 1998        | integer |
    And I query a "profile" with the object id
    Then response field "type" is "profile"
    And response field "attrs.displayName" is "Jean Dupont"
    And response field "attrs.yearOfBirth" is "1998"