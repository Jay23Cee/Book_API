package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Book struct {
	Title  string `gorm:"Title"`
	Author string

	ID  uint `gorm:"primaryKey"`
	Key uint

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Thing struct {
	Value string
	book  Book
}

func readBooks(w http.ResponseWriter, r *http.Request) {
	var book Book
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.First(&book, 1)                  // find book with integer primary key
	db.First(&book, "Title = ?", "D42") // find book with code D42
	// fmt.Fprintf(w, " \n READING %v", book)
	// w.Write([]byte(book.Title))

}

func getBooks(w http.ResponseWriter, r *http.Request) {
	var book Book
	mymap := make(map[int]Book)
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	rows, err := db.Model(&Book{}).Rows()

	defer rows.Close()
	if err != nil {
		panic(err)
	}
	count := 0
	for rows.Next() {

		db.ScanRows(rows, &book)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		mymap[count] = book

		count = count + 1

	}
	e, err := json.Marshal(mymap)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write([]byte(e))

	// db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	// if err != nil {
	// 	panic("failed to connect database")
	// }

	// db.First(&book, 1)                  // find book with integer primary key
	// db.First(&book, "Title = ?", "D42") // find book with code D42
	// w.Write([]byte(book.Title))

}

func deletebook(w http.ResponseWriter, r *http.Request) {

	jsonMap := make(map[string]Book)
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal([]byte(body), &jsonMap)
	var temp Book
	temp = jsonMap["book"]
	fmt.Fprintf(w, " HERE AT DELETE %v,", temp.ID)

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Update - update book's price to 200

	// Update - update multiple fields

	fmt.Fprintf(w, "UPDATED ---%v", db)

	db.Delete(&temp)
}
func addbooks(w http.ResponseWriter, r *http.Request) {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Book{})

	body, err := ioutil.ReadAll(r.Body)
	// var d struct{ Result map[string][]Book }

	fmt.Fprintf(w, " bODY %v", body)

	var t Book

	err = json.Unmarshal(body, &t)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, " \n Unmarshall t %v", t)

	jsonMap := make(map[string]Book)

	err = json.Unmarshal([]byte(body), &jsonMap)
	book := jsonMap["book"]
	fmt.Fprintf(w, " \n JSON MAP Unmarshall t %v", jsonMap["book"])

	db.Create(&book)
	w.Write([]byte("Adding Book"))
	fmt.Fprintf(w, " \n BOOK ADDED %v", book)

}

func editbook(w http.ResponseWriter, r *http.Request) {

	jsonMap := make(map[string]Book)
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal([]byte(body), &jsonMap)
	var temp Book
	temp = jsonMap["book"]
	fmt.Fprintf(w, " HERE AT EDIT %v,", temp)

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Update - update book's price to 200

	// Update - update multiple fields
	db.Model(&temp).Where("id = ?", temp.ID).Updates(map[string]interface{}{"Title": temp.Title, "Author": temp.Author})
	fmt.Fprintf(w, "UPDATED ---%v", db)

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

	r.Post("/add", addbooks)
	r.Post("/edit", editbook)
	r.Get("/read", getBooks)
	r.Post("/delete", deletebook)
	// RESTy routes for "articles" resource

	// Subrouters:

	// Mount the admin sub-router
	fmt.Print("ACTIVE")
	http.ListenAndServe(":3333", r)

}
