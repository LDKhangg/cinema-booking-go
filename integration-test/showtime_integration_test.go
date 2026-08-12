package integration_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/LDKhangg/cinema-booking-go/internal/movie"
	"github.com/LDKhangg/cinema-booking-go/internal/server"
	"github.com/LDKhangg/cinema-booking-go/internal/showtime"
	"github.com/LDKhangg/cinema-booking-go/internal/theater"
)

func TestShowtimeRepositoryIntegration_CreateAndQueryByMovieRoomAndTheater(t *testing.T) {
	db := newTestDB(t)
	resetAllTables(t, db)

	movieRepo := movie.NewRepository(db)
	theaterRepo := theater.NewRepository(db)
	showtimeRepo := showtime.NewRepository(db)
	ctx := context.Background()

	movieID, roomID, theaterID, targetDate := seedShowtimeDependencies(t, ctx, movieRepo, theaterRepo)

	startTime := targetDate.Add(2 * time.Hour)
	endTime := startTime.Add(2 * time.Hour)
	st := &showtime.Showtime{
		MovieID:   movieID,
		RoomID:    roomID,
		StartTime: startTime,
		EndTime:   endTime,
		Price:     95000,
		CreatedAt: startTime.Add(-30 * time.Minute),
		UpdatedAt: startTime.Add(-30 * time.Minute),
	}

	if err := showtimeRepo.Create(ctx, st); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if st.ID <= 0 {
		t.Fatalf("Create() did not set ID, got %d", st.ID)
	}

	byID, err := showtimeRepo.GetByID(ctx, st.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if byID.MovieID != movieID || byID.RoomID != roomID {
		t.Fatalf("GetByID() = %+v, want movieID=%d roomID=%d", byID, movieID, roomID)
	}

	byMovie, err := showtimeRepo.GetByMovieID(ctx, movieID, targetDate)
	if err != nil {
		t.Fatalf("GetByMovieID() error = %v", err)
	}
	if len(byMovie) != 1 || byMovie[0].ID != st.ID {
		t.Fatalf("GetByMovieID() = %+v, want one showtime with ID %d", byMovie, st.ID)
	}

	byRoom, err := showtimeRepo.GetByRoomID(ctx, roomID, targetDate)
	if err != nil {
		t.Fatalf("GetByRoomID() error = %v", err)
	}
	if len(byRoom) != 1 || byRoom[0].ID != st.ID {
		t.Fatalf("GetByRoomID() = %+v, want one showtime with ID %d", byRoom, st.ID)
	}

	theaterLookupRepo, ok := showtimeRepo.(interface {
		GetByTheaterID(ctx context.Context, theaterID int, targetDate time.Time) ([]showtime.Showtime, error)
	})
	if !ok {
		t.Fatal("showtime repository does not implement GetByTheaterID")
	}

	byTheater, err := theaterLookupRepo.GetByTheaterID(ctx, theaterID, targetDate)
	if err != nil {
		t.Fatalf("GetByTheaterID() error = %v", err)
	}
	if len(byTheater) != 1 || byTheater[0].ID != st.ID {
		t.Fatalf("GetByTheaterID() = %+v, want one showtime with ID %d", byTheater, st.ID)
	}
}

func TestShowtimeHTTPIntegration_CreateAndQueryFlow(t *testing.T) {
	db := newTestDB(t)
	resetAllTables(t, db)

	router := server.SetupRouter(db)
	movieID, roomID, theaterID, targetDate := seedShowtimeDependencies(
		t,
		context.Background(),
		movie.NewRepository(db),
		theater.NewRepository(db),
	)

	startTime := targetDate.Add(3 * time.Hour)
	endTime := startTime.Add(2 * time.Hour)

	createBody := performJSONRequest(t, router, http.MethodPost, "/showtimes", map[string]any{
		"movie_id":   movieID,
		"room_id":    roomID,
		"start_time": startTime.Format(time.RFC3339),
		"end_time":   endTime.Format(time.RFC3339),
		"price":      120000,
	})

	var createdShowtime showtime.Showtime
	decodeJSONResponse(t, createBody, &createdShowtime)
	if createdShowtime.ID <= 0 {
		t.Fatalf("create showtime response ID = %d, want > 0", createdShowtime.ID)
	}

	getBody := performRequest(t, router, http.MethodGet, fmt.Sprintf("/showtimes/%d", createdShowtime.ID), nil, http.StatusOK)
	var fetched showtime.Showtime
	decodeJSONResponse(t, getBody, &fetched)
	if fetched.ID != createdShowtime.ID {
		t.Fatalf("GET showtime ID = %d, want %d", fetched.ID, createdShowtime.ID)
	}

	byMovieBody := performRequest(t, router, http.MethodGet, fmt.Sprintf("/showtimes?movie_id=%d&date=%s", movieID, targetDate.Format("2006-01-02")), nil, http.StatusOK)
	var movieShowtimes []showtime.Showtime
	decodeJSONResponse(t, byMovieBody, &movieShowtimes)
	if len(movieShowtimes) != 1 || movieShowtimes[0].ID != createdShowtime.ID {
		t.Fatalf("GET /showtimes by movie = %+v, want one showtime with ID %d", movieShowtimes, createdShowtime.ID)
	}

	byTheaterBody := performRequest(t, router, http.MethodGet, fmt.Sprintf("/showtimes?theater_id=%d&date=%s", theaterID, targetDate.Format("2006-01-02")), nil, http.StatusOK)
	var theaterShowtimes []showtime.Showtime
	decodeJSONResponse(t, byTheaterBody, &theaterShowtimes)
	if len(theaterShowtimes) != 1 || theaterShowtimes[0].ID != createdShowtime.ID {
		t.Fatalf("GET /showtimes by theater = %+v, want one showtime with ID %d", theaterShowtimes, createdShowtime.ID)
	}
}

func TestShowtimeHTTPIntegration_RejectsOverlappingShowtime(t *testing.T) {
	db := newTestDB(t)
	resetAllTables(t, db)

	router := server.SetupRouter(db)
	movieID, roomID, _, targetDate := seedShowtimeDependencies(
		t,
		context.Background(),
		movie.NewRepository(db),
		theater.NewRepository(db),
	)

	startTime := targetDate.Add(5 * time.Hour)
	endTime := startTime.Add(2 * time.Hour)

	_ = performJSONRequest(t, router, http.MethodPost, "/showtimes", map[string]any{
		"movie_id":   movieID,
		"room_id":    roomID,
		"start_time": startTime.Format(time.RFC3339),
		"end_time":   endTime.Format(time.RFC3339),
		"price":      90000,
	})

	overlapBody := performJSONRequestWithStatus(t, router, http.MethodPost, "/showtimes", map[string]any{
		"movie_id":   movieID,
		"room_id":    roomID,
		"start_time": startTime.Add(30 * time.Minute).Format(time.RFC3339),
		"end_time":   endTime.Add(30 * time.Minute).Format(time.RFC3339),
		"price":      100000,
	}, http.StatusInternalServerError)

	if strings.TrimSpace(overlapBody) != `{"error":"lịch chiếu bị trùng với suất chiếu khác trong cùng phòng"}` {
		t.Fatalf("overlap body = %s, want overlap error", strings.TrimSpace(overlapBody))
	}
}
func seedShowtimeDependencies(t *testing.T, ctx context.Context, movieRepo movie.Repository, theaterRepo theater.Repository) (int, int, int, time.Time) {
	t.Helper()

	targetDate := time.Now().Add(24 * time.Hour).UTC().Truncate(24 * time.Hour)
	m := &movie.Movie{
		Title:       "The Dark Knight",
		Description: "Batman",
		ReleaseDate: targetDate,
		ClosedDate:  targetDate.Add(7 * 24 * time.Hour),
		IsPublished: true,
		CreatedAt:   targetDate.Add(-2 * time.Hour),
		UpdatedAt:   targetDate.Add(-2 * time.Hour),
	}
	if err := movieRepo.Create(ctx, m); err != nil {
		t.Fatalf("seed movie error = %v", err)
	}

	th := &theater.Theater{
		Name:      "Galaxy",
		Address:   "1 Nguyen Hue",
		City:      "HCM",
		CreatedAt: targetDate.Add(-2 * time.Hour),
		UpdatedAt: targetDate.Add(-2 * time.Hour),
	}
	if err := theaterRepo.CreateTheater(ctx, th); err != nil {
		t.Fatalf("seed theater error = %v", err)
	}

	rm := &theater.Room{
		TheatherID: th.ID,
		Name:       "Room Showtime",
		TotalSeats: 20,
		RoomType:   "2D",
		CreatedAt:  targetDate.Add(-2 * time.Hour),
		UpdatedAt:  targetDate.Add(-2 * time.Hour),
	}
	if err := theaterRepo.CreateRoom(ctx, rm); err != nil {
		t.Fatalf("seed room error = %v", err)
	}

	return m.ID, rm.ID, th.ID, targetDate
}
