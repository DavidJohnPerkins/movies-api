package api

import (
	"dperkins/movies-api/config"
	"dperkins/movies-api/store"

	//"github.com/kashifsoofi/blog-code-samples/movies-api-with-go-chi-and-memory-store/config"
	//"github.com/kashifsoofi/blog-code-samples/movies-api-with-go-chi-and-memory-store/store"

	chi "github.com/go-chi/chi/v5"
)

type Server struct {
	cfg    config.HTTPServer
	store  store.Interface
	router *chi.Mux
}

func NewServer(cfg config.HTTPServer, store store.Interface) *Server {
	srv := &Server{
		cfg:    cfg,
		store:  store,
		router: chi.NewRouter(),
	}

	srv.routes()

	return srv
}
