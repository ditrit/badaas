package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	_ "github.com/lib/pq"
)

func Authorized(w http.ResponseWriter, r *http.Request) {
	sub := r.URL.Query().Get("sub")
	dom := r.URL.Query().Get("dom")
	obj := r.URL.Query().Get("obj")
	act := r.URL.Query().Get("act")

	fmt.Printf("%v, %v, %v, %v ", sub, dom, obj, act)

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

	// If you want to load the policy from a Cockroach database, uncomment the above code and comment the line below
	e = GetEnforcerFromFiles()

	e.AddNamedDomainMatchingFunc("g", "KeyMatch2", util.KeyMatch2)

	http.HandleFunc("/authorized", Authorized)

	fmt.Println("Ready")

	log.Fatal(http.ListenAndServe("127.0.0.1:5556", nil))
}

func GetEnforcerFromFiles() *casbin.Enforcer {
	e, err := casbin.NewEnforcer("model.conf", "policy.csv")
	if err != nil {
		panic(err)
	}

	return e
}
