package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/aftaab60/e-library-api/internal/db_manager"
	"github.com/aftaab60/e-library-api/models"
)

type LoanRepositoryDB struct {
	DB *db_manager.DB
}

func NewLoanRepositoryDB(db *db_manager.DB) *LoanRepositoryDB {
	return &LoanRepositoryDB{
		DB: db,
	}
}

func (l *LoanRepositoryDB) GetLoan(ctx context.Context, title string, borrowerName string) (*models.Loan, error) {
	query := `
        SELECT l.id, l.book_id, l.borrower_name, l.loan_date, l.return_date, l.is_returned
        FROM loans l
        JOIN books b ON l.book_id = b.id
        WHERE b.title = $1 AND l.borrower_name = $2 AND l.is_returned = FALSE
    `
	row := l.DB.GetRecord(ctx, query, title, borrowerName)

	var loan models.Loan
	if err := row.Scan(&loan.Id, &loan.BookId, &loan.BorrowerName, &loan.LoanDate, &loan.ReturnDate, &loan.IsReturn); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLoanNotFound
		}
		return nil, err
	}
	return &loan, nil
}

func (l *LoanRepositoryDB) CreateLoan(ctx context.Context, title string, loan *models.Loan) (*models.Loan, error) {
	insertQuery := `
        INSERT INTO loans (book_id, borrower_name, loan_date, return_date, is_returned)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, book_id, borrower_name, loan_date, return_date, is_returned
    `
	row := l.DB.CreateRecord(ctx, insertQuery, loan.BookId, loan.BorrowerName, loan.LoanDate, loan.ReturnDate, loan.IsReturn)
	if row == nil {
		return nil, sql.ErrNoRows
	}

	var insertedLoan models.Loan
	err := row.Scan(&insertedLoan.Id, &insertedLoan.BookId, &insertedLoan.BorrowerName, &insertedLoan.LoanDate, &insertedLoan.ReturnDate, &insertedLoan.IsReturn)
	if err != nil {
		return nil, fmt.Errorf("error creating loan for title %s: %w", title, err)
	}
	return &insertedLoan, nil
}

func (l *LoanRepositoryDB) UpdateLoan(ctx context.Context, title string, borrowerName string, loanUpdate *models.LoanUpdate) (*models.Loan, error) {
	var bookID int
	query := "SELECT id FROM books WHERE title = $1"
	row := l.DB.GetRecord(ctx, query, title)
	err := row.Scan(&bookID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("book with title %s not found", title)
		}
		return nil, fmt.Errorf("error fetching book id: %w", err)
	}

	updateQuery := `
        UPDATE loans
        SET return_date = $1, is_returned = $2
        WHERE book_id = $3 AND borrower_name = $4 AND is_returned = FALSE
        RETURNING id, book_id, borrower_name, loan_date, return_date, is_returned
    `
	row = l.DB.UpdateRecord(ctx, updateQuery, loanUpdate.ReturnDate, loanUpdate.IsReturn, bookID, borrowerName)
	if row == nil {
		return nil, sql.ErrNoRows
	}

	var updatedLoan models.Loan
	err = row.Scan(&updatedLoan.Id, &updatedLoan.BookId, &updatedLoan.BorrowerName, &updatedLoan.LoanDate, &updatedLoan.ReturnDate, &updatedLoan.IsReturn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLoanNotFound
		}
		return nil, fmt.Errorf("error updating loan for title %s: %w", title, err)
	}

	return &updatedLoan, nil
}

func (l *LoanRepositoryDB) DeleteLoan(ctx context.Context, title string, borrowerName string) error {
	// Step 1: Get the book_id for the given title
	var bookID int
	query := "SELECT id FROM books WHERE title = $1"
	row := l.DB.GetRecord(ctx, query, title)
	err := row.Scan(&bookID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("book with title %s not found", title)
		}
		return fmt.Errorf("error fetching book id: %w", err)
	}

	loanQuery := `
        SELECT id
        FROM loans
        WHERE book_id = $1 AND borrower_name = $2 AND is_returned = FALSE
    `
	row = l.DB.GetRecord(ctx, loanQuery, bookID, borrowerName)

	var loanID int
	err = row.Scan(&loanID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no active loan found for borrower %s with book title %s", borrowerName, title)
		}
		return fmt.Errorf("error fetching loan record: %w", err)
	}

	deleteQuery := "DELETE FROM loans WHERE id = $1"
	_, err = l.DB.DeleteRecord(ctx, deleteQuery, loanID)
	if err != nil {
		return fmt.Errorf("error deleting loan: %w", err)
	}

	return nil
}
