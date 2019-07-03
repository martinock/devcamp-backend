package main

import (
	"database/sql"
)

// args used for this application
type args struct {
	port int
}

// handler object used to handle the HTTP API
type handler struct {
	db              *sql.DB
	notFoundHandler *notFoundHandler
}
type notFoundHandler struct{}

// user struct for database query
type user struct {
	id   int
	name string
}
