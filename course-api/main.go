package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/course-api/internal/pkg/database"
	"github.com/course-api/internal/pkg/server"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("This is going to be the server")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbURL := os.Getenv("DATABASE_URL")
	log.Println("Starting server....")
	router := mux.NewRouter()
	db := database.NewCoursesDBSession(dbURL)
	s := server.NewApiServer(":6060", router, db)
	ctx := context.Background()
	s.Run(ctx)

}
