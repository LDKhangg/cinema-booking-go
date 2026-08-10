package theater

import (
	"context"
	"time"
)

type Theater struct {
	ID        int       `json:"jd"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	City      string    `json:"city"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Room struct {
	ID         int       `json:"id"`
	TheatherID int       `json:"theather_id"`
	Name       string    `json:"name"`
	TotalSeats int       `json:"total_seats"`
	RoomType   string    `json:"room_type"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Seat struct {
	ID       int    `json:"id"`
	RoomID   int    `json:"room_id"`
	RowLine  string `json:"row_line"`
	Number   int    `json:"number"`
	SeatType string `json:"seat_type"`
}

type Repository interface {
	GetTheaters(ctx context.Context) ([]Theater, error)
	GetTheaterByID(ctx context.Context, id int) (Theater, error)
	CreateTheater(ctx context.Context, t *Theater) error

	GetRoomsByTheaterID(ctx context.Context, theaterID int) ([]Room, error)
	CreateRoom(ctx context.Context, r *Room) error

	GetSeatsByRoomID(ctx context.Context, roomID int) ([]Seat, error)
	CreateSeats(ctx context.Context, seats []Seat) error
}

type Service interface {
	GetTheaters(ctx context.Context) ([]Theater, error)
	GetTheaterByID(ctx context.Context, id int) (Theater, error)
	CreateTheater(ctx context.Context, t *Theater) error

	GetRoomsByTheaterID(ctx context.Context, theaterID int) ([]Room, error)
	CreateRoom(ctx context.Context, r *Room) error

	GetSeatsByRoomID(ctx context.Context, roomID int) ([]Seat, error)
	CreateSeats(ctx context.Context, roomID int, rows int, seatsPerRow int) error
}
