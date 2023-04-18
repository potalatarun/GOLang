package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var err error

func homePage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "index.html")
}
func signUp(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "signup.html")
	}
	displayname := req.FormValue("displayname")
	email := req.FormValue("email")
	password := req.FormValue("password")

	var user string
	err := db.QueryRow("SELECT email FROM users WHERE email=?", email).Scan(&user)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, "Server error, ubnable to create your account.", 500)
			return
		}

		_, err = db.Exec("INSERT INTO users(displayname,email,password) VALUES(?,?,?)", displayname, email, hashedPassword)
		if err != nil {
			http.Error(res, "Server error, ubnable to create your account.", 500)
			return
		}
		res.Write([]byte("User created!"))
		return
	case err != nil:
		http.Error(res, "Server error, unable to create your account.", 500)
		return
	default:
		http.Redirect(res, req, "/", 301)
	}
}

func login(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "login.html")
		return
	}

	email := req.FormValue("email")
	password := req.FormValue("password")

	var databasedisplayname string
	var databaseemail string
	var databasepassword string

	err := db.QueryRow("SELECT displayname,email,password FROM users WHERE email=?", email).Scan(&databasedisplayname, &databaseemail, &databasepassword)

	if err != nil {
		http.Redirect(res, req, "/login", 301)
		return
	}
	fmt.Fprintf(res, "this is password %s", password)
	fmt.Fprintf(res, "this is hashpassowrd %s", databasepassword)
	err = bcrypt.CompareHashAndPassword([]byte(databasepassword), []byte(password))
	if err != nil {
		// fmt.Fprintf(res, "error 2")
		http.Redirect(res, req, "/login", 301)
	}

	res.Write([]byte("Hello" + databasedisplayname))
}
func main() {
	// creating the database
	db, err = sql.Open("mysql", "root:Tharun@123@(127.0.0.1:3306)/user")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("/", homePage)
	http.HandleFunc("/signup", signUp)
	http.HandleFunc("/login", login)
	http.ListenAndServe(":8080", nil)
}
