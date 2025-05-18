package postgres

// Package postgres provides database connection management and CRUD operations for the "people" table.
// It includes functions for inserting, retrieving, updating, and deleting records.

import (
	"TestRest/external"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-default:"localhost"`
	Port     uint16 `yaml:"port" env:"POSTGRES_PORT" env-default:"5432"`
	Username string `yaml:"username" env:"USERNAME" env-default:"root"`
	Password string `yaml:"password" env:"PASSWORD" env-default:"qwerty"`
	Database string `yaml:"database" env:"DATABASE" env-default:"postgres"`

	MinConns int32 `yaml:"min_conns" env:"MIN_CONNS" env-default:"5"`
	MaxConns int32 `yaml:"max_conns" env:"MAX_CONNS" env-default:"10"`
}

// New creates a new PostgreSQL connection pool.
// @Summary Create PostgreSQL connection pool
// @Description New initializes a new PostgreSQL connection pool using the provided configuration.
// @Tags postgres
// @Param config body postgres.Config true "PostgreSQL configuration"
// @Success 200 {object} pgxpool.Pool "Database connection pool"
// @Failure 500 {string} string "Failed to create connection pool"
// @Router /postgres/new [post]
func New(ctx context.Context, config Config) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse pool config: %w", err)
	}

	poolConfig.MinConns = config.MinConns
	poolConfig.MaxConns = config.MaxConns

	conn, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %w", err)
	}

	return conn, nil
}

// Person represents a person entity in the database.
// @Description Person entity
// @ID Person
// @Param id int true "Person ID"
// @Param name string true "Person's name"
// @Param surname string true "Person's surname"
// @Param patronymic string false "Person's patronymic"
// @Param age int true "Person's age"
// @Param nationality string true "Person's nationality"
// @Param gender string true "Person's gender"
type Person struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Nationality string `json:"nationality"`
	Gender      string `json:"gender"`
}

// InsertPerson inserts a new person into the database.
// @Summary Insert person
// @Description InsertPerson adds a new person to the database with details fetched from external APIs.
// @Tags postgres
// @Param name query string true "Person's name"
// @Param surname query string true "Person's surname"
// @Param patronymic query string false "Person's patronymic"
// @Success 200 {object} postgres.Person "Inserted person"
// @Failure 500 {string} string "Failed to insert person"
// @Router /postgres/person [post]
func InsertPerson(ctx context.Context, db *pgxpool.Pool, name string, surname string, patronymic string) (*Person, error) {
	age, err := external.GetAge(name)
	if err != nil {
		return nil, err
	}

	gender, err := external.GetGender(name)
	if err != nil {
		return nil, err
	}

	nationality, err := external.GetNationality(name)
	if err != nil {
		return nil, err
	}

	var person Person
	query := `
		INSERT INTO people (name, surname, patronymic, age, nationality, gender)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, surname, patronymic, age, nationality, gender
	`
	err = db.QueryRow(ctx, query, name, surname, patronymic, age, nationality, gender).Scan(
		&person.ID, &person.Name, &person.Surname, &person.Patronymic, &person.Age, &person.Nationality, &person.Gender,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert and retrieve person: %w", err)
	}

	return &person, nil
}

// GetPerson retrieves a person's details by ID.
// @Summary Get person
// @Description GetPerson fetches a person's details from the database by their ID.
// @Tags postgres
// @Param id query int true "Person ID"
// @Success 200 {object} postgres.Person "Person details"
// @Failure 500 {string} string "Failed to retrieve person"
// @Router /postgres/person [get]
func GetPerson(ctx context.Context, db *pgxpool.Pool, id int) (*Person, error) {
	query := `
		SELECT id, name, surname, patronymic, age, nationality, gender
		FROM people
		WHERE id = $1
	`
	var p Person
	err := db.QueryRow(ctx, query, id).Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Nationality, &p.Gender)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve person: %w", err)
	}
	return &p, nil
}

// DeletePerson deletes a person from the database by ID.
// @Summary Delete person
// @Description DeletePerson removes a person from the database by their ID.
// @Tags postgres
// @Param id query int true "Person ID"
// @Success 200 {string} string "Person deleted"
// @Failure 500 {string} string "Failed to delete person"
// @Router /postgres/person [delete]
func DeletePerson(ctx context.Context, db *pgxpool.Pool, id int) error {
	query := `
		DELETE FROM people
		WHERE id = $1
	`
	_, err := db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete person: %w", err)
	}
	return nil
}

// UpdatePerson updates a person's details in the database.
// @Summary Update person
// @Description UpdatePerson updates a person's details in the database by their ID.
// @Tags postgres
// @Param id query int true "Person ID"
// @Param name query string true "Person's name"
// @Param surname query string true "Person's surname"
// @Param patronymic query string false "Person's patronymic"
// @Success 200 {object} postgres.Person "Updated person"
// @Failure 500 {string} string "Failed to update person"
// @Router /postgres/person [put]
func UpdatePerson(ctx context.Context, db *pgxpool.Pool, id int, name, surname, patronymic string) (*Person, error) {
	query := `
		UPDATE people
		SET name = $1, surname = $2, patronymic = $3
		WHERE id = $4
		RETURNING id, name, surname, patronymic, age, nationality, gender
	`
	var p Person
	err := db.QueryRow(ctx, query, name, surname, patronymic, id).Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Nationality, &p.Gender)
	if err != nil {
		return nil, fmt.Errorf("failed to update person: %w", err)
	}
	return &p, nil
}
