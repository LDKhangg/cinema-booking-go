package integration_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/LDKhangg/cinema-booking-go/migrations"
	"github.com/LDKhangg/cinema-booking-go/pkg/database"
	"github.com/pressly/goose/v3"
)

const defaultTestDSN = "postgres://postgres:123456@localhost:5432/cinema_db?sslmode=disable"

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("TEST_POSTGRES_DSN")
	if dsn == "" {
		dsn = defaultTestDSN
	}

	db, err := database.NewPostgresDB(dsn)
	if err != nil {
		t.Fatalf("NewPostgresDB() error = %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	goose.SetBaseFS(migrations.FS)
	if err := goose.Up(db, "."); err != nil {
		t.Fatalf("goose.Up() error = %v", err)
	}

	return db
}

func resetTheaterTables(t *testing.T, db *sql.DB) {
	t.Helper()

	query := `TRUNCATE TABLE seats, rooms, theaters RESTART IDENTITY CASCADE`
	if _, err := db.Exec(query); err != nil {
		t.Fatalf("reset theater tables error = %v", err)
	}
}

func resetAllTables(t *testing.T, db *sql.DB) {
	t.Helper()

	query := `TRUNCATE TABLE showtimes, seats, rooms, theaters, movies RESTART IDENTITY CASCADE`
	if _, err := db.Exec(query); err != nil {
		t.Fatalf("reset all tables error = %v", err)
	}
}

func performJSONRequest(t *testing.T, handler http.Handler, method string, path string, payload any) string {
	t.Helper()
	return performJSONRequestWithStatus(t, handler, method, path, payload, http.StatusCreated)
}

func performJSONRequestWithStatus(t *testing.T, handler http.Handler, method string, path string, payload any, wantStatus int) string {
	t.Helper()

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	return performRequest(t, handler, method, path, bytes.NewReader(body), wantStatus)
}

func performRequest(t *testing.T, handler http.Handler, method string, path string, body *bytes.Reader, wantStatus int) string {
	t.Helper()

	var requestBody *bytes.Reader
	if body == nil {
		requestBody = bytes.NewReader(nil)
	} else {
		requestBody = body
	}

	req := httptest.NewRequest(method, path, requestBody)
	if requestBody.Len() > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != wantStatus {
		t.Fatalf("%s %s status = %d, want %d, body = %s", method, path, rec.Code, wantStatus, strings.TrimSpace(rec.Body.String()))
	}

	return rec.Body.String()
}

func decodeJSONResponse(t *testing.T, body string, dst any) {
	t.Helper()

	if err := json.Unmarshal([]byte(body), dst); err != nil {
		t.Fatalf("json.Unmarshal() error = %v; body = %s", err, body)
	}
}
