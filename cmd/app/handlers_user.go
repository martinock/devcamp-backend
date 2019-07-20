package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

// userGET a method to get user given userID params in URL
func (h *handler) GetUserByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	query := "SELECT * FROM users WHERE id = " + param.ByName("userID")
	rows, err := h.db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}

	var users []user

	for rows.Next() {
		user := user{}
		err := rows.Scan(
			&user.id,
			&user.name,
		)
		if err != nil {
			log.Println(err)
			continue
		}
		users = append(users, user)
	}

	bytes, err := json.Marshal(users)
	if err != nil {
		log.Println(err)
		return
	}

	renderJSON(w, bytes, http.StatusOK)
}

// InsertUser a function to insert user data (id, name) to DB
func (h *handler) InsertUser(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

}

// EditUserByID a function to change user data (name) in DB with given params (id, name)
func (h *handler) EditUserByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

}

// DeleteUserData a function to remove user data from DB given the userID
func (h *handler) DeleteUserByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

}

// GetMultipleUsers a function to get multiple users row in a single request
func (h *handler) GetMultipleUsers(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

}
