package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gorilla/mux"
	_ "github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var db *sql.DB
var err error

func signupPage(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	if req.Method != "POST" {
		http.ServeFile(res, req, "signup.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	var user string

	err := db.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, "<h1>Server error, unable to create your account.</h1>", 500)
			return
		}

		_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashedPassword)
		if err != nil {
			http.Error(res, "<h1>Server error, unable to create your account.</h1>", 500)
			return
		}

		res.Write([]byte("<h1>Successfully Registered, Go to <a href=\"/login\">Login Page</a>"))
	case err != nil:
		http.Error(res, "<h1>Server error, unable to create your account.</h1>", 500)
		return
	default:
		http.Redirect(res, req, "/", 301)
	}
}
// login
func loginPage(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	if req.Method != "POST" {
		http.ServeFile(res, req, "login.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	var databaseUsername string
	var databasePassword string

	err := db.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&databaseUsername, &databasePassword)

	//if err != nil {
	//	http.Redirect(res, req, "/login", 500)
	//	return
	//}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil{
		http.Redirect(res, req, "/login", 301)
		return
	}

	if err != bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password)){
		http.Redirect(res, req, "/UserNotFound", 500)
		return
	}

	res.Write([]byte("<h1>Hello " + databaseUsername + ", Welcome the website</h1>"))
}

func notF(res http.ResponseWriter, req *http.Request){
	res.Header().Set("Content-Type", "text/html")
	http.ServeFile(res, req,"logerr.html")
}

// home page
func homePage(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	http.ServeFile(res, req, "index.html")
}

// contact page
func contact(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "contact.html")
}
// faq page
func faq(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r,"faq.html")
}
// not found
func notFound(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Println("/404")
	fmt.Fprint(w, "<h1>Sorry, but we couldn't find the page you're looking for")
}

func main() {
	// database handler
	db, err = sql.Open("mysql", "root:@/go")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	// page handler
	http.HandleFunc("/signup", signupPage)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/", homePage)
	http.HandleFunc("/contact", contact)
	http.HandleFunc("/faq", faq)
	http.HandleFunc("/UserNotFound", notF)
	http.ListenAndServe(":8080", nil)
}
