Feature: Saving and querying objects in the database using the BaDorm
  Scenario: Created objects can be queried by a property
    Given a sale object exists for product "1", code "1" and description "sale1"
    And a sale object exists for product "2", code "2" and description "sale2"
    When I query all sale objects with conditions
      | key  | value | type    |
      | code | 1     | integer |
    Then there are "1" "sale" objects
    And there is a sale object with attributes
      | key         | value | type   |
      | Code        | 1     | float  |
      | Description | sale1 | string |

  Scenario: Created objects can be queried by multiple properties
    Given a sale object exists for product "1", code "1" and description "sale1"
    And a sale object exists for product "2", code "2" and description "sale2"
    When I query all sale objects with conditions
      | key         | value | type    |
      | code        | 1     | integer |
      | description | sale1 | string  |
    Then there are "1" "sale" objects
    And there is a sale object with attributes
      | key         | value | type   |
      | Code        | 1     | float  |
      | Description | sale1 | string |

  Scenario: Created objects can be queried doing joins
    Given a sale object exists for product "1", code "1" and description "sale1"
    And a sale object exists for product "2", code "2" and description "sale2"
    When I query all sale objects with conditions
      | key     | value        | type    |
      | code    | 1            | integer |
      | Product | {"int": "1"} | json    |
    Then there are "1" "sale" objects
    And there is a sale object with attributes
      | key         | value | type   |
      | Code        | 1     | float  |
      | Description | sale1 | string |