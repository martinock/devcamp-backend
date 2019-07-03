package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

func (h *handler) userGET(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

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

func (h *handler) userPOST(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

}

func (h *handler) userPUT(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

}

func (h *handler) userDELETE(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

}

func (h *handler) usersGET(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

}
