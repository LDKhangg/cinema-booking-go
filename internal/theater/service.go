package theather

import (
	"context"
	"errors"
	"time"
)

type theaterService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &theaterService{repo: repo}
}

func (s *theaterService) GetTheaters(ctx context.Context) ([]Theater, error) {
	return s.repo.GetTheaters(ctx)
}

func (s *theaterService) GetTheaterByID(ctx context.Context, id int) (Theater, error) {
	if id <= 0 {
		return Theater{}, errors.New("ID rạp chiếu không hợp lệ")
	}
	return s.repo.GetTheaterByID(ctx, id)
}

func (s *theaterService) CreateTheater(ctx context.Context, t *Theater) error {
	if t.Name == "" {
		return errors.New("tên rạp không được để trống")
	}
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	return s.repo.CreateTheater(ctx, t)
}

func (s *theaterService) GetRoomsByTheaterID(ctx context.Context, theaterID int) ([]Room, error) {
	if theaterID <= 0 {
		return nil, errors.New("ID rạp chiếu không hợp lệ")
	}
	return s.repo.GetRoomsByTheaterID(ctx, theaterID)
}

func (s *theaterService) CreateRoom(ctx context.Context, r *Room) error {
	if r.TheatherID <= 0 {
		return errors.New("ID rạp chiếu không hợp lệ")
	}
	if r.Name == "" {
		return errors.New("tên phòng chiếu không được để trống")
	}
	now := time.Now()
	r.CreatedAt = now
	r.UpdatedAt = now
	return s.repo.CreateRoom(ctx, r)
}

func (s *theaterService) GetSeatsByRoomID(ctx context.Context, roomID int) ([]Seat, error) {
	if roomID <= 0 {
		return nil, errors.New("ID phòng chiếu không hợp lệ")
	}
	return s.repo.GetSeatsByRoomID(ctx, roomID)
}

func (s *theaterService) CreateSeats(ctx context.Context, roomID int, rows int, seatsPerRow int) error {
	if roomID <= 0 {
		return errors.New("ID phòng chiếu không hợp lệ")
	}
	if rows <= 0 || seatsPerRow <= 0 {
		return errors.New("số hàng và số ghế mỗi hàng phải lớn hơn 0")
	}

	var seats []Seat
	for i := 0; i < rows; i++ {
		rowChar := string(rune('A' + i))
		for j := 1; j <= seatsPerRow; j++ {
			seatType := "standard"
			if i >= rows-2 {
				seatType = "vip"
			}
			seats = append(seats, Seat{
				RoomID:   roomID,
				RowLine:  rowChar,
				Number:   j,
				SeatType: seatType,
			})
		}
	}

	return s.repo.CreateSeats(ctx, seats)
}
