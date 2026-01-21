package handlers

import (
	"libraryAPI/db"
	"libraryAPI/models"

	"time"

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
	var loans []BookLoanResponse

	rows, err := db.DB.Query("SELECT l.id, l.borrow_date, l.return_date, u.id, u.first_name, u.last_name, b.id, b.title FROM loans l JOIN users u ON u.id=l.user_id JOIN books b ON b.id=l.book_id")

	if err != nil {

		c.JSON(500, gin.H{"error": "Something went wrong in query."})
		return
	}

	for rows.Next() {
		var loan BookLoanResponse

		if err := rows.Scan(&loan.ID, &loan.BorrowDate, &loan.ReturnDate, &loan.User.ID, &loan.User.FirstName, &loan.User.LastName, &loan.Book.ID, &loan.Book.Title); err != nil {
			c.JSON(500, gin.H{"error": "Something went wrong in scan."})
		}
		loans = append(loans, loan)
	}

	if err := rows.Err(); err != nil {

		c.JSON(500, gin.H{"error": "Something went wrong."})
		return
	}

	c.JSON(200, loans)
}

func NewLoan(c *gin.Context) {
	var req struct {
		BookID string `json:"book_id"`
		UserID string `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	loan := models.BookLoan{
		ID:         uuid.New(),
		BookID:     uuid.MustParse(req.BookID),
		UserID:     uuid.MustParse(req.UserID),
		BorrowDate: time.Now(),
	}

	rows, err := db.DB.Exec("UPDATE books SET available=available-1 WHERE id=$1 AND available > 0", loan.BookID)

	if no, _ := rows.RowsAffected(); no == 0 {
		c.JSON(400, gin.H{"message": "Book not found or not available."})
		return
	}

	if err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong."})
		return
	}

	db.DB.Exec("INSERT INTO loans (id,user_id,book_id,borrow_date) VALUES ($1,$2,$3,$4)", loan.ID, loan.UserID, loan.BookID, loan.BorrowDate)

	c.JSON(201, loan)

}

func ReturnLoan(c *gin.Context) {
	id := c.Param("id")

	loan, _ := db.DB.Exec("UPDATE loans SET return_date=$2 WHERE id=$1 AND return_date IS NULL", uuid.MustParse(id), time.Now())

	if no, _ := loan.RowsAffected(); no == 0 {
		c.JSON(400, gin.H{"message": "Book was already returned."})
		return
	}

	book, _ := db.DB.Exec("UPDATE books SET available=available+1 WHERE id=(SELECT book_id FROM loans WHERE id=$1)", uuid.MustParse(id))

	if no, _ := book.RowsAffected(); no == 0 {
		c.JSON(400, gin.H{"message": "Book was already returned."})
		return
	}

	c.JSON(200, gin.H{"status": "Book returned."})

}
