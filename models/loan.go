package models

import (
	"time"

	"github.com/google/uuid"
)

type BookLoan struct {
    ID         uuid.UUID  `json:"id"`
    UserID     uuid.UUID  `json:"user_id"`
    BookID     uuid.UUID  `json:"book_id"`
    BorrowDate time.Time  `json:"borrow_date"`
    ReturnDate *time.Time `json:"return_date,omitempty"`
}
