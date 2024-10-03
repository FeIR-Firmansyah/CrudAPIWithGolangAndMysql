// main.go
package main

import (
	"SimpleCrud/routes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//connect to anekazoo  database
	dataSourceName := os.Getenv("DATABASE_URL")
	database, connectionError := sql.Open("mysql", dataSourceName)

	if connectionError != nil {
		log.Fatalf("error: %s", connectionError.Error())
	}
	defer database.Close()

	// Check if the database answer, if not then create one.
	databasePingError := database.Ping()
	if databasePingError != nil {
		log.Printf("databasePingError: %s", databasePingError.Error())

		//creating new anekazoo database
		dataSourceName = os.Getenv("DATABASE_ROOT_URL")
		database, connectionError = sql.Open("mysql", dataSourceName)                                   //making connection to "127.../" cause its root directory
		if runCreatingDatabaseError := RunCreatingDatabase(database); runCreatingDatabaseError != nil { // run RunCreatingDatabase in root mysql directory
			log.Fatalf("runCreatingDatabaseError: %v", runCreatingDatabaseError.Error())
		}

		//connect to newly created database
		dataSourceName = os.Getenv("DATABASE_URL")
		database, connectionError = sql.Open("mysql", dataSourceName)
		if connectionError != nil {
			log.Fatalf("error: %s", connectionError.Error())
		}
		defer database.Close()

	}
	log.Printf("Database Connected Successfully")

	//create or retrieving tables from database
	if runCreatingTableError := RunCreatingTable(database); runCreatingTableError != nil {
		log.Fatalf("runCreatingTableError: %v", runCreatingTableError.Error())
	}
	log.Printf("Creating/Retrieving Table Successfully")

	// Register the handlers for different URLs
	http.HandleFunc("/", routes.HandleRoot)
	http.HandleFunc("/animal", func(responseWriter http.ResponseWriter, request *http.Request) {
		routes.HandleAnimal(responseWriter, request, database)
	})
	http.HandleFunc("/animal/", func(responseWriter http.ResponseWriter, request *http.Request) {
		routes.HandleAnimalByID(responseWriter, request, database)
	})

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func RunCreatingDatabase(database *sql.DB) error {
	//get and execute CreatingDatabase.sql migration file
	currentWorkingDirectory, directoryError := os.Getwd()
	if directoryError != nil {
		return directoryError
	}
	sqlFile := currentWorkingDirectory + "\\migration\\CreatingDatabase.sql"
	sqlScript, fileReadError := ioutil.ReadFile(sqlFile)
	if fileReadError != nil {
		return fileReadError
	}

	log.Printf("Running SQL Script: %s", string(sqlScript))
	_, executionError := database.Exec(string(sqlScript))
	return executionError
}

func RunCreatingTable(database *sql.DB) error {
	//get and execute CreatingDatabase.sql migration file
	currentWorkingDirectory, directoryError := os.Getwd()
	if directoryError != nil {
		return directoryError
	}
	sqlFile := currentWorkingDirectory + "\\migration\\CreatingTable.sql"
	sqlScript, fileReadError := ioutil.ReadFile(sqlFile)
	if fileReadError != nil {
		return fileReadError
	}

	log.Printf("Running SQL Script: %s", string(sqlScript))
	_, executionError := database.Exec(string(sqlScript))
	return executionError
}
