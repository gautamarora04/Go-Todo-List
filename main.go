// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-chi/chi/v5/middleware"
// 	"github.com/thedevsaddam/renderer"
// 	"gopkg.in/mgo.v2"
// 	"gopkg.in/mgo.v2/bson"
// )

// var rnd *renderer.Render
// var db *mgo.Database

// const (
// 	hostname       string = "localhost:27017"
// 	dbname         string = "demo_todo"
// 	collectionName string = "todo"
// 	port           string = ":9000"
// )

// type (
// 	todoModel struct {
// 		ID        bson.ObjectId `bson:"_id, omitempty"`
// 		Title     string        `bson:"title"`
// 		Completed bool          `bson:"completed"`
// 		CreatedAt time.Time     `bson:"cretedat"`
// 	}
// 	todo struct {
// 		ID        string    `json:"id"`
// 		Title     string    `json:"title"`
// 		Completed bool      `json:"completed"`
// 		CreatedAt time.Time `json:"cretedat"`
// 	}
// )

// func checkErr(err error) {
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func init() {
// 	rnd = renderer.New()
// 	session, err := mgo.Dial(hostname)
// 	checkErr(err)
// 	session.SetMode(mgo.Monotonic, true)
// 	db = session.DB(dbname)
// }

// func todohandlers() http.Handler {
// 	rg := chi.NewRouter() // new group router
// 	rg.Group(func(r chi.Router) {
// 		r.Get("/", fetchtodos)
// 		r.Post("/", createtodos)
// 		r.Put("/{id}", updatetodos)
// 		r.Delete("/{id}", deletetodos)
// 	})
// 	return rg
// }

// func main() {
// 	r := chi.NewRouter()
// 	r.Use(middleware.Logger)
// 	r.Get("/", homehandler)
// 	r.Mount("/todo", todohandlers())

// 	server := &http.Server{
// 		Addr:         port,
// 		Handler:      r,
// 		ReadTimeout:  time.Minute,
// 		WriteTimeout: time.Minute,
// 		IdleTimeout:  time.Minute,
// 	}
// 	if err := server.ListenAndServe(); err != nil {
// 		fmt.Println("Server Error", err)
// 	}

// }
package main

import (
	"fmt"
	"time"
)

func main() {
	// Without goroutine
	// printNumbers(3)

	// With goroutine
	go printNumbers(3)

	// Give the goroutine some time to execute before the program exits
	time.Sleep(time.Second)
}

func printNumbers(n int) {
	for i := 1; i <= n; i++ {
		fmt.Println(i)
		time.Sleep(time.Second)
	}
}
