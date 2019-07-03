package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

// User a simple user struct for database query
type User struct {
	id   int
	name string
}

func home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Welcome!\n")
}

func hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	getUser(w, ps.ByName("userID"))
}

func getUser(w http.ResponseWriter, userID string) {
	conn := "postgres://postgres:postgres@db:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	queryUser := fmt.Sprintf("SELECT * FROM users WHERE id = %s", userID)

	rows, dbError := db.Query(queryUser)
	if dbError != nil {
		log.Fatal(dbError)
		return
	}

	count := 0
	for rows.Next() {
		user := User{}
		err := rows.Scan(
			&user.id,
			&user.name,
		)

		fmt.Fprintf(w, "Hello, %s!\n", user.name)

		if err != nil {
			log.Fatal(err)
			return
		}
		count++
	}
	fmt.Fprintf(w, "%d data found\n", count)
}

func main() {
	portPtr := flag.Int("port", 3000, "port number for your apps")

	router := httprouter.New()
	router.GET("/", home)
	router.GET("/hello/:userID", hello)

	fmt.Printf("Apps served on :%d\n", *portPtr)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *portPtr), router))
}
