package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LDKhangg/cinema-booking-go/internal/server"
	"github.com/LDKhangg/cinema-booking-go/pkg/database"
)

func main() {
	dsn := "postgres://postgres:123456@localhost:5432/cinema_db?sslmode=disable"
	db, err := database.NewPostgresDB(dsn)
	if err != nil {
		log.Fatalf("Không thể kết nối DB: %v", err)
	}

	defer db.Close()
	fmt.Println("Connected DB")
	route := server.SetupRouter(db)
	port := ":8080"
	fmt.Printf("Server running on port: %s", port)

	if err := http.ListenAndServe(port, route); err != nil {
		log.Fatalf("Error when start server %v", err)
	}
}
