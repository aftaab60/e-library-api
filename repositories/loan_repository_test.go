package repositories

import (
	"github.com/aftaab60/e-library-api/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLoanRepository_CreateLoan(t *testing.T) {
	repo := NewLoanRepository()

	loanDetail := &models.LoanDetail{
		NameOfBorrower: "user1",
		LoanDate:       time.Now(),
		ReturnDate:     time.Now().AddDate(0, 0, 28),
	}

	t.Run("Create new loan", func(t *testing.T) {
		createdLoan, err := repo.CreateLoan("book1", loanDetail)
		assert.NoError(t, err)
		assert.NotNil(t, createdLoan)
		assert.Equal(t, "user1", createdLoan.NameOfBorrower)
	})

	t.Run("Fail to create duplicate loan", func(t *testing.T) {
		_, err := repo.CreateLoan("book1", loanDetail)
		assert.Error(t, err)
		assert.Equal(t, ErrExistingActiveLoan, err)
	})
}

func TestLoanRepository_GetLoan(t *testing.T) {
	repo := NewLoanRepository()

	loanDetail := &models.LoanDetail{
		NameOfBorrower: "user2",
		LoanDate:       time.Now(),
		ReturnDate:     time.Now().AddDate(0, 0, 28),
	}
	_, err := repo.CreateLoan("book2", loanDetail)
	assert.NoError(t, err)

	t.Run("Get existing loan", func(t *testing.T) {
		loan, err := repo.GetLoan("book2", "user2")
		assert.NoError(t, err)
		assert.NotNil(t, loan)
		assert.Equal(t, "user2", loan.NameOfBorrower)
	})

	t.Run("Get non-existent loan", func(t *testing.T) {
		loan, err := repo.GetLoan("book2", "user_xyz")
		assert.Error(t, err)
		assert.Nil(t, loan)
		assert.Equal(t, ErrLoanNotFound, err)
	})
}

func TestLoanRepository_UpdateLoan(t *testing.T) {
	repo := NewLoanRepository()

	loanDetail := &models.LoanDetail{
		NameOfBorrower: "user3",
		LoanDate:       time.Now(),
		ReturnDate:     time.Now().AddDate(0, 0, 28),
	}
	_, err := repo.CreateLoan("book3", loanDetail)
	assert.NoError(t, err)

	t.Run("Update existing loan", func(t *testing.T) {
		newReturnDate := time.Now().AddDate(0, 0, 35) // Extend loan
		updatedLoan, err := repo.UpdateLoan("book3", &models.LoanDetail{
			NameOfBorrower: "user3",
			LoanDate:       loanDetail.LoanDate,
			ReturnDate:     newReturnDate,
		})

		assert.NoError(t, err)
		assert.NotNil(t, updatedLoan)
		assert.Equal(t, newReturnDate, updatedLoan.ReturnDate)
	})

	t.Run("Fail to update non-existent loan", func(t *testing.T) {
		_, err := repo.UpdateLoan("book3", &models.LoanDetail{NameOfBorrower: "user_xyz"})
		assert.Error(t, err)
		assert.Equal(t, ErrLoanNotFound, err)
	})
}

func TestLoanRepository_DeleteLoan(t *testing.T) {
	repo := NewLoanRepository()

	loanDetail := &models.LoanDetail{
		NameOfBorrower: "user4",
		LoanDate:       time.Now(),
		ReturnDate:     time.Now().AddDate(0, 0, 28),
	}
	_, _ = repo.CreateLoan("book4", loanDetail)

	t.Run("Delete existing loan", func(t *testing.T) {
		err := repo.DeleteLoan("book4", "user4")
		assert.NoError(t, err)

		// Verify loan is removed
		loan, err := repo.GetLoan("book4", "user4")
		assert.Error(t, err)
		assert.Nil(t, loan)
	})

	t.Run("Fail to delete non-existent loan", func(t *testing.T) {
		err := repo.DeleteLoan("book4", "user_xyz")
		assert.Error(t, err)
		assert.Equal(t, ErrLoanNotFound, err)
	})
}
