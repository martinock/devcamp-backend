package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func initFlags(args *args) {
	port := flag.Int("port", 3000, "port number for your apps")
	args.port = *port
}

func initHandler(handler *handler) error {

	// Initialize SQL DB
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/?sslmode=disable")
	if err != nil {
		return err
	}
	handler.db = db

	// Handler for 404 pages
	handler.notFoundHandler = &notFoundHandler{}

	return nil
}

func initRouter(router *httprouter.Router, handler *handler) {

	router.GET("/", handler.index)

	// Single user API
	router.GET("/user/:userID", handler.GetUserByID)
	router.POST("/user", handler.InsertUser)
	router.PUT("/user/:userID", handler.EditUserByID)
	router.DELETE("/user/:userID", handler.DeleteUserByID)

	// Batch user API
	router.GET("/users", handler.GetMultipleUsers)

	// Single book API
	router.GET("/book/:bookID", handler.GetBookByID)
	router.POST("/book", handler.InsertBook)
	router.PUT("/book/:bookID", handler.EditBook)
	router.DELETE("/book/:bookID", handler.DeleteBookByID)

	// Batch book API
	router.POST("/books", handler.InsertMultipleBooks)

	// Lending API
	router.POST("/lend/:bookID", handler.LendBook)

	router.NotFound = handler.notFoundHandler
}

func main() {
	args := new(args)
	initFlags(args)

	handler := new(handler)
	if err := initHandler(handler); err != nil {
		panic(err)
	}

	router := httprouter.New()
	initRouter(router, handler)

	fmt.Printf("Apps served on :%d\n", args.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", args.port), router))
}
