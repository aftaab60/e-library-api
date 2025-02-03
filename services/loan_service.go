package services

import (
	"errors"
	"github.com/aftaab60/e-library-api/models"
	"github.com/aftaab60/e-library-api/repositories"
	"log"
	"time"
)

type LoanService struct {
	LoanRepository repositories.ILoanRepository
	BookRepository repositories.IBookRepository
}

// NewLoanService uses interface so that we can switch between in-memory and actual pgsql repo data easily
func NewLoanService(loanRepository repositories.ILoanRepository, bookRepository repositories.IBookRepository) LoanService {
	return LoanService{
		LoanRepository: loanRepository,
		BookRepository: bookRepository,
	}
}

func (s *LoanService) GetLoanDetailByTitleAndBorrower(title string, borrowerName string) (*models.LoanDetail, error) {
	LoanDetail, err := s.LoanRepository.GetLoan(title, borrowerName)
	if err != nil {
		if errors.Is(err, repositories.ErrLoanNotFound) {
			log.Printf("Loan '%s' not found", title)
		} else {
			log.Printf("error getting Loan from repository: %v", err)
		}
		return nil, err
	}
	return LoanDetail, nil
}

var ErrExistingLoanFound = errors.New("existing loan found")

var ErrNoAvailableCopiesFound = errors.New("no available copies found")

func (s *LoanService) BorrowBook(title string, borrowerName string) (*models.LoanDetail, error) {
	//check existing loan
	LoanDetail, err := s.LoanRepository.GetLoan(title, borrowerName)
	if err != nil && !errors.Is(err, repositories.ErrLoanNotFound) {
		return nil, err
	}
	if LoanDetail != nil {
		return nil, ErrExistingLoanFound
	}

	//check book and availability
	book, err := s.BookRepository.GetBook(title)
	if err != nil {
		log.Printf("error getting book: %v", err)
		return nil, err
	}
	if book.AvailableCopies == 0 {
		return nil, ErrNoAvailableCopiesFound
	}

	if _, err = s.BookRepository.UpdateBook(title, book.AvailableCopies-1); err != nil {
		log.Printf("error updating book available copies: %v", err)
		return nil, err
	}

	LoanDetail, err = s.LoanRepository.CreateLoan(title, &models.LoanDetail{
		NameOfBorrower: borrowerName,
		LoanDate:       time.Now(),
		ReturnDate:     time.Now().AddDate(0, 0, 28),
	})
	if err != nil {
		log.Printf("error creating Loan from repository: %v", err)
		return nil, err
	}
	return LoanDetail, nil
}

func (s *LoanService) ExtendLoan(title string, borrowerName string) (*models.LoanDetail, error) {
	LoanDetail, err := s.LoanRepository.GetLoan(title, borrowerName)
	if err != nil {
		if errors.Is(err, repositories.ErrLoanNotFound) {
			log.Printf("Loan '%s' not found", title)
		} else {
			log.Printf("error getting Loan from repository: %v", err)
		}
	}

	//extend 3 more weeks
	updatedLoanDetail, err := s.LoanRepository.UpdateLoan(title, &models.LoanDetail{
		NameOfBorrower: LoanDetail.NameOfBorrower,
		LoanDate:       LoanDetail.LoanDate,
		ReturnDate:     LoanDetail.ReturnDate.AddDate(0, 0, 21),
	})
	if err != nil {
		log.Printf("error updating Loan from repository: %v", err)
		return nil, err
	}
	return updatedLoanDetail, nil
}

func (s *LoanService) ReturnBook(title string, borrowerName string) error {
	err := s.LoanRepository.DeleteLoan(title, borrowerName)
	if err != nil {
		if errors.Is(err, repositories.ErrLoanNotFound) {
			log.Printf("Loan '%s' not found", title)
		}
		return err
	}

	book, err := s.BookRepository.GetBook(title)
	if err != nil {
		return err
	}
	if _, err := s.BookRepository.UpdateBook(title, book.AvailableCopies+1); err != nil {
		log.Printf("error updating book available copies: %v", err)
		return err
	}

	return nil
}
