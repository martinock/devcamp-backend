package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// GetBookByID a function to get a single book given it's ID
func (h *handler) GetBookByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

}

// InsertBook a function to insert book to DB
func (h *handler) InsertBook(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

}

// EditBook a function to change book data in DB, with given params
func (h *handler) EditBook(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

}

// DeleteBookByID a function to remove book data from DB, given bookID
func (h *handler) DeleteBookByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

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
