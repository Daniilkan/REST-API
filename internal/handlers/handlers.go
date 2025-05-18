package handlers

// Package handlers provides HTTP handler functions for the REST API.
// It includes handlers for CRUD operations on the "people" table and initializes the database connection for handlers.

import (
	"TestRest/pkg/postgres"
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"strconv"
)

var db *pgxpool.Pool
var ctx context.Context

func InitHandlers(database *pgxpool.Pool, context context.Context) {
	db = database
	ctx = context
}

// GetInfo retrieves a person's information by ID.
// @Summary Get person info
// @Description Get a person's details by their ID
// @Tags people
// @Accept json
// @Produce json
// @Param id query int true "Person ID"
// @Success 200 {object} postgres.Person
// @Failure 400 {string} string "Invalid ID parameter"
// @Failure 500 {string} string "Failed to get person"
// @Router /get [get]
func GetInfo(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	idStr := params.Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid id parameter"))
		return
	}
	person, err := postgres.GetPerson(ctx, db, id)
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
// @Description Delete a person by their ID
// @Tags people
// @Param id query int true "Person ID"
// @Success 200 {string} string "Deleted person by ID"
// @Failure 400 {string} string "Invalid ID parameter"
// @Failure 500 {string} string "Failed to delete person"
// @Router /delete [delete]
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	idStr := params.Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid id parameter"))
		return
	}

	err = postgres.DeletePerson(ctx, db, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to delete person"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deleted person by id: " + idStr))
}

// InsertPerson inserts a new person into the database.
// @Summary Insert person
// @Description Add a new person to the database
// @Tags people
// @Param name query string true "Person's name"
// @Param surname query string true "Person's surname"
// @Param patronymic query string false "Person's patronymic"
// @Success 200 {object} postgres.Person
// @Failure 500 {string} string "Failed to insert person"
// @Router /post [post]
func InsertPerson(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	name := params.Get("name")
	surname := params.Get("surname")
	patronymic := params.Get("patronymic")

	person, err := postgres.InsertPerson(ctx, db, name, surname, patronymic)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to insert person"))
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
// @Description Update a person's details by their ID
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
	params := r.URL.Query()

	idStr := params.Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid id parameter"))
		return
	}

	name := params.Get("name")
	surname := params.Get("surname")
	patronymic := params.Get("patronymic")

	person, err := postgres.UpdatePerson(ctx, db, id, name, surname, patronymic)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to update person"))
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
