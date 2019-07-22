package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// GetBookByID a function to get a single book given it's ID
func (h *handler) GetBookByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	query := "SELECT id, title, author, isbn, stock FROM books WHERE id = " + param.ByName("bookID")
	rows, err := h.db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}

	var books []Book

	for rows.Next() {
		book := Book{}
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.ISBN,
			&book.Stock,
		)
		if err != nil {
			log.Println(err)
			continue
		}
		books = append(books, book)
	}

	bytes, err := json.Marshal(books)
	if err != nil {
		log.Println(err)
		return
	}

	renderJSON(w, bytes, http.StatusOK)
}

// InsertBook a function to insert book to DB
func (h *handler) InsertBook(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	// read json body
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		renderJSON(w, []byte(`
			message: "Failed to read body"
		`), http.StatusBadRequest)
		return
	}

	// parse json body
	var book Book
	err = json.Unmarshal(body, &book)
	if err != nil {
		log.Println(err)
		return
	}

	// executing insert query
	query := "INSERT INTO books (id, title, author, isbn, stock) VALUES ($1, $2, $3, $4, $5)"
	_, err = h.db.Exec(query, book.ID, book.Title, book.Author, book.ISBN, book.Stock)
	if err != nil {
		log.Println(err)
		return
	}

	renderJSON(w, []byte(`
	{
		status: "success",
		message: "Insert book success!"
	}
	`), http.StatusOK)
}

// EditBook a function to change book data in DB, with given params
func (h *handler) EditBook(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	bookID := param.ByName("bookID")
	// read json body
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		renderJSON(w, []byte(`
			message: "Failed to read body"
		`), http.StatusBadRequest)
		return
	}

	// parse json body
	var book Book
	err = json.Unmarshal(body, &book)
	if err != nil {
		log.Println(err)
		return
	}

	query := "UPDATE books SET title = %s, author = $s, isbn = %s, stock = %d WHERE id = %s"
	_, err = h.db.Exec(fmt.Sprintf(query, book.Title, book.Author, book.ISBN, book.Stock, bookID))
	if err != nil {
		log.Println(err)
		return
	}

	renderJSON(w, []byte(`
	{
		status: "success",
		message: "Update book success!"
	}
	`), http.StatusOK)
}

// DeleteBookByID a function to remove book data from DB, given bookID
func (h *handler) DeleteBookByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	bookID := param.ByName("bookID")

	query := "DELETE FROM books WHERE id = %s"
	_, err := h.db.Exec(fmt.Sprintf(query, bookID))
	if err != nil {
		log.Println(err)
		return
	}

	renderJSON(w, []byte(`
	{
		status: "success",
		message: "Delete book success!"
	}
	`), http.StatusOK)
}

// InsertMultipleBooks a function to insert multiple book data, given file of books data
func (h *handler) InsertMultipleBooks(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	var buffer bytes.Buffer

	file, header, err := r.FormFile("books")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	// get file name
	name := strings.Split(header.Filename, ".")
	log.Printf("Received a file with name = %s\n", name[0])

	// copy file to buffer
	io.Copy(&buffer, file)

	contents := buffer.String()

	// TODO: Do something with the file's contents, instead of just print it
	fmt.Println(contents)

	buffer.Reset()
}

// LendBook a function to record book lending in DB and update book stock in book tables
func (h *handler) LendBook(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

}
