package internal

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
func (h *Handler) GetBookByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	query := "SELECT id, title, author, isbn, stock FROM books WHERE id = " + param.ByName("bookID")
	rows, err := h.DB.Query(query)
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
func (h *Handler) InsertBook(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
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
	_, err = h.DB.Exec(query, book.ID, book.Title, book.Author, book.ISBN, book.Stock)
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
func (h *Handler) EditBook(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
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

	query := "UPDATE books SET title = '%s', author = '%s', isbn = '%s', stock = %d WHERE id = %s"
	_, err = h.DB.Exec(fmt.Sprintf(query, book.Title, book.Author, book.ISBN, book.Stock, bookID))
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
func (h *Handler) DeleteBookByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	bookID := param.ByName("bookID")

	query := "DELETE FROM books WHERE id = %s"
	_, err := h.DB.Exec(fmt.Sprintf(query, bookID))
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
func (h *Handler) InsertMultipleBooks(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	var buffer bytes.Buffer

	query := "INSERT INTO books (id, title, author, isbn, stock) VALUES ($1, $2, $3, $4, $5)"

	file, header, err := r.FormFile("books")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	// get file name
	name := strings.Split(header.Filename, ".")
	if name[1] != "csv" {
		log.Println("File format not supported")
		return
	}
	log.Printf("Received a file with name = %s\n", name[0])

	// copy file to buffer
	io.Copy(&buffer, file)

	contents := buffer.String()

	// Split contents to rows
	rows := strings.Split(contents, "\n")
	for i, row := range rows {
		// skip title
		if i == 0 {
			continue
		}
		// Split rows to column
		columns := strings.Split(row, ",")

		_, err = h.DB.Exec(query, columns[0], columns[1], columns[2], columns[3], columns[4])
		if err != nil {
			log.Println(err)
			return
		}
	}

	buffer.Reset()

	renderJSON(w, []byte(`
	{
		status: "success",
		message: "Update book success!"
	}
	`), http.StatusOK)
}

// LendBook a function to record book lending in DB and update book stock in book tables
func (h *Handler) LendBook(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	bookID := param.ByName("bookID")
	bookStockQuery := "SELECT stock FROM books WHERE id = %s"

	// Read userID
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		renderJSON(w, []byte(`
			message: "Failed to read body"
		`), http.StatusBadRequest)
		return
	}

	// parse json body
	var request LendRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Println(err)
		return
	}

	// Get book stock from DB
	rows, err := h.DB.Query(fmt.Sprintf(bookStockQuery, bookID))
	if err != nil {
		log.Println(err)
		return
	}

	var bookStock int
	for rows.Next() {
		err := rows.Scan(&bookStock)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	// Insert Book to Lend tables
	insertBookLendingQuery := "INSERT INTO lend (user_id, book_id) VALUES ($1, $2)"
	_, err = h.DB.Exec(insertBookLendingQuery, request.UserID, bookID)
	if err != nil {
		log.Println(err)
		return
	}

	// Update Book stock query
	updateStockQuery := "UPDATE books SET stock = $1 WHERE id = $2"
	_, err = h.DB.Exec(updateStockQuery, (bookStock - 1), bookID)
	if err != nil {
		log.Println(err)
		return
	}

	renderJSON(w, []byte(`
	{
		status: "success",
		message: "Your order has been recorded!"
	}
	`), http.StatusOK)
}

// GetMultipleBooks a function to get multiple books in a single request
func (h *Handler) GetMultipleBooks(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	// TODO
}
