package database

import (
	"fmt"
)

type DuplicateKeyError struct {
	Id string
}

func (e *DuplicateKeyError) Error() string {
	return fmt.Sprintf("duplicate couse id", e.Id)
}
