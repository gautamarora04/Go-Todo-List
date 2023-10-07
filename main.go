package main

import (
	"encoding/json"
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
		ID        bson.ObjectId `bson:"_id,omitempty"`
		Title     string        `bson:"title"`
		Completed bool          `bson:"completed"`
		CreatedAt time.Time     `bson:"createAt"`
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

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if err := rnd.Template(w, http.StatusOK, []string{"static/home.tpl"}, nil); err != nil {
		log.Fatal(err)
	}
}

func fetchTodos(w http.ResponseWriter, r *http.Request) {
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
		"data": todolist,
	})

}

// func fetchTodos(w http.ResponseWriter, r *http.Request) {
// 	todos := []todoModel{}

// 	if err := db.C(collectionName).
// 		Find(bson.M{}).
// 		All(&todos); err != nil {
// 		rnd.JSON(w, http.StatusProcessing, renderer.M{
// 			"message": "Failed to fetch todo",
// 			"error":   err,
// 		})
// 		return
// 	}

// 	todoList := []todo{}
// 	for _, t := range todos {
// 		todoList = append(todoList, todo{
// 			ID:        t.ID.Hex(),
// 			Title:     t.Title,
// 			Completed: t.Completed,
// 			CreatedAt: t.CreatedAt,
// 		})
// 	}

// 	rnd.JSON(w, http.StatusOK, renderer.M{
// 		"data": todoList,
// 	})
// }

func createTodo(w http.ResponseWriter, r *http.Request) {
	var data todo
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Error while creating",
			"error":   err,
		})
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
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to save todo",
			"error":   err,
		})
		return
	}
	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "Todo Created Successfully!!",
		"todo_id": datamodel.ID.Hex(),
	})
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))

	if bson.IsObjectIdHex(id) == false {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "ID not valid!!",
		})
		return
	}

	if err := db.C(collectionName).RemoveId(bson.ObjectIdHex(id)); err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to delete Todo",
		})
		return
	}
	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "Todo Deleted Successfully !!",
	})

}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))

	if bson.IsObjectIdHex(id) == false {
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "ID not valid!!",
		})
		return
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

func todoHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Get("/", fetchTodos)
		r.Post("/", createTodo)
		r.Put("/{id}", updateTodo)
		r.Delete("/{id}", deleteTodo)
	})
	return rg
}
func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", homeHandler)

	r.Mount("/todo", todoHandlers())

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Listening on port ", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("listen: %s\n", err)
	}

}
