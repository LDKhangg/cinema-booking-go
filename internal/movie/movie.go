package movie

import (
	"context"
	"time"
)

type Movie struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ReleaseDate time.Time `json:"release_date"`
	ClosedDate  time.Time `json:"closed_date"`
	IsPublished bool      `json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Repository interface {
	GetPublishedMovies(ctx context.Context) ([]Movie, error)
	GetByID(ctx context.Context, id int) (Movie, error)
	Create(ctx context.Context, m *Movie) error
	Update(ctx context.Context, m *Movie) error
	Delete(ctx context.Context, id int) error
}
