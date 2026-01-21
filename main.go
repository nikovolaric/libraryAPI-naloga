package main

import (
	"libraryAPI/db"
	"libraryAPI/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect()

	router := gin.Default()

	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.GET("/users", handlers.GetAllUsers)
	router.POST("/users", handlers.CreateUser)
	router.GET("/users/:id", handlers.GetOneUser)
	router.GET("/users/:id/loans", handlers.GetUserLoans)

	router.GET("/books", handlers.GetAllBooks)
	router.POST("/books", handlers.CreateBook)
	router.GET("/books/available", handlers.GetAllAvailableBooks)

	router.GET("/loans", handlers.GetAllBookLoans)
	router.POST("/loans/new", handlers.NewLoan)
	router.POST("/loans/return/:id", handlers.ReturnLoan)

	router.Run(":8000")
}
