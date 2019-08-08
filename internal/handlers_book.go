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
	"time"

	"github.com/julienschmidt/httprouter"
)

// GetBookByID a function to get a single book given it's ID
func (h *Handler) GetBookByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	query := fmt.Sprintf("SELECT id, title, author, isbn, stock FROM books WHERE id = %s", param.ByName("bookID"))
	rows, err := h.DB.Query(query)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

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
	query := fmt.Sprintf("INSERT INTO books (id, title, author, isbn, stock) VALUES (%d, '%s', '%s', '%s', %d)", book.ID, book.Title, book.Author, book.ISBN, book.Stock)
	_, err = h.DB.Query(query)
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

	query := fmt.Sprintf("UPDATE books SET title = '%s', author = '%s', isbn = '%s', stock = %d WHERE id = %s", book.Title, book.Author, book.ISBN, book.Stock, bookID)
	_, err = h.DB.Query(query)
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

	query := fmt.Sprintf("DELETE FROM books WHERE id = %s", bookID)
	_, err := h.DB.Exec(query)
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
	rows := strings.Split(contents, "\n")[1:]

	h.RunWorkerPool(rows)

	_, err = h.DB.Query("DELETE FROM BOOKS WHERE ID > 10")
	if err != nil {
		log.Panic("Failed to delete", err)
	}

	start := time.Now()
	for _, row := range rows[1:] {
		// Split rows to column
		columns := strings.Split(row, ",")

		// check for any wrong input
		if len(columns) != 5 {
			continue
		}

		query := fmt.Sprintf("INSERT INTO books (id, title, author, isbn, stock) VALUES (%s, '%s', '%s', '%s', %s) ON CONFLICT(id) DO NOTHING", columns[0], columns[1], columns[2], columns[3], columns[4])
		_, err := h.DB.Exec(query)
		if err != nil {
			log.Println(err)
			log.Panic(h.DB.Stats())
			continue
		}
	}
	log.Println("[Without worker] Finished with time", time.Since(start).Seconds(), "s")

	buffer.Reset()

	renderJSON(w, []byte(`
	{
		status: "success",
		message: "Insert book success!"
	}
	`), http.StatusOK)
}

// LendBook a function to record book lending in DB and update book stock in book tables
func (h *Handler) LendBook(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
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

	// Insert Book to Lend tables
	insertBookLendingQuery := fmt.Sprintf("INSERT INTO lend (user_id, book_id) VALUES (%d, %d)", request.UserID, request.BookID)
	_, err = h.DB.Query(insertBookLendingQuery)
	if err != nil {
		log.Println(err)
		return
	}

	// Update Book stock query
	updateStockQuery := fmt.Sprintf("UPDATE books SET stock = stock - 1 WHERE id = %d", request.BookID)
	_, err = h.DB.Query(updateStockQuery)
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
