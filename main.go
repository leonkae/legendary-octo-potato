package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var db *sql.DB

func main() {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		log.Fatal("MYSQL_DSN not set")
	}

	var err error

	for i := 1; i <= 15; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
		}

		if err == nil {
			log.Println("✅ Connected to MySQL")
			break
		}

		log.Printf("⏳ Waiting for MySQL... (%d/15)", i)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("❌ Could not connect to MySQL:", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(100)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/users", usersHandler)

	log.Println("🚀 Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		rows, err := db.Query("SELECT id, name FROM users")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var u User
			rows.Scan(&u.ID, &u.Name)
			users = append(users, u)
		}
		json.NewEncoder(w).Encode(users)

	case http.MethodPost:
		var u User
		json.NewDecoder(r.Body).Decode(&u)

		res, err := db.Exec("INSERT INTO users (name) VALUES (?)", u.Name)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		id, _ := res.LastInsertId()
		u.ID = int(id)
		json.NewEncoder(w).Encode(u)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
