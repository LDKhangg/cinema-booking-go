package theater

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeTheaterService struct {
	getTheatersCalled bool
	getTheatersResult []Theater
	getTheatersErr    error

	getTheaterByIDCalled bool
	getTheaterByIDInput  int
	getTheaterByIDResult Theater
	getTheaterByIDErr    error

	createTheaterCalled bool
	createTheaterInput  *Theater
	createTheaterErr    error

	getRoomsCalled       bool
	getRoomsTheaterID    int
	getRoomsResult       []Room
	getRoomsByTheaterErr error

	createRoomCalled bool
	createRoomInput  *Room
	createRoomErr    error

	createSeatsCalled      bool
	createSeatsRoomID      int
	createSeatsRows        int
	createSeatsSeatsPerRow int
	createSeatsErr         error

	getSeatsCalled  bool
	getSeatsRoomID  int
	getSeatsResult  []Seat
	getSeatsByIDErr error
}

func (f *fakeTheaterService) GetTheaters(ctx context.Context) ([]Theater, error) {
	f.getTheatersCalled = true
	return f.getTheatersResult, f.getTheatersErr
}

func (f *fakeTheaterService) GetTheaterByID(ctx context.Context, id int) (Theater, error) {
	f.getTheaterByIDCalled = true
	f.getTheaterByIDInput = id
	return f.getTheaterByIDResult, f.getTheaterByIDErr
}

func (f *fakeTheaterService) CreateTheater(ctx context.Context, t *Theater) error {
	f.createTheaterCalled = true
	f.createTheaterInput = t
	return f.createTheaterErr
}

func (f *fakeTheaterService) GetRoomsByTheaterID(ctx context.Context, theaterID int) ([]Room, error) {
	f.getRoomsCalled = true
	f.getRoomsTheaterID = theaterID
	return f.getRoomsResult, f.getRoomsByTheaterErr
}

func (f *fakeTheaterService) CreateRoom(ctx context.Context, r *Room) error {
	f.createRoomCalled = true
	f.createRoomInput = r
	return f.createRoomErr
}

func (f *fakeTheaterService) GetSeatsByRoomID(ctx context.Context, roomID int) ([]Seat, error) {
	f.getSeatsCalled = true
	f.getSeatsRoomID = roomID
	return f.getSeatsResult, f.getSeatsByIDErr
}

func (f *fakeTheaterService) CreateSeats(ctx context.Context, roomID int, rows int, seatsPerRow int) error {
	f.createSeatsCalled = true
	f.createSeatsRoomID = roomID
	f.createSeatsRows = rows
	f.createSeatsSeatsPerRow = seatsPerRow
	return f.createSeatsErr
}

func TestHandlerGetTheaters(t *testing.T) {
	service := &fakeTheaterService{getTheatersResult: []Theater{{ID: 1, Name: "CGV"}}}
	h := NewHandler(service)
	req := httptest.NewRequest(http.MethodGet, "/theaters", nil)
	rec := httptest.NewRecorder()

	h.GetTheaters(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if !service.getTheatersCalled {
		t.Fatal("expected GetTheaters to be called")
	}

	var got []Theater
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("expected valid json body, got %v", err)
	}

	if len(got) != 1 || got[0].Name != "CGV" {
		t.Fatalf("expected one theater named CGV, got %+v", got)
	}
}

func TestHandlerGetTheatersServiceError(t *testing.T) {
	service := &fakeTheaterService{getTheatersErr: errors.New("db down")}
	h := NewHandler(service)
	req := httptest.NewRequest(http.MethodGet, "/theaters", nil)
	rec := httptest.NewRecorder()

	h.GetTheaters(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}

	if strings.TrimSpace(rec.Body.String()) != `{"error":"db down"}` {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestHandlerGetTheaterByIDBadPath(t *testing.T) {
	service := &fakeTheaterService{}
	h := NewHandler(service)
	req := httptest.NewRequest(http.MethodGet, "/theaters/abc", nil)
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	h.GetTheaterByID(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	if service.getTheaterByIDCalled {
		t.Fatal("expected GetTheaterByID not to be called")
	}
}

func TestHandlerGetTheaterByIDSuccess(t *testing.T) {
	service := &fakeTheaterService{getTheaterByIDResult: Theater{ID: 2, Name: "Lotte"}}
	h := NewHandler(service)
	req := httptest.NewRequest(http.MethodGet, "/theaters/2", nil)
	req.SetPathValue("id", "2")
	rec := httptest.NewRecorder()

	h.GetTheaterByID(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if !service.getTheaterByIDCalled || service.getTheaterByIDInput != 2 {
		t.Fatal("expected GetTheaterByID to be called with id 2")
	}

	var got Theater
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("expected valid json body, got %v", err)
	}

	if got.ID != 2 || got.Name != "Lotte" {
		t.Fatalf("unexpected theater: %+v", got)
	}
}

func TestHandlerGetTheaterByIDServiceError(t *testing.T) {
	service := &fakeTheaterService{getTheaterByIDErr: errors.New("not found")}
	h := NewHandler(service)
	req := httptest.NewRequest(http.MethodGet, "/theaters/2", nil)
	req.SetPathValue("id", "2")
	rec := httptest.NewRecorder()

	h.GetTheaterByID(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}

func TestHandlerCreateTheaterBadJSON(t *testing.T) {
	service := &fakeTheaterService{}
	h := NewHandler(service)
	req := httptest.NewRequest(http.MethodPost, "/theaters", strings.NewReader("{"))
	rec := httptest.NewRecorder()

	h.CreateTheater(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	if service.createTheaterCalled {
		t.Fatal("expected CreateTheater not to be called")
	}
}

func TestHandlerCreateTheaterSuccess(t *testing.T) {
	service := &fakeTheaterService{}
	h := NewHandler(service)
	body := `{"name":"CGV Vincom","address":"123 Street","city":"HCM"}`
	req := httptest.NewRequest(http.MethodPost, "/theaters", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateTheater(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}

	if !service.createTheaterCalled {
		t.Fatal("expected CreateTheater to be called")
	}

	if service.createTheaterInput == nil || service.createTheaterInput.Name != "CGV Vincom" {
		t.Fatalf("unexpected create theater input: %+v", service.createTheaterInput)
	}
}

func TestHandlerCreateTheaterServiceError(t *testing.T) {
	service := &fakeTheaterService{createTheaterErr: errors.New("insert failed")}
	h := NewHandler(service)
	body := `{"name":"CGV Vincom","address":"123 Street","city":"HCM"}`
	req := httptest.NewRequest(http.MethodPost, "/theaters", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateTheater(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
}

func TestHandlerGetRoomsByTheaterIDBadPath(t *testing.T) {
	service := &fakeTheaterService{}
	h := NewHandler(service)
	req := httptest.NewRequest(http.MethodGet, "/theaters/abc/rooms", nil)
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	h.GetRoomsByTheaterID(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestHandlerGetRoomsByTheaterIDSuccess(t *testing.T) {
	service := &fakeTheaterService{getRoomsResult: []Room{{ID: 1, Name: "Room 1"}}}
	h := NewHandler(service)
	req := httptest.NewRequest(http.MethodGet, "/theaters/3/rooms", nil)
	req.SetPathValue("id", "3")
	rec := httptest.NewRecorder()

	h.GetRoomsByTheaterID(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if !service.getRoomsCalled || service.getRoomsTheaterID != 3 {
		t.Fatal("expected GetRoomsByTheaterID to be called with theater id 3")
	}
}

func TestHandlerGetRoomsByTheaterIDServiceError(t *testing.T) {
	service := &fakeTheaterService{getRoomsByTheaterErr: errors.New("query rooms failed")}
	h := NewHandler(service)
	req := httptest.NewRequest(http.MethodGet, "/theaters/3/rooms", nil)
	req.SetPathValue("id", "3")
	rec := httptest.NewRecorder()

	h.GetRoomsByTheaterID(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}

	if strings.TrimSpace(rec.Body.String()) != `{"error":"query rooms failed"}` {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestHandlerCreateRoomBadJSON(t *testing.T) {
	service := &fakeTheaterService{}
	h := NewHandler(service)
	req := httptest.NewRequest(http.MethodPost, "/theaters/rooms", strings.NewReader("{"))
	rec := httptest.NewRecorder()

	h.CreateRoom(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	if service.createRoomCalled {
		t.Fatal("expected CreateRoom not to be called")
	}
}

func TestHandlerCreateRoomSuccess(t *testing.T) {
	service := &fakeTheaterService{}
	h := NewHandler(service)
	body := `{"theather_id":1,"name":"Room 1","total_seats":50,"room_type":"2D"}`
	req := httptest.NewRequest(http.MethodPost, "/theaters/rooms", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateRoom(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}

	if !service.createRoomCalled {
		t.Fatal("expected CreateRoom to be called")
	}

	if service.createRoomInput == nil || service.createRoomInput.Name != "Room 1" {
		t.Fatalf("unexpected create room input: %+v", service.createRoomInput)
	}
}

func TestHandlerCreateRoomServiceError(t *testing.T) {
	service := &fakeTheaterService{createRoomErr: errors.New("insert room failed")}
	h := NewHandler(service)
	body := `{"theather_id":1,"name":"Room 1","total_seats":50,"room_type":"2D"}`
	req := httptest.NewRequest(http.MethodPost, "/theaters/rooms", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateRoom(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}

	if strings.TrimSpace(rec.Body.String()) != `{"error":"insert room failed"}` {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestHandlerCreateSeatsBadPath(t *testing.T) {
	service := &fakeTheaterService{}
	h := NewHandler(service)
	body := `{"rows":3,"seats_per_row":2}`
	req := httptest.NewRequest(http.MethodPost, "/rooms/abc/seats", strings.NewReader(body))
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	h.CreateSeats(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	if service.createSeatsCalled {
		t.Fatal("expected CreateSeats not to be called")
	}
}

func TestHandlerCreateSeatsBadJSON(t *testing.T) {
	service := &fakeTheaterService{}
	h := NewHandler(service)
	req := httptest.NewRequest(http.MethodPost, "/rooms/10/seats", strings.NewReader("{"))
	req.SetPathValue("id", "10")
	rec := httptest.NewRecorder()

	h.CreateSeats(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	if service.createSeatsCalled {
		t.Fatal("expected CreateSeats not to be called")
	}
}

func TestHandlerCreateSeatsSuccess(t *testing.T) {
	service := &fakeTheaterService{}
	h := NewHandler(service)
	body := `{"rows":3,"seats_per_row":2}`
	req := httptest.NewRequest(http.MethodPost, "/rooms/10/seats", strings.NewReader(body))
	req.SetPathValue("id", "10")
	rec := httptest.NewRecorder()

	h.CreateSeats(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}

	if !service.createSeatsCalled {
		t.Fatal("expected CreateSeats to be called")
	}

	if service.createSeatsRoomID != 10 || service.createSeatsRows != 3 || service.createSeatsSeatsPerRow != 2 {
		t.Fatalf("unexpected create seats input: roomID=%d rows=%d seatsPerRow=%d", service.createSeatsRoomID, service.createSeatsRows, service.createSeatsSeatsPerRow)
	}

	if strings.TrimSpace(rec.Body.String()) != `{"message":"Tạo ghế thành công"}` {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestHandlerCreateSeatsServiceError(t *testing.T) {
	service := &fakeTheaterService{createSeatsErr: errors.New("create seats failed")}
	h := NewHandler(service)
	body := `{"rows":3,"seats_per_row":2}`
	req := httptest.NewRequest(http.MethodPost, "/rooms/10/seats", strings.NewReader(body))
	req.SetPathValue("id", "10")
	rec := httptest.NewRecorder()

	h.CreateSeats(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}

	if strings.TrimSpace(rec.Body.String()) != `{"error":"create seats failed"}` {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestHandlerGetSeatsByRoomIDBadPath(t *testing.T) {
	service := &fakeTheaterService{}
	h := NewHandler(service)
	req := httptest.NewRequest(http.MethodGet, "/rooms/abc/seats", nil)
	req.SetPathValue("id", "abc")
	rec := httptest.NewRecorder()

	h.GetSeatsByRoomID(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	if service.getSeatsCalled {
		t.Fatal("expected GetSeatsByRoomID not to be called")
	}
}

func TestHandlerGetSeatsByRoomIDSuccess(t *testing.T) {
	service := &fakeTheaterService{getSeatsResult: []Seat{{RoomID: 10, RowLine: "A", Number: 1, SeatType: "standard"}}}
	h := NewHandler(service)
	req := httptest.NewRequest(http.MethodGet, "/rooms/10/seats", nil)
	req.SetPathValue("id", "10")
	rec := httptest.NewRecorder()

	h.GetSeatsByRoomID(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if !service.getSeatsCalled || service.getSeatsRoomID != 10 {
		t.Fatal("expected GetSeatsByRoomID to be called with room id 10")
	}
}

func TestHandlerGetSeatsByRoomIDServiceError(t *testing.T) {
	service := &fakeTheaterService{getSeatsByIDErr: errors.New("query seats failed")}
	h := NewHandler(service)
	req := httptest.NewRequest(http.MethodGet, "/rooms/10/seats", nil)
	req.SetPathValue("id", "10")
	rec := httptest.NewRecorder()

	h.GetSeatsByRoomID(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}

	if strings.TrimSpace(rec.Body.String()) != `{"error":"query seats failed"}` {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
