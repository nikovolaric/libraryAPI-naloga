package handlers

import (
	"libraryAPI/db"
	"libraryAPI/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateUser(c *gin.Context) {
	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	user := models.User{
		ID:        uuid.New(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	db.DB.Exec(
		"INSERT INTO users (id, first_name, last_name) VALUES ($1, $2, $3)",
		user.ID, user.FirstName, user.LastName)

	c.JSON(201, user)
}

func GetAllUsers(c *gin.Context) {

	var users []models.User

	rows, err := db.DB.Query("SELECT * FROM users")

	if err != nil {

		c.JSON(500, gin.H{"error": "Something went wrong."})
		return
	}

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName); err != nil {
			c.JSON(500, gin.H{"error": "Something went wrong."})
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {

		c.JSON(500, gin.H{"error": "Something went wrong."})
		return
	}

	c.JSON(200, users)
}

func GetOneUser(c *gin.Context) {
	var user models.User

	id := uuid.MustParse(c.Param("id"))

	err := db.DB.QueryRow("SELECT id, first_name, last_name FROM users WHERE id = $1", id).Scan(&user.ID, &user.FirstName, &user.LastName)

	if err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	c.JSON(200, user)
}

type Loan struct {
	ID         uuid.UUID  `json:"id"`
	BorrowDate time.Time  `json:"borrow_date"`
	ReturnDate *time.Time `json:"return_date"`

	Book struct {
		ID    uuid.UUID `json:"id"`
		Title string    `json:"title"`
	} `json:"book"`
}

func GetUserLoans(c *gin.Context) {
	var loans []Loan
	id := uuid.MustParse(c.Param("id"))

	rows, err := db.DB.Query("SELECT l.id, l.borrow_date, l.return_date, b.id, b.title FROM loans l JOIN books b ON b.id=l.book_id WHERE l.user_id=$1", id)

	if err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	for rows.Next() {
		var loan Loan

		if err = rows.Scan(&loan.ID, &loan.BorrowDate, &loan.ReturnDate, &loan.Book.ID, &loan.Book.Title); err != nil {
			c.JSON(500, gin.H{"error": "Something went wrong."})
		}

		loans = append(loans, loan)
	}

	c.JSON(200, loans)
}
