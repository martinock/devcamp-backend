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
	router.GET("/user/:userID", handler.userGET)
	router.POST("/user", handler.userPOST)
	router.PUT("/user/:userID", handler.userPUT)
	router.DELETE("/user/:userID", handler.userDELETE)

	// Batch user API
	router.GET("/users", handler.usersGET)

	// Single book API
	router.GET("/book/:bookID", handler.bookGET)
	router.POST("/book", handler.bookPOST)
	router.PUT("/book/:bookID", handler.bookPUT)
	router.DELETE("/book/:bookID", handler.bookDELETE)

	// Batch book API
	router.GET("/books", handler.booksGET)

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
