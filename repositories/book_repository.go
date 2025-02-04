package repositories

import (
	"context"
	"errors"
	"github.com/aftaab60/e-library-api/models"
	"sync"
)

type IBookRepository interface {
	GetBook(ctx context.Context, title string) (*models.Book, error)
	UpdateBook(ctx context.Context, title string, availableQuantity int) (*models.Book, error)
}

type BookRepository struct {
	books map[string]*models.Book
	mutex sync.RWMutex
}

func NewBookRepository() *BookRepository {
	repo := &BookRepository{
		books: make(map[string]*models.Book),
	}
	repo.initBookRepository()
	return repo
}

// initialise some books by default at launch
func (br *BookRepository) initBookRepository() {
	br.mutex.Lock()
	defer br.mutex.Unlock()

	books := []models.Book{
		{Id: 1, Title: "book1", AvailableCopies: 5},
		{Id: 2, Title: "book2", AvailableCopies: 3},
		{Id: 3, Title: "book3", AvailableCopies: 1},
		{Id: 4, Title: "book4", AvailableCopies: 0},
	}
	for _, book := range books {
		br.books[book.Title] = &book
	}
}

// ErrBookNotFound is returned when a book is not found
var ErrBookNotFound = errors.New("book not found")

func (br *BookRepository) GetBook(ctx context.Context, title string) (*models.Book, error) {
	book, ok := br.books[title]
	if !ok {
		return nil, ErrBookNotFound
	}
	return book, nil
}

func (br *BookRepository) UpdateBook(ctx context.Context, title string, availableCopies int) (*models.Book, error) {
	book, ok := br.books[title]
	if !ok {
		return nil, ErrBookNotFound
	}
	book.AvailableCopies = availableCopies
	br.books[title] = book
	return book, nil
}
