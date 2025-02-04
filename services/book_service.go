package services

import (
	"context"
	"errors"
	"github.com/aftaab60/e-library-api/models"
	"github.com/aftaab60/e-library-api/repositories"
	"log"
)

type BookService struct {
	bookRepository repositories.IBookRepository
}

// NewBookService uses interface so that we can switch between in-memory and actual pgsql repo data easily
func NewBookService(bookRepository repositories.IBookRepository) BookService {
	return BookService{bookRepository: bookRepository}
}

func (s *BookService) GetBookByTitle(ctx context.Context, title string) (*models.BookDetail, error) {
	book, err := s.bookRepository.GetBook(ctx, title)
	if err != nil {
		if errors.Is(err, repositories.ErrBookNotFound) {
			log.Printf("book '%s' not found", title)
		} else {
			log.Printf("error getting book from repository: %v", err)
		}
		return nil, err
	}

	return &models.BookDetail{
		Title:           book.Title,
		AvailableCopies: book.AvailableCopies,
	}, nil
}
