package main

import (
	"github.com/aftaab60/e-library-api/repositories"
	"github.com/aftaab60/e-library-api/routes"
	"github.com/aftaab60/e-library-api/services"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery()) // to recover from panics in execution

	setupRoutes(r)
	if err := r.Run(":3000"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func setupRoutes(r *gin.Engine) {
	bookRepository := repositories.NewBookRepository()
	loanRepository := repositories.NewLoanRepository()

	bookRoute := routes.NewBookRoute(services.NewBookService(bookRepository))
	loanRoute := routes.NewLoanRoute(services.NewLoanService(loanRepository, bookRepository))

	r.GET("/book/:title", bookRoute.GetBookByTitle)
	r.POST("/borrow", loanRoute.BorrowBook)
	r.POST("/extend", loanRoute.ExtendLoan)
	r.POST("/return", loanRoute.ReturnBook)
}
