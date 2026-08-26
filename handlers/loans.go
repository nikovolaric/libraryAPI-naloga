package handlers

import (
	"time"

	"libraryAPI/db"
	"libraryAPI/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookLoanResponse struct {
	ID         uuid.UUID  `json:"id"`
	BorrowDate time.Time  `json:"borrow_date"`
	ReturnDate *time.Time `json:"return_date"`

	User struct {
		ID        uuid.UUID `json:"id"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
	} `json:"user"`

	Book struct {
		ID    uuid.UUID `json:"id"`
		Title string    `json:"title"`
	} `json:"book"`
}

func GetAllBookLoans(c *gin.Context) {
	loans := []BookLoanResponse{}

	rows, err := db.DB.QueryContext(c.Request.Context(),
		`SELECT l.id, l.borrow_date, l.return_date,
		        u.id, u.first_name, u.last_name,
		        b.id, b.title
		 FROM loans l
		 JOIN users u ON u.id = l.user_id
		 JOIN books b ON b.id = l.book_id`)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not fetch loans"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var loan BookLoanResponse
		if err := rows.Scan(
			&loan.ID, &loan.BorrowDate, &loan.ReturnDate,
			&loan.User.ID, &loan.User.FirstName, &loan.User.LastName,
			&loan.Book.ID, &loan.Book.Title,
		); err != nil {
			c.JSON(500, gin.H{"error": "could not read loan"})
			return
		}
		loans = append(loans, loan)
	}

	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": "could not fetch loans"})
		return
	}

	c.JSON(200, loans)
}

func NewLoan(c *gin.Context) {
	var req struct {
		BookID string `json:"book_id" binding:"required"`
		UserID string `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	bookID, err := uuid.Parse(req.BookID)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid book_id"})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user_id"})
		return
	}

	ctx := c.Request.Context()

	// Decrementing availability and inserting the loan must be atomic:
	// a crash between the two would otherwise leak a copy.
	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not start transaction"})
		return
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx,
		"UPDATE books SET available = available - 1 WHERE id = $1 AND available > 0",
		bookID)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not update book"})
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		c.JSON(400, gin.H{"error": "book not found or not available"})
		return
	}

	loan := models.BookLoan{
		ID:         uuid.New(),
		BookID:     bookID,
		UserID:     userID,
		BorrowDate: time.Now(),
	}

	if _, err := tx.ExecContext(ctx,
		"INSERT INTO loans (id, user_id, book_id, borrow_date) VALUES ($1, $2, $3, $4)",
		loan.ID, loan.UserID, loan.BookID, loan.BorrowDate); err != nil {
		// Most likely an unknown user_id (foreign key violation).
		c.JSON(400, gin.H{"error": "could not create loan, check user_id"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(500, gin.H{"error": "could not commit transaction"})
		return
	}

	c.JSON(201, loan)
}

func ReturnLoan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid loan id"})
		return
	}

	ctx := c.Request.Context()

	// Closing the loan and returning the copy must be atomic.
	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not start transaction"})
		return
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx,
		"UPDATE loans SET return_date = $2 WHERE id = $1 AND return_date IS NULL",
		id, time.Now())
	if err != nil {
		c.JSON(500, gin.H{"error": "could not update loan"})
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		c.JSON(400, gin.H{"error": "loan not found or already returned"})
		return
	}

	if _, err := tx.ExecContext(ctx,
		"UPDATE books SET available = available + 1 WHERE id = (SELECT book_id FROM loans WHERE id = $1)",
		id); err != nil {
		c.JSON(500, gin.H{"error": "could not update book"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(500, gin.H{"error": "could not commit transaction"})
		return
	}

	c.JSON(200, gin.H{"status": "book returned"})
}
