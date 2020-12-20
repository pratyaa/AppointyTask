package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
//	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

type Person struct {
	id                  int                `json:"id,omitempty"`
	name               string              `json:"name,omitempty" `
	dob                string              `json:"dob,omitempty" `
	phoneNumber		   string			   `json:"phoneNumber,omitempty" `
    email              string              `json:"email,omitempty" `
	CreationTimestamp  time.Time           `json:"timestamp,omitempty" `
	
}

type Contact struct{
	id1      int       `json:"id1,omitempty"`  
	id2      int       `json:"id2,omitempty"`
    contact  string    `json:"contact,omitempty"`
}

func createUserEndpoint(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("content-type", "application/json")
	var person Person
	_ = json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database("TracingContactDB").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}

 func getUserEndpoint(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("content-type", "application/json")
	id, ok := request.URL.Query()["id"]
	var person Person
	collection := client.Database("TracingContactDB").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Person{ID: id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}	
	json.NewEncoder(response).Encode(person)

 }


func createContact(response http.ResponseWriter, request *http.Request) {

	response.Header().Set("content-type", "application/json")
	var contact Contact
	_ = json.NewDecoder(request.Body).Decode(&contact)
	collection := client.Database("TracingContactDB").Collection("contact")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, contact)
	json.NewEncoder(response).Encode(result)
}



func handleRequests() {
    http.HandleFunc("/users", createUserEndpoint)
	http.HandleFunc("/users/<id>", getUserEndpoint)
	http.HandleFunc("/contact", createContact)
//	http.HandleFunc("/contacts?user=<user id>&infection_timestamp=<timestamp>")

   log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {

	// Establishing connection to MongoDB
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)

	handleRequests()
}