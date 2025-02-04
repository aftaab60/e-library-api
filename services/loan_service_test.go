package services

import (
	"context"
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

	ctx := context.Background()
	// Add test book data
	_, err := bookRepo.UpdateBook(ctx, "book1", 3)
	assert.NoError(t, err) // Set 3 available copies

	t.Run("Successfully borrow a book", func(t *testing.T) {
		loan, err := loanService.BorrowBook(ctx, "book1", "borrower1")

		assert.NoError(t, err)
		assert.NotNil(t, loan)
		assert.Equal(t, "borrower1", loan.NameOfBorrower)

		// Ensure book copies reduced
		updatedBook, _ := bookRepo.GetBook(ctx, "book1")
		assert.Equal(t, 2, updatedBook.AvailableCopies)
	})

	t.Run("Fail to borrow when no copies left", func(t *testing.T) {
		_, err := bookRepo.UpdateBook(ctx, "book2", 0)
		assert.NoError(t, err) // Set 0 copies

		loan, err := loanService.BorrowBook(ctx, "book2", "borrower2")
		assert.Error(t, err)
		assert.Nil(t, loan)
		assert.Equal(t, ErrNoAvailableCopiesFound, err)
	})
}

func TestLoanService_ExtendLoan(t *testing.T) {
	bookRepo := repositories.NewBookRepository()
	loanRepo := repositories.NewLoanRepository()
	loanService := NewLoanService(loanRepo, bookRepo)

	ctx := context.Background()
	currTime := time.Now()
	_, err := loanRepo.CreateLoan(ctx, "book1", &models.Loan{
		BorrowerName: "borrower2",
		LoanDate:     currTime.AddDate(0, 0, -31),
		ReturnDate:   currTime.AddDate(0, 0, -21),
		BookId:       2,
		IsReturn:     false,
		Id:           2,
	})
	assert.NoError(t, err)

	t.Run("Successfully extend a loan", func(t *testing.T) {
		loan, err := loanService.ExtendLoan(ctx, "book1", "borrower2")
		assert.NoError(t, err)
		assert.NotNil(t, loan)
		assert.Equal(t, currTime.Unix(), loan.ReturnDate.Unix()) // Extended by 21 days
	})
}

func TestLoanService_ReturnBook(t *testing.T) {
	bookRepo := repositories.NewBookRepository()
	loanRepo := repositories.NewLoanRepository()
	loanService := NewLoanService(loanRepo, bookRepo)
	ctx := context.Background()

	// Add a book and loan
	bookRepo.UpdateBook(ctx, "book1", 2)
	loanRepo.CreateLoan(ctx, "book1", &models.Loan{
		BorrowerName: "borrower3",
		LoanDate:     time.Now(),
		ReturnDate:   time.Now().AddDate(0, 0, 28),
		BookId:       2,
		IsReturn:     false,
		Id:           3,
	})

	t.Run("Successfully return a book", func(t *testing.T) {
		err := loanService.ReturnBook(ctx, "book1", "borrower3")
		assert.NoError(t, err)

		// Ensure book copies increased
		updatedBook, _ := bookRepo.GetBook(ctx, "book1")
		assert.Equal(t, 3, updatedBook.AvailableCopies)
	})

	t.Run("Fail to return a non-existent loan", func(t *testing.T) {
		err := loanService.ReturnBook(ctx, "book2", "borrower2")

		assert.Error(t, err)
		assert.Equal(t, repositories.ErrLoanNotFound, err)
	})
}
