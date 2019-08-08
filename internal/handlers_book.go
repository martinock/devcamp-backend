package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// GetBookByID a function to get a single book given it's ID
func (h *Handler) GetBookByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	var books []Book
	var output []byte
	key := "Book_" + param.ByName("bookID")

	data, err := h.MemCache.Get(key)
	if err != nil {
		log.Printf("Got Error , err : %+v", err)
	}

	if data != nil {
		output, err = json.Marshal(data)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		query := fmt.Sprintf("SELECT id, title, author, isbn, stock FROM books WHERE id = %s", param.ByName("bookID"))
		rows, err := h.DB.Query(query)
		if err != nil {
			log.Println(err)
			return
		}

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
		output, err = json.Marshal(books)
		if err != nil {
			log.Println(err)
			return
		}

		h.MemCache.Set(key, books)
	}

	renderJSON(w, output, http.StatusOK)
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
	rows := strings.Split(contents, "\n")
	for i, row := range rows {
		// skip title
		if i == 0 {
			continue
		}
		// Split rows to column
		columns := strings.Split(row, ",")
		query := "INSERT INTO books (id, title, author, isbn, stock) VALUES ($1, $2, $3, $4, $5)"
		_, err = h.DB.Query(query, columns[0], columns[1], columns[2], columns[3], columns[4])
		if err != nil {
			log.Println(err)
			continue
		}

	}

	buffer.Reset()

	renderJSON(w, []byte(`
	{
		status: "success",
		message: "Insert book success!"
	}
	`), http.StatusOK)
}

// InsertMultipleBooksWithBatchingProcess a function to insert multiple book data, given file of books data
func (h *Handler) InsertMultipleBooksWithBatchingProcess(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
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
	rows := strings.Split(contents, "\n")
	totaltask := len(rows)
	maxBatch := 2
	for start := 0; start < totaltask; start += maxBatch {
		end := start + maxBatch

		if end > totaltask {
			end = totaltask
		}

		err := h.insertBatchProcess(rows[start:end])
		if err != nil {
			log.Println("err batch : ", err)
		}
	}

	buffer.Reset()

	renderJSON(w, []byte(`
	{
		status: "success",
		message: "Insert book success!"
	}
	`), http.StatusOK)
}

func (h *Handler) insertBatchProcess(data []string) (err error) {
	sqlStr := "INSERT INTO books(id, title, author, isbn, stock) VALUES "
	vals := []interface{}{}
	val := "(?, ?, ?, ?, ?)"

	var values []string

	for _, row := range data {
		values = append(values, val)
		columns := strings.Split(row, ",")
		vals = append(vals, columns[0], columns[1], columns[2], columns[3], columns[4])
		log.Println("data added : ", row)
	}
	//trim the last ,
	sqlStr = sqlStr + strings.Join(values, ",")
	log.Println(sqlStr)
	log.Println(vals)
	// prepare the statement
	sqlStr = h.ReplaceSQL(sqlStr, "?")
	stmt, err := h.DB.Prepare(sqlStr)
	if err != nil {
		log.Println("err prepare")
		return
	}

	//format all vals at once
	res, err := stmt.Exec(vals...)
	if err != nil {
		log.Println("err exec")
		return
	}

	log.Println(res)

	return
}

func (h *Handler) ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
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
