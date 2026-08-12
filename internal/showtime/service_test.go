package showtime

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/LDKhangg/cinema-booking-go/pkg/apperror"
)

type fakeShowtimeRepo struct {
	createErr         error
	getByIDResult     Showtime
	getByIDErr        error
	getByMovieResult  []Showtime
	getByMovieErr     error
	getByRoomResult   []Showtime
	getByRoomErr      error
	getByTheaterResult []Showtime
	getByTheaterErr   error

	createCalled      bool
	getByIDCalled     bool
	getByMovieCalled  bool
	getByRoomCalled   bool
	getByTheaterCalled bool

	createInput       *Showtime
	getByIDInput      int
	getByMovieInput   int
	getByRoomInput    int
	getByTheaterInput int
}

type repoWithoutTheaterLookup struct{}

func (f *fakeShowtimeRepo) Create(ctx context.Context, s *Showtime) error {
	f.createCalled = true
	f.createInput = s
	return f.createErr
}

func (f *fakeShowtimeRepo) GetByID(ctx context.Context, id int) (Showtime, error) {
	f.getByIDCalled = true
	f.getByIDInput = id
	return f.getByIDResult, f.getByIDErr
}

func (f *fakeShowtimeRepo) GetByMovieID(ctx context.Context, movieID int, targetDate time.Time) ([]Showtime, error) {
	f.getByMovieCalled = true
	f.getByMovieInput = movieID
	return f.getByMovieResult, f.getByMovieErr
}

func (f *fakeShowtimeRepo) GetByRoomID(ctx context.Context, roomID int, targetDate time.Time) ([]Showtime, error) {
	f.getByRoomCalled = true
	f.getByRoomInput = roomID
	return f.getByRoomResult, f.getByRoomErr
}

func (f *fakeShowtimeRepo) GetByTheaterID(ctx context.Context, theaterID int, targetDate time.Time) ([]Showtime, error) {
	f.getByTheaterCalled = true
	f.getByTheaterInput = theaterID
	return f.getByTheaterResult, f.getByTheaterErr
}

func (repoWithoutTheaterLookup) Create(ctx context.Context, s *Showtime) error { return nil }

func (repoWithoutTheaterLookup) GetByID(ctx context.Context, id int) (Showtime, error) {
	return Showtime{}, nil
}

func (repoWithoutTheaterLookup) GetByMovieID(ctx context.Context, movieID int, targetDate time.Time) ([]Showtime, error) {
	return nil, nil
}

func (repoWithoutTheaterLookup) GetByRoomID(ctx context.Context, roomID int, targetDate time.Time) ([]Showtime, error) {
	return nil, nil
}

func TestShowtimeServiceCreateShowtimeValidationReturnsBadRequest(t *testing.T) {
	svc := &showtimeService{repo: &fakeShowtimeRepo{}}

	err := svc.CreateShowtime(context.Background(), &Showtime{})

	assertShowtimeAppError(t, err, "bad_request", "ID phim không hợp lệ")
}

func TestShowtimeServiceCreateShowtimeConflictReturnsConflict(t *testing.T) {
	now := time.Now().Add(2 * time.Hour)
	repo := &fakeShowtimeRepo{
		getByRoomResult: []Showtime{{
			RoomID:    1,
			StartTime: now,
			EndTime:   now.Add(2 * time.Hour),
		}},
	}
	svc := &showtimeService{repo: repo}

	err := svc.CreateShowtime(context.Background(), &Showtime{
		MovieID:   1,
		RoomID:    1,
		StartTime: now.Add(30 * time.Minute),
		EndTime:   now.Add(3 * time.Hour),
		Price:     90000,
	})

	assertShowtimeAppError(t, err, "conflict", "lịch chiếu bị trùng với suất chiếu khác trong cùng phòng")
	if repo.createCalled {
		t.Fatal("expected Create not to be called")
	}
}

func TestShowtimeServiceGetShowtimeByIDNotFoundReturnsAppError(t *testing.T) {
	repo := &fakeShowtimeRepo{getByIDErr: ErrShowtimeNotFound}
	svc := &showtimeService{repo: repo}

	_, err := svc.GetShowtimeByID(context.Background(), 99)

	assertShowtimeAppError(t, err, "not_found", "cannot find showtime")
}

func TestShowtimeServiceCreateShowtimeRepositoryErrorReturnsInternal(t *testing.T) {
	repoErr := errors.New("db down")
	repo := &fakeShowtimeRepo{createErr: repoErr}
	svc := &showtimeService{repo: repo}
	now := time.Now().Add(2 * time.Hour)

	err := svc.CreateShowtime(context.Background(), &Showtime{
		MovieID:   1,
		RoomID:    2,
		StartTime: now,
		EndTime:   now.Add(2 * time.Hour),
		Price:     100000,
	})

	appErr := assertShowtimeAppError(t, err, "internal_error", "không thể tạo suất chiếu")
	if !errors.Is(appErr, repoErr) {
		t.Fatalf("expected wrapped repo error, got %v", appErr)
	}
}

func TestShowtimeServiceGetShowtimesByTheaterUnsupportedRepositoryReturnsInternal(t *testing.T) {
	svc := &showtimeService{repo: repoWithoutTheaterLookup{}}

	_, err := svc.GetShowtimesByTheater(context.Background(), 1, time.Now())

	assertShowtimeAppError(t, err, "internal_error", "repository không hỗ trợ tra cứu theo rạp")
}

func assertShowtimeAppError(t *testing.T, err error, wantCode string, wantMessage string) *apperror.AppError {
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
