package showtime

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var ErrShowtimeNotFound = errors.New("cannot find showtime")

type postgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, s *Showtime) error {
	query := `
		INSERT INTO showtimes (movie_id, room_id, start_time, end_time, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	err := r.db.QueryRowContext(
		ctx, query,
		s.MovieID, s.RoomID, s.StartTime, s.EndTime, s.Price, s.CreatedAt, s.UpdatedAt,
	).Scan(&s.ID)
	if err != nil {
		return fmt.Errorf("cannot create showtime: %w", err)
	}
	return nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int) (Showtime, error) {
	query := `
		SELECT id, movie_id, room_id, start_time, end_time, price, created_at, updated_at
		FROM showtimes
		WHERE id = $1
	`
	var s Showtime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.MovieID, &s.RoomID, &s.StartTime, &s.EndTime, &s.Price, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Showtime{}, ErrShowtimeNotFound
		}
		return Showtime{}, fmt.Errorf("error getting showtime by id: %w", err)
	}
	return s, nil
}

func (r *postgresRepository) GetByMovieID(ctx context.Context, movieID int, targetDate time.Time) ([]Showtime, error) {
	query := `
		SELECT id, movie_id, room_id, start_time, end_time, price, created_at, updated_at
		FROM showtimes
		WHERE movie_id = $1 AND DATE(start_time) = DATE($2)
		ORDER BY start_time
	`
	rows, err := r.db.QueryContext(ctx, query, movieID, targetDate)
	if err != nil {
		return nil, fmt.Errorf("error getting showtimes by movie: %w", err)
	}
	defer rows.Close()

	var showtimes []Showtime
	for rows.Next() {
		var s Showtime
		if err := rows.Scan(&s.ID, &s.MovieID, &s.RoomID, &s.StartTime, &s.EndTime, &s.Price, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		showtimes = append(showtimes, s)
	}
	return showtimes, rows.Err()
}

func (r *postgresRepository) GetByRoomID(ctx context.Context, roomID int, targetDate time.Time) ([]Showtime, error) {
	query := `
		SELECT s.id, s.movie_id, s.room_id, s.start_time, s.end_time, s.price, s.created_at, s.updated_at
		FROM showtimes s
		WHERE s.room_id = $1 AND DATE(s.start_time) = DATE($2)
		ORDER BY s.start_time
	`
	rows, err := r.db.QueryContext(ctx, query, roomID, targetDate)
	if err != nil {
		return nil, fmt.Errorf("error getting showtimes by room: %w", err)
	}
	defer rows.Close()

	var showtimes []Showtime
	for rows.Next() {
		var s Showtime
		if err := rows.Scan(&s.ID, &s.MovieID, &s.RoomID, &s.StartTime, &s.EndTime, &s.Price, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		showtimes = append(showtimes, s)
	}
	return showtimes, rows.Err()
}

func (r *postgresRepository) GetByTheaterID(ctx context.Context, theaterID int, targetDate time.Time) ([]Showtime, error) {
	query := `
		SELECT s.id, s.movie_id, s.room_id, s.start_time, s.end_time, s.price, s.created_at, s.updated_at
		FROM showtimes s
		JOIN rooms r ON s.room_id = r.id
		WHERE r.theater_id = $1 AND DATE(s.start_time) = DATE($2)
		ORDER BY s.start_time
	`
	rows, err := r.db.QueryContext(ctx, query, theaterID, targetDate)
	if err != nil {
		return nil, fmt.Errorf("error getting showtimes by theater: %w", err)
	}
	defer rows.Close()

	var showtimes []Showtime
	for rows.Next() {
		var s Showtime
		if err := rows.Scan(&s.ID, &s.MovieID, &s.RoomID, &s.StartTime, &s.EndTime, &s.Price, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		showtimes = append(showtimes, s)
	}
	return showtimes, rows.Err()
}
