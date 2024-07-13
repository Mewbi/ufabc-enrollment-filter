package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"pdf-parser/config"
	"pdf-parser/types"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	conn *sql.DB
}

func NewSQLite() (*SQLite, error) {
	conf := config.Get()
	db, err := sql.Open(
		"sqlite3",
		fmt.Sprintf("%s%s?_foreign_keys=on&cache=%s", conf.Database.Type, conf.Database.Address, conf.Database.Cache),
	)
	if err != nil {
		return nil, err
	}

	// Set the maximum number of open connections
	db.SetMaxOpenConns(conf.Database.MaxConn)

	// Ping to check if the database connection is established
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	repo := &SQLite{
		conn: db,
	}

	err = repo.migrate(conf.Database.Schema)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (s *SQLite) migrate(filepath string) error {
	// Read the schema file
	schema, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	// Execute the SQL statements from the schema file
	_, err = s.conn.Exec(string(schema))
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLite) SaveEnrollment(ctx context.Context, enrollment *types.Enrollment) error {
	sql := `INSERT INTO enrollments (id, name, url, content, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := s.conn.
		ExecContext(ctx, sql,
			enrollment.ID,
			enrollment.Name,
			enrollment.URL,
			enrollment.Content,
			enrollment.CreatedAt,
		)
	return err
}

func (s *SQLite) GetEnrollmentByURL(ctx context.Context, url string) (*types.Enrollment, error) {
	var enrollment types.Enrollment
	sqlQuery := `SELECT id, name, url, content, created_at FROM enrollments WHERE url = $1`
	err := s.conn.QueryRowContext(ctx, sqlQuery, url).Scan(
		&enrollment.ID,
		&enrollment.Name,
		&enrollment.URL,
		&enrollment.Content,
		&enrollment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &enrollment, nil
}

func (s *SQLite) GetEnrollmentByID(ctx context.Context, id string) (*types.Enrollment, error) {
	var enrollment types.Enrollment
	sqlQuery := `SELECT id, name, url, content, created_at FROM enrollments WHERE id = $1`
	err := s.conn.QueryRowContext(ctx, sqlQuery, id).Scan(
		&enrollment.ID,
		&enrollment.Name,
		&enrollment.URL,
		&enrollment.Content,
		&enrollment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &enrollment, nil
}

func (s *SQLite) ListEnrollments(ctx context.Context) ([]types.Enrollment, error) {
	enrollments := make([]types.Enrollment, 0)

	sqlQuery := `SELECT id, name, created_at FROM enrollments`
	rows, err := s.conn.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var enrollment types.Enrollment
		err = rows.Scan(
			&enrollment.ID,
			&enrollment.Name,
			&enrollment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		enrollments = append(enrollments, enrollment)
	}
	return enrollments, nil
}
