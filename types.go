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

// User struct for database query
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
