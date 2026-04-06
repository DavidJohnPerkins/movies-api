package store

import (
	"fmt"

	"github.com/google/uuid"
)

type DuplicateKeyError struct {
	ID uuid.UUID
}

type RecordNotFoundError struct{}

func (e *DuplicateKeyError) Error() string {
	return fmt.Sprintf("Duplicate movie id: %v", e.ID)
}

func (e *RecordNotFoundError) Error() string {
	return "Record not found"
}
