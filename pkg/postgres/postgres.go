package postgres

import (
	"TestRest/pkg/logger"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"strings"
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

type Person struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Nationality string `json:"nationality"`
	Gender      string `json:"gender"`
}

func InsertPerson(ctx context.Context, db *pgxpool.Pool, name string, surname string, patronymic string, age int, gender string, nationality string) (*Person, error) {
	var person Person
	query := `
		INSERT INTO people (name, surname, patronymic, age, nationality, gender)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, surname, patronymic, age, nationality, gender
	`
	err := db.QueryRow(ctx, query, name, surname, patronymic, age, nationality, gender).Scan(
		&person.ID, &person.Name, &person.Surname, &person.Patronymic, &person.Age, &person.Nationality, &person.Gender,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert and retrieve person: %w", err)
	}

	logger.GetLoggerFromContext(ctx).Info(ctx, "Inserted person", zap.Int("id", person.ID), zap.String("name", person.Name), zap.String("surname", person.Surname), zap.String("patronymic", person.Patronymic), zap.Int("age", age), zap.String("gender", gender), zap.String("nationality", nationality))
	return &person, nil
}

func GetPerson(ctx context.Context, db *pgxpool.Pool, id int, name, surname, patronymic string, age int, gender, nationality string) ([]Person, error) {
	query := `
		SELECT id, name, surname, patronymic, age, nationality, gender
		FROM people
	`
	conditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if id != 0 {
		conditions = append(conditions, fmt.Sprintf("id = $%d", argIndex))
		args = append(args, id)
		argIndex++
	}
	if name != "" {
		conditions = append(conditions, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, name)
		argIndex++
	}
	if surname != "" {
		conditions = append(conditions, fmt.Sprintf("surname = $%d", argIndex))
		args = append(args, surname)
		argIndex++
	}
	if patronymic != "" {
		conditions = append(conditions, fmt.Sprintf("patronymic = $%d", argIndex))
		args = append(args, patronymic)
		argIndex++
	}
	if age != 0 {
		conditions = append(conditions, fmt.Sprintf("age = $%d", argIndex))
		args = append(args, age)
		argIndex++
	}
	if gender != "" {
		conditions = append(conditions, fmt.Sprintf("gender = $%d", argIndex))
		args = append(args, gender)
		argIndex++
	}
	if nationality != "" {
		conditions = append(conditions, fmt.Sprintf("nationality = $%d", argIndex))
		args = append(args, nationality)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve persons: %w", err)
	}
	defer rows.Close()

	var persons []Person
	for rows.Next() {
		var p Person
		if err := rows.Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Nationality, &p.Gender); err != nil {
			return nil, fmt.Errorf("failed to scan person: %w", err)
		}
		persons = append(persons, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	logger.GetLoggerFromContext(ctx).Info(ctx, "Retrieved persons", zap.Int("count", len(persons)), zap.Any("persons", persons))
	return persons, nil
}

func DeletePerson(ctx context.Context, db *pgxpool.Pool, id int) error {
	query := `
		DELETE FROM people
		WHERE id = $1
	`
	_, err := db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete person: %w", err)
	}
	logger.GetLoggerFromContext(ctx).Info(ctx, "Deleted person", zap.Int("id", id))
	return nil
}

func UpdatePerson(ctx context.Context, db *pgxpool.Pool, id int, name, surname, patronymic string, age int, gender string, nationality string) (*Person, error) {
	query := `
		UPDATE people
		SET name = $1, surname = $2, patronymic = $3, age = $4, gender = $5, nationality = $6
		WHERE id = $7
		RETURNING id, name, surname, patronymic, age, nationality, gender
	`
	var p Person
	err := db.QueryRow(ctx, query, name, surname, patronymic, age, gender, nationality, id).Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Nationality, &p.Gender)
	if err != nil {
		return nil, fmt.Errorf("failed to update person: %w", err)
	}
	logger.GetLoggerFromContext(ctx).Info(ctx, "Updated person", zap.Int("id", p.ID), zap.String("name", p.Name), zap.String("surname", p.Surname), zap.String("patronymic", p.Patronymic), zap.Int("age", age), zap.String("gender", gender), zap.String("nationality", nationality))
	return &p, nil
}
