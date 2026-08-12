package showtime

import (
	"context"
	"errors"
	"time"

	"github.com/LDKhangg/cinema-booking-go/pkg/apperror"
)

type showtimeService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &showtimeService{repo: repo}
}

func (s *showtimeService) CreateShowtime(ctx context.Context, st *Showtime) error {
	if st.MovieID <= 0 {
		return apperror.BadRequest("ID phim không hợp lệ")
	}
	if st.RoomID <= 0 {
		return apperror.BadRequest("ID phòng chiếu không hợp lệ")
	}
	if st.StartTime.Before(time.Now()) {
		return apperror.BadRequest("thời gian chiếu không được ở trong quá khứ")
	}
	if !st.EndTime.After(st.StartTime) {
		return apperror.BadRequest("thời gian kết thúc phải sau thời gian bắt đầu")
	}
	if st.Price <= 0 {
		return apperror.BadRequest("giá vé phải lớn hơn 0")
	}

	existing, err := s.repo.GetByRoomID(ctx, st.RoomID, st.StartTime)
	if err != nil {
		return apperror.Internal("không thể kiểm tra lịch chiếu hiện có", err)
	}

	for _, ex := range existing {
		if st.StartTime.Before(ex.EndTime) && st.EndTime.After(ex.StartTime) {
			return apperror.Conflict("lịch chiếu bị trùng với suất chiếu khác trong cùng phòng")
		}
	}

	now := time.Now()
	st.CreatedAt = now
	st.UpdatedAt = now

	if err := s.repo.Create(ctx, st); err != nil {
		return apperror.Internal("không thể tạo suất chiếu", err)
	}

	return nil
}

func (s *showtimeService) GetShowtimeByID(ctx context.Context, id int) (Showtime, error) {
	if id <= 0 {
		return Showtime{}, apperror.BadRequest("ID suất chiếu không hợp lệ")
	}

	st, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrShowtimeNotFound) {
			return Showtime{}, apperror.NotFound("cannot find showtime")
		}
		return Showtime{}, apperror.Internal("không thể lấy thông tin suất chiếu", err)
	}

	return st, nil
}

func (s *showtimeService) GetShowtimesByMovie(ctx context.Context, movieID int, date time.Time) ([]Showtime, error) {
	if movieID <= 0 {
		return nil, apperror.BadRequest("ID phim không hợp lệ")
	}

	showtimes, err := s.repo.GetByMovieID(ctx, movieID, date)
	if err != nil {
		return nil, apperror.Internal("không thể lấy danh sách suất chiếu theo phim", err)
	}

	return showtimes, nil
}

func (s *showtimeService) GetShowtimesByTheater(ctx context.Context, theaterID int, date time.Time) ([]Showtime, error) {
	if theaterID <= 0 {
		return nil, apperror.BadRequest("ID rạp chiếu không hợp lệ")
	}
	// Extended repository supports GetByTheaterID
	type theaterRepo interface {
		GetByTheaterID(ctx context.Context, theaterID int, targetDate time.Time) ([]Showtime, error)
	}
	if tr, ok := s.repo.(theaterRepo); ok {
		showtimes, err := tr.GetByTheaterID(ctx, theaterID, date)
		if err != nil {
			return nil, apperror.Internal("không thể lấy danh sách suất chiếu theo rạp", err)
		}
		return showtimes, nil
	}

	return nil, apperror.Internal("repository không hỗ trợ tra cứu theo rạp", nil)
}
