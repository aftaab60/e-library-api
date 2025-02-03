package repositories

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBookRepository_GetBook(t *testing.T) {
	repo := NewBookRepository()

	t.Run("Get existing book", func(t *testing.T) {
		book, err := repo.GetBook("book1")
		assert.NoError(t, err)
		assert.NotNil(t, book)
		assert.Equal(t, "book1", book.Title)
		assert.Equal(t, 5, book.AvailableCopies)
	})

	t.Run("Get non-existent book", func(t *testing.T) {
		book, err := repo.GetBook("nonexistent")
		assert.Error(t, err)
		assert.Nil(t, book)
		assert.Equal(t, ErrBookNotFound, err)
	})
}

func TestBookRepository_UpdateBook(t *testing.T) {
	repo := NewBookRepository()

	t.Run("Update existing book", func(t *testing.T) {
		updatedBook, err := repo.UpdateBook("book1", 2)
		assert.NoError(t, err)
		assert.NotNil(t, updatedBook)
		assert.Equal(t, 2, updatedBook.AvailableCopies)

		// Verify the book was updated
		book, err := repo.GetBook("book1")
		assert.NoError(t, err)
		assert.Equal(t, 2, book.AvailableCopies)
	})

	t.Run("Update non-existent book", func(t *testing.T) {
		updatedBook, err := repo.UpdateBook("nonexistent", 2)
		assert.Error(t, err)
		assert.Nil(t, updatedBook)
		assert.Equal(t, ErrBookNotFound, err)
	})
}
