package movie

import (
	"context"
	"errors"
	"time"
)

type movieService struct {
	repo Repository
}

func (s *movieService) CreateMovie(ctx context.Context, m *Movie) error {
	if m.Title == "" {
		return errors.New("tên phim không được để trống")
	}
	if !m.ClosedDate.IsZero() && m.ClosedDate.Before(m.ReleaseDate) {
		return errors.New("ngày kết thúc chiếu không được trước ngày khởi chiếu")
	}

	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now

	return s.repo.Create(ctx, m)
}

func (m *movieService) DeleteMovie(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("ID phim không hợp lệ")
	}

	existingMovie, err := m.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingMovie.IsPublished {
		return errors.New("không thể xóa phim đang trong trạng thái công chiếu, vui lòng gỡ xuống trước")
	}

	return m.repo.Delete(ctx, id)
}

func (m *movieService) GetMovieById(ctx context.Context, id int) (Movie, error) {
	if id <= 0 {
		return Movie{}, errors.New("ID phim không hợp lệ")
	}

	return m.repo.GetByID(ctx, id)
}

func (m *movieService) GetMovies(ctx context.Context) ([]Movie, error) {
	return m.repo.GetPublishedMovies(ctx)
}

func (mb *movieService) UpdateMovie(ctx context.Context, m *Movie) error {
	if m.ID <= 0 {
		return errors.New("ID phim không hợp lệ")
	}

	if m.Title == "" {
		return errors.New("tên phim không được để trống")
	}

	if !m.ClosedDate.IsZero() && m.ClosedDate.Before(m.ReleaseDate) {
		return errors.New("ngày kết thúc chiếu không được trước ngày khởi chiếu")
	}

	m.UpdatedAt = time.Now()

	return mb.repo.Update(ctx, m)
}

func NewService(repo Repository) Service {
	return &movieService{
		repo: repo,
	}
}
