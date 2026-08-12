package movie

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/LDKhangg/cinema-booking-go/pkg/apperror"
)

type fakeMovieService struct {
	createMovieErr error
	getMovieErr    error
	updateMovieErr error
	deleteMovieErr error
	getMoviesErr   error
	getMovieResult Movie
	getMoviesResult []Movie
}

func (f *fakeMovieService) CreateMovie(ctx context.Context, m *Movie) error { return f.createMovieErr }
func (f *fakeMovieService) GetMovieById(ctx context.Context, id int) (Movie, error) {
	return f.getMovieResult, f.getMovieErr
}
func (f *fakeMovieService) UpdateMovie(ctx context.Context, m *Movie) error { return f.updateMovieErr }
func (f *fakeMovieService) GetMovies(ctx context.Context) ([]Movie, error) { return f.getMoviesResult, f.getMoviesErr }
func (f *fakeMovieService) DeleteMovie(ctx context.Context, id int) error { return f.deleteMovieErr }

func TestHandlerCreateMovieServiceBadRequest(t *testing.T) {
	h := NewHandler(&fakeMovieService{createMovieErr: apperror.BadRequest("tên phim không được để trống")})
	req := httptest.NewRequest(http.MethodPost, "/movies", strings.NewReader(`{"title":""}`))
	rec := httptest.NewRecorder()

	h.CreateMovie(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if strings.TrimSpace(rec.Body.String()) != `{"code":"bad_request","error":"tên phim không được để trống"}` {
		t.Fatalf("unexpected body: %s", strings.TrimSpace(rec.Body.String()))
	}
}

func TestHandlerDeleteMovieServiceConflict(t *testing.T) {
	h := NewHandler(&fakeMovieService{deleteMovieErr: apperror.Conflict("không thể xóa phim đang trong trạng thái công chiếu, vui lòng gỡ xuống trước")})
	req := httptest.NewRequest(http.MethodDelete, "/movies/1", nil)
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()

	h.DeleteMovie(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}
	if strings.TrimSpace(rec.Body.String()) != `{"code":"conflict","error":"không thể xóa phim đang trong trạng thái công chiếu, vui lòng gỡ xuống trước"}` {
		t.Fatalf("unexpected body: %s", strings.TrimSpace(rec.Body.String()))
	}
}

func TestHandlerGetMovieByIDServiceNotFound(t *testing.T) {
	h := NewHandler(&fakeMovieService{getMovieErr: apperror.NotFound("can not find movie")})
	req := httptest.NewRequest(http.MethodGet, "/movies/1", nil)
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()

	h.GetMovieById(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
	if strings.TrimSpace(rec.Body.String()) != `{"code":"not_found","error":"can not find movie"}` {
		t.Fatalf("unexpected body: %s", strings.TrimSpace(rec.Body.String()))
	}
}
