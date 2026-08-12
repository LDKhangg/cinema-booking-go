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
)

func TestMovieRepositoryIntegration_CreateUpdateGetAndDelete(t *testing.T) {
	db := newTestDB(t)
	resetAllTables(t, db)

	repo := movie.NewRepository(db)
	ctx := context.Background()

	releaseDate := time.Date(2026, time.January, 10, 9, 0, 0, 0, time.UTC)
	closedDate := releaseDate.Add(72 * time.Hour)
	createdAt := releaseDate.Add(-24 * time.Hour)

	m := &movie.Movie{
		Title:       "Interstellar",
		Description: "Space exploration",
		ReleaseDate: releaseDate,
		ClosedDate:  closedDate,
		IsPublished: true,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}

	if err := repo.Create(ctx, m); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if m.ID <= 0 {
		t.Fatalf("Create() did not set ID, got %d", m.ID)
	}

	fetched, err := repo.GetByID(ctx, m.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if fetched.Title != m.Title || fetched.Description != m.Description {
		t.Fatalf("GetByID() = %+v, want title=%q description=%q", fetched, m.Title, m.Description)
	}

	updatedAt := createdAt.Add(2 * time.Hour)
	m.Title = "Interstellar Remastered"
	m.Description = "Updated description"
	m.UpdatedAt = updatedAt

	if err := repo.Update(ctx, m); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	updatedMovie, err := repo.GetByID(ctx, m.ID)
	if err != nil {
		t.Fatalf("GetByID() after update error = %v", err)
	}
	if updatedMovie.Title != m.Title || updatedMovie.Description != m.Description {
		t.Fatalf("updated movie = %+v, want title=%q description=%q", updatedMovie, m.Title, m.Description)
	}

	publishedMovies, err := repo.GetPublishedMovies(ctx)
	if err != nil {
		t.Fatalf("GetPublishedMovies() error = %v", err)
	}
	if len(publishedMovies) != 1 {
		t.Fatalf("GetPublishedMovies() len = %d, want 1", len(publishedMovies))
	}
	if publishedMovies[0].ID != m.ID {
		t.Fatalf("GetPublishedMovies()[0].ID = %d, want %d", publishedMovies[0].ID, m.ID)
	}

	if err := repo.Delete(ctx, m.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	if _, err := repo.GetByID(ctx, m.ID); err == nil {
		t.Fatal("GetByID() after delete error = nil, want not found")
	}
}

func TestMovieHTTPIntegration_CreateUpdateGetListAndDeleteFlow(t *testing.T) {
	db := newTestDB(t)
	resetAllTables(t, db)

	router := server.SetupRouter(db)

	releaseDate := time.Date(2026, time.February, 1, 10, 0, 0, 0, time.UTC)
	closedDate := releaseDate.Add(96 * time.Hour)

	createBody := performJSONRequest(t, router, http.MethodPost, "/movies", map[string]any{
		"title":        "Dune Part Two",
		"description":  "Sci-fi epic",
		"release_date": releaseDate.Format(time.RFC3339),
		"closed_date":  closedDate.Format(time.RFC3339),
		"is_published": true,
	})

	var createdMovie movie.Movie
	decodeJSONResponse(t, createBody, &createdMovie)
	if createdMovie.ID <= 0 {
		t.Fatalf("create movie response ID = %d, want > 0", createdMovie.ID)
	}
	if createdMovie.Title != "Dune Part Two" {
		t.Fatalf("create movie title = %q, want %q", createdMovie.Title, "Dune Part Two")
	}

	getBody := performRequest(t, router, http.MethodGet, fmt.Sprintf("/movies/%d", createdMovie.ID), nil, http.StatusOK)
	var fetchedMovie movie.Movie
	decodeJSONResponse(t, getBody, &fetchedMovie)
	if fetchedMovie.ID != createdMovie.ID {
		t.Fatalf("GET movie ID = %d, want %d", fetchedMovie.ID, createdMovie.ID)
	}

	updateBody := performJSONRequestWithStatus(t, router, http.MethodPut, fmt.Sprintf("/movies/%d", createdMovie.ID), map[string]any{
		"title":        "Dune Part Two Extended",
		"description":  "Sci-fi epic extended cut",
		"release_date": releaseDate.Format(time.RFC3339),
		"closed_date":  closedDate.Add(24 * time.Hour).Format(time.RFC3339),
		"is_published": true,
	}, http.StatusOK)
	if strings.TrimSpace(updateBody) != `{"message":"Cập nhật phim thành công"}` {
		t.Fatalf("update movie body = %s, want success message", strings.TrimSpace(updateBody))
	}

	updatedGetBody := performRequest(t, router, http.MethodGet, fmt.Sprintf("/movies/%d", createdMovie.ID), nil, http.StatusOK)
	var updatedMovie movie.Movie
	decodeJSONResponse(t, updatedGetBody, &updatedMovie)
	if updatedMovie.Title != "Dune Part Two Extended" {
		t.Fatalf("updated movie title = %q, want %q", updatedMovie.Title, "Dune Part Two Extended")
	}

	listBody := performRequest(t, router, http.MethodGet, "/movies", nil, http.StatusOK)
	var movies []movie.Movie
	decodeJSONResponse(t, listBody, &movies)
	if len(movies) != 1 {
		t.Fatalf("GET /movies len = %d, want 1", len(movies))
	}
	if movies[0].ID != createdMovie.ID {
		t.Fatalf("GET /movies first ID = %d, want %d", movies[0].ID, createdMovie.ID)
	}

	unpublishBody := performJSONRequestWithStatus(t, router, http.MethodPut, fmt.Sprintf("/movies/%d", createdMovie.ID), map[string]any{
		"title":        "Dune Part Two Extended",
		"description":  "Sci-fi epic extended cut",
		"release_date": releaseDate.Format(time.RFC3339),
		"closed_date":  closedDate.Add(24 * time.Hour).Format(time.RFC3339),
		"is_published": false,
	}, http.StatusOK)
	if strings.TrimSpace(unpublishBody) != `{"message":"Cập nhật phim thành công"}` {
		t.Fatalf("unpublish movie body = %s, want success message", strings.TrimSpace(unpublishBody))
	}

	deleteBody := performRequest(t, router, http.MethodDelete, fmt.Sprintf("/movies/%d", createdMovie.ID), nil, http.StatusOK)
	if strings.TrimSpace(deleteBody) != `{"message":"Xóa phim thành công"}` {
		t.Fatalf("delete movie body = %s, want success message", strings.TrimSpace(deleteBody))
	}

	missingBody := performRequest(t, router, http.MethodGet, fmt.Sprintf("/movies/%d", createdMovie.ID), nil, http.StatusNotFound)
	if strings.TrimSpace(missingBody) != `{"error":"can not find movie"}` {
		t.Fatalf("missing movie body = %s, want not found error", strings.TrimSpace(missingBody))
	}
}
