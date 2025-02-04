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
	//From below 2 lines, select either in-memory or pgsql db repository.
	//Implementation are on interfaces hence same service works in both cases.
	
	bookRepository := repositories.NewBookRepository()
	loanRepository := repositories.NewLoanRepository()
	//bookRepository := repositories.NewBookRepositoryDB(db_manager.InitPgsqlConnection())
	//loanRepository := repositories.NewLoanRepositoryDB(db_manager.InitPgsqlConnection())

	bookRoute := routes.NewBookRoute(services.NewBookService(bookRepository))
	loanRoute := routes.NewLoanRoute(services.NewLoanService(loanRepository, bookRepository))

	r.GET("/book/:title", bookRoute.GetBookByTitle)
	r.POST("/borrow", loanRoute.BorrowBook)
	r.POST("/extend", loanRoute.ExtendLoan)
	r.POST("/return", loanRoute.ReturnBook)
}
