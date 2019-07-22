package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	var users []User

	for rows.Next() {
		user := User{}
		err := rows.Scan(
			&user.ID,
			&user.Name,
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
	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Println(err)
		return
	}

	// executing insert query
	query := "INSERT INTO users (id, name) VALUES ($1, $2)"
	_, err = h.db.Exec(query, user.ID, user.Name)
	if err != nil {
		log.Println(err)
		return
	}

	renderJSON(w, []byte(`
	{
		status: "success",
		message: "Insert user success!"
	}
	`), http.StatusOK)
}

// EditUserByID a function to change user data (name) in DB with given params (id, name)
func (h *handler) EditUserByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	userID := param.ByName("userID")
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
	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Println(err)
		return
	}

	query := "UPDATE users SET name = %s WHERE id = %s"
	_, err = h.db.Exec(fmt.Sprintf(query, user.Name, userID))
	if err != nil {
		log.Println(err)
		return
	}

	renderJSON(w, []byte(`
	{
		status: "success",
		message: "Update user success!"
	}
	`), http.StatusOK)
}

// DeleteUserData a function to remove user data from DB given the userID
func (h *handler) DeleteUserByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	userID := param.ByName("userID")

	query := "DELETE FROM users WHERE id = %s"
	_, err := h.db.Exec(fmt.Sprintf(query, userID))
	if err != nil {
		log.Println(err)
		return
	}

	renderJSON(w, []byte(`
	{
		status: "success",
		message: "Delete user success!"
	}
	`), http.StatusOK)
}

// GetMultipleUsers a function to get multiple users row in a single request
func (h *handler) GetMultipleUsers(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

}
