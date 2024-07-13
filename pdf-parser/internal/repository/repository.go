package repository

import (
	"context"

	"pdf-parser/types"
)

type Repository interface {
	GetEnrollmentByURL(ctx context.Context, url string) (*types.Enrollment, error)
	GetEnrollmentByID(ctx context.Context, id string) (*types.Enrollment, error)
	ListEnrollments(ctx context.Context) ([]types.Enrollment, error)
	SaveEnrollment(ctx context.Context, enrollment *types.Enrollment) error
}

type Databases struct {
	sql *SQLite
}

func New() (*Databases, error) {
	sql, err := NewSQLite()
	if err != nil {
		return nil, err
	}

	return &Databases{
		sql: sql,
	}, nil
}
