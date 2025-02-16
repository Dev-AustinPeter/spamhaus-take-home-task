package api

import (
	"log"
	"net/http"

	urlHlr "github.com/Dev-AustinPeter/spamhaus-take-home-task/handler/url"
	"github.com/Dev-AustinPeter/spamhaus-take-home-task/middleware"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		addr: addr,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	rateLimiter := middleware.NewRateLimiter()

	urlHandler := urlHlr.NewHandler()
	urlHandler.RegisterRoutes(subrouter, rateLimiter)

	log.Println("[INFO]: Listening on port", s.addr)
	return http.ListenAndServe(s.addr, router)
}
