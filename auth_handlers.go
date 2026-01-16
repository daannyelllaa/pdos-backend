package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func register(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		AgencyID       int    `json:"agency_id"`
		FirstName      string `json:"first_name"`
		LastName       string `json:"last_name"`
		MiddleName     string `json:"middle_name"`
		Gender         string `json:"gender"`
		DateOfBirth    string `json:"date_of_birth"`
		PassportNumber string `json:"passport_number"`
		Email          string `json:"email"`
		Phone          string `json:"phone"`
		Username       string `json:"username"`
		Password       string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	tx, err := DB.Begin()
	if err != nil {
		http.Error(w, "Transaction error", 500)
		return
	}

	// 1️⃣ Insert worker
	res, err := tx.Exec(`
		INSERT INTO workers (
			agency_id, first_name, last_name, middle_name,
			gender, date_of_birth, passport_number, email, phone
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		input.AgencyID,
		input.FirstName,
		input.LastName,
		input.MiddleName,
		input.Gender,
		input.DateOfBirth,
		input.PassportNumber,
		input.Email,
		input.Phone,
	)

	if err != nil {
		tx.Rollback()
		http.Error(w, "Worker registration failed", 400)
		return
	}

	workerID, _ := res.LastInsertId()

	// 2️⃣ Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Password error", 500)
		return
	}

	// 3️⃣ Insert user
	_, err = tx.Exec(`
		INSERT INTO users (worker_id, username, email, password_hash)
		VALUES (?, ?, ?, ?)`,
		workerID,
		input.Username,
		input.Email,
		string(hash),
	)

	if err != nil {
		tx.Rollback()
		http.Error(w, "User creation failed", 400)
		return
	}

	tx.Commit()

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Registration successful",
	})
}

func login(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)

	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	json.NewDecoder(r.Body).Decode(&input)

	var user User
	err := DB.QueryRow(`
		SELECT id, worker_id, password_hash
		FROM users
		WHERE username = ?
	`, input.Username).Scan(&user.ID, &user.WorkerID, &user.PasswordHash)

	if err != nil {
		http.Error(w, "Invalid credentials", 401)
		return
	}

	if bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(input.Password),
	) != nil {
		http.Error(w, "Invalid credentials", 401)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Login successful",
		"worker_id": user.WorkerID,
	})
}
