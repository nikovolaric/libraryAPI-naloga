package handlers

import (
	"libraryAPI/db"
	"libraryAPI/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AvailableBook struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Available int       `json:"available"`
}

func CreateBook(c *gin.Context) {
	var req struct {
		Title string `json:"title" binding:"required"`
		Total int    `json:"total" binding:"required,gte=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	book := models.Book{
		ID:        uuid.New(),
		Title:     req.Title,
		Total:     req.Total,
		Available: req.Total,
	}

	if _, err := db.DB.ExecContext(c.Request.Context(),
		"INSERT INTO books (id, title, total, available) VALUES ($1, $2, $3, $4)",
		book.ID, book.Title, book.Total, book.Available); err != nil {
		c.JSON(500, gin.H{"error": "could not create book"})
		return
	}

	c.JSON(201, book)
}

func GetAllBooks(c *gin.Context) {
	books := []models.Book{}

	rows, err := db.DB.QueryContext(c.Request.Context(),
		"SELECT id, title, total, available FROM books")
	if err != nil {
		c.JSON(500, gin.H{"error": "could not fetch books"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Total, &book.Available); err != nil {
			c.JSON(500, gin.H{"error": "could not read book"})
			return
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": "could not fetch books"})
		return
	}

	c.JSON(200, books)
}

func GetAllAvailableBooks(c *gin.Context) {
	books := []AvailableBook{}

	rows, err := db.DB.QueryContext(c.Request.Context(),
		"SELECT id, title, available FROM books WHERE available > 0")
	if err != nil {
		c.JSON(500, gin.H{"error": "could not fetch books"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var book AvailableBook
		if err := rows.Scan(&book.ID, &book.Title, &book.Available); err != nil {
			c.JSON(500, gin.H{"error": "could not read book"})
			return
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": "could not fetch books"})
		return
	}

	c.JSON(200, books)
}
