package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LDKhangg/cinema-booking-go/config"
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
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.NewPostgresDB(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal("failed to connect db ", err)
	}
	defer db.Close()

	goose.SetBaseFS(migrations.FS)

	if err := goose.Up(db, "."); err != nil {
		log.Fatalf("Lỗi khi chạy Migration: %v", err)
	}
	fmt.Println("Đã kiểm tra và chạy Migration (nếu có) thành công!")

	fmt.Println("Connected DB")
	router := server.SetupRouter(db)

	log.Println("Server đang chạy tại", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, router); err != nil {
		log.Fatal("Khởi động server thất bại: ", err)
	}
}

