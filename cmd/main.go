package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/LDKhangg/cinema-booking-go/internal/movie"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dsn := "postgres://postgres:123456@localhost:5432/cinema_db?sslmode=disable"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Không thể kết nối DB: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal("DB Không phản hồi: %v", err)
	}
	fmt.Println("connected")

	movieRepo := movie.NewRepository(db)
	movieService := movie.NewService(movieRepo)
	movieHandler := movie.NewHandler(movieService)

	http.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			movieHandler.GetMovies(w, r)
		case http.MethodPost:
			movieHandler.CreateMovie(w, r)
		default:
			http.Error(w, "Method not allow", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/movies/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			movieHandler.GetMovieById(w, r)
		case http.MethodPut:
			movieHandler.UpdateMovie(w, r)
		case http.MethodDelete:
			movieHandler.DeleteMovie(w, r)
		default:
			http.Error(w, "Method not support", http.StatusMethodNotAllowed)
		}
	})

	port := ":8080"
	fmt.Printf("Server running on port: %s", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Error when start server %v", err)
	}
}
