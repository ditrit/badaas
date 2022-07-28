package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/casbin/casbin/util"
	"github.com/casbin/casbin/v2"
	_ "github.com/lib/pq"
)

func Authorized(w http.ResponseWriter, r *http.Request) {
	sub := r.URL.Query().Get("sub")
	dom := r.URL.Query().Get("dom")
	obj := r.URL.Query().Get("obj")
	act := r.URL.Query().Get("act")

	if res, _ := e.Enforce(sub, dom, obj, act); res {
		response := "Access authorized"
		fmt.Println(response)
		w.Write([]byte(response))
	} else {
		response := "Access denied"
		fmt.Println(response)
		w.Write([]byte(response))
	}
}

var e *casbin.Enforcer

func main() {
	/*
		// Connect to the database first
		db, err := sql.Open("postgres", "user=root password=postgres host=localhost port=26257 sslmode=disable dbname=rbac")
		if err != nil {
			panic(err)
		}
		if err = db.Ping(); err != nil {
			panic(err)
		}

		a, err := sqladapter.NewAdapter(db, "postgres", "policy")

		if err != nil {
			panic(err)
		}

		e, err = casbin.NewEnforcer("rbac_with_domains_model.conf", a)
	*/

	// If you want to load the policy from a Cockroach database, decomment the above code and comment the two lines below
	err := *new(error)
	e, err = casbin.NewEnforcer("rbac_with_domains_model.conf", "rbac_with_domains_policy.csv")

	if err != nil {
		panic(err)
	}

	e.AddNamedDomainMatchingFunc("g", "KeyMatch2", util.KeyMatch2)

	http.HandleFunc("/authorized", Authorized)

	fmt.Println("Ready")

	log.Fatal(http.ListenAndServe("127.0.0.1:5556", nil))
}
