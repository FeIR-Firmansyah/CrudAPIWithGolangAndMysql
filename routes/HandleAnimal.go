package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Animal struct updated to include Class and Legs fields
type Animal struct {
	ID    int
	Name  string
	Class string
	Legs  int
}

func HandleAnimal(responseWriter http.ResponseWriter, request *http.Request, database *sql.DB) {
	switch request.Method {
	case http.MethodGet:
		handleAnimalGet(database, responseWriter, request)
	case http.MethodPost:
		handleAnimalPost(database, responseWriter, request)
	case http.MethodPut:
		handleAnimalPut(database, responseWriter, request)
	case http.MethodDelete:
		handleAnimalDelete(database, responseWriter, request)
	default:
		http.Error(responseWriter, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleAnimalPost(database *sql.DB, responseWriter http.ResponseWriter, request *http.Request) {
	fmt.Println("Handling POST request in /Animal")
	//read request json body
	body, readBodyError := io.ReadAll(request.Body)
	if readBodyError != nil {
		http.Error(responseWriter, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	// create key/value pairs
	var animalData map[string]interface{}
	if err := json.Unmarshal(body, &animalData); err != nil {
		http.Error(responseWriter, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	// Create the Animal struct from the data
	animal := Animal{}
	if name, ok := animalData["name"].(string); ok {
		animal.Name = name
	}
	if class, ok := animalData["class"].(string); ok {
		animal.Class = class
	}
	if legs, ok := animalData["legs"].(float64); ok { // JSON numbers are decoded as float64
		animal.Legs = int(legs)
	}

	// Check if an animal with the same name, class, and legs already exists in the database
	checkIfIdExistQuery := "SELECT id FROM animal WHERE name = ? AND class = ? AND legs = ?"
	var existingID int
	err := database.QueryRow(checkIfIdExistQuery, animal.Name, animal.Class, animal.Legs).Scan(&existingID)
	if err != nil && err != sql.ErrNoRows {
		// Handle any errors other than "not found"
		http.Error(responseWriter, "Error checking for existing animal", http.StatusInternalServerError)
		return
	}

	// if we have an existing animal then response with http.StatusConflict
	if existingID != 0 {
		http.Error(responseWriter, "Exact match already exists", http.StatusConflict)
		return
	}

	// Insert the animal into the database
	query := "INSERT INTO animal (name, class, legs) VALUES (?, ?, ?)"
	_, queryError := database.Exec(query, animal.Name, animal.Class, animal.Legs)
	if queryError != nil {
		http.Error(responseWriter, "Error inserting animal into database", http.StatusInternalServerError)
		return
	}

	// Respond with success
	responseWriter.WriteHeader(http.StatusCreated)
	fmt.Fprintf(responseWriter, "Animal %s of class %s with %d legs added successfully!", animal.Name, animal.Class, animal.Legs)
}

func handleAnimalGet(database *sql.DB, responseWriter http.ResponseWriter, request *http.Request) {
	// Query the database to select all animals
	fmt.Println("Handling GET request in /Animal")
	query := "SELECT id, name, class, legs FROM animal"
	rows, err := database.Query(query)
	if err != nil {
		http.Error(responseWriter, "Error retrieving animals", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// read the results
	var animals []Animal
	for rows.Next() {
		var animal Animal
		if err := rows.Scan(&animal.ID, &animal.Name, &animal.Class, &animal.Legs); err != nil {
			http.Error(responseWriter, "Error scanning animals", http.StatusInternalServerError)
			return
		}
		animals = append(animals, animal)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		http.Error(responseWriter, "Error encountered during row iteration", http.StatusInternalServerError)
		return
	}

	// if theres no rows returned http.StatusNotFound
	if len(animals) == 0 {
		http.Error(responseWriter, "No animals found", http.StatusNotFound)
		return
	}

	// append all data into one sigle output, then send it
	output := ""
	for _, animal := range animals {
		output += fmt.Sprintf("ID: %d, Name: %s, Class: %s, Legs: %d\n", animal.ID, animal.Name, animal.Class, animal.Legs)
	}
	responseWriter.WriteHeader(http.StatusOK)
	fmt.Fprintf(responseWriter, "%s", output)
}

func handleAnimalPut(database *sql.DB, responseWriter http.ResponseWriter, request *http.Request) {
	fmt.Println("Handling PUT request in /Animal")

	// get request json body
	body, readBodyError := io.ReadAll(request.Body)
	if readBodyError != nil {
		http.Error(responseWriter, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	// Parse the request body into a struct (Animal)
	var animalData Animal
	if jsonUnmarshalError := json.Unmarshal(body, &animalData); jsonUnmarshalError != nil {
		http.Error(responseWriter, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	// Check if an animal with the given ID exists in the database
	checkIfIdExistQuery := "SELECT id FROM animal WHERE id = ?"
	var existingID int
	err := database.QueryRow(checkIfIdExistQuery, animalData.ID).Scan(&existingID)
	if err != nil && err != sql.ErrNoRows {
		// Handle any errors other than "not found"
		http.Error(responseWriter, "Error checking for existing animal", http.StatusInternalServerError)
		return
	}

	// if the id exist, Update the animal
	if existingID != 0 {
		// Update the existing animal
		query := "UPDATE animal SET name = ?, class = ?, legs = ? WHERE id = ?"
		_, updateError := database.Exec(query, animalData.Name, animalData.Class, animalData.Legs, animalData.ID)
		if updateError != nil {
			http.Error(responseWriter, "Error updating animal in database", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(responseWriter, "Update on %s (ID: %d) succeeded!", animalData.Name, animalData.ID)
		return
	}

	//if id doesnt exist then insert/create as new data rows
	query := "INSERT INTO animal (name, class, legs) VALUES (?, ?, ?)"
	_, queryError := database.Exec(query, animalData.Name, animalData.Class, animalData.Legs)
	if queryError != nil {
		http.Error(responseWriter, "Error inserting animal into database", http.StatusInternalServerError)
		return
	}

	responseWriter.WriteHeader(http.StatusOK)
	fmt.Fprintf(responseWriter, "Data with %d doesnt exist, added as new data instead", animalData.ID)
}

func handleAnimalDelete(database *sql.DB, responseWriter http.ResponseWriter, request *http.Request) {
	fmt.Println("Handling DELETE request in /Animal")

	// get request json body
	body, readBodyError := io.ReadAll(request.Body)
	if readBodyError != nil {
		http.Error(responseWriter, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	//pasing request body into struct (Animal)
	var animalData Animal
	if jsonUnmarshalError := json.Unmarshal(body, &animalData); jsonUnmarshalError != nil {
		http.Error(responseWriter, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	// check if id exist
	checkIfIdExistQuery := "SELECT id FROM animal WHERE id = ?"
	var existingID int
	err := database.QueryRow(checkIfIdExistQuery, animalData.ID).Scan(&existingID)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no rows are returned, the ID does not exist. response with http.StatusNotFound
			http.Error(responseWriter, fmt.Sprintf("Animal with ID %d does not exist", animalData.ID), http.StatusNotFound)
			return
		}
		// Handle other possible errors
		http.Error(responseWriter, "Error checking for existing animal", http.StatusInternalServerError)
		return
	}

	//delete animal rows that match the given id
	query := "DELETE FROM animal WHERE id = ?"
	_, deleteError := database.Exec(query, animalData.ID)
	if deleteError != nil {
		http.Error(responseWriter, "Error deleting animal from database", http.StatusInternalServerError)
		return
	}

	responseWriter.WriteHeader(http.StatusOK)
	fmt.Fprintf(responseWriter, "Animal with ID %d successfully deleted.", animalData.ID)
}
