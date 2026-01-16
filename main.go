package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Connect to database
	if err := ConnectDB(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Routes
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/agencies", getAgencies)

	// Railway-compatible port handling
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback for local development
	}

	log.Println("Server running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
