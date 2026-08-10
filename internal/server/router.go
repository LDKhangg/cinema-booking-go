package server

import (
	"database/sql"
	"net/http"

	"github.com/LDKhangg/cinema-booking-go/internal/movie"
	"github.com/LDKhangg/cinema-booking-go/internal/showtime"
	"github.com/LDKhangg/cinema-booking-go/internal/theater"
)

func SetupRouter(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	// Movie routes
	movieRepo := movie.NewRepository(db)
	movieService := movie.NewService(movieRepo)
	movieHandler := movie.NewHandler(movieService)

	mux.HandleFunc("GET /movies", movieHandler.GetMovies)
	mux.HandleFunc("POST /movies", movieHandler.CreateMovie)
	mux.HandleFunc("GET /movies/{id}", movieHandler.GetMovieById)
	mux.HandleFunc("PUT /movies/{id}", movieHandler.UpdateMovie)
	mux.HandleFunc("DELETE /movies/{id}", movieHandler.DeleteMovie)

	// Theater routes
	theaterRepo := theather.NewRepository(db)
	theaterService := theather.NewService(theaterRepo)
	theaterHandler := theather.NewHandler(theaterService)

	mux.HandleFunc("GET /theaters", theaterHandler.GetTheaters)
	mux.HandleFunc("POST /theaters", theaterHandler.CreateTheater)
	mux.HandleFunc("GET /theaters/{id}", theaterHandler.GetTheaterByID)
	mux.HandleFunc("GET /theaters/{id}/rooms", theaterHandler.GetRoomsByTheaterID)
	mux.HandleFunc("POST /theaters/rooms", theaterHandler.CreateRoom)
	mux.HandleFunc("POST /rooms/{id}/seats", theaterHandler.CreateSeats)
	mux.HandleFunc("GET /rooms/{id}/seats", theaterHandler.GetSeatsByRoomID)

	// Showtime routes
	showtimeRepo := showtime.NewRepository(db)
	showtimeService := showtime.NewService(showtimeRepo)
	showtimeHandler := showtime.NewHandler(showtimeService)

	mux.HandleFunc("GET /showtimes", showtimeHandler.GetShowtimes)
	mux.HandleFunc("POST /showtimes", showtimeHandler.CreateShowtime)
	mux.HandleFunc("GET /showtimes/{id}", showtimeHandler.GetShowtimeByID)

	return mux
}
