package showtime

import (
	"context"
	"errors"
	"time"
)

type showtimeService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &showtimeService{repo: repo}
}

func (s *showtimeService) CreateShowtime(ctx context.Context, st *Showtime) error {
	if st.MovieID <= 0 {
		return errors.New("ID phim không hợp lệ")
	}
	if st.RoomID <= 0 {
		return errors.New("ID phòng chiếu không hợp lệ")
	}
	if st.StartTime.Before(time.Now()) {
		return errors.New("thời gian chiếu không được ở trong quá khứ")
	}
	if !st.EndTime.After(st.StartTime) {
		return errors.New("thời gian kết thúc phải sau thời gian bắt đầu")
	}
	if st.Price <= 0 {
		return errors.New("giá vé phải lớn hơn 0")
	}

	existing, err := s.repo.GetByRoomID(ctx, st.RoomID, st.StartTime)
	if err == nil {
		for _, ex := range existing {
			if st.StartTime.Before(ex.EndTime) && st.EndTime.After(ex.StartTime) {
				return errors.New("lịch chiếu bị trùng với suất chiếu khác trong cùng phòng")
			}
		}
	}

	now := time.Now()
	st.CreatedAt = now
	st.UpdatedAt = now

	return s.repo.Create(ctx, st)
}

func (s *showtimeService) GetShowtimeByID(ctx context.Context, id int) (Showtime, error) {
	if id <= 0 {
		return Showtime{}, errors.New("ID suất chiếu không hợp lệ")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *showtimeService) GetShowtimesByMovie(ctx context.Context, movieID int, date time.Time) ([]Showtime, error) {
	if movieID <= 0 {
		return nil, errors.New("ID phim không hợp lệ")
	}
	return s.repo.GetByMovieID(ctx, movieID, date)
}

func (s *showtimeService) GetShowtimesByTheater(ctx context.Context, theaterID int, date time.Time) ([]Showtime, error) {
	if theaterID <= 0 {
		return nil, errors.New("ID rạp chiếu không hợp lệ")
	}
	// Extended repository supports GetByTheaterID
	type theaterRepo interface {
		GetByTheaterID(ctx context.Context, theaterID int, targetDate time.Time) ([]Showtime, error)
	}
	if tr, ok := s.repo.(theaterRepo); ok {
		return tr.GetByTheaterID(ctx, theaterID, date)
	}
	return nil, errors.New("repository không hỗ trợ tra cứu theo rạp")
}
