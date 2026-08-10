package showtime

import (
	"context"
	"time"
)

type Showtime struct {
	ID        int       `json:"id"`
	MovieID   int       `json:"movie_id"`
	RoomID    int       `json:"room_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Repository interface {
	Create(ctx context.Context, s *Showtime) error
	GetByID(ctx context.Context, id int) (Showtime, error)

	GetByMovieID(ctx context.Context, movieID int, targetDate time.Time) ([]Showtime, error)
	GetByRoomID(ctx context.Context, roomID int, targetDate time.Time) ([]Showtime, error)
}

type Service interface {
	CreateShowtime(ctx context.Context, s *Showtime) error
	GetShowtimeByID(ctx context.Context, id int) (Showtime, error)

	GetShowtimesByMovie(ctx context.Context, movieID int, date time.Time) ([]Showtime, error)
	GetShowtimesByTheater(ctx context.Context, theaterID int, date time.Time) ([]Showtime, error)
}
