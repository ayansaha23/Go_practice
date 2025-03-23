package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/course-api/internal/pkg/models"
	"github.com/gorilla/mux"
)

func getCourses() *[]models.Course {
	courses := []models.Course{
		{
			Id:         1,
			Name:       "ruby on rails",
			Price:      500.00,
			Technology: []string{"ruby", "rails", "docker"},
		},
		{
			Id:         2,
			Name:       "go",
			Price:      500.00,
			Technology: []string{"go", "gin", "docker"},
		},
	}
	return &courses
}

// crud routes

func SetUpRoutes(r *mux.Router) {
	r.HandleFunc("/", Homelander).Methods("GET")
	r.HandleFunc("/courses", showCourses).Methods("GET")
	r.HandleFunc("/course", createCourse).Methods("POST")
	r.HandleFunc("/courses/{id:[0-9]+}", showCourse).Methods("GET")
	r.HandleFunc("/courses/{id:[0-9]+}", updateCourse).Methods("PUT")
	r.HandleFunc("/courses/{id:[0-9]+}", deleteCourse).Methods("DELETE")
}

func Homelander(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h2>Welcome to home page</h2>"))
	return
}

func showCourses(w http.ResponseWriter, r *http.Request) {
	// initially iteration without db
	w.Header().Set("Content-Type", "application/json")
	courses := getCourses()
	json.NewEncoder(w).Encode(courses)
}

func createCourse(w http.ResponseWriter, r *http.Request) {
	var receivedCourse models.Course
	// reading the body io
	err := json.NewDecoder(r.Body).Decode(&receivedCourse)
	if err == io.EOF {
		log.Println("Payload not provided:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode("please provide payload")
		return
	}

	if err != nil {
		log.Println("Could not parse request Bod:", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("problem with payload")
		return
	}

	if receivedCourse.IsEmpty() {
		log.Println("Empty payload")
		w.WriteHeader(http.StatusExpectationFailed)
		json.NewEncoder(w).Encode("Payload is empty")
		return
	} else {
		// here the payload is not empty
		// append it to the slice and return  the new object
		existingCourses := getCourses()
		*existingCourses = append(*existingCourses, receivedCourse)
		for _, item := range *existingCourses {
			log.Println(item)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(receivedCourse)
	}
}

func findCourse(receivedId int, existingCourses *[]models.Course) (models.Course, int, error) {
	for i, item := range *existingCourses {
		if item.Id == receivedId {
			return item, i, nil
		}
	}
	return models.Course{}, 0, errors.New("no course found")
}

func showCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	val, ok := params["id"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("No id provided")
	} else {
		// consider id provided
		receivedId, _ := strconv.Atoi(val)
		existingCourses := getCourses()
		course, _, err := findCourse(receivedId, existingCourses)
		if err != nil {

			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		json.NewEncoder(w).Encode(course)
	}
}

func updateCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	val, ok := params["id"]

	// check if params["id"] is present or not
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("no id provided")
		return
	}
	// need to validate the body of the request receivd
	var receivedCourse models.Course
	err := json.NewDecoder(r.Body).Decode(&receivedCourse)
	if err == io.EOF {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode("no payload provided")
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode("invalid payload provided")
		return
	}

	if receivedCourse.IsEmpty() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode("payload contain empty values")
		return
	}
	existingCourses := getCourses()
	receivedId, _ := strconv.Atoi(val)

	found, _, err := findCourse(receivedId, existingCourses)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	if found.Id != 0 {
		// means we found the course and the zero value is no the default
		found.Id = receivedCourse.Id
		found.Name = receivedCourse.Name
		copy(found.Technology, receivedCourse.Technology)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(found)
		return
	}
}

func deleteCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// check if params["id"] is present
	params := mux.Vars(r)

	val, ok := params["id"]

	if !ok {
		// id is not provided
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("no id provided")
		return
	}

	// convert id from to i and get the specific course
	receivedId, _ := strconv.Atoi(val)
	existingCourses := getCourses()
	_, index, err := findCourse(receivedId, existingCourses)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("no course found with the provided id")
		return
	}

	*existingCourses = append((*existingCourses)[:index], (*existingCourses)[index+1:]...)
	json.NewEncoder(w).Encode(*existingCourses)

}
