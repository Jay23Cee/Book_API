package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Post struct {
	Title string `bson:"title,omitempty"`
	Body  string `bson:"body,omitempty"`
}

type Book struct {
	Title  string `bson:"Title,omitempty"`
	Author string `bson:"Author,omitempty"`

	ID string `json:"ID,omitempty" bson:"_id,omitempty"`

	// CreatedAt time.Time
	// UpdatedAt time.Time
	// DeletedAt gorm.DeletedAt `gorm:"index"`
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

	//	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	mymap := make(map[int]Book)
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://mongo:LOsLH6a40mcR0QzB@cluster0.esomu.mongodb.net/?retryWrites=true&w=majority")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	col := client.Database("BookAPI").Collection("book")
	var results []Book

	// Get a MongoDB document using the FindOne() method
	cursor, err := col.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	count := 0
	//fmt.Fprintf(w, "collec %v", len(results))
	for _, result := range results {
		//	fmt.Println(result)
		//	fmt.Fprintf(w, "result: %v", result)
		fmt.Println(result.ID)

		mymap[count] = result

		count = count + 1

	}

	e, err := json.Marshal(mymap)
	if err != nil {
		panic(err)

	}
	w.Write([]byte(e))

}

func deletebook(w http.ResponseWriter, r *http.Request) {

	jsonMap := make(map[string]Book)
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal([]byte(body), &jsonMap)
	var temp Book
	temp = jsonMap["book"]
	fmt.Fprintf(w, " HERE AT DELETE %v,", temp.ID)

	clientOptions := options.Client().
		ApplyURI("mongodb+srv://mongo:LOsLH6a40mcR0QzB@cluster0.esomu.mongodb.net/?retryWrites=true&w=majority")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	collection := client.Database("BookAPI").Collection("book")
	fmt.Fprintf(w, " LOGGED IN %v", collection.Database())

	objid, err := primitive.ObjectIDFromHex(temp.ID)
	if err != nil {
		panic(err)
	}

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objid})
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "\nDeleted %v", result.DeletedCount)

}

func addbooks(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	jsonMap := make(map[string]Book)

	err = json.Unmarshal([]byte(body), &jsonMap)
	book := jsonMap["book"]

	fmt.Println(book)
	if err != nil {
		panic(err)
	}

	clientOptions := options.Client().
		ApplyURI("mongodb+srv://mongo:LOsLH6a40mcR0QzB@cluster0.esomu.mongodb.net/?retryWrites=true&w=majority")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	collection := client.Database("BookAPI").Collection("book")
	fmt.Fprintf(w, " LOGGED IN %v", collection.Database())

	doc := bson.D{{"Title", book.Title}, {"Author", book.Author}, {"_id", primitive.NewObjectID()}}
	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Inserted document with id: %v\n", result.InsertedID)
}

func CheckError(e error) {
	panic("unimplemented")
}

func editbook(w http.ResponseWriter, r *http.Request) {

	jsonMap := make(map[string]Book)
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal([]byte(body), &jsonMap)
	var book Book
	book = jsonMap["book"]
	fmt.Fprintf(w, " HERE AT EDIT %v,", book)

	clientOptions := options.Client().
		ApplyURI("mongodb+srv://mongo:LOsLH6a40mcR0QzB@cluster0.esomu.mongodb.net/?retryWrites=true&w=majority")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	collection := client.Database("BookAPI").Collection("book")
	fmt.Fprintf(w, " LOGGED IN %v", book.Title)

	id, _ := primitive.ObjectIDFromHex(book.ID)
	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{{"Title", book.Title}, {"Author", book.Author}}},
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)

}
func main() {

	// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://mongo:LOsLH6a40mcR0QzB@cluster0.esomu.mongodb.net/?retryWrites=true&w=majority"))

	// if err != nil {
	// 	panic(err)
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	// err = client.Connect((ctx))
	// if err != nil {
	// 	panic(err)
	// }

	// defer client.Disconnect(ctx)

	// err = client.Ping(ctx, readpref.Primary())
	// if err != nil {
	// 	panic(err)
	// }

	// databases, err := client.ListDatabaseNames(ctx, bson.M{})
	// if err != nil {
	// 	panic(err)

	//////////

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
