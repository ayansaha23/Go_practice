package models

import (
	"encoding/json"
	"errors"
	"log"
)

// no space between json and fields
type Course struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Price      float64  `json:"price"`
	Technology []string `json:"technology"`
}

// during insertion of records the slice of strings need to be converted into string for saving
type CourseDatabase struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Technology string  `json:"technology"`
}

type CreateCourseParams struct {
	Name       string   `db:"name"`
	Price      float64  `db:"price"`
	Technology []string `db:"technology"`
}

type UpdateCourseParams struct {
	Name       string   `db:"name"`
	Price      float64  `db:"price"`
	Technology []string `db:"technology"`
}

func (c *CreateCourseParams) IsEmpty() bool {

	return c.Name == "" || c.Price == 0 || c.Technology == nil
}
func (c *UpdateCourseParams) IsEmpty() bool {

	return c.Name == "" || c.Price == 0 || c.Technology == nil
}
func (c Course) IsEmpty() bool {
	return c.Name == "" || c.Price == 0 || c.Technology == nil
}

func ConvertToSlice(fieldVal string) ([]string, error) {
	// here the string accepted is "["python","django","celery"]"
	// should return slice of string
	var stringSlice []string
	err := json.Unmarshal([]byte(fieldVal), &stringSlice)

	if err != nil {
		log.Println("")
		return []string{}, errors.New("Unable to marshal technology field")
	}
	return stringSlice, nil

}
