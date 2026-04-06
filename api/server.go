package api

import (
	"dperkins/movies-api/config"
	"dperkins/movies-api/store"
	//"github.com/kashifsoofi/blog-code-samples/movies-api-with-go-chi-and-memory-store/config"
	//"github.com/kashifsoofi/blog-code-samples/movies-api-with-go-chi-and-memory-store/store"
)

type Server struct {
	cfg   config.HTTPServer
	store store.Interface
}
