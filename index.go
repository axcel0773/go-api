// Go package
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "123"
	DB_NAME     = "movies"
)

// DB set up
func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)

	checkErr(err)

	return DB
}

type Movie struct {
	MovieID   string `json:"movieid"`
	MovieName string `json:moviename`
}

type JsonResponse struct {
	Type    string  `json:"type"`
	Data    []Movie `json:"data"`
	Message string  `json:message`
}

// Main function
func main() {

	//init the mux router
	router := mux.NewRouter()

	//Route handles & endpoints

	//Get all movies
	router.HandleFunc("/movies/", GetMovies).Methods("GET")

	//Create a movie
	router.HandleFunc("/movies/", CreateMovie).Methods("POST")

	//Delete a specific movie by the movieID
	router.HandleFunc("/movies/{movieid}", DeleteMovie).Methods("DELETE")

	//Delete all movies
	router.HandleFunc("/movies/", DeleteMovies).Methods("DELETE")

	//serve the app
	fmt.Println("Server at 8080")
	log.Fatal(http.ListenAndServe(":8000", router))
}

// Function for handling messages
func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

// Function for handling erros
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

//Get all movies

// response and request handlers
func GetMovies(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	printMessage("Getting movies...")

	//Get all movies from movies table that don't have movieID = "1"
	rows, err := db.Query("SELECT * FROM movies")

	//check erros
	checkErr(err)

	//var response []JsonResponse
	var movies []Movie

	//for each movie
	for rows.Next() {
		var id int
		var movieID string
		var movieName string

		err = rows.Scan(&id, &movieID, &movieName)

		//check the erros
		checkErr(err)
	}
}
