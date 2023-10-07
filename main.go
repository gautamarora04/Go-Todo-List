package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/thedevsaddam/renderer"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var rnd *renderer.Render
var db *mgo.Database

const (
	hostname       string = "localhost:27017"
	dbname         string = "demo_todo"
	collectionName string = "todo"
	port           string = ":9000"
)

type (
	todoModel struct {
		ID        bson.ObjectId `bson:"_id, omitempty"`
		Title     string        `bson:"title"`
		Completed bool          `bson:"completed"`
		CreatedAt time.Time     `bson:"cretedat"`
	}
	todo struct {
		ID        string    `json:"id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"`
		CreatedAt time.Time `json:"cretedat"`
	}
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rnd = renderer.New()
	session, err := mgo.Dial(hostname)
	checkErr(err)
	session.SetMode(mgo.Monotonic, true)
	db = session.DB(dbname)
}

func homehandler(w http.ResponseWriter, r *http.Request) {
	if err := rnd.Template(w, http.StatusOK, []string{"static/home.tpl"}, nil); err != nil {
		log.Fatal(err)
	}
}

func todohandlers() http.Handler {
	rg := chi.NewRouter() // new group router
	rg.Group(func(r chi.Router) {
		r.Get("/", fetchtodos)
		r.Post("/", createtodos)
		r.Put("/{id}", updatetodos)
		r.Delete("/{id}", deletetodos)
	})
	return rg
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", homehandler)
	r.Mount("/todo", todohandlers())

	server := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
		IdleTimeout:  time.Minute,
	}
	fmt.Println("Starting the server on the port: ", port)
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Server Error", err)
	}

}
