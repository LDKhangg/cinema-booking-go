package movie

import (
	"context"
	"errors"
	"time"

	"github.com/LDKhangg/cinema-booking-go/pkg/apperror"
)

type movieService struct {
	repo Repository
}

func (s *movieService) CreateMovie(ctx context.Context, m *Movie) error {
	if m.Title == "" {
		return apperror.BadRequest("tên phim không được để trống")
	}
	if !m.ClosedDate.IsZero() && m.ClosedDate.Before(m.ReleaseDate) {
		return apperror.BadRequest("ngày kết thúc chiếu không được trước ngày khởi chiếu")
	}

	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now

	if err := s.repo.Create(ctx, m); err != nil {
		return apperror.Internal("không thể tạo phim", err)
	}

	return nil
}

func (m *movieService) DeleteMovie(ctx context.Context, id int) error {
	if id <= 0 {
		return apperror.BadRequest("ID phim không hợp lệ")
	}

	existingMovie, err := m.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrMovieNotFound) {
			return apperror.NotFound("can not find movie")
		}
		return apperror.Internal("không thể lấy thông tin phim", err)
	}

	if existingMovie.IsPublished {
		return apperror.Conflict("không thể xóa phim đang trong trạng thái công chiếu, vui lòng gỡ xuống trước")
	}

	if err := m.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, ErrMovieNotFound) {
			return apperror.NotFound("can not find movie")
		}
		return apperror.Internal("không thể xóa phim", err)
	}

	return nil
}

func (m *movieService) GetMovieById(ctx context.Context, id int) (Movie, error) {
	if id <= 0 {
		return Movie{}, apperror.BadRequest("ID phim không hợp lệ")
	}

	movie, err := m.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrMovieNotFound) {
			return Movie{}, apperror.NotFound("can not find movie")
		}
		return Movie{}, apperror.Internal("không thể lấy thông tin phim", err)
	}

	return movie, nil
}

func (m *movieService) GetMovies(ctx context.Context) ([]Movie, error) {
	movies, err := m.repo.GetPublishedMovies(ctx)
	if err != nil {
		return nil, apperror.Internal("không thể lấy danh sách phim", err)
	}

	return movies, nil
}

func (mb *movieService) UpdateMovie(ctx context.Context, m *Movie) error {
	if m.ID <= 0 {
		return apperror.BadRequest("ID phim không hợp lệ")
	}

	if m.Title == "" {
		return apperror.BadRequest("tên phim không được để trống")
	}

	if !m.ClosedDate.IsZero() && m.ClosedDate.Before(m.ReleaseDate) {
		return apperror.BadRequest("ngày kết thúc chiếu không được trước ngày khởi chiếu")
	}

	m.UpdatedAt = time.Now()

	if err := mb.repo.Update(ctx, m); err != nil {
		if errors.Is(err, ErrMovieNotFound) {
			return apperror.NotFound("can not find movie")
		}
		return apperror.Internal("không thể cập nhật phim", err)
	}

	return nil
}

func NewService(repo Repository) Service {
	return &movieService{
		repo: repo,
	}
}
