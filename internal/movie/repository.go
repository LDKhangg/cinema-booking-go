package movie

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrMovieNotFound = errors.New("Can not find movie")

type postgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{
		db: db,
	}
}

func (p *postgresRepository) Create(ctx context.Context, m *Movie) error {
	query := `INSERT INTO movies(
		title,
		description,
		release_date,
		closed_date,
		is_published,
		created_at, 
		updated_at)
					VALUES($1,$2,$3,$4,$5,$6,$7)					
					RETURNING id
`
	err := p.db.QueryRowContext(
		ctx,
		query,
		m.Title,
		m.Description,
		m.ReleaseDate,
		m.ClosedDate,
		m.IsPublished,
		m.CreatedAt,
		m.UpdatedAt,
	).Scan(&m.ID)
	if err != nil {
		return fmt.Errorf("Cant not create movie: %w", err)
	}
	return nil
}

func (p *postgresRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM movies WHERE id = $1`
	result, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("Can not remove movie: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrMovieNotFound
	}
	return nil
}

func (p *postgresRepository) GetByID(ctx context.Context, id int) (Movie, error) {
	query := `
		SELECT id, title, description, release_date, closed_date, is_published, created_at, updated_at
		FROM movies
		WHERE id = $1
	`
	var m Movie
	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID,
		&m.Title,
		&m.Description,
		&m.ReleaseDate,
		&m.ClosedDate,
		&m.IsPublished,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Movie{}, ErrMovieNotFound
		}
		return Movie{}, fmt.Errorf("invalid query claues: %w", err)
	}
	return m, nil
}

func (p *postgresRepository) GetPublishedMovies(ctx context.Context) ([]Movie, error) {
	query := `
	SELECT id, title, description, release_date, closed_date, is_published, created_at, updated_at
	FROM movies
	WHERE is_published = true
	ORDER BY release_date
`
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error get list: %w", err)
	}
	defer rows.Close()
	var movies []Movie
	for rows.Next() {
		var m Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.Description, &m.ReleaseDate, &m.ClosedDate, &m.IsPublished, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return movies, nil
}

// Update implements [Repository].
func (p *postgresRepository) Update(ctx context.Context, m *Movie) error {
	query := `
		UPDATE movies
		SET title = $1, description = $2, release_date = $3, closed_date = $4, is_published = $5, updated_at = $6
		WHERE id = $7
	`
	result, err := p.db.ExecContext(
		ctx, query,
		m.Title, m.Description, m.ReleaseDate, m.ClosedDate, m.IsPublished, m.UpdatedAt, m.ID,
	)
	if err != nil {
		return fmt.Errorf("không thể cập nhật phim: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrMovieNotFound
	}

	return nil
}
