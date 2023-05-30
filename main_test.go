package main

import (
	"fmt"
	"testing"
)

type Data struct {
	Name    string
	Owner   string
	Domain  string
	Creator string
}

type CasbinTest struct {
	Dom    string
	Sub    string
	Obj    Data
	Act    string
	Expect bool
}

func TestCasbin(t *testing.T) {
	e = GetEnforcerFromFiles()
	test2Data := Data{
		Name:    "test2",
		Owner:   "Pierre",
		Domain:  "Orness",
		Creator: "Xavier",
	}
	data1 := Data{
		Name:    "data1",
		Owner:   "Xavier",
		Domain:  "Orness",
		Creator: "Pierre",
	}
	data2 := Data{
		Name:    "data2",
		Owner:   "Xavier",
		Domain:  "Orness",
		Creator: "Pierre",
	}
	data3 := Data{
		Name:    "data3",
		Owner:   "Xavier",
		Domain:  "Orness",
		Creator: "Pierre",
	}
	tests := []CasbinTest{
		{
			Dom: "Orness",
			Sub: "Pierre",
			Obj: Data{
				Name:    "test",
				Owner:   "Xavier",
				Domain:  "Orness",
				Creator: "Pierre",
			},
			Act:    "write",
			Expect: true,
		},
		{
			Dom:    "Orness",
			Sub:    "Pierre",
			Obj:    test2Data,
			Act:    "write",
			Expect: false,
		},
		{
			Dom:    "Orness",
			Sub:    "Pierre",
			Obj:    test2Data,
			Act:    "read",
			Expect: true,
		},
		{
			Dom: "Orness",
			Sub: "Pierre",
			Obj: Data{
				Name:    "test3",
				Owner:   "Xavier",
				Domain:  "Ditrit",
				Creator: "Pierre",
			},
			Act:    "write",
			Expect: false,
		},
		// ESTE no anda
		{
			Dom:    "domain1.sub2",
			Sub:    "Vincent",
			Obj:    data2,
			Act:    "exec",
			Expect: true,
		},
		{
			Dom:    "domain1.sub2",
			Sub:    "Vincent",
			Obj:    data3,
			Act:    "exec",
			Expect: false,
		},
		{
			Dom:    "domain1.sub1",
			Sub:    "Vincent",
			Obj:    data2,
			Act:    "exec",
			Expect: false,
		},
		// ESTE no anda
		{
			Dom:    "domain1.sub1",
			Sub:    "devaut",
			Obj:    data1,
			Act:    "exec",
			Expect: true,
		},
		// ESTE no anda
		{
			Dom:    "domain1.sub2",
			Sub:    "devaut",
			Obj:    data2,
			Act:    "exec",
			Expect: true,
		},
		// ESTE no anda
		{
			Dom:    "Orness",
			Sub:    "super",
			Obj:    data1,
			Act:    "exec",
			Expect: true,
		},
	}

	for _, test := range tests {
		allowed, err := e.Enforce(test.Dom, test.Sub, test.Obj, test.Act)
		if err != nil {
			t.Error("Error while executing test", err)
		} else {
			if allowed {
				fmt.Printf("Access granted for user %s to resource %s with permission %s in domain %s \n", test.Sub, test.Obj, test.Act, test.Dom)
			} else {
				fmt.Printf("Access denied for user %s to resource %s with permission %s in domain %s \n", test.Sub, test.Obj, test.Act, test.Dom)
			}

			if allowed != test.Expect {
				t.Errorf("Expected %t but got %t", test.Expect, allowed)
			}
		}
	}
}
