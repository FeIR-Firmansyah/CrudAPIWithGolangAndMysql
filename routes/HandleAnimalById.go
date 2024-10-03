package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	StringConvert "strconv"
	"strings"
)

func HandleAnimalByID(responseWriter http.ResponseWriter, request *http.Request, database *sql.DB) {
	switch request.Method {
	case http.MethodGet, http.MethodHead:
		handleRequestWithGetMethod(responseWriter, request, database)
	default:
		http.Error(responseWriter, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleRequestWithGetMethod(responseWriter http.ResponseWriter, request *http.Request, database *sql.DB) {
	// Extract the ID from the URL path
	fmt.Println("Handling GET request in /Animal/")
	parts := strings.Split(request.URL.Path, "/") //parts [0] is empty
	idInString := parts[2]
	id, err := StringConvert.Atoi(idInString)
	if err != nil {
		http.Error(responseWriter, "Invalid animal ID", http.StatusBadRequest)
		return
	}

	// Query the database for the animal with the given ID
	var animal Animal
	query := "SELECT name, class, legs FROM animal WHERE id = ?"
	row := database.QueryRow(query, id)
	err = row.Scan(&animal.Name, &animal.Class, &animal.Legs)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(responseWriter, "Animal not found", http.StatusNotFound)
		} else {
			http.Error(responseWriter, "Error retrieving animal", http.StatusInternalServerError)
		}
		return
	}

	// Respond with the animal details
	responseWriter.WriteHeader(http.StatusOK)
	fmt.Fprintf(responseWriter, "Animal found: Name: %s, Class: %s, Legs: %d", animal.Name, animal.Class, animal.Legs)
}
