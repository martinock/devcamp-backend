package internal

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// GetBookByID a function to get a single book given it's ID
func (h *Handler) GetBookByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	// TODO: Implement this. Query = SELECT id, title, author, isbn, stock FROM books WHERE id = <bookID>
}

// InsertBook a function to insert book to DB
func (h *Handler) InsertBook(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	// TODO: Implement this. Query = INSERT INTO books (id, title, author, isbn, stock) VALUES (<id>, '<title>', '<author>', '<isbn>', <stock>)
	// read json body

	// parse json body

	// executing insert query

}

// EditBook a function to change book data in DB, with given params
func (h *Handler) EditBook(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	// TODO: Implement this. Query = UPDATE books SET title = '<title>', author = '<author>', isbn = '<isbn>', stock = <stock> WHERE id = <id>
	// read json body

	// parse json body

	// executing update query
}

// DeleteBookByID a function to remove book data from DB, given bookID
func (h *Handler) DeleteBookByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	// TODO: implement this. Query = DELETE FROM books WHERE id = <id>
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

	// TODO: uncomment this when implementing
	// contents := buffer.String()

	// Split contents to rows
	// TODO: uncomment this when implementing
	// rows := strings.Split(contents, "\n")

	// TODO: iterate csv rows here.

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
	// TODO: implement this.
	// Get stock query = SELECT stock FROM books WHERE id = <bookID>
	// Insert Book Lending query = INSERT INTO lend (user_id, book_id) VALUES (<userID>, <bookID>)
	// Update stock query = UPDATE books SET stock = <newStock> WHERE id = <bookID>

	// Read userID

	// parse json body

	// Get book stock from DB

	// Insert Book to Lend tables

	// Update Book stock query
}
