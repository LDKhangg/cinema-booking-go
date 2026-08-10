package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/LDKhangg/cinema-booking-go/docs"
	"github.com/LDKhangg/cinema-booking-go/internal/server"
	"github.com/LDKhangg/cinema-booking-go/migrations"
	"github.com/LDKhangg/cinema-booking-go/pkg/database"
	"github.com/pressly/goose/v3"
)

// @title Cinema Booking API
// @version 1.0
// @description API quản lý đặt vé xem phim: phim, rạp chiếu, phòng chiếu, ghế và suất chiếu
// @host localhost:8080
// @BasePath /
func main() {
	dsn := "postgres://postgres:123456@localhost:5432/cinema_db?sslmode=disable"
	db, err := database.NewPostgresDB(dsn)
	if err != nil {
		log.Fatalf("Không thể kết nối DB: %v", err)
	}

	goose.SetBaseFS(migrations.FS)

	if err := goose.Up(db, "."); err != nil {
		log.Fatalf("Lỗi khi chạy Migration: %v", err)
	}
	fmt.Println("Đã kiểm tra và chạy Migration (nếu có) thành công!")

	defer db.Close()
	fmt.Println("Connected DB")
	route := server.SetupRouter(db)
	port := ":8080"
	fmt.Printf("Server running on port: %s", port)

	if err := http.ListenAndServe(port, route); err != nil {
		log.Fatalf("Error when start server %v", err)
	}
}