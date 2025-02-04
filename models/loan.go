package models

import (
	"errors"
	"time"
)

type Loan struct {
	Id           int       `json:"id"`
	BookId       int       `json:"book_id"`
	BorrowerName string    `json:"borrower_name"`
	LoanDate     time.Time `json:"loan_date"`
	ReturnDate   time.Time `json:"return_date"`
	IsReturn     bool      `json:"is_return"`
}

type LoanDetail struct {
	NameOfBorrower string    `json:"name_of_borrower"`
	LoanDate       time.Time `json:"loan_date"`
	ReturnDate     time.Time `json:"return_date"`
}

type LoanRequest struct {
	Title        string `json:"title"`
	BorrowerName string `json:"borrower_name"`
}

func (b *LoanRequest) Validate() error {
	if len(b.Title) == 0 {
		return errors.New("missing title")
	}
	if len(b.BorrowerName) == 0 {
		return errors.New("missing borrower_name")
	}
	//other validations as needed...
	return nil
}

// LoanUpdate allowed fields that can be updated
type LoanUpdate struct {
	ReturnDate *time.Time `json:"return_date"`
	IsReturn   *bool      `json:"is_return"`
}
