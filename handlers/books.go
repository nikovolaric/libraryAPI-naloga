package handlers

import (
	"libraryAPI/db"
	"libraryAPI/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateBook(c *gin.Context) {
	var req struct {
		Title string `json:"title"`
		Total int    `json:"total"`
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

	db.DB.Exec("INSERT INTO books (id, title, total, available) VALUES ($1, $2, $3, $4)", book.ID, book.Title, book.Total, book.Available)

	c.JSON(201, book)

}

func GetAllBooks(c *gin.Context) {
	var books []models.Book

	rows, err := db.DB.Query("SELECT * FROM books")

	if err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong."})
		return
	}

	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Total, &book.Available); err != nil {
			c.JSON(500, gin.H{"error": "Something went wrong."})
			return

		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong."})
		return
	}

	c.JSON(200, books)
}

func GetAllAvailableBooks(c *gin.Context) {
	var books []struct {
		ID        uuid.UUID
		Title     string
		Available int
	}

	rows, err := db.DB.Query("SELECT id, title, available FROM books WHERE available != 0")

	if err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong."})
		return
	}

	for rows.Next() {
		var book struct {
			ID        uuid.UUID
			Title     string
			Available int
		}
		if err := rows.Scan(&book.ID, &book.Title, &book.Available); err != nil {
			c.JSON(500, gin.H{"error": "Something went wrong."})

		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong."})
		return
	}

	c.JSON(200, books)

}
