package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if err := ConnectDB(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/agencies", getAgencies)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
