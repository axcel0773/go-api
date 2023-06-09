// Go package
package main

import (
	"database/sql"
	"encoding/json"
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

// Get all movies
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

		movies = append(movies, Movie{MovieID: movieID, MovieName: movieName})
	}

	var response = JsonResponse{Type: "success", Data: movies}

	json.NewEncoder(w).Encode(response)
}

// Create a movie
// response and request handlers
func CreateMovie(w http.ResponseWriter, r *http.Request) {
	movieId := r.FormValue("movieid")
	movieName := r.FormValue("moviename")

	var response = JsonResponse{}

	if movieId == "" || movieName == "" {
		response = JsonResponse{Type: "error", Message: "You are missing movieID or movieName parameter."}
	} else {
		db := setupDB()

		printMessage("Inserting movie into DB")

		fmt.Println("Inserting new movie with ID: " + movieId + " and name: " + movieName)

		var lastInsertID int

		err := db.QueryRow("INSERT INTO movies(movieID, movieName) VALUES ($1, $2) returning id;", movieId, movieName).Scan(&lastInsertID)

		//check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The movie has been inserted successfully!"}
	}
	json.NewEncoder(w).Encode(response)
}

// Delete a movie
// response and request handlers
func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	movieId := params["movieid"]

	var response = JsonResponse{}

	if movieId == "" {
		response = JsonResponse{Type: "error", Message: "You are missing movieID parameter."}
	} else {
		db := setupDB()

		printMessage("Deleting movie from Db")

		_, err := db.Exec("DELETE FROM movies WHERE movieID = $1", movieId)

		//check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The movie has been deleted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

// Delete all records
// response and request handlers
func DeleteMovies(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	printMessage("Deleting all movies...")

	_, err := db.Exec("DELETE FROM movies")

	//check errors
	checkErr(err)

	printMessage("All movies have been deleted successfully!")

	var response = JsonResponse{Type: "success", Message: "All movies have been deleted succesfully!"}

	json.NewEncoder(w).Encode(response)
}
