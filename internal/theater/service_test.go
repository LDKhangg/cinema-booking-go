package theater

import (
	"context"
	"errors"
	"testing"
)

type fakeTheaterRepo struct {
	createSeatsCalled bool
	receivedSeats     []Seat
	createSeatsErr    error
}

func (f *fakeTheaterRepo) GetTheaters(ctx context.Context) ([]Theater, error) {
	panic("unexpected call to GetTheaters")
}

func (f *fakeTheaterRepo) GetTheaterByID(ctx context.Context, id int) (Theater, error) {
	panic("unexpected call to GetTheaterByID")
}

func (f *fakeTheaterRepo) CreateTheater(ctx context.Context, t *Theater) error {
	panic("unexpected call to CreateTheater")
}

func (f *fakeTheaterRepo) GetRoomsByTheaterID(ctx context.Context, theaterID int) ([]Room, error) {
	panic("unexpected call to GetRoomsByTheaterID")
}

func (f *fakeTheaterRepo) CreateRoom(ctx context.Context, r *Room) error {
	panic("unexpected call to CreateRoom")
}

func (f *fakeTheaterRepo) GetSeatsByRoomID(ctx context.Context, roomID int) ([]Seat, error) {
	panic("unexpected call to GetSeatsByRoomID")
}

func (f *fakeTheaterRepo) CreateSeats(ctx context.Context, seats []Seat) error {
	f.createSeatsCalled = true
	f.receivedSeats = seats
	return f.createSeatsErr
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
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if err.Error() != tt.wantErr {
				t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
			}

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
	if err == nil {
		t.Fatal("expected error, got nil")
	}

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
