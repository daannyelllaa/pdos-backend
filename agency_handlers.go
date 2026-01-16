package main

import (
	"encoding/json"
	"net/http"
)

func getAgencies(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := DB.Query(`
		SELECT id, agency_name
		FROM agencies
		ORDER BY agency_name ASC
	`)
	if err != nil {
		http.Error(w, "Failed to fetch agencies", 500)
		return
	}
	defer rows.Close()

	type Agency struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	var agencies []Agency

	for rows.Next() {
		var a Agency
		rows.Scan(&a.ID, &a.Name)
		agencies = append(agencies, a)
	}

	json.NewEncoder(w).Encode(agencies)
}
