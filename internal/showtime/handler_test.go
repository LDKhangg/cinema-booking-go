package showtime

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/LDKhangg/cinema-booking-go/pkg/apperror"
)

type fakeShowtimeService struct {
	createErr         error
	getByIDErr        error
	getByMovieErr     error
	getByTheaterErr   error
	getByIDResult     Showtime
	getByMovieResult  []Showtime
	getByTheaterResult []Showtime
}

func (f *fakeShowtimeService) CreateShowtime(ctx context.Context, s *Showtime) error { return f.createErr }
func (f *fakeShowtimeService) GetShowtimeByID(ctx context.Context, id int) (Showtime, error) {
	return f.getByIDResult, f.getByIDErr
}
func (f *fakeShowtimeService) GetShowtimesByMovie(ctx context.Context, movieID int, date time.Time) ([]Showtime, error) {
	return f.getByMovieResult, f.getByMovieErr
}
func (f *fakeShowtimeService) GetShowtimesByTheater(ctx context.Context, theaterID int, date time.Time) ([]Showtime, error) {
	return f.getByTheaterResult, f.getByTheaterErr
}

func TestHandlerCreateShowtimeServiceBadRequest(t *testing.T) {
	h := NewHandler(&fakeShowtimeService{createErr: apperror.BadRequest("ID phim không hợp lệ")})
	req := httptest.NewRequest(http.MethodPost, "/showtimes", strings.NewReader(`{"movie_id":0}`))
	rec := httptest.NewRecorder()

	h.CreateShowtime(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if strings.TrimSpace(rec.Body.String()) != `{"code":"bad_request","error":"ID phim không hợp lệ"}` {
		t.Fatalf("unexpected body: %s", strings.TrimSpace(rec.Body.String()))
	}
}

func TestHandlerCreateShowtimeServiceConflict(t *testing.T) {
	h := NewHandler(&fakeShowtimeService{createErr: apperror.Conflict("lịch chiếu bị trùng với suất chiếu khác trong cùng phòng")})
	req := httptest.NewRequest(http.MethodPost, "/showtimes", strings.NewReader(`{"movie_id":1,"room_id":1}`))
	rec := httptest.NewRecorder()

	h.CreateShowtime(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}
	if strings.TrimSpace(rec.Body.String()) != `{"code":"conflict","error":"lịch chiếu bị trùng với suất chiếu khác trong cùng phòng"}` {
		t.Fatalf("unexpected body: %s", strings.TrimSpace(rec.Body.String()))
	}
}

func TestHandlerGetShowtimeByIDServiceNotFound(t *testing.T) {
	h := NewHandler(&fakeShowtimeService{getByIDErr: apperror.NotFound("cannot find showtime")})
	req := httptest.NewRequest(http.MethodGet, "/showtimes/1", nil)
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()

	h.GetShowtimeByID(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
	if strings.TrimSpace(rec.Body.String()) != `{"code":"not_found","error":"cannot find showtime"}` {
		t.Fatalf("unexpected body: %s", strings.TrimSpace(rec.Body.String()))
	}
}

func TestHandlerGetShowtimesByMovieServiceError(t *testing.T) {
	h := NewHandler(&fakeShowtimeService{getByMovieErr: apperror.Internal("không thể lấy danh sách suất chiếu theo phim", nil)})
	req := httptest.NewRequest(http.MethodGet, "/showtimes?movie_id=1&date=2026-01-01", nil)
	rec := httptest.NewRecorder()

	h.GetShowtimes(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	if strings.TrimSpace(rec.Body.String()) != `{"code":"internal_error","error":"không thể lấy danh sách suất chiếu theo phim"}` {
		t.Fatalf("unexpected body: %s", strings.TrimSpace(rec.Body.String()))
	}
}

func TestHandlerGetShowtimesByTheaterServiceError(t *testing.T) {
	h := NewHandler(&fakeShowtimeService{getByTheaterErr: apperror.Internal("không thể lấy danh sách suất chiếu theo rạp", nil)})
	req := httptest.NewRequest(http.MethodGet, "/showtimes?theater_id=1&date=2026-01-01", nil)
	rec := httptest.NewRecorder()

	h.GetShowtimes(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	if strings.TrimSpace(rec.Body.String()) != `{"code":"internal_error","error":"không thể lấy danh sách suất chiếu theo rạp"}` {
		t.Fatalf("unexpected body: %s", strings.TrimSpace(rec.Body.String()))
	}
}
