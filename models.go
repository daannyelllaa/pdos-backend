package main

type Agency struct {
	ID         uint
	AgencyName string
}

type Worker struct {
	ID                 uint
	AgencyID           uint
	FirstName          string
	LastName           string
	PassportNumber     string
	RegistrationStatus string
}

type User struct {
	ID           uint
	WorkerID     uint
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    string
}
