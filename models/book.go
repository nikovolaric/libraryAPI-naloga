package models

import "github.com/google/uuid"

type Book struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Total     int       `json:"total"`
	Available int       `json:"available"`
}
