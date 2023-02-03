package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	sessions "github.com/kataras/go-sessions"
	//"golang.org/x/crypto/bcrypt"
	 //"os"
	 //"go get -u github.com/satori/go.uuid@master"
)

var db *sql.DB
var err error

type user struct {
	ID        int  `json :"id"`
	Username  string `json : "username"`
	FirstName string `json : "firstname"`
	LastName  string `json : "lastname"`
	Password  string `json : "password"`
	Email     string `json : "email"`
}

func connect_db() {
	db, err = sql.Open("mysql", "abhay:Admin@123@tcp(127.0.0.1:3306)/project")

	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
}

func routes() {
	http.HandleFunc("/register", register)
	//http.HandleFunc("/login", login)
	//http.HandleFunc("/logout", logout)
}

func main() {
	connect_db()
	routes()

	defer db.Close()

	fmt.Println("Server running on port :8080")
	http.ListenAndServe(":8080", nil)
}

type statusRes struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

func checkErr(w http.ResponseWriter, r *http.Request, err error) bool {
	if err != nil {

		fmt.Println(r.Host + r.URL.Path)

		http.Redirect(w, r, r.Host+r.URL.Path, 301)
		return false
	}

	return true
}

func QueryUser(username string) user {
	var users = user{}
	err = db.QueryRow(`
		SELECT id, 
		username, 
		first_name, 
		last_name, 
		password,
		email
		FROM users WHERE username=?
		`, username).
		Scan(
			&users.ID,
			&users.Username,
			&users.FirstName,
			&users.LastName,
			&users.Password,
			&users.Email,
		)
	return users
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "register.html")
		return
	}
}
	username := r.FormValue("username")
	email := r.FormValue("email")
	first_name := r.FormValue("first_name")
	last_name := r.FormValue("last_name")
	password := r.FormValue("password")
	fmt.Println(username)
	users := QueryUser(first_name)

	if (user{}) == users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		if len(hashedPassword) != 0 && checkErr(w, r, err) {
			stmt, err := db.Prepare("INSERT INTO user SET username=?, password=?, first_name=?, last_name=?, email=?")
			if err == nil {
				_, err := stmt.Exec(&username, &password, &first_name, &last_name, &email)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				res := statusRes{Status: 200, Msg: "Success"}
				json.NewEncoder(w).Encode(res)
			}
		} 
	}	
		else {
		res := statusRes{Status: 400, Msg: "Method Must be post"}
		json.NewEncoder(w).Encode(res)
	}


func login(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) != 0 && checkErr(w, r, err) {
		//check session if avaliabel
		res := statusRes{Status: 200, Msg: " login session avliabe"}
		json.NewEncoder(w).Encode(res)
	}
	if r.Method != "POST" {
		res := statusRes{Status: 400, Msg: "Method Must be post"}
		json.NewEncoder(w).Encode(res)
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	users := QueryUser(username)

	// compare password
	var password_tes = bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(password))

	if password_tes == nil {
		//login success
		session := sessions.Start(w, r)
		session.Set("username", users.Username)
		session.Set("name", users.FirstName)
		res := statusRes{Status: 200, Msg: "Successfull login"}
		json.NewEncoder(w).Encode(res)
	} else {
		//login failed
		res := statusRes{Status: 400, Msg: "Invalid login"}
		json.NewEncoder(w).Encode(res)
	}

}

func logout(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	session.Clear()
	sessions.Destroy(w, r)
	res := statusRes{Status: 200, Msg: "Successfull logout"}
	json.NewEncoder(w).Encode(res)
}
