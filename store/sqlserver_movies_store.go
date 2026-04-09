package store

import (
	"context"
	"database/sql"
	"log"
	"strings"

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
	r, err := s.dbx.QueryxContext(
		ctx, `
		EXEC dbo.r_movie;
	`)

	if err != nil {
		return nil, err
	}
	defer r.Close()

	for r.Next() {
		var m Movie
		if err := r.StructScan(&m); err != nil {
			return nil, err
		}
		movies = append(movies, m)
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
	r, err := s.dbx.QueryxContext(
		ctx, `
		EXEC dbo.r_movie_by_id @id=@id;
	`, sql.Named("id", id))
	if err != nil {
		return Movie{}, err
	}
	defer r.Close()

	if r.Next() {
		if err := r.StructScan(&movie); err != nil {
			return Movie{}, err
		}
	} else {
		return Movie{}, sql.ErrNoRows
	}

	return movie, nil
}

func (s *SqlServerMoviesStore) Create(ctx context.Context, jsonBody string) error {
	err := s.connect(ctx)
	if err != nil {
		return err
	}
	defer s.close()

	_, err = s.dbx.ExecContext(
		ctx,
		`EXEC dbo.c_movie @json = @json`,
		sql.Named("json", jsonBody),
	)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return &DuplicateKeyError{err}
		}
		return err
	}
	return nil
}

// func (s *SqlServerMoviesStore) Create(ctx context.Context, createMovieParams CreateMovieParams) error {
// 	err := s.connect(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	defer s.close()

// 	movie := Movie{
// 		ID:          createMovieParams.ID,
// 		Title:       createMovieParams.Title,
// 		Director:    createMovieParams.Director,
// 		ReleaseDate: createMovieParams.ReleaseDate,
// 		TicketPrice: createMovieParams.TicketPrice,
// 		CreatedAt:   time.Now().UTC(),
// 		UpdatedAt:   time.Now().UTC(),
// 	}
// 	log.Printf("movie: %v", movie)

// 	if _, err := s.dbx.NamedExecContext(
// 		ctx,
// 		`INSERT INTO dbo.movies
// 			(Id, Title, Director, ReleaseDate, TicketPrice, CreatedAt, UpdatedAt)
// 		VALUES
// 			(:Id, :Title, :Director, :ReleaseDate, :TicketPrice, :CreatedAt, :UpdatedAt)`,
// 		movie); err != nil {
// 		if strings.Contains(err.Error(), "Cannot insert duplicate key") {
// 			return &DuplicateKeyError{ID: createMovieParams.ID}
// 		}
// 		return err
// 	}

// 	return nil
// }

func (s *SqlServerMoviesStore) Update(ctx context.Context, jsonBody string) error {
	err := s.connect(ctx)
	if err != nil {
		return err
	}
	defer s.close()

	_, err = s.dbx.ExecContext(
		ctx,
		`EXEC dbo.u_movie @json = @json`,
		sql.Named("json", jsonBody),
	)
	if err != nil {
		if strings.Contains(err.Error(), "operation failed") {
			return &RecordNotFoundError{err}
		}
		return err
	}
	return nil
}

// func (s *SqlServerMoviesStore) Update(ctx context.Context, id uuid.UUID, updateMovieParams UpdateMovieParams) error {
// 	err := s.connect(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	defer s.close()

// 	movie := Movie{
// 		ID:          id,
// 		Title:       updateMovieParams.Title,
// 		Director:    updateMovieParams.Director,
// 		ReleaseDate: updateMovieParams.ReleaseDate,
// 		TicketPrice: updateMovieParams.TicketPrice,
// 		UpdatedAt:   time.Now().UTC(),
// 	}

// 	if _, err := s.dbx.NamedExecContext(
// 		ctx,
// 		`UPDATE dbo.movies
// 		SET Title = :Title, Director = :Director, ReleaseDate = :ReleaseDate, TicketPrice = :TicketPrice, UpdatedAt = :UpdatedAt
// 		WHERE Id = :Id`,
// 		movie); err != nil {
// 		return err
// 	}

// 	return nil
// }

func (s *SqlServerMoviesStore) Delete(ctx context.Context, jsonBody string) error {
	err := s.connect(ctx)
	if err != nil {
		return err
	}
	defer s.close()
	log.Printf("jsonBody: %v", jsonBody)
	_, err = s.dbx.ExecContext(
		ctx,
		`EXEC dbo.d_movie @json = @json`,
		sql.Named("json", jsonBody),
	)
	if err != nil {
		if strings.Contains(err.Error(), "operation failed") {
			return &RecordNotFoundError{err}
		}
		return err
	}
	return nil
}
