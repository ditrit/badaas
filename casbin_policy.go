package main

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

type Data struct {
	Name    string
	Owner   string
	Domain  string
	Creator string
}

func main() {
	model, err := model.NewModelFromFile("model.conf")
	if err != nil {
		panic(err)
	}
	adapter := fileadapter.NewAdapter("policy.csv")
	policy := adapter
	e, err := casbin.NewEnforcer(model, policy)
	if err != nil {
		panic(err)
	}
	//e.AddNamedMatchingFunc("g2", "KeyMatch2", util.KeyMatch)
	//e.AddNamedDomainMatchingFunc("g3", "KeyMatch2", util.KeyMatch)
	tests := []struct {
		Sub    string
		Obj    Data
		Act    string
		Dom    string
		Expect bool
	}{
		{
			Sub: "Pierre",
			Obj: Data{
				Name:    "test",
				Owner:   "Xavier",
				Domain:  "Orness",
				Creator: "Pierre",
			},
			Act:    "write",
			Dom:    "Orness",
			Expect: true,
		},
		{
			Sub: "Pierre",
			Obj: Data{
				Name:    "test2",
				Owner:   "Pierre",
				Domain:  "Orness",
				Creator: "Xavier",
			},
			Act:    "write",
			Dom:    "Orness",
			Expect: false,
		},
		{
			Sub: "Pierre",
			Obj: Data{
				Name:    "test2",
				Owner:   "Pierre",
				Domain:  "Orness",
				Creator: "Xavier",
			},
			Act:    "read",
			Dom:    "Orness",
			Expect: true,
		},
		{
			Sub: "Pierre",
			Obj: Data{
				Name:    "test3",
				Owner:   "Xavier",
				Domain:  "Ditrit",
				Creator: "Pierre",
			},
			Act:    "write",
			Dom:    "Orness",
			Expect: false,
		},
		{
			Sub: "Vincent",
			Obj: Data{
				Name:    "data2",
				Owner:   "Xavier",
				Domain:  "Orness",
				Creator: "Pierre",
			},
			Act:    "exec",
			Dom:    "domain1.sub2",
			Expect: true,
		},
		{
			Sub: "Vincent",
			Obj: Data{
				Name:    "data3",
				Owner:   "Xavier",
				Domain:  "Orness",
				Creator: "Pierre",
			},
			Act:    "exec",
			Dom:    "domain1.sub2",
			Expect: false,
		},
		{
			Sub: "Vincent",
			Obj: Data{
				Name:    "data2",
				Owner:   "Xavier",
				Domain:  "Orness",
				Creator: "Pierre",
			},
			Act:    "exec",
			Dom:    "domain1.sub1",
			Expect: false,
		},
		{
			Sub: "devaut",
			Obj: Data{
				Name:    "data1",
				Owner:   "Xavier",
				Domain:  "Orness",
				Creator: "Pierre",
			},
			Act:    "exec",
			Dom:    "domain1.sub1",
			Expect: true,
		},
		{
			Sub: "devaut",
			Obj: Data{
				Name:    "data2",
				Owner:   "Xavier",
				Domain:  "Orness",
				Creator: "Pierre",
			},
			Act:    "exec",
			Dom:    "domain1.sub2",
			Expect: true,
		},
		{
			Sub: "super",
			Obj: Data{
				Name:    "data1",
				Owner:   "Xavier",
				Domain:  "Orness",
				Creator: "Pierre",
			},
			Act:    "exec",
			Dom:    "Orness",
			Expect: true,
		},
		{
			Sub: "new_group",
			Obj: Data{
				Name:    "Ditrit",
				Owner:   "Xavier",
				Domain:  "Orness",
				Creator: "Pierre",
			},
			Act:    "play",
			Dom:    "domain1",
			Expect: true,
		},
		{
			Sub: "Vincent",
			Obj: Data{
				Name:    "Ditrit",
				Owner:   "Xavier",
				Domain:  "Orness",
				Creator: "Pierre",
			},
			Act:    "play",
			Dom:    "domain1",
			Expect: true,
		},
	}

	for _, test := range tests {
		allowed, err := e.Enforce(test.Dom, test.Sub, test.Obj, test.Act)
		if err != nil {
			panic(err)
		}

		if allowed == test.Expect {
			fmt.Printf("Test Succeeded for user %s to resource %s with permission %s in domain %s \n", test.Sub, test.Obj, test.Act, test.Dom)
		} else {
			fmt.Printf("Test Failed for user %s to resource %s with permission %s in domain %s \n", test.Sub, test.Obj, test.Act, test.Dom)
		}
	}
}
