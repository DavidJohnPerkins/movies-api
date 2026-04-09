package store

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	ID          uuid.UUID `db:"Id"`
	Title       string
	Director    string
	ReleaseDate time.Time
	TicketPrice float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
type MovieJson struct {
	ID          uuid.UUID `db:"Id"`
	Title       string
	Director    string
	ReleaseDate time.Time
	TicketPrice float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Interface interface {
	GetAll(ctx context.Context) ([]Movie, error)
	GetByID(ctx context.Context, id uuid.UUID) (Movie, error)
	Create(ctx context.Context, jsonString string) error
	Update(ctx context.Context, jsonString string) error
	Delete(ctx context.Context, jsonString string) error
}
