package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Book struct {
	Title  string
	Author string
}

func addbooks(w http.ResponseWriter, r *http.Request) {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Book{})

	// Create
	db.Create(&Book{Title: "D42", Author: "LEGEND MAKER"})

	// // Read
	// var book Book
	// db.First(&book, 1)                 // find book with integer primary key
	// db.First(&book, "code = ?", "D42") // find book with code D42

	// Delete - delete book
	//	db.Delete(&book, 1)

}

func editbook(w http.ResponseWriter, r *http.Request) {
	var book Book
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Book{})

	// Update - update book's price to 200
	db.Model(&book).Update("Price", 200)
	// Update - update multiple fields
	db.Model(&book).Updates(Book{Title: "D42", Author: "LEGEND MAKER"}) // non-zero fields
	db.Model(&book).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

}
func main() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	r.Get("/add", addbooks)
	r.Get("/edit", editbook)
	// RESTy routes for "articles" resource

	// Subrouters:

	// Mount the admin sub-router
	fmt.Print("ACTIVE")
	http.ListenAndServe(":3333", r)

}
