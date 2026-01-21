package main

import (
	"database/sql"
	"fmt"
	"os"
)

var DB *sql.DB

func ConnectDB() error {
	pswd := os.Getenv("MYSQL_PASSWORD")
	dsn := fmt.Sprintf(
		"root:%s@tcp(127.0.0.1:3306)/pdos_system?charset=utf8mb4",
		pswd,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	DB = db
	fmt.Println("Successfully connected and verified database!")
	return nil
}
