package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

//Student is a struct for JSON data entry
type Student struct {
	ID    int    `json: id`
	Name  string `json:"name"`
	Range int    `json: range`
}

var people []Student

//GetPeople - getting the data from file data.json
func GetPeople() []Student {
	raw, err := ioutil.ReadFile("./data.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var data []Student
	json.Unmarshal(raw, &data)
	return data
}

//GetAll - is endpoint for GET request for People
func GetAll(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(people)
}

//GetByID - is endpoint for GET request for Student
func GetByID(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, item := range people {
		i, _ := strconv.Atoi(params["id"])
		if item.ID == i {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Student{})
}

//Create - is endpoint for POST request to add a Student
func Create(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var person Student
	_ = json.NewDecoder(req.Body).Decode(&person)
	i, _ := strconv.Atoi(params["id"])
	r, _ := strconv.Atoi(params["range"])
	person.ID = i
	person.Name = params["name"]
	person.Range = r
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

//Delete - is endpoint for DELETE request
func Delete(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, item := range people {
		i, _ := strconv.Atoi(params["id"])
		if item.ID == i {
			people = append(people[:index], people[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(people)
}

func main() {
	router := mux.NewRouter()
	people = GetPeople()
	router.HandleFunc("/people", GetAll).Methods("GET")
	router.HandleFunc("/people/{id}", GetByID).Methods("GET")
	router.HandleFunc("/people/{id}/{name}/{range}", Create).Methods("POST")
	router.HandleFunc("/people/{id}", Delete).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":5555", router))
}
