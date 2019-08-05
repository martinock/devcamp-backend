package internal

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// GetUserByID a method to get user given userID params in URL
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	// TODO: implement this. Query = SELECT * FROM users WHERE id = <userID>
}

// InsertUser a function to insert user data (id, name) to DB
func (h *Handler) InsertUser(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	// TODO: implement this. Query = INSERT INTO users (id, name) VALUES (<userID>, '<name>')
	// read json body

	// parse json body

	// executing insert query
}

// EditUserByID a function to change user data (name) in DB with given params (id, name)
func (h *Handler) EditUserByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	// TODO: implement this. Query = UPDATE users SET name = '<name>' WHERE id = <userID>
	// read json body

	// parse json body
}

// DeleteUserByID a function to remove user data from DB given the userID
func (h *Handler) DeleteUserByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	// TODO: implement this. Query = DELETE FROM users WHERE id = <userID>
}
