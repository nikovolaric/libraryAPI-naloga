package handlers

import (
	"database/sql"
	"errors"
	"time"

	"libraryAPI/db"
	"libraryAPI/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateUser(c *gin.Context) {
	var req struct {
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
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

	if _, err := db.DB.ExecContext(c.Request.Context(),
		"INSERT INTO users (id, first_name, last_name) VALUES ($1, $2, $3)",
		user.ID, user.FirstName, user.LastName); err != nil {
		c.JSON(500, gin.H{"error": "could not create user"})
		return
	}

	c.JSON(201, user)
}

func GetAllUsers(c *gin.Context) {
	users := []models.User{}

	rows, err := db.DB.QueryContext(c.Request.Context(),
		"SELECT id, first_name, last_name FROM users")
	if err != nil {
		c.JSON(500, gin.H{"error": "could not fetch users"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName); err != nil {
			c.JSON(500, gin.H{"error": "could not read user"})
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": "could not fetch users"})
		return
	}

	c.JSON(200, users)
}

func GetOneUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user id"})
		return
	}

	var user models.User
	err = db.DB.QueryRowContext(c.Request.Context(),
		"SELECT id, first_name, last_name FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.FirstName, &user.LastName)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "could not fetch user"})
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
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user id"})
		return
	}

	loans := []Loan{}

	rows, err := db.DB.QueryContext(c.Request.Context(),
		`SELECT l.id, l.borrow_date, l.return_date, b.id, b.title
		 FROM loans l
		 JOIN books b ON b.id = l.book_id
		 WHERE l.user_id = $1`, id)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not fetch loans"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var loan Loan
		if err := rows.Scan(&loan.ID, &loan.BorrowDate, &loan.ReturnDate, &loan.Book.ID, &loan.Book.Title); err != nil {
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
