package theater

import (
	"context"
	"errors"
	"time"

	"github.com/LDKhangg/cinema-booking-go/pkg/apperror"
)

type theaterService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &theaterService{repo: repo}
}

func (s *theaterService) GetTheaters(ctx context.Context) ([]Theater, error) {
	theaters, err := s.repo.GetTheaters(ctx)
	if err != nil {
		return nil, apperror.Internal("không thể lấy danh sách rạp", err)
	}

	return theaters, nil
}

func (s *theaterService) GetTheaterByID(ctx context.Context, id int) (Theater, error) {
	if id <= 0 {
		return Theater{}, apperror.BadRequest("ID rạp chiếu không hợp lệ")
	}

	theater, err := s.repo.GetTheaterByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrTheaterNotFound) {
			return Theater{}, apperror.NotFound("không tìm thấy rạp chiếu")
		}
		return Theater{}, apperror.Internal("không thể lấy thông tin rạp chiếu", err)
	}

	return theater, nil
}

func (s *theaterService) CreateTheater(ctx context.Context, t *Theater) error {
	if t.Name == "" {
		return apperror.BadRequest("tên rạp không được để trống")
	}
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	if err := s.repo.CreateTheater(ctx, t); err != nil {
		return apperror.Internal("không thể tạo rạp mới", err)
	}

	return nil
}

func (s *theaterService) GetRoomsByTheaterID(ctx context.Context, theaterID int) ([]Room, error) {
	if theaterID <= 0 {
		return nil, apperror.BadRequest("ID rạp chiếu không hợp lệ")
	}

	rooms, err := s.repo.GetRoomsByTheaterID(ctx, theaterID)
	if err != nil {
		return nil, apperror.Internal("không thể lấy danh sách phòng chiếu", err)
	}

	return rooms, nil
}

func (s *theaterService) CreateRoom(ctx context.Context, r *Room) error {
	if r.TheatherID <= 0 {
		return apperror.BadRequest("ID rạp chiếu không hợp lệ")
	}
	if r.Name == "" {
		return apperror.BadRequest("tên phòng chiếu không được để trống")
	}
	now := time.Now()
	r.CreatedAt = now
	r.UpdatedAt = now
	if err := s.repo.CreateRoom(ctx, r); err != nil {
		return apperror.Internal("không thể tạo phòng chiếu mới", err)
	}

	return nil
}

func (s *theaterService) GetSeatsByRoomID(ctx context.Context, roomID int) ([]Seat, error) {
	if roomID <= 0 {
		return nil, apperror.BadRequest("ID phòng chiếu không hợp lệ")
	}

	seats, err := s.repo.GetSeatsByRoomID(ctx, roomID)
	if err != nil {
		return nil, apperror.Internal("không thể lấy danh sách ghế", err)
	}

	return seats, nil
}

func (s *theaterService) CreateSeats(ctx context.Context, roomID int, rows int, seatsPerRow int) error {
	if roomID <= 0 {
		return apperror.BadRequest("ID phòng chiếu không hợp lệ")
	}
	if rows <= 0 || seatsPerRow <= 0 {
		return apperror.BadRequest("số hàng và số ghế mỗi hàng phải lớn hơn 0")
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

	if err := s.repo.CreateSeats(ctx, seats); err != nil {
		return apperror.Internal("không thể tạo ghế", err)
	}

	return nil
}
