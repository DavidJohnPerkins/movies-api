package store

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/microsoft/go-mssqldb"
)

const driverName = "sqlserver"

type SqlServerMoviesStore struct {
	databaseUrl string
	dbx         *sqlx.DB
}

func NewSqlServerMoviesStore(databaseUrl string) *SqlServerMoviesStore {
	return &SqlServerMoviesStore{
		databaseUrl: databaseUrl,
	}
}

func noOpMapper(s string) string {
	return s
}

func (s *SqlServerMoviesStore) connect(ctx context.Context) error {
	dbx, err := sqlx.ConnectContext(ctx, driverName, s.databaseUrl)
	if err != nil {
		log.Printf("DB connect failed: %v", err)
		return err
	}

	dbx.MapperFunc(noOpMapper)
	s.dbx = dbx
	return nil
}

func (s *SqlServerMoviesStore) close() error {
	return s.dbx.Close()
}

func (s *SqlServerMoviesStore) GetAll(ctx context.Context) ([]Movie, error) {
	err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer s.close()

	var movies []Movie
	if err := s.dbx.SelectContext(
		ctx,
		&movies,
		`SELECT
			Id, Title, Director, ReleaseDate, TicketPrice, CreatedAt, UpdatedAt 
			FROM dbo.movies`); err != nil {
		return nil, err
	}
	return movies, nil

}

func (s *SqlServerMoviesStore) GetByID(ctx context.Context, id uuid.UUID) (Movie, error) {
	err := s.connect(ctx)
	if err != nil {
		return Movie{}, err
	}
	defer s.close()

	var movie Movie
	if err := s.dbx.GetContext(
		ctx,
		&movie,
		`SELECT Id, Title, Director, ReleaseDate, TicketPrice, CreatedAt, UpdatedAt
		FROM dbo.movies
		WHERE Id = @id`,
		sql.Named("id", id)); err != nil {
		if err != sql.ErrNoRows {
			return Movie{}, err
		}
		return Movie{}, &RecordNotFoundError{}
	}

	return movie, nil
}

func (s *SqlServerMoviesStore) Create(ctx context.Context, createMovieParams CreateMovieParams) error {
	err := s.connect(ctx)
	if err != nil {
		return err
	}
	defer s.close()

	movie := Movie{
		ID:          createMovieParams.ID,
		Title:       createMovieParams.Title,
		Director:    createMovieParams.Director,
		ReleaseDate: createMovieParams.ReleaseDate,
		TicketPrice: createMovieParams.TicketPrice,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if _, err := s.dbx.NamedExecContext(
		ctx,
		`INSERT INTO dbo.movies
			(Id, Title, Director, ReleaseDate, TicketPrice, CreatedAt, UpdatedAt)
		VALUES
			(:Id, :Title, :Director, :ReleaseDate, :TicketPrice, :CreatedAt, :UpdatedAt)`,
		movie); err != nil {
		if strings.Contains(err.Error(), "Cannot insert duplicate key") {
			return &DuplicateKeyError{ID: createMovieParams.ID}
		}
		return err
	}

	return nil
}

func (s *SqlServerMoviesStore) Update(ctx context.Context, id uuid.UUID, updateMovieParams UpdateMovieParams) error {
	err := s.connect(ctx)
	if err != nil {
		return err
	}
	defer s.close()

	movie := Movie{
		ID:          id,
		Title:       updateMovieParams.Title,
		Director:    updateMovieParams.Director,
		ReleaseDate: updateMovieParams.ReleaseDate,
		TicketPrice: updateMovieParams.TicketPrice,
		UpdatedAt:   time.Now().UTC(),
	}
	log.Printf("movie: %v", movie)
	if _, err := s.dbx.NamedExecContext(
		ctx,
		`UPDATE dbo.movies
		SET Title = :Title, Director = :Director, ReleaseDate = :ReleaseDate, TicketPrice = :TicketPrice, UpdatedAt = :UpdatedAt
		WHERE Id = :Id`,
		movie); err != nil {
		return err
	}

	return nil
}

func (s *SqlServerMoviesStore) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.connect(ctx)
	if err != nil {
		return err
	}
	defer s.close()

	if _, err := s.dbx.ExecContext(
		ctx,
		`DELETE FROM dbo.movies
		WHERE id = @id`, sql.Named("id", id)); err != nil {
		return err
	}

	return nil
}
