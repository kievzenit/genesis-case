package models

import (
	"time"

	"github.com/google/uuid"
)

type ConfirmationEmail struct {
	Id           int
	ToAddress    string
	Token        uuid.UUID
	Completed    bool
	Attempts     int
	NextTryAfter time.Time
}
