// Code generated by badctl v0.0.0, DO NOT EDIT.
package integrationtests

import badorm "github.com/ditrit/badaas/badorm"

func PersonNameCondition(v string) badorm.WhereCondition[Person] {
	return badorm.WhereCondition[Person]{
		Field: "name",
		Value: v,
	}
}
