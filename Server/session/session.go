package session

import (
	"github.com/google/uuid"
)

func NewSessionID() *string {
	u := uuid.New().String()
	return &u
}
