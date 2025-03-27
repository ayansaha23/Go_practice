package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/course-api/internal/pkg/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// crud routes

func (s *ApiServer) Homelander(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("<h2>Welcome to home page</h2>"))
	return
}

func (s *ApiServer) showCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	courses, err := s.Db.GetAll(r.Context())
	if err != nil {
		log.Println("err in fetching courses:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("oops something went wrong")
		return
	}
	json.NewEncoder(w).Encode(courses)
}

func (s *ApiServer) createCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var receivedCourse models.CreateCourseParams
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
		log.Println("Could not parse request Body:", err)
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
		newCourse, err := s.Db.Create(r.Context(), receivedCourse)
		if err != nil {
			log.Printf("DB error when creating", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("oops something went wrong")
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newCourse)
	}
}

func (s *ApiServer) showCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	val, ok := params["id"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("No id provided")
	} else {
		// consider id provided
		id, err := uuid.Parse(val)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("could not read id from input")
			return
		}
		course, err := s.Db.GetByID(r.Context(), id)
		if err != nil {

			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		json.NewEncoder(w).Encode(course)
	}
}

func (s *ApiServer) updateCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	val, ok := params["id"]

	// check if params["id"] is present or not
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("no id provided")
		return
	}
	receivedId, err := uuid.Parse(val)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("could not read id from input")
		return
	}
	// need to validate the body of the request receivd
	var receivedCourse models.UpdateCourseParams
	err = json.NewDecoder(r.Body).Decode(&receivedCourse)
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

	course, err := s.Db.Update(r.Context(), receivedId, receivedCourse)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(course)

}

func (s *ApiServer) deleteCourse(w http.ResponseWriter, r *http.Request) {
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

	receivedId, err := uuid.Parse(val)
	if err != nil {
		log.Println("error in parsing received ID")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("could not parse Id")
		return
	}
	err = s.Db.Delete(r.Context(), receivedId)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode("failed to delete record")
		return
	}
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode("record deleted sucessfully")
	return
}
