package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/course-api/internal/pkg/database"
	"github.com/gorilla/mux"
)

type ApiServer struct {
	Addr    string
	Handler *mux.Router
	Db      *database.CoursesDBSession
}

func NewApiServer(addr string, handler *mux.Router, db *database.CoursesDBSession) *ApiServer {
	return &ApiServer{
		Addr:    addr,
		Handler: handler,
		Db:      db,
	}
}

func (s *ApiServer) Run(ctx context.Context) error {
	log.Println("Server starting on:", s.Addr)
	server := &http.Server{
		Addr:    s.Addr,
		Handler: s.Handler,
	}

	s.SetUpRoutes()
	s.Handler.Use(loggingMiddleware)
	log.Println("router set")
	err := s.Db.Ping(ctx)
	if err != nil {
		panic("DB server not connected")
	}
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

func (s *ApiServer) SetUpRoutes() {
	s.Handler.HandleFunc("/", s.Homelander).Methods("GET")
	s.Handler.HandleFunc("/courses", s.showCourses).Methods("GET")
	s.Handler.HandleFunc("/course", s.createCourse).Methods("POST")
	s.Handler.HandleFunc("/courses/{id}", s.showCourse).Methods("GET")
	s.Handler.HandleFunc("/courses/{id}", s.updateCourse).Methods("PUT")
	s.Handler.HandleFunc("/courses/{id}", s.deleteCourse).Methods("DELETE")
}
