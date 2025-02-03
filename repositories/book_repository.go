package repositories

import (
	"errors"
	"github.com/aftaab60/e-library-api/models"
	"sync"
)

type IBookRepository interface {
	GetBook(title string) (*models.BookDetail, error)
	UpdateBook(title string, availableQuantity int) (*models.BookDetail, error)
}

type BookRepository struct {
	books map[string]*models.BookDetail
	mutex sync.RWMutex
}

func NewBookRepository() *BookRepository {
	repo := &BookRepository{
		books: make(map[string]*models.BookDetail),
	}
	repo.initBookRepository()
	return repo
}

// initialise some books by default at launch
func (br *BookRepository) initBookRepository() {
	br.mutex.Lock()
	defer br.mutex.Unlock()

	books := []models.BookDetail{
		{Title: "book1", AvailableCopies: 5},
		{Title: "book2", AvailableCopies: 3},
		{Title: "book3", AvailableCopies: 1},
		{Title: "book4", AvailableCopies: 0},
	}
	for _, book := range books {
		br.books[book.Title] = &book
	}
}

// ErrBookNotFound is returned when a book is not found
var ErrBookNotFound = errors.New("book not found")

func (br *BookRepository) GetBook(title string) (*models.BookDetail, error) {
	book, ok := br.books[title]
	if !ok {
		return nil, ErrBookNotFound
	}
	return book, nil
}

func (br *BookRepository) UpdateBook(title string, availableCopies int) (*models.BookDetail, error) {
	book, ok := br.books[title]
	if !ok {
		return nil, ErrBookNotFound
	}
	book.AvailableCopies = availableCopies
	br.books[title] = book
	return book, nil
}
