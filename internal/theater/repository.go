package theather

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrTheaterNotFound = errors.New("không tìm thấy rạp chiếu")
	ErrRoomNotFound    = errors.New("không tìm thấy phòng chiếu")
	ErrSeatNotFound    = errors.New("không tìm thấy ghế")
)

type postgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) GetTheaters(ctx context.Context) ([]Theater, error) {
	query := `
		SELECT id, name, address, city, created_at, updated_at
		FROM theaters
		ORDER BY id DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy danh sách rạp: %w", err)
	}
	defer rows.Close()

	var theaters []Theater
	for rows.Next() {
		var t Theater
		if err := rows.Scan(&t.ID, &t.Name, &t.Address, &t.City, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("không thể đọc dữ liệu rạp: %w", err)
		}
		theaters = append(theaters, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return theaters, nil
}

func (r *postgresRepository) GetTheaterByID(ctx context.Context, id int) (Theater, error) {
	query := `
		SELECT id, name, address, city, created_at, updated_at
		FROM theaters
		WHERE id = $1
	`
	var t Theater
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.Name, &t.Address, &t.City, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Theater{}, ErrTheaterNotFound
		}
		return Theater{}, fmt.Errorf("không thể lấy rạp theo ID: %w", err)
	}
	return t, nil
}

func (r *postgresRepository) CreateTheater(ctx context.Context, t *Theater) error {
	query := `
		INSERT INTO theaters (name, address, city, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := r.db.QueryRowContext(
		ctx, query,
		t.Name, t.Address, t.City, t.CreatedAt, t.UpdatedAt,
	).Scan(&t.ID)
	if err != nil {
		return fmt.Errorf("không thể tạo rạp mới: %w", err)
	}
	return nil
}

func (r *postgresRepository) GetRoomsByTheaterID(ctx context.Context, theaterID int) ([]Room, error) {
	query := `
		SELECT id, theater_id, name, total_seats, room_type, created_at, updated_at
		FROM rooms
		WHERE theater_id = $1
		ORDER BY id DESC
	`
	rows, err := r.db.QueryContext(ctx, query, theaterID)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy danh sách phòng chiếu: %w", err)
	}
	defer rows.Close()

	var rooms []Room
	for rows.Next() {
		var rm Room
		if err := rows.Scan(&rm.ID, &rm.TheatherID, &rm.Name, &rm.TotalSeats, &rm.RoomType, &rm.CreatedAt, &rm.UpdatedAt); err != nil {
			return nil, fmt.Errorf("không thể đọc dữ liệu phòng: %w", err)
		}
		rooms = append(rooms, rm)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (r *postgresRepository) CreateRoom(ctx context.Context, rm *Room) error {
	query := `
		INSERT INTO rooms (theater_id, name, total_seats, room_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err := r.db.QueryRowContext(
		ctx, query,
		rm.TheatherID, rm.Name, rm.TotalSeats, rm.RoomType, rm.CreatedAt, rm.UpdatedAt,
	).Scan(&rm.ID)
	if err != nil {
		return fmt.Errorf("không thể tạo phòng chiếu mới: %w", err)
	}
	return nil
}

func (r *postgresRepository) GetSeatsByRoomID(ctx context.Context, roomID int) ([]Seat, error) {
	query := `
		SELECT id, room_id, row_line, number, seat_type
		FROM seats
		WHERE room_id = $1
		ORDER BY row_line, number
	`
	rows, err := r.db.QueryContext(ctx, query, roomID)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy danh sách ghế: %w", err)
	}
	defer rows.Close()

	var seats []Seat
	for rows.Next() {
		var s Seat
		if err := rows.Scan(&s.ID, &s.RoomID, &s.RowLine, &s.Number, &s.SeatType); err != nil {
			return nil, fmt.Errorf("không thể đọc dữ liệu ghế: %w", err)
		}
		seats = append(seats, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return seats, nil
}

func (r *postgresRepository) CreateSeats(ctx context.Context, seats []Seat) error {
	if len(seats) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("không thể bắt đầu transaction tạo ghế: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO seats (room_id, row_line, number, seat_type)
		VALUES ($1, $2, $3, $4)
	`)
	if err != nil {
		return fmt.Errorf("không thể chuẩn bị câu lệnh tạo ghế: %w", err)
	}
	defer stmt.Close()

	for _, s := range seats {
		_, err := stmt.ExecContext(ctx, s.RoomID, s.RowLine, s.Number, s.SeatType)
		if err != nil {
			return fmt.Errorf("không thể chèn ghế: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("không thể commit transaction tạo ghế: %w", err)
	}

	return nil
}
