package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Staring server!")
	registerHandlers(http.NewServeMux())
}

type Response struct {
	Message string `json:"message"`
}

type Greeter interface {
	Greet() string
}

type Person struct {
	Name string
}

func (p Person) Greet() string {
	return fmt.Sprintf("Hello, my name is %s", p.Name)
}

func registerHandlers(m *http.ServeMux) {

	m.HandleFunc("/getPersons", getPersonsHandler)
	m.HandleFunc("/getGreetings", getGreetingsHandler)
	http.ListenAndServe(":8080", m)
}

func getPersons() []Person {
	person1 := Person{Name: "John Doe"}
	person2 := Person{Name: "Jane Doe"}
	persons := []Person{person1, person2}
	return persons
}

func getPersonsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getPersons())
}

func getGreetingsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	persons := getPersons()
	greetings := make([]string, len(persons))
	for i, person := range persons {
		greetings[i] = person.Greet()
	}
	json.NewEncoder(w).Encode(greetings)
}
