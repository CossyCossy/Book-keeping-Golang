package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Person struct {
	gorm.Model

	Name  string
	Email string `gorm:"type:varchar(100);unique_index"`
	Books []Book
}

type Book struct {
	gorm.Model

	Title      string
	Author     string
	CallNumber int `gorm:"unique_index"`
	PersonID   int
}

/*
	when creating a database with dummy data
*/

// var (
// 	person = &Person{Name: "Cosmas", Email: "saxonmuti@gmail.com"}
// 	books  = []Book{
// 		{Title: "Pirates of the Caribbean", Author: "T.S.Eliot", CallNumber: 1234, PersonID: 1},
// 		{Title: "Jack Ryan", Author: "Lint Master", CallNumber: 1334, PersonID: 1},
// 	}
// )

var db *gorm.DB
var err error

func main() {

	/*
		loading environmental variables
	*/
	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbport := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbName := os.Getenv("NAME")
	password := "zoom20$$" //os.Getenv("PASSWD")

	/*
		Database connection string
	*/
	dbURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, dbport, user, dbName, password)

	/*
		opening database connection
	*/
	db, err = gorm.Open(dialect, dbURI)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully connected to database")
	}

	/*
		close connection to database when main function finishes
	*/
	defer db.Close()

	/*
		Migrating migrations to the database if they have not been created
	*/
	db.AutoMigrate(&Person{}, &Book{})

	/*
		when creating a database with dummy data
	*/
	// db.Create(person)
	// for idx := range books {
	// 	db.Create(&books[idx])
	// }

	/*
		API endpoints
	*/
	router := mux.NewRouter()

	router.HandleFunc("/create/person", createPerson).Methods("POST")
	router.HandleFunc("/people", getPeople).Methods("GET")
	router.HandleFunc("/person/{id}", getPerson).Methods("GET")
	router.HandleFunc("/delete/person/{id}", deletePerson).Methods("DELETE")
	router.HandleFunc("/update/person/{id}", updatePerson).Methods("PATCH")

	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/create/book", createBook).Methods("POST")
	router.HandleFunc("/book/{id}", getBook).Methods("GET")
	router.HandleFunc("/delete/book/{id}", deleteBook).Methods("DELETE")
	router.HandleFunc("/update/book/{id}", updateBook).Methods("PATCH")

	log.Fatal(http.ListenAndServe(":8080", router))
}

/*
	API controllers
*/

/*People API controlllers*/
func getPeople(w http.ResponseWriter, r *http.Request) {

	var people []Person

	db.Find(&people)

	json.NewEncoder(w).Encode(people)
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var person Person
	var books []Book

	db.First(&person, id)
	db.Model(&person).Related(&books)

	person.Books = books

	if person.ID == 0 {
		json.NewEncoder(w).Encode(person)
		return
	}

	json.NewEncoder(w).Encode(person)
}

func createPerson(w http.ResponseWriter, r *http.Request) {
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)

	//db.Create(&person)
	createdPerson := db.Create(&person)
	err = createdPerson.Error

	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(person)
	}
}

func deletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var person Person

	db.First(&person, id)

	if person.ID == 0 {
		json.NewEncoder(w).Encode(person)
		return
	}

	db.Delete(&person)
	json.NewEncoder(w).Encode(person)
}

func updatePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var person Person

	db.First(&person, id)

	if person.ID == 0 {
		json.NewEncoder(w).Encode(person)
		return
	}

	_ = json.NewDecoder(r.Body).Decode(&person)

	db.Save(&person)
	json.NewEncoder(w).Encode(person)
}

/*Books API controlllers*/
func getBooks(w http.ResponseWriter, r *http.Request) {

	var books []Book

	db.Find(&books)

	json.NewEncoder(w).Encode(books)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	createdBook := db.Create(&book)
	err = createdBook.Error

	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(book)
	}
}

func getBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var book Book

	db.First(&book, id)

	if book.ID == 0 {
		json.NewEncoder(w).Encode(book)
		return
	}

	json.NewEncoder(w).Encode(book)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var book Book

	db.First(&book, id)

	if book.ID == 0 {
		json.NewEncoder(w).Encode(book)
		return
	}

	db.Delete(&book)
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var book Book

	db.First(&book, id)

	if book.ID == 0 {
		json.NewEncoder(w).Encode(book)
		return
	}

	_ = json.NewDecoder(r.Body).Decode(&book)

	db.Save(&book)
	json.NewEncoder(w).Encode(book)
}
