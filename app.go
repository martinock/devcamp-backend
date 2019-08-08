package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"

	"github.com/tokopedia/devcamp-backend/internal"

	kcache "github.com/koding/cache"
)

const (
	//TTL for books
	booksTTL = 5
)

func initFlags(args *internal.Args) {
	port := flag.Int("port", 3000, "port number for your apps")
	args.Port = *port
}

func initHandler(handler *internal.Handler) error {

	// Initialize SQL DB
	db, err := sql.Open("postgres", "postgres://postgres:postgres@127.0.0.1:5432/?sslmode=disable")
	if err != nil {
		return err
	}

	// create a cache with TTL
	cache := kcache.NewMemoryWithTTL(booksTTL * time.Second)
	// start garbage collection for expired keys
	cache.StartGC(time.Millisecond * 10)

	handler.DB = db
	handler.MemCache = cache

	return nil
}

func initRouter(router *httprouter.Router, handler *internal.Handler) {

	router.GET("/", handler.Index)

	// Single user API
	router.GET("/user/:userID", handler.GetUserByID)
	router.POST("/user", handler.InsertUser)
	router.PUT("/user/:userID", handler.EditUserByID)
	router.DELETE("/user/:userID", handler.DeleteUserByID)

	// Single book API
	router.GET("/book/:bookID", handler.GetBookByID)
	router.POST("/book", handler.InsertBook)
	router.PUT("/book/:bookID", handler.EditBook)
	router.DELETE("/book/:bookID", handler.DeleteBookByID)

	// Batch book API
	router.POST("/books", handler.InsertMultipleBooks)

	// Batch book API
	router.POST("/books/batch", handler.InsertMultipleBooksWithBatchingProcess)

	// Lending API
	router.POST("/lend", handler.LendBook)

	// `httprouter` library uses `ServeHTTP` method for it's 404 pages
	router.NotFound = handler
}

func main() {
	args := new(internal.Args)
	initFlags(args)

	handler := new(internal.Handler)
	if err := initHandler(handler); err != nil {
		panic(err)
	}

	router := httprouter.New()
	initRouter(router, handler)

	fmt.Printf("Apps served on :%d\n", args.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", args.Port), router))
}
