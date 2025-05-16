package models

import "github.com/google/uuid"

type Subscription struct {
	Id        int
	Token     uuid.UUID
	Confirmed bool
	Email     string
	City      string
	Frequency Frequency
}
