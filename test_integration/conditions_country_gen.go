// Code generated by badctl v0.0.0, DO NOT EDIT.
package integrationtests

import badorm "github.com/ditrit/badaas/badorm"

func CountryNameCondition(v string) badorm.WhereCondition[Country] {
	return badorm.WhereCondition[Country]{
		Field: "name",
		Value: v,
	}
}