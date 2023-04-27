package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type server struct {
	db *sql.DB
}

func database() server {
	database, _ := sql.Open("sqlite3", "database.db")
	server := server{db: database}
	return server
}

func (s *server) register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	_, err = s.db.Exec("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)", name, email, password)
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *server) users(w http.ResponseWriter, r *http.Request) {
	var users []User

	result, _ := s.db.Query("SELECT * FROM users;")

	for result.Next() {
		var user User
		err := result.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	templ, _ := template.ParseFiles("static/users.html")
	templ.Execute(w, users)
}

func (s *server) delete(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	_, _ = s.db.Exec("DELETE FROM users WHERE id=$1", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func update(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	data := map[string]interface{}{"id": id}
	templ, _ := template.ParseFiles("static/update.html")
	templ.Execute(w, data)
}

func (s *server) updateFinal(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	_, _ = s.db.Exec("UPDATE users set name=$1, email=$2, password=$3 where id=$4", name, email, password, id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	s := database()
	defer s.db.Close()

	http.Handle("/", http.FileServer(http.Dir("./static")))
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/register", s.register)
	http.HandleFunc("/users", s.users)
	http.HandleFunc("/delete", s.delete)
	http.HandleFunc("/update", update)
	http.HandleFunc("/updateFinal", s.updateFinal)
	fmt.Println("Server is running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
