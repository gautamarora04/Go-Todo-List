package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
		CreatedAt time.Time     `bson:"createdAt"`
	}
	todo struct {
		ID        string    `json:"id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"`
		CreatedAt time.Time `json:"created_at"`
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

func fetchtodos(w http.ResponseWriter, r *http.Request) {
	todos := []todoModel{}
	if err := db.C(collectionName).Find(bson.M{}).All(&todos); err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch todo",
			"Error":   err,
		})
		return
	}
	todolist := []todo{}

	for _, t := range todos {
		todolist = append(todolist, todo{
			ID:        t.ID.Hex(),
			Title:     t.Title,
			Completed: t.Completed,
			CreatedAt: t.CreatedAt,
		})
	}
	rnd.JSON(w, http.StatusOK, renderer.M{
		"Todolist": todolist,
	})

}

func createtodos(w http.ResponseWriter, r *http.Request) {
	var data todo
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		rnd.JSON(w, http.StatusProcessing, err)
		return
	}
	if data.Title == "" {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "Please enter the title",
		})
		return
	}
	datamodel := todoModel{
		ID:        bson.NewObjectId(),
		Title:     data.Title,
		Completed: false,
		CreatedAt: time.Now(),
	}
	if err := db.C(collectionName).Insert(&datamodel); err != nil {
		rnd.JSON(w, http.StatusBadRequest, err)
		return
	}
	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "Todo Created Successfully!!",
		"todo_id": datamodel.ID.Hex(),
	})
}

func deletetodos(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))

	if bson.IsObjectIdHex(id) == false {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "ID not valid!!",
		})
		return
	}

	if err := db.C(collectionName).RemoveId(id); err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to delete Todo",
		})
		return
	}
	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "Todo Deleted Successfully !!",
	})

}

func updatetodos(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))

	if bson.IsObjectIdHex(id) == false {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "ID not valid!!",
		})
	}

	var data todo

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		rnd.JSON(w, http.StatusProcessing, err)
		return
	}

	if err := db.C(collectionName).Update(
		bson.M{"_id": bson.ObjectIdHex(id)},
		bson.M{"title": data.Title, "completed": data.Completed},
	); err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to Update",
		})
		return
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
