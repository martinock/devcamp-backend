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

// Book struct for database query
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	ISBN   string `json:"isbn"`
	Stock  int64  `json:"stock"`
}

// LendRequest struct for receiving lend http request
type LendRequest struct {
	UserID int `json:"user_id"`
}
