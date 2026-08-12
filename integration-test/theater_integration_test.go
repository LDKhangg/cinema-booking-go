package integration_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/LDKhangg/cinema-booking-go/internal/server"
	"github.com/LDKhangg/cinema-booking-go/internal/theater"
)

func TestTheaterRepositoryIntegration_CreateTheaterRoomAndSeats(t *testing.T) {
	db := newTestDB(t)
	resetTheaterTables(t, db)

	repo := theater.NewRepository(db)
	ctx := context.Background()

	createdAt := time.Now().UTC().Truncate(time.Second)
	theatre := &theater.Theater{
		Name:      "CGV Su Van Hanh",
		Address:   "11 Su Van Hanh",
		City:      "HCM",
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}

	if err := repo.CreateTheater(ctx, theatre); err != nil {
		t.Fatalf("CreateTheater() error = %v", err)
	}
	if theatre.ID <= 0 {
		t.Fatalf("CreateTheater() did not set ID, got %d", theatre.ID)
	}

	fetchedTheater, err := repo.GetTheaterByID(ctx, theatre.ID)
	if err != nil {
		t.Fatalf("GetTheaterByID() error = %v", err)
	}
	if fetchedTheater.Name != theatre.Name {
		t.Fatalf("GetTheaterByID() name = %q, want %q", fetchedTheater.Name, theatre.Name)
	}

	room := &theater.Room{
		TheatherID: theatre.ID,
		Name:       "Room 1",
		TotalSeats: 6,
		RoomType:   "2D",
		CreatedAt:  createdAt,
		UpdatedAt:  createdAt,
	}

	if err := repo.CreateRoom(ctx, room); err != nil {
		t.Fatalf("CreateRoom() error = %v", err)
	}
	if room.ID <= 0 {
		t.Fatalf("CreateRoom() did not set ID, got %d", room.ID)
	}

	rooms, err := repo.GetRoomsByTheaterID(ctx, theatre.ID)
	if err != nil {
		t.Fatalf("GetRoomsByTheaterID() error = %v", err)
	}
	if len(rooms) != 1 {
		t.Fatalf("GetRoomsByTheaterID() len = %d, want 1", len(rooms))
	}
	if rooms[0].ID != room.ID {
		t.Fatalf("GetRoomsByTheaterID() room ID = %d, want %d", rooms[0].ID, room.ID)
	}

	seatsToCreate := []theater.Seat{
		{RoomID: room.ID, RowLine: "A", Number: 1, SeatType: "standard"},
		{RoomID: room.ID, RowLine: "A", Number: 2, SeatType: "standard"},
		{RoomID: room.ID, RowLine: "B", Number: 1, SeatType: "vip"},
	}

	if err := repo.CreateSeats(ctx, seatsToCreate); err != nil {
		t.Fatalf("CreateSeats() error = %v", err)
	}

	seats, err := repo.GetSeatsByRoomID(ctx, room.ID)
	if err != nil {
		t.Fatalf("GetSeatsByRoomID() error = %v", err)
	}
	if len(seats) != len(seatsToCreate) {
		t.Fatalf("GetSeatsByRoomID() len = %d, want %d", len(seats), len(seatsToCreate))
	}

	for i, want := range seatsToCreate {
		got := seats[i]
		if got.RoomID != want.RoomID || got.RowLine != want.RowLine || got.Number != want.Number || got.SeatType != want.SeatType {
			t.Fatalf("seat[%d] = %+v, want %+v", i, got, want)
		}
	}
}

func TestTheaterHTTPIntegration_CreateRoomAndSeatsFlow(t *testing.T) {
	db := newTestDB(t)
	resetTheaterTables(t, db)

	router := server.SetupRouter(db)

	theaterBody := performJSONRequest(t, router, http.MethodPost, "/theaters", map[string]any{
		"name":    "Lotte Cinema",
		"address": "469 Nguyen Huu Tho",
		"city":    "Da Nang",
	})

	var createdTheater theater.Theater
	decodeJSONResponse(t, theaterBody, &createdTheater)
	if createdTheater.ID <= 0 {
		t.Fatalf("create theater response ID = %d, want > 0", createdTheater.ID)
	}

	roomBody := performJSONRequest(t, router, http.MethodPost, "/theaters/rooms", map[string]any{
		"theather_id": createdTheater.ID,
		"name":        "IMAX",
		"total_seats": 6,
		"room_type":   "3D",
	})

	var createdRoom theater.Room
	decodeJSONResponse(t, roomBody, &createdRoom)
	if createdRoom.ID <= 0 {
		t.Fatalf("create room response ID = %d, want > 0", createdRoom.ID)
	}

	messageBody := performJSONRequest(t, router, http.MethodPost, fmt.Sprintf("/rooms/%d/seats", createdRoom.ID), map[string]any{
		"rows":          3,
		"seats_per_row": 2,
	})
	if strings.TrimSpace(messageBody) != `{"message":"Tạo ghế thành công"}` {
		t.Fatalf("create seats body = %s, want success message", strings.TrimSpace(messageBody))
	}

	seatsBody := performRequest(t, router, http.MethodGet, fmt.Sprintf("/rooms/%d/seats", createdRoom.ID), nil, http.StatusOK)
	var seats []theater.Seat
	decodeJSONResponse(t, seatsBody, &seats)

	if len(seats) != 6 {
		t.Fatalf("GET seats len = %d, want 6", len(seats))
	}

	checks := []struct {
		index    int
		rowLine  string
		number   int
		seatType string
	}{
		{index: 0, rowLine: "A", number: 1, seatType: "standard"},
		{index: 1, rowLine: "A", number: 2, seatType: "standard"},
		{index: 2, rowLine: "B", number: 1, seatType: "vip"},
		{index: 3, rowLine: "B", number: 2, seatType: "vip"},
		{index: 4, rowLine: "C", number: 1, seatType: "vip"},
		{index: 5, rowLine: "C", number: 2, seatType: "vip"},
	}

	for _, check := range checks {
		got := seats[check.index]
		if got.RowLine != check.rowLine || got.Number != check.number || got.SeatType != check.seatType {
			t.Fatalf("seat[%d] = %+v, want row=%s number=%d type=%s", check.index, got, check.rowLine, check.number, check.seatType)
		}
	}
}
