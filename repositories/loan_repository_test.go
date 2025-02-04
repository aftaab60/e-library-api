package repositories

import (
	"context"
	"github.com/aftaab60/e-library-api/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLoanRepository_CreateLoan(t *testing.T) {
	repo := NewLoanRepository()
	ctx := context.Background()

	loanDetail := &models.Loan{
		BookId:       1,
		BorrowerName: "user1",
		LoanDate:     time.Now(),
		ReturnDate:   time.Now().AddDate(0, 0, 28),
		IsReturn:     false,
	}

	t.Run("Create new loan", func(t *testing.T) {
		createdLoan, err := repo.CreateLoan(ctx, "book1", loanDetail)
		assert.NoError(t, err)
		assert.NotNil(t, createdLoan)
		assert.Equal(t, "user1", createdLoan.BorrowerName)
	})

	t.Run("Fail to create duplicate loan", func(t *testing.T) {
		_, err := repo.CreateLoan(ctx, "book1", loanDetail)
		assert.Error(t, err)
		assert.Equal(t, ErrExistingActiveLoan, err)
	})
}

func TestLoanRepository_GetLoan(t *testing.T) {
	repo := NewLoanRepository()
	ctx := context.Background()

	loanDetail := &models.Loan{
		BookId:       1,
		BorrowerName: "user2",
		LoanDate:     time.Now(),
		ReturnDate:   time.Now().AddDate(0, 0, 28),
		IsReturn:     false,
	}
	_, err := repo.CreateLoan(ctx, "book2", loanDetail)
	assert.NoError(t, err)

	t.Run("Get existing loan", func(t *testing.T) {
		loan, err := repo.GetLoan(ctx, "book2", "user2")
		assert.NoError(t, err)
		assert.NotNil(t, loan)
		assert.Equal(t, "user2", loan.BorrowerName)
	})

	t.Run("Get non-existent loan", func(t *testing.T) {
		loan, err := repo.GetLoan(ctx, "book2", "user_xyz")
		assert.Error(t, err)
		assert.Nil(t, loan)
		assert.Equal(t, ErrLoanNotFound, err)
	})
}

func TestLoanRepository_UpdateLoan(t *testing.T) {
	repo := NewLoanRepository()
	ctx := context.Background()

	loanDetail := &models.Loan{
		LoanDate:     time.Now(),
		ReturnDate:   time.Now().AddDate(0, 0, 28),
		IsReturn:     false,
		BorrowerName: "user3",
		BookId:       3,
	}
	_, err := repo.CreateLoan(ctx, "book3", loanDetail)
	assert.NoError(t, err)

	t.Run("Update existing loan", func(t *testing.T) {
		newReturnDate := time.Now().AddDate(0, 0, 35) // Extend loan
		updatedLoan, err := repo.UpdateLoan(ctx, "book3", "user3", &models.LoanUpdate{
			ReturnDate: &newReturnDate,
		})
		assert.NoError(t, err)
		assert.NotNil(t, updatedLoan)
		assert.Equal(t, newReturnDate, updatedLoan.ReturnDate)
	})

	currTime := time.Now()
	t.Run("Fail to update non-existent loan", func(t *testing.T) {
		_, err := repo.UpdateLoan(ctx, "book3", "user30", &models.LoanUpdate{
			ReturnDate: &currTime,
		})
		assert.Error(t, err)
		assert.Equal(t, ErrLoanNotFound, err)
	})
}

func TestLoanRepository_DeleteLoan(t *testing.T) {
	repo := NewLoanRepository()
	ctx := context.Background()

	loanDetail := &models.Loan{
		LoanDate:     time.Now(),
		ReturnDate:   time.Now().AddDate(0, 0, 28),
		IsReturn:     false,
		BorrowerName: "user4",
		BookId:       1,
	}
	_, _ = repo.CreateLoan(ctx, "book4", loanDetail)

	t.Run("Delete existing loan", func(t *testing.T) {
		err := repo.DeleteLoan(ctx, "book4", "user4")
		assert.NoError(t, err)

		// Verify loan is removed
		loan, err := repo.GetLoan(ctx, "book4", "user4")
		assert.Error(t, err)
		assert.Nil(t, loan)
	})

	t.Run("Fail to delete non-existent loan", func(t *testing.T) {
		err := repo.DeleteLoan(ctx, "book4", "user_xyz")
		assert.Error(t, err)
		assert.Equal(t, ErrLoanNotFound, err)
	})
}
