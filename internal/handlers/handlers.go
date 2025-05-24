package handlers

// Package handlers provides HTTP handler functions for the REST API.
// It includes handlers for CRUD operations on the "people" table and initializes the database connection for handlers.

import (
	"TestRest/external"
	"TestRest/pkg/postgres"
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

var db *pgxpool.Pool
var ctx context.Context

func InitHandlers(database *pgxpool.Pool, context context.Context) {
	db = database
	ctx = context
}

// GetInfo retrieves a person's information by ID.
// @Summary Get person info
// @Description GetInfo Get a person's details by their ID
// @Tags people
// @Accept json
// @Produce json
// @Param id query int true "Person ID"
// @Success 200 {object} postgres.Person
// @Failure 400 {string} string "Invalid ID parameter"
// @Failure 500 {string} string "Failed to get person"
// @Router /get [get]
func GetInfo(w http.ResponseWriter, r *http.Request) {
	var params struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Surname     string `json:"surname"`
		Patronymic  string `json:"patronymic"`
		Age         int    `json:"age"`
		Gender      string `json:"gender"`
		Nationality string `json:"nationality"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	person, err := postgres.GetPerson(ctx, db, params.ID, params.Name, params.Surname, params.Patronymic, params.Age, params.Gender, params.Nationality)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to get person"))
		return
	}

	response, err := json.Marshal(person)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to process person data"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// DeletePerson deletes a person by ID.
// @Summary Delete person
// @Description DeletePerson Delete a person by their ID
// @Tags people
// @Param id query int true "Person ID"
// @Success 200 {string} string "Deleted person by ID"
// @Failure 400 {string} string "Invalid ID parameter"
// @Failure 500 {string} string "Failed to delete person"
// @Router /delete [delete]
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	var params struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err := postgres.GetPerson(ctx, db, params.ID, "", "", "", 0, "", "")
	if err != nil {
		if err.Error() == "no person found" {
			http.Error(w, "Person not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve person", http.StatusInternalServerError)
		return
	}

	err = postgres.DeletePerson(ctx, db, params.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to delete person"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deleted person by id: " + string(params.ID)))
}

// InsertPerson inserts a new person into the database.
// @Summary Insert person
// @Description InsertPerson Add a new person to the database
// @Tags people
// @Param name query string true "Person's name"
// @Param surname query string true "Person's surname"
// @Param patronymic query string false "Person's patronymic"
// @Success 200 {object} postgres.Person
// @Failure 500 {string} string "Failed to insert person"
// @Router /post [post]
func InsertPerson(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Name       string `json:"name"`
		Surname    string `json:"surname"`
		Patronymic string `json:"patronymic"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	name := params.Name
	surname := params.Surname
	patronymic := params.Patronymic
	age, err := external.GetAge(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to insert person - unable to fetch age"))
		return
	}
	gender, err := external.GetGender(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to insert person - unable to fetch gender"))
		return
	}
	nationality, err := external.GetNationality(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to insert person - unable to fetch nationality"))
		return
	}

	person, err := postgres.InsertPerson(ctx, db, name, surname, patronymic, age, gender, nationality)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to insert person - error in database"))
		return
	}

	response, err := json.Marshal(person)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to process person data"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// UpdatePerson updates an existing person's information.
// @Summary Update person
// @Description UpdatePerson Update a person's details by their ID
// @Tags people
// @Param id query int true "Person ID"
// @Param name query string true "Person's name"
// @Param surname query string true "Person's surname"
// @Param patronymic query string false "Person's patronymic"
// @Success 200 {object} postgres.Person
// @Failure 400 {string} string "Invalid ID parameter"
// @Failure 500 {string} string "Failed to update person"
// @Router /put [put]
func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	var params struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Surname     string `json:"surname"`
		Patronymic  string `json:"patronymic"`
		Age         int    `json:"age"`
		Gender      string `json:"gender"`
		Nationality string `json:"nationality"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id := params.ID

	persons, err := postgres.GetPerson(ctx, db, id, "", "", "", 0, "", "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to retrieve person"))
		return
	}

	if len(persons) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID or multiple persons found"))
		return
	}
	p := persons[0]

	name := params.Name
	if name == "" {
		name = p.Name
	}
	surname := params.Surname
	if surname == "" {
		surname = p.Surname
	}
	patronymic := params.Patronymic
	if patronymic == "" {
		patronymic = p.Patronymic
	}
	age := params.Age
	if age == 0 {
		age = p.Age
	}
	gender := params.Gender
	if gender == "" {
		gender = p.Gender
	}
	nationality := params.Nationality
	if nationality == "" {
		nationality = p.Nationality
	}

	updatedPerson, err := postgres.UpdatePerson(ctx, db, id, name, surname, patronymic, age, gender, nationality)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to update person"))
		return
	}

	response, err := json.Marshal(updatedPerson)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to process person data"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
