package movie

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/LDKhangg/cinema-booking-go/pkg/apperror"
)

type fakeMovieRepo struct {
	createErr        error
	getByIDMovie     Movie
	getByIDErr       error
	getPublished     []Movie
	getPublishedErr  error
	updateErr        error
	deleteErr        error
	createCalled     bool
	deleteCalled     bool
	updateCalled     bool
	getByIDCalled    bool
	publishedCalled  bool
}

func (f *fakeMovieRepo) GetPublishedMovies(ctx context.Context) ([]Movie, error) {
	f.publishedCalled = true
	return f.getPublished, f.getPublishedErr
}

func (f *fakeMovieRepo) GetByID(ctx context.Context, id int) (Movie, error) {
	f.getByIDCalled = true
	return f.getByIDMovie, f.getByIDErr
}

func (f *fakeMovieRepo) Create(ctx context.Context, m *Movie) error {
	f.createCalled = true
	return f.createErr
}

func (f *fakeMovieRepo) Update(ctx context.Context, m *Movie) error {
	f.updateCalled = true
	return f.updateErr
}

func (f *fakeMovieRepo) Delete(ctx context.Context, id int) error {
	f.deleteCalled = true
	return f.deleteErr
}

func TestMovieServiceCreateMovieValidationReturnsBadRequest(t *testing.T) {
	svc := &movieService{repo: &fakeMovieRepo{}}
	movie := &Movie{}

	err := svc.CreateMovie(context.Background(), movie)

	assertAppError(t, err, "bad_request", "tên phim không được để trống")
}

func TestMovieServiceDeleteMoviePublishedReturnsConflict(t *testing.T) {
	repo := &fakeMovieRepo{getByIDMovie: Movie{ID: 1, IsPublished: true}}
	svc := &movieService{repo: repo}

	err := svc.DeleteMovie(context.Background(), 1)

	assertAppError(t, err, "conflict", "không thể xóa phim đang trong trạng thái công chiếu, vui lòng gỡ xuống trước")
	if repo.deleteCalled {
		t.Fatal("expected delete not to be called")
	}
}

func TestMovieServiceGetMovieByIDNotFoundReturnsAppError(t *testing.T) {
	repo := &fakeMovieRepo{getByIDErr: ErrMovieNotFound}
	svc := &movieService{repo: repo}

	_, err := svc.GetMovieById(context.Background(), 99)

	assertAppError(t, err, "not_found", "can not find movie")
}

func TestMovieServiceCreateMovieRepositoryErrorReturnsInternal(t *testing.T) {
	repoErr := errors.New("db down")
	repo := &fakeMovieRepo{createErr: repoErr}
	svc := &movieService{repo: repo}

	err := svc.CreateMovie(context.Background(), &Movie{
		Title:       "Dune",
		ReleaseDate: time.Now(),
	})

	appErr := assertAppError(t, err, "internal_error", "không thể tạo phim")
	if !errors.Is(appErr, repoErr) {
		t.Fatalf("expected wrapped repo error, got %v", appErr)
	}
}

func assertAppError(t *testing.T, err error, wantCode string, wantMessage string) *apperror.AppError {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	appErr, ok := err.(*apperror.AppError)
	if !ok {
		t.Fatalf("expected *apperror.AppError, got %T", err)
	}

	if appErr.Code != wantCode {
		t.Fatalf("expected code %q, got %q", wantCode, appErr.Code)
	}

	if appErr.Message != wantMessage {
		t.Fatalf("expected message %q, got %q", wantMessage, appErr.Message)
	}

	return appErr
}
