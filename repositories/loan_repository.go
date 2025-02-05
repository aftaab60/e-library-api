package repositories

import (
	"context"
	"errors"
	"github.com/aftaab60/e-library-api/models"
	"sync"
)

type ILoanRepository interface {
	GetLoan(ctx context.Context, title string, borrowerName string) (*models.Loan, error)
	CreateLoan(ctx context.Context, title string, loanDetail *models.Loan) (*models.Loan, error)
	UpdateLoan(ctx context.Context, title string, borrowerName string, loanUpdate *models.LoanUpdate) (*models.Loan, error)
	DeleteLoan(ctx context.Context, title string, borrowerName string) error
}

type LoanRepository struct {
	//book_id: All loans of this book title. Value can also be map[borrower]Loan but keeping slice for simplicity
	loans map[string][]models.Loan
	mutex sync.RWMutex
}

func NewLoanRepository() *LoanRepository {
	return &LoanRepository{
		loans: make(map[string][]models.Loan),
	}
}

// ErrLoanNotFound is returned when a book is not found
var ErrLoanNotFound = errors.New("loan not found")

// ErrExistingActiveLoan is returned when a book is not found
var ErrExistingActiveLoan = errors.New("existing active loan")

func (l *LoanRepository) GetLoan(ctx context.Context, title string, borrowerName string) (*models.Loan, error) {
	loanDetails, exists := l.loans[title]
	if !exists {
		return nil, ErrLoanNotFound
	}
	for _, loanDetail := range loanDetails {
		if loanDetail.BorrowerName == borrowerName && !loanDetail.IsReturn {
			return &loanDetail, nil
		}
	}
	return nil, ErrLoanNotFound
}

func (l *LoanRepository) CreateLoan(ctx context.Context, title string, loanDetail *models.Loan) (*models.Loan, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	loanDetails, exists := l.loans[title]
	if exists {
		for _, loan := range loanDetails {
			if loanDetail.BorrowerName == loan.BorrowerName && !loan.IsReturn {
				return nil, ErrExistingActiveLoan
			}
		}
	} else {
		loanDetails = make([]models.Loan, 0)
	}

	loanDetail.Id = len(loanDetails) + 1 //incremental id
	loanDetails = append(loanDetails, *loanDetail)
	l.loans[title] = loanDetails

	return &loanDetails[len(loanDetails)-1], nil
}

func (l *LoanRepository) UpdateLoan(ctx context.Context, title string, borrowerName string, loanUpdate *models.LoanUpdate) (*models.Loan, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	loanDetails, exists := l.loans[title]
	if !exists || loanUpdate == nil {
		return nil, ErrLoanNotFound
	}

	var updatedLoan *models.Loan
	for i, loan := range loanDetails {
		if loan.BorrowerName == borrowerName {
			//update values if not null
			if loanUpdate.ReturnDate != nil {
				loanDetails[i].ReturnDate = *loanUpdate.ReturnDate
			}
			if loanUpdate.IsReturn != nil {
				loanDetails[i].IsReturn = *loanUpdate.IsReturn
			}
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

func (l *LoanRepository) DeleteLoan(ctx context.Context, title string, borrowerName string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	loanDetails, exists := l.loans[title]
	if !exists {
		return ErrLoanNotFound
	}

	for i, loanDetail := range loanDetails {
		if loanDetail.BorrowerName == borrowerName {
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
