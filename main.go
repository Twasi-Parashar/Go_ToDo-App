package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Task struct {
	ID          int
	Title       string
	Description string
	DueDate     string
}

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
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/add-task", addTaskHandler)
	http.HandleFunc("/delete-task", deleteTaskHandler)
	http.HandleFunc("/logout", logoutHandler)

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
	http.Redirect(w, r, "/dashboard?username="+username, http.StatusSeeOther)
}

// DashBoard Handler
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Missing username", http.StatusBadRequest)
		return
	}

	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	rows, err := db.Query("SELECT id,title,description,due_date FROM tasks WHERE user_id = ? ORDER BY due_date", userID)
	if err != nil {
		http.Error(w, "Error fetching tasks: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		var dueDate string
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &dueDate); err != nil {
			http.Error(w, "Error reading task: "+err.Error(), http.StatusInternalServerError)
			return
		}
		t.DueDate = dueDate
		tasks = append(tasks, t)
	}

	//Render Dashboard
	fmt.Fprint(w, `<html><head><title>Dashboard</title></head><body>`)
	fmt.Fprintf(w, `<h1>Welcome, %s</h1><a href="/logout">Logout</a><br><br>`, username)
	fmt.Fprintf(w, `<form method="POST" action="/add-task">
		<input type="hidden" name="username" value="%s">
		<input name="title" placeholder="Title" required>
		<input name="description" placeholder="Description">
		<input type="date" name="due_date" required>
		<button type="submit">Add Task</button>
	</form><br><ul>`, username)

	//Uploading data
	for _, task := range tasks {
		fmt.Fprintf(w, `<li><strong>%s</strong> - %s (Due: %s) 
			<a href="/delete-task?id=%d&username=%s">Delete</a></li>`,
			task.Title, task.Description, task.DueDate, task.ID, username)
	}

	fmt.Fprint(w, `</ul></body></html>`)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	title := r.FormValue("title")
	description := r.FormValue("description")
	dueDate := r.FormValue("due_date")

	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO tasks (user_id,title,description,due_date) VALUES(?,?,?,?)", userID, title, description, dueDate)
	if err != nil {
		http.Error(w, "Error adding task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard?username="+username, http.StatusSeeOther)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	username := r.URL.Query().Get("username")

	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM tasks WHERE id = ? AND user_id = ?", id, userID)
	if err != nil {
		http.Error(w, "Error deleting task:"+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/dashboard?username="+username, http.StatusSeeOther)
}

// Logout
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Since you're not using sessions, logout just redirects to login
	http.Redirect(w, r, "/login.html", http.StatusSeeOther)
}
