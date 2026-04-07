package api

import (
	"dperkins/movies-api/store"
	"errors"
	"net/http"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type movieResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Director    string    `json:"director"`
	ReleaseDate time.Time `json:"release_date"`
	TicketPrice float64   `json:"ticket_price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateMovieRequest struct {
	Title       string    `json:"title"`
	Director    string    `json:"director"`
	ReleaseDate time.Time `json:"release_date"`
	TicketPrice float64   `json:"ticket_price"`
}

type CreateMovieRequest struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Director    string    `json:"director"`
	ReleaseDate time.Time `json:"release_date"`
	TicketPrice float64   `json:"ticket_price"`
}

func NewMovieResponse(m store.Movie) movieResponse {
	return movieResponse{
		ID:          m.ID,
		Title:       m.Title,
		Director:    m.Director,
		ReleaseDate: m.ReleaseDate,
		TicketPrice: m.TicketPrice,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func NewMovieListResponse(movies []store.Movie) []render.Renderer {
	list := []render.Renderer{}
	for _, movie := range movies {
		mr := NewMovieResponse(movie)
		list = append(list, mr)
	}
	return list
}

func (mr movieResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) handleGetMovie(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		render.Render(w, r, ErrBadRequest)
	}

	movie, err := s.store.GetByID(r.Context(), id)
	if err != nil {
		var rnfErr *store.RecordNotFoundError
		if errors.As(err, &rnfErr) {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	mr := movieResponse(movie)
	render.Render(w, r, mr)
}

func (s *Server) handleListMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := s.store.GetAll(r.Context())
	if err != nil {
		render.Render(w, r, ErrInternalServerError)
		return
	}

	render.RenderList(w, r, NewMovieListResponse(movies))
}

func (mr *CreateMovieRequest) Bind(r *http.Request) error {
	return nil
}

func (s *Server) handleCreateMovie(w http.ResponseWriter, r *http.Request) {
	data := &CreateMovieRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	createMovieParams := store.CreateMovieParams{
		ID:          data.ID,
		Title:       data.Title,
		Director:    data.Director,
		ReleaseDate: data.ReleaseDate,
		TicketPrice: data.TicketPrice,
	}

	err := s.store.Create(r.Context(), createMovieParams)
	if err != nil {
		var dupKeyError *store.DuplicateKeyError
		if errors.As(err, &dupKeyError) {
			render.Render(w, r, ErrConflict(err))
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	w.WriteHeader(200)
	w.Write(nil)
}

func (mr *UpdateMovieRequest) Bind(r *http.Request) error {
	return nil
}

func (s *Server) handleUpdateMovie(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	data := &UpdateMovieRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	updateMovieParams := store.UpdateMovieParams{
		Title:       data.Title,
		Director:    data.Director,
		ReleaseDate: data.ReleaseDate,
		TicketPrice: data.TicketPrice,
	}

	err = s.store.Update(r.Context(), id, updateMovieParams)
	if err != nil {
		var rnfErr *store.RecordNotFoundError
		if errors.As(err, &rnfErr) {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	w.WriteHeader(200)
	w.Write(nil)
}

func (s *Server) handleDeleteMovie(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	err = s.store.Delete(r.Context(), id)
	if err != nil {
		var rnfErr *store.RecordNotFoundError
		if errors.As(err, &rnfErr) {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	w.WriteHeader(200)
	w.Write(nil)
}
