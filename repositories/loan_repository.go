package repositories

import (
	"errors"
	"github.com/aftaab60/e-library-api/models"
	"sync"
)

type ILoanRepository interface {
	GetLoan(title string, borrowerName string) (*models.LoanDetail, error)
	CreateLoan(title string, loanDetail *models.LoanDetail) (*models.LoanDetail, error)
	UpdateLoan(title string, loanDetail *models.LoanDetail) (*models.LoanDetail, error)
	DeleteLoan(title string, borrowerName string) error
}

type LoanRepository struct {
	loans map[string][]models.LoanDetail
	mutex sync.RWMutex
}

func NewLoanRepository() *LoanRepository {
	return &LoanRepository{
		loans: make(map[string][]models.LoanDetail),
	}
}

// ErrLoanNotFound is returned when a book is not found
var ErrLoanNotFound = errors.New("loan not found")

// ErrExistingActiveLoan is returned when a book is not found
var ErrExistingActiveLoan = errors.New("existing active loan")

func (l *LoanRepository) GetLoan(title string, borrowerName string) (*models.LoanDetail, error) {
	loanDetails, exists := l.loans[title]
	if !exists {
		return nil, ErrLoanNotFound
	}
	for _, loanDetail := range loanDetails {
		if loanDetail.NameOfBorrower == borrowerName {
			return &loanDetail, nil
		}
	}
	return nil, ErrLoanNotFound
}

func (l *LoanRepository) CreateLoan(title string, loanDetail *models.LoanDetail) (*models.LoanDetail, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	loanDetails, exists := l.loans[title]
	if exists {
		for _, loan := range loanDetails {
			if loanDetail.NameOfBorrower == loan.NameOfBorrower {
				return nil, ErrExistingActiveLoan
			}
		}
	} else {
		loanDetails = make([]models.LoanDetail, 0)
	}

	loanDetails = append(loanDetails, *loanDetail)
	l.loans[title] = loanDetails

	return &loanDetails[len(loanDetails)-1], nil
}

func (l *LoanRepository) UpdateLoan(title string, loanDetail *models.LoanDetail) (*models.LoanDetail, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	loanDetails, exists := l.loans[title]
	if !exists || loanDetail == nil {
		return nil, ErrLoanNotFound
	}

	var updatedLoan *models.LoanDetail
	for i, loan := range loanDetails {
		if loan.NameOfBorrower == loanDetail.NameOfBorrower {
			loanDetails[i] = *loanDetail
			updatedLoan = &loanDetails[i]
			break
		}
	}

	if updatedLoan == nil {
		return nil, ErrLoanNotFound
	}
	l.loans[title] = loanDetails

	return updatedLoan, nil
}

func (l *LoanRepository) DeleteLoan(title string, borrowerName string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	loanDetails, exists := l.loans[title]
	if !exists {
		return ErrLoanNotFound
	}

	for i, loanDetail := range loanDetails {
		if loanDetail.NameOfBorrower == borrowerName {
			updatedLoanDetails := append(loanDetails[:i], loanDetails[i+1:]...)
			l.loans[title] = updatedLoanDetails

			if len(updatedLoanDetails) == 0 {
				delete(l.loans, title)
			}
			return nil
		}
	}
	return ErrLoanNotFound
}
