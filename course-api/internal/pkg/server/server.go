package server

import (
	"log"
	"net/http"
	"time"

	"github.com/course-api/internal/pkg/handler"
	"github.com/gorilla/mux"
)

type ApiServer struct {
	Addr    string
	Handler *mux.Router
}

func NewApiServer(addr string, handler *mux.Router) *ApiServer {
	return &ApiServer{
		Addr:    addr,
		Handler: handler,
	}
}

func (s *ApiServer) Run() error {
	log.Println("Server starting on:", s.Addr)
	server := &http.Server{
		Addr:    s.Addr,
		Handler: s.Handler,
	}
	handler.SetUpRoutes(s.Handler)
	s.Handler.Use(loggingMiddleware)
	log.Println("router set")
	return server.ListenAndServe()
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		log.Printf("Request is %s ,%s", r.URL, r.URL.Path)
		next.ServeHTTP(w, r)
		duration := time.Since(startTime)
		log.Printf("Response: %s %s - %v", r.Method, r.URL.Path, duration)
	})
}
