package services

import (
	"testing"
	"time"

	"github.com/aftaab60/e-library-api/models"
	"github.com/aftaab60/e-library-api/repositories"
	"github.com/stretchr/testify/assert"
)

func TestLoanService_BorrowBook(t *testing.T) {
	bookRepo := repositories.NewBookRepository()
	loanRepo := repositories.NewLoanRepository()
	loanService := NewLoanService(loanRepo, bookRepo)

	// Add test book data
	_, err := bookRepo.UpdateBook("book1", 3)
	assert.NoError(t, err) // Set 3 available copies

	t.Run("Successfully borrow a book", func(t *testing.T) {
		loan, err := loanService.BorrowBook("book1", "borrower1")

		assert.NoError(t, err)
		assert.NotNil(t, loan)
		assert.Equal(t, "borrower1", loan.NameOfBorrower)

		// Ensure book copies reduced
		updatedBook, _ := bookRepo.GetBook("book1")
		assert.Equal(t, 2, updatedBook.AvailableCopies)
	})

	t.Run("Fail to borrow when no copies left", func(t *testing.T) {
		_, err := bookRepo.UpdateBook("book2", 0)
		assert.NoError(t, err) // Set 0 copies

		loan, err := loanService.BorrowBook("book2", "borrower2")
		assert.Error(t, err)
		assert.Nil(t, loan)
		assert.Equal(t, ErrNoAvailableCopiesFound, err)
	})
}

func TestLoanService_ExtendLoan(t *testing.T) {
	bookRepo := repositories.NewBookRepository()
	loanRepo := repositories.NewLoanRepository()
	loanService := NewLoanService(loanRepo, bookRepo)

	currTime := time.Now()
	_, err := loanRepo.CreateLoan("book1", &models.LoanDetail{
		NameOfBorrower: "borrower2",
		LoanDate:       currTime.AddDate(0, 0, -31),
		ReturnDate:     currTime.AddDate(0, 0, -21),
	})
	assert.NoError(t, err)

	t.Run("Successfully extend a loan", func(t *testing.T) {
		loan, err := loanService.ExtendLoan("book1", "borrower2")
		assert.NoError(t, err)
		assert.NotNil(t, loan)
		assert.Equal(t, currTime.Unix(), loan.ReturnDate.Unix()) // Extended by 21 days
	})
}

func TestLoanService_ReturnBook(t *testing.T) {
	bookRepo := repositories.NewBookRepository()
	loanRepo := repositories.NewLoanRepository()
	loanService := NewLoanService(loanRepo, bookRepo)

	// Add a book and loan
	bookRepo.UpdateBook("book1", 2)
	loanRepo.CreateLoan("book1", &models.LoanDetail{
		NameOfBorrower: "borrower3",
		LoanDate:       time.Now(),
		ReturnDate:     time.Now().AddDate(0, 0, 28),
	})

	t.Run("Successfully return a book", func(t *testing.T) {
		err := loanService.ReturnBook("book1", "borrower3")
		assert.NoError(t, err)

		// Ensure book copies increased
		updatedBook, _ := bookRepo.GetBook("book1")
		assert.Equal(t, 3, updatedBook.AvailableCopies)
	})

	t.Run("Fail to return a non-existent loan", func(t *testing.T) {
		err := loanService.ReturnBook("book2", "borrower2")

		assert.Error(t, err)
		assert.Equal(t, repositories.ErrLoanNotFound, err)
	})
}
