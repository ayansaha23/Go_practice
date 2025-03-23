package main

import (
	"fmt"
	"log"

	"github.com/gorilla/mux"

	"github.com/course-api/internal/pkg/server"
)

func main() {
	fmt.Println("This is going to be the server")
	// create a server
	// create a router
	// create a controller
	log.Println("Starting server....")
	router := mux.NewRouter()
	s := server.NewApiServer(":6060", router)
	s.Run()

}
