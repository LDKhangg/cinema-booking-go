package theater

import (
	"context"
	"errors"
	"testing"

	"github.com/LDKhangg/cinema-booking-go/pkg/apperror"
)

func assertAppError(t *testing.T, err error, wantCode string, wantMessage string) {
	t.Helper()

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected *apperror.AppError, got %T (%v)", err, err)
	}

	if appErr.Code != wantCode {
		t.Fatalf("expected code %q, got %q", wantCode, appErr.Code)
	}

	if appErr.Message != wantMessage {
		t.Fatalf("expected message %q, got %q", wantMessage, appErr.Message)
	}
}

type fakeTheaterRepo struct {
	getTheatersResult []Theater
	getTheatersErr    error

	getTheaterByIDCalled bool
	getTheaterByIDInput  int
	getTheaterByIDResult Theater
	getTheaterByIDErr    error

	createTheaterCalled bool
	receivedTheater     *Theater
	createTheaterErr    error

	getRoomsCalled        bool
	getRoomsTheaterID     int
	getRoomsResult        []Room
	getRoomsByTheaterErr  error

	createRoomCalled bool
	receivedRoom     *Room
	createRoomErr    error

	getSeatsCalled  bool
	getSeatsRoomID  int
	getSeatsResult  []Seat
	getSeatsByIDErr error

	createSeatsCalled bool
	receivedSeats     []Seat
	createSeatsErr    error
}

func (f *fakeTheaterRepo) GetTheaters(ctx context.Context) ([]Theater, error) {
	return f.getTheatersResult, f.getTheatersErr
}

func (f *fakeTheaterRepo) GetTheaterByID(ctx context.Context, id int) (Theater, error) {
	f.getTheaterByIDCalled = true
	f.getTheaterByIDInput = id
	return f.getTheaterByIDResult, f.getTheaterByIDErr
}

func (f *fakeTheaterRepo) CreateTheater(ctx context.Context, t *Theater) error {
	f.createTheaterCalled = true
	f.receivedTheater = t
	return f.createTheaterErr
}

func (f *fakeTheaterRepo) GetRoomsByTheaterID(ctx context.Context, theaterID int) ([]Room, error) {
	f.getRoomsCalled = true
	f.getRoomsTheaterID = theaterID
	return f.getRoomsResult, f.getRoomsByTheaterErr
}

func (f *fakeTheaterRepo) CreateRoom(ctx context.Context, r *Room) error {
	f.createRoomCalled = true
	f.receivedRoom = r
	return f.createRoomErr
}

func (f *fakeTheaterRepo) GetSeatsByRoomID(ctx context.Context, roomID int) ([]Seat, error) {
	f.getSeatsCalled = true
	f.getSeatsRoomID = roomID
	return f.getSeatsResult, f.getSeatsByIDErr
}

func (f *fakeTheaterRepo) CreateSeats(ctx context.Context, seats []Seat) error {
	f.createSeatsCalled = true
	f.receivedSeats = seats
	return f.createSeatsErr
}

func TestTheaterServiceGetTheaters(t *testing.T) {
	want := []Theater{{ID: 1, Name: "CGV"}}
	repo := &fakeTheaterRepo{getTheatersResult: want}
	service := &theaterService{repo: repo}

	got, err := service.GetTheaters(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(got) != len(want) {
		t.Fatalf("expected %d theaters, got %d", len(want), len(got))
	}

	if got[0].ID != want[0].ID || got[0].Name != want[0].Name {
		t.Fatalf("expected first theater %+v, got %+v", want[0], got[0])
	}
}

func TestTheaterServiceGetTheaterByIDValidation(t *testing.T) {
	repo := &fakeTheaterRepo{}
	service := &theaterService{repo: repo}

	_, err := service.GetTheaterByID(context.Background(), 0)
	assertAppError(t, err, "bad_request", "ID rạp chiếu không hợp lệ")

	if repo.getTheaterByIDCalled {
		t.Fatal("expected GetTheaterByID not to be called")
	}
}

func TestTheaterServiceGetTheaterByIDSuccess(t *testing.T) {
	want := Theater{ID: 7, Name: "Lotte"}
	repo := &fakeTheaterRepo{getTheaterByIDResult: want}
	service := &theaterService{repo: repo}

	got, err := service.GetTheaterByID(context.Background(), 7)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !repo.getTheaterByIDCalled {
		t.Fatal("expected GetTheaterByID to be called")
	}

	if repo.getTheaterByIDInput != 7 {
		t.Fatalf("expected repo to receive id 7, got %d", repo.getTheaterByIDInput)
	}

	if got.ID != want.ID || got.Name != want.Name {
		t.Fatalf("expected theater %+v, got %+v", want, got)
	}
}

func TestTheaterServiceCreateTheaterSuccess(t *testing.T) {
	repo := &fakeTheaterRepo{}
	service := &theaterService{repo: repo}
	theater := &Theater{Name: "CGV Vincom", Address: "123 Street", City: "HCM"}

	err := service.CreateTheater(context.Background(), theater)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !repo.createTheaterCalled {
		t.Fatal("expected CreateTheater to be called")
	}

	if repo.receivedTheater != theater {
		t.Fatal("expected repo to receive the same theater pointer")
	}

	if theater.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be set")
	}

	if theater.UpdatedAt.IsZero() {
		t.Fatal("expected UpdatedAt to be set")
	}

	if !theater.CreatedAt.Equal(theater.UpdatedAt) {
		t.Fatal("expected CreatedAt and UpdatedAt to be equal on create")
	}
}

func TestTheaterServiceCreateTheaterValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   *Theater
		wantErr string
	}{
		{
			name:    "empty name",
			input:   &Theater{Name: ""},
			wantErr: "tên rạp không được để trống",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeTheaterRepo{}
			service := &theaterService{repo: repo}

			err := service.CreateTheater(context.Background(), tt.input)
			assertAppError(t, err, "bad_request", tt.wantErr)

			if repo.createTheaterCalled {
				t.Fatal("expected CreateTheater not to be called")
			}
		})
	}
}

func TestTheaterServiceCreateTheaterRepoError(t *testing.T) {
	repoErr := errors.New("database is down")
	repo := &fakeTheaterRepo{createTheaterErr: repoErr}
	service := &theaterService{repo: repo}
	theater := &Theater{Name: "CGV Vincom"}

	err := service.CreateTheater(context.Background(), theater)
	assertAppError(t, err, "internal_error", "không thể tạo rạp mới")

	if !errors.Is(err, repoErr) {
		t.Fatalf("expected error %v, got %v", repoErr, err)
	}

	if !repo.createTheaterCalled {
		t.Fatal("expected CreateTheater to be called")
	}

	if theater.CreatedAt.IsZero() || theater.UpdatedAt.IsZero() {
		t.Fatal("expected timestamps to be set before repo error")
	}
}

func TestTheaterServiceGetRoomsByTheaterIDValidation(t *testing.T) {
	repo := &fakeTheaterRepo{}
	service := &theaterService{repo: repo}

	_, err := service.GetRoomsByTheaterID(context.Background(), 0)
	assertAppError(t, err, "bad_request", "ID rạp chiếu không hợp lệ")

	if repo.getRoomsCalled {
		t.Fatal("expected GetRoomsByTheaterID not to be called")
	}
}

func TestTheaterServiceGetRoomsByTheaterIDSuccess(t *testing.T) {
	want := []Room{{ID: 1, TheatherID: 2, Name: "Room 1"}}
	repo := &fakeTheaterRepo{getRoomsResult: want}
	service := &theaterService{repo: repo}

	got, err := service.GetRoomsByTheaterID(context.Background(), 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !repo.getRoomsCalled {
		t.Fatal("expected GetRoomsByTheaterID to be called")
	}

	if repo.getRoomsTheaterID != 2 {
		t.Fatalf("expected repo to receive theater id 2, got %d", repo.getRoomsTheaterID)
	}

	if len(got) != len(want) {
		t.Fatalf("expected %d rooms, got %d", len(want), len(got))
	}
}

func TestTheaterServiceCreateRoomSuccess(t *testing.T) {
	repo := &fakeTheaterRepo{}
	service := &theaterService{repo: repo}
	room := &Room{TheatherID: 1, Name: "Room 1", TotalSeats: 50, RoomType: "2D"}

	err := service.CreateRoom(context.Background(), room)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !repo.createRoomCalled {
		t.Fatal("expected CreateRoom to be called")
	}

	if repo.receivedRoom != room {
		t.Fatal("expected repo to receive the same room pointer")
	}

	if room.CreatedAt.IsZero() || room.UpdatedAt.IsZero() {
		t.Fatal("expected room timestamps to be set")
	}

	if !room.CreatedAt.Equal(room.UpdatedAt) {
		t.Fatal("expected room timestamps to be equal on create")
	}
}

func TestTheaterServiceCreateRoomValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   *Room
		wantErr string
	}{
		{
			name:    "invalid theater id",
			input:   &Room{TheatherID: 0, Name: "Room 1", TotalSeats: 50, RoomType: "2D"},
			wantErr: "ID rạp chiếu không hợp lệ",
		},
		{
			name:    "empty room name",
			input:   &Room{TheatherID: 1, Name: "", TotalSeats: 50, RoomType: "2D"},
			wantErr: "tên phòng chiếu không được để trống",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeTheaterRepo{}
			service := &theaterService{repo: repo}

			err := service.CreateRoom(context.Background(), tt.input)
			assertAppError(t, err, "bad_request", tt.wantErr)

			if repo.createRoomCalled {
				t.Fatal("expected CreateRoom not to be called")
			}
		})
	}
}

func TestTheaterServiceCreateRoomRepoError(t *testing.T) {
	repoErr := errors.New("insert room failed")
	repo := &fakeTheaterRepo{createRoomErr: repoErr}
	service := &theaterService{repo: repo}
	room := &Room{TheatherID: 1, Name: "Room 1"}

	err := service.CreateRoom(context.Background(), room)
	assertAppError(t, err, "internal_error", "không thể tạo phòng chiếu mới")

	if !errors.Is(err, repoErr) {
		t.Fatalf("expected error %v, got %v", repoErr, err)
	}

	if !repo.createRoomCalled {
		t.Fatal("expected CreateRoom to be called")
	}
}

func TestTheaterServiceGetSeatsByRoomIDValidation(t *testing.T) {
	repo := &fakeTheaterRepo{}
	service := &theaterService{repo: repo}

	_, err := service.GetSeatsByRoomID(context.Background(), 0)
	assertAppError(t, err, "bad_request", "ID phòng chiếu không hợp lệ")

	if repo.getSeatsCalled {
		t.Fatal("expected GetSeatsByRoomID not to be called")
	}
}

func TestTheaterServiceGetSeatsByRoomIDSuccess(t *testing.T) {
	want := []Seat{{RoomID: 3, RowLine: "A", Number: 1, SeatType: "standard"}}
	repo := &fakeTheaterRepo{getSeatsResult: want}
	service := &theaterService{repo: repo}

	got, err := service.GetSeatsByRoomID(context.Background(), 3)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !repo.getSeatsCalled {
		t.Fatal("expected GetSeatsByRoomID to be called")
	}

	if repo.getSeatsRoomID != 3 {
		t.Fatalf("expected repo to receive room id 3, got %d", repo.getSeatsRoomID)
	}

	if len(got) != len(want) {
		t.Fatalf("expected %d seats, got %d", len(want), len(got))
	}
}

func TestTheaterServiceCreateSeatsSuccess(t *testing.T) {
	repo := &fakeTheaterRepo{}
	service := &theaterService{repo: repo}

	err := service.CreateSeats(context.Background(), 10, 3, 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !repo.createSeatsCalled {
		t.Fatal("expected CreateSeats to be called")
	}

	if len(repo.receivedSeats) != 6 {
		t.Fatalf("expected 6 seats, got %d", len(repo.receivedSeats))
	}

	first := repo.receivedSeats[0]
	if first.RoomID != 10 {
		t.Fatalf("expected first seat room id 10, got %d", first.RoomID)
	}
	if first.RowLine != "A" {
		t.Fatalf("expected first seat row A, got %s", first.RowLine)
	}
	if first.Number != 1 {
		t.Fatalf("expected first seat number 1, got %d", first.Number)
	}
	if first.SeatType != "standard" {
		t.Fatalf("expected first seat type standard, got %s", first.SeatType)
	}

	third := repo.receivedSeats[2]
	if third.RowLine != "B" {
		t.Fatalf("expected third seat row B, got %s", third.RowLine)
	}
	if third.Number != 1 {
		t.Fatalf("expected third seat number 1, got %d", third.Number)
	}
	if third.SeatType != "vip" {
		t.Fatalf("expected third seat type vip, got %s", third.SeatType)
	}

	last := repo.receivedSeats[5]
	if last.RowLine != "C" {
		t.Fatalf("expected last seat row C, got %s", last.RowLine)
	}
	if last.Number != 2 {
		t.Fatalf("expected last seat number 2, got %d", last.Number)
	}
	if last.SeatType != "vip" {
		t.Fatalf("expected last seat type vip, got %s", last.SeatType)
	}
}

func TestTheaterServiceCreateSeatsValidation(t *testing.T) {
	tests := []struct {
		name        string
		roomID      int
		rows        int
		seatsPerRow int
		wantErr     string
	}{
		{
			name:        "invalid room id",
			roomID:      0,
			rows:        3,
			seatsPerRow: 2,
			wantErr:     "ID phòng chiếu không hợp lệ",
		},
		{
			name:        "invalid rows",
			roomID:      10,
			rows:        0,
			seatsPerRow: 2,
			wantErr:     "số hàng và số ghế mỗi hàng phải lớn hơn 0",
		},
		{
			name:        "invalid seats per row",
			roomID:      10,
			rows:        3,
			seatsPerRow: 0,
			wantErr:     "số hàng và số ghế mỗi hàng phải lớn hơn 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeTheaterRepo{}
			service := &theaterService{repo: repo}

			err := service.CreateSeats(context.Background(), tt.roomID, tt.rows, tt.seatsPerRow)
			assertAppError(t, err, "bad_request", tt.wantErr)

			if repo.createSeatsCalled {
				t.Fatal("expected CreateSeats not to be called")
			}
		})
	}
}

func TestTheaterServiceCreateSeatsRepoError(t *testing.T) {
	repoErr := errors.New("database is down")
	repo := &fakeTheaterRepo{createSeatsErr: repoErr}
	service := &theaterService{repo: repo}

	err := service.CreateSeats(context.Background(), 10, 3, 2)
	assertAppError(t, err, "internal_error", "không thể tạo ghế")

	if !errors.Is(err, repoErr) {
		t.Fatalf("expected error %v, got %v", repoErr, err)
	}

	if !repo.createSeatsCalled {
		t.Fatal("expected CreateSeats to be called")
	}

	if len(repo.receivedSeats) != 6 {
		t.Fatalf("expected generated seats to still be passed to repo, got %d", len(repo.receivedSeats))
	}
}
