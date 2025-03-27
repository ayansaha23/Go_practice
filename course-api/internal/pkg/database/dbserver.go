package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"strings"

	"github.com/course-api/internal/pkg/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CoursesDBSession struct {
	DatabaseUrl string
	dbx         *sqlx.DB
}

// SelectContext - used to receive multiple rows from db.Returns []T(slice of structs)
// GetContext
const driverName = "mysql"

type Interface interface {
	GetAll(ctx context.Context) ([]models.Course, error)
	GetByID(ctx context.Context, id uuid.UUID) (models.Course, error)
	Create(ctx context.Context, createParams models.CreateCourseParams) error
	Update(ctx context.Context, id uuid.UUID, updateParams models.UpdateCourseParams) (models.Course, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

func NewCoursesDBSession(url string) *CoursesDBSession {
	return &CoursesDBSession{
		DatabaseUrl: url,
	}
}

func (s *CoursesDBSession) Ping(ctx context.Context) error {
	err := s.connect(ctx)
	if err != nil {
		log.Println("Could not connect to db")
	}
	defer s.close()
	err = s.dbx.PingContext(ctx)
	if err != nil {
		log.Println("PING failed:", err)
		return err
	}
	return nil
}

func (s *CoursesDBSession) connect(ctx context.Context) error {
	dbx, err := sqlx.ConnectContext(ctx, driverName, s.DatabaseUrl)
	fmt.Println("inside connect method....")
	if err != nil {
		fmt.Println("err is", err)
		return err
	}
	s.dbx = dbx
	return nil
}

func (s *CoursesDBSession) close() {
	s.dbx.Close()
}

func (s *CoursesDBSession) Create(ctx context.Context, Params models.CreateCourseParams) (models.Course, error) {
	err := s.connect(ctx)
	if err != nil {
		log.Println("Could not connect to db")
		return models.Course{}, err
	}
	defer s.close()
	query := `INSERT INTO courses(id,name,price,technology) VALUES(:id, :name, :price, :technology)`
	uuidGenerated := uuid.New()
	// since technology field is a slice need to store this in json encoded way(serialization)
	technologyJson, err := json.Marshal(Params.Technology)

	if err != nil {
		log.Println("json marshal err")
		return models.Course{}, errors.New("technology field parse error")
	}

	c := models.CourseDatabase{
		Id:         uuidGenerated.String(),
		Name:       Params.Name,
		Price:      Params.Price,
		Technology: string(technologyJson),
	}

	result, err := s.dbx.NamedExecContext(ctx, query, c)

	if err != nil {
		if strings.Contains(err.Error(), "Error 1062") {
			return models.Course{}, &DuplicateKeyError{Id: c.Id}
		}
		return models.Course{}, err
	}
	rowsAffected, _ := result.RowsAffected()
	fmt.Println("rows updated:", rowsAffected)
	if rowsAffected > 0 {
		// convert into struct that can be returned
		course := models.Course{
			Id:         c.Id,
			Name:       c.Name,
			Technology: Params.Technology,
			Price:      c.Price,
		}
		return course, nil
	}
	return models.Course{}, nil
}

func (s *CoursesDBSession) GetAll(ctx context.Context) ([]models.Course, error) {
	err := s.connect(ctx)
	if err != nil {
		log.Println("could not connext to db")
		return nil, err
	}
	defer s.close()
	var coursesDatabase []models.CourseDatabase
	var courses []models.Course
	query := `SELECT id, name, price, technology FROM courses`
	// as technology is stored as json encoded , needed to convert this into []string.
	// a temp struct to hold values retrieved from db
	// set the db.course
	err = s.dbx.SelectContext(ctx, &coursesDatabase, query)
	if err != nil {
		return nil, err
	}
	// iterate over every record and setup the fields of struct
	for _, item := range coursesDatabase {
		tech, err := models.ConvertToSlice(item.Technology)
		if err != nil {
			log.Println("unmarshsal err:", tech)
			return nil, err
		}
		course := models.Course{
			Id:         item.Id,
			Name:       item.Name,
			Price:      item.Price,
			Technology: tech,
		}
		courses = append(courses, course)
	}
	return courses, nil
}

func (s *CoursesDBSession) GetByID(ctx context.Context, id uuid.UUID) (models.Course, error) {
	err := s.connect(ctx)
	if err != nil {
		log.Println("could not connect to db")
		return models.Course{}, err
	}
	defer s.close()
	var courseRow models.CourseDatabase
	query := `SeLect id, name, price,technology from courses where id=?`
	err = s.dbx.GetContext(ctx, &courseRow, query, id)
	if err != nil {
		log.Println("Error fetching course  err:", err)
		return models.Course{}, err
	}
	// convert
	technologyVal, err := models.ConvertToSlice(courseRow.Technology)
	if err != nil {
		log.Println("Failed to convert Technology field for id:", id)
	}
	fetchedCourse := models.Course{
		Id:         courseRow.Id,
		Name:       courseRow.Name,
		Price:      courseRow.Price,
		Technology: technologyVal,
	}
	return fetchedCourse, nil
}

func (s *CoursesDBSession) Update(ctx context.Context, id uuid.UUID, updateParams models.UpdateCourseParams) (models.Course, error) {
	// accepts uuid and all other params. Updates
	err := s.connect(ctx)
	if err != nil {
		log.Println("could not connect to db")
		return models.Course{}, err
	}
	defer s.close()
	query := `UPDATE courses SET name = :name, price = :price, technology = :technology where id = :id`
	// convert updated courseParams into CourseDatabase struct
	technologyBytes, err := json.Marshal(updateParams.Technology)
	if err != nil {
		log.Println(err)
		return models.Course{}, err
	}
	courseData := models.CourseDatabase{
		Id:         id.String(),
		Name:       updateParams.Name,
		Price:      updateParams.Price,
		Technology: string(technologyBytes),
	}
	result, err := s.dbx.NamedExecContext(ctx, query, courseData)
	if err != nil {
		log.Println("Error in updating:", err)
		return models.Course{}, err
	}
	log.Println(result.RowsAffected())
	updatedCourse, err := s.GetByID(ctx, id)

	if err != nil {
		log.Println("error in fetching updated record", err)
		return models.Course{}, err
	}
	return updatedCourse, nil

}

func (s *CoursesDBSession) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.connect(ctx)
	if err != nil {
		log.Println("could not connect to db")
		return err
	}
	defer s.close()
	query := `DELETE from courses where id=?`
	result, err := s.dbx.ExecContext(ctx, query, id.String())
	if err != nil {
		log.Println("error in deleting", err)
		return err
	}
	log.Println(result)
	return nil
}
