package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/gorilla/sessions"
)

// Struct holding user's credentials
type Credentials struct {
	Username string
	Password string
}

// Struct holding information about data
type Data struct {
	Name    string
	Owner   string
	Domain  string
	Creator string
}

// Define a map to store valid usernames and passwords (for demo purposes, maybe implementing hash later when
// session management is implemented the way it is in Badaas)
var validUsers = map[string]string{
	"john": "password1",
	"jane": "password2",
}

var (
	// Session store
	store = sessions.NewCookieStore([]byte("saqs0bxcrepSVajybhaQssej1p8nVWq4lP1mX37F8xk="))
)

// Define the main handler for the login page
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method == "POST" {
		r.ParseForm()

		// Getting credentials from the form
		username := r.Form.Get("username")
		password := r.Form.Get("password")

		// Checking if the provided credentials are valid
		if validUsers[username] == password {
			// Create a new session
			session, _ := store.Get(r, "session-name")
			// Set the username in the session
			session.Values["username"] = username
			// Save the session
			session.Save(r, w)

			// Redirecting the user to the domain selection page
			http.Redirect(w, r, "/domain.html", http.StatusFound)
			return
		}
	}

	// Ask the user to log in if the credentials are invalid (or if the method is not POST)
	tmpl := template.Must(template.ParseFiles("login.html"))
	tmpl.Execute(w, nil)
}

// Domain selection page
func domainHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the username from the session
	session, _ := store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok {
		// Redirect to login if the session is not found or username is not stored
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Check if the request method is POST
	if r.Method == "POST" {
		// Parse the form to get the selected domain and other parameters
		r.ParseForm()
		dom := r.Form.Get("domain")
		act := r.Form.Get("action")

		// Get the obj value from the form
		objValue := r.Form.Get("object")

		// Convert the obj value from string to Data struct
		var obj Data
		err := json.Unmarshal([]byte(objValue), &obj)
		fmt.Printf("Parsed object: %+v\n", obj)
		if err != nil {
			http.Error(w, "Invalid object data", http.StatusInternalServerError)
			return
		}
		fmt.Printf("Type of username: %T\n", username)
		fmt.Printf("Type of domain: %T\n", dom)
		fmt.Printf("Type of obj: %T\n", obj)
		fmt.Printf("Type of act: %T\n", act)
		fmt.Printf("username is %s, object is %s, domain is %s, action is %s, ", username, obj, dom, act)
		// test test test test yes yes yes yes shkribibim pam pam pam pam
		fmt.Println("Before enforcing policies")
		// Initialize the Casbin enforcer
		model, err := model.NewModelFromFile("model.conf")
		if err != nil {
			http.Error(w, "Internal Server Error loading model", http.StatusInternalServerError)
			return
		}
		adapter := fileadapter.NewAdapter("policy.csv")
		e, err := casbin.NewEnforcer(model, adapter)
		if err != nil {
			http.Error(w, "Internal Server Error creating enforcer", http.StatusInternalServerError)
			return
		}

		// to test the object please use the test.txt file

		// Enforce the Casbin policy
		allowed, err := e.Enforce(dom, username, obj, act)
		if err != nil {
			http.Error(w, "Internal Server Error enforcing policy", http.StatusInternalServerError)
			return
		}

		// test after work party I'm hungry and coding
		fmt.Println("After enforcing policies")

		// Store the enforcement result in the session
		session.Values["enforcement_result"] = allowed
		session.Save(r, w)

		// Redirect to the result page
		http.Redirect(w, r, "/result", http.StatusFound)
		return
	}

	// Render the domain.html page
	tmpl := template.Must(template.ParseFiles("domain.html"))
	tmpl.Execute(w, nil)
}

// Result page
func resultHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the username from the session
	session, _ := store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		// Redirect to login if the session is not found or username is not stored
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Retrieve the enforcement result from the session
	allowed, ok := session.Values["enforcement_result"].(bool)
	if !ok {
		// Redirect to domain selection if enforcement result is not found in session
		http.Redirect(w, r, "/domain.html", http.StatusFound)
		return
	}

	// Render the response template with the access status message
	tmpl := template.Must(template.ParseFiles("response.html"))
	if allowed {
		statusMessage := fmt.Sprintf("Access granted for user %s", username)
		tmpl.Execute(w, statusMessage)
	} else {
		statusMessage := fmt.Sprintf("Access denied for user %s", username)
		tmpl.Execute(w, statusMessage)
	}
}

func main() {
	// session configuration
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600, // 1 hour for the cookie to expire
		HttpOnly: true,
	}
	// Routing
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/domain.html", domainHandler)
	http.HandleFunc("/result", resultHandler)

	// Server Start
	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// Home page handler "/"
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Checking whether the user is logged in or not
	session, _ := store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if ok && username != "" {
		http.Redirect(w, r, "/domain.html", http.StatusFound)
		return
	}

	// Redirect to "/login" if the user is not logged in
	http.Redirect(w, r, "/login", http.StatusFound)
}
