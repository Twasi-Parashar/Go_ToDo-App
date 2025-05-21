package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	//Database connection
	var err error
	db, err = sql.Open("mysql", "root:Mohit2607@@tcp(127.0.0.1:330)/todoapp")
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	// Serve static files (e.g., HTML, CSS, JS files)
	fs := http.FileServer(http.Dir("./")) // serve files from the current directory
	http.Handle("/", fs)

	//Routes
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)

	//Setting server
	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Registration handler
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//Getting data for registration
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	//Inserting data into database
	_, err := db.Exec("INSERT INTO users (username,password) VALUES(?, ?)", username, password)
	if err != nil {
		http.Error(w, "Could not register user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//Success
	http.Redirect(w, r, "/login.html", http.StatusSeeOther)
}

// Login Handler
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//Login data
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	//Retrieving the stored password
	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPassword)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//Validating
	if storedPassword != password {
		http.Error(w, "Invaalid username or password", http.StatusUnauthorized)
		return
	}

	//Login to dashboard
	http.Redirect(w, r, "/dashboard.html", http.StatusSeeOther)

}
