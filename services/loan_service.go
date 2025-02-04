package services

import (
	"context"
	"errors"
	"github.com/aftaab60/e-library-api/internal/db_manager"
	"github.com/aftaab60/e-library-api/models"
	"github.com/aftaab60/e-library-api/repositories"
	"log"
	"time"
)

type LoanService struct {
	LoanRepository repositories.ILoanRepository
	BookRepository repositories.IBookRepository
	TxDB           db_manager.ItxDB
}

// NewLoanService uses interface so that we can switch between in-memory and actual pgsql repo data easily
func NewLoanService(loanRepository repositories.ILoanRepository, bookRepository repositories.IBookRepository) LoanService {
	return LoanService{
		LoanRepository: loanRepository,
		BookRepository: bookRepository,
	}
}

func (s *LoanService) GetLoanDetailByTitleAndBorrower(ctx context.Context, title string, borrowerName string) (*models.Loan, error) {
	loan, err := s.LoanRepository.GetLoan(ctx, title, borrowerName)
	if err != nil {
		if errors.Is(err, repositories.ErrLoanNotFound) {
			log.Printf("Loan '%s' not found", title)
		} else {
			log.Printf("error getting Loan from repository: %v", err)
		}
		return nil, err
	}
	return loan, nil
}

var ErrExistingLoanFound = errors.New("existing loan found")

var ErrNoAvailableCopiesFound = errors.New("no available copies found")

func (s *LoanService) BorrowBook(ctx context.Context, title string, borrowerName string) (*models.LoanDetail, error) {
	//check existing loan
	loan, err := s.LoanRepository.GetLoan(ctx, title, borrowerName)
	if err != nil && !errors.Is(err, repositories.ErrLoanNotFound) {
		return nil, err
	}
	if loan != nil {
		return nil, ErrExistingLoanFound
	}

	//check book and availability
	book, err := s.BookRepository.GetBook(ctx, title)
	if err != nil {
		log.Printf("error getting book: %v", err)
		return nil, err
	}
	if book.AvailableCopies == 0 {
		return nil, ErrNoAvailableCopiesFound
	}

	//book and loan, both should be part of atomic operation and need to run in a transaction
	if err = db_manager.WrapInTransaction(ctx, s.TxDB, func(ctx context.Context) error {
		if _, err = s.BookRepository.UpdateBook(ctx, title, book.AvailableCopies-1); err != nil {
			log.Printf("error updating book available copies: %v", err)
			return err
		}

		loan, err = s.LoanRepository.CreateLoan(ctx, title, &models.Loan{
			BookId:       book.Id,
			BorrowerName: borrowerName,
			LoanDate:     time.Now(),
			ReturnDate:   time.Now().AddDate(0, 0, 28),
			IsReturn:     false,
		})
		if err != nil {
			log.Printf("error creating loan from repository: %v", err)
			return err
		}
		return nil
	}, nil); err != nil {
		log.Printf("error running book and loan update transaction: %v", err)
		return nil, err
	}

	/*
		if _, err = s.BookRepository.UpdateBook(ctx, title, book.AvailableCopies-1); err != nil {
			log.Printf("error updating book available copies: %v", err)
			return nil, err
		}

		loan, err = s.LoanRepository.CreateLoan(ctx, title, &models.Loan{
			BookId:       book.Id,
			BorrowerName: borrowerName,
			LoanDate:     time.Now(),
			ReturnDate:   time.Now().AddDate(0, 0, 28),
			IsReturn:     false,
		})
		if err != nil {
			log.Printf("error creating loan from repository: %v", err)
			return nil, err
		}
	*/

	return &models.LoanDetail{
		NameOfBorrower: loan.BorrowerName,
		LoanDate:       loan.LoanDate,
		ReturnDate:     loan.ReturnDate,
	}, nil
}

func (s *LoanService) ExtendLoan(ctx context.Context, title string, borrowerName string) (*models.LoanDetail, error) {
	loan, err := s.LoanRepository.GetLoan(ctx, title, borrowerName)
	if err != nil {
		if errors.Is(err, repositories.ErrLoanNotFound) {
			log.Printf("Loan '%s' not found", title)
		} else {
			log.Printf("error getting Loan from repository: %v", err)
		}
	}

	//extend 3 more weeks
	t := loan.ReturnDate.AddDate(0, 0, 21)
	updatedLoanDetail, err := s.LoanRepository.UpdateLoan(ctx, title, borrowerName, &models.LoanUpdate{
		ReturnDate: &t,
	})
	if err != nil {
		log.Printf("error updating Loan from repository: %v", err)
		return nil, err
	}
	return &models.LoanDetail{
		NameOfBorrower: updatedLoanDetail.BorrowerName,
		LoanDate:       updatedLoanDetail.LoanDate,
		ReturnDate:     updatedLoanDetail.ReturnDate,
	}, nil
}

func (s *LoanService) ReturnBook(ctx context.Context, title string, borrowerName string) error {
	book, err := s.BookRepository.GetBook(ctx, title)
	if err != nil {
		return err
	}

	t := time.Now()
	isReturn := true
	updatedLoan, err := s.LoanRepository.UpdateLoan(ctx, title, borrowerName, &models.LoanUpdate{
		ReturnDate: &t,
		IsReturn:   &isReturn,
	})
	if err != nil {
		if errors.Is(err, repositories.ErrLoanNotFound) {
			log.Printf("Loan '%s' not found", title)
		}
		return err
	}

	if _, err := s.BookRepository.UpdateBook(ctx, title, book.AvailableCopies+1); err != nil {
		log.Printf("error updating book available copies: %v", err)
		return err
	}

	log.Printf("book has been returned, loanId: %d, booktitle: %s, borrowerName: %s\n", updatedLoan.Id, title, borrowerName)
	return nil
}
