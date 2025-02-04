package repositories

import (
	"context"
	"github.com/aftaab60/e-library-api/internal/db_manager"
	"github.com/aftaab60/e-library-api/models"
)

type BookRepositoryDB struct {
	DB *db_manager.DB
}

func NewBookRepositoryDB(db *db_manager.DB) *BookRepositoryDB {
	return &BookRepositoryDB{DB: db}
}

func (br *BookRepositoryDB) GetBook(ctx context.Context, title string) (*models.Book, error) {
	query := "SELECT id, title, available_copies FROM books WHERE title = $1"
	result := br.DB.GetRecord(ctx, query, title)

	var bookDetail models.Book
	if err := result.Scan(&bookDetail.Id, &bookDetail.Title, &bookDetail.AvailableCopies); err != nil {
		return nil, err
	}
	return &bookDetail, nil
}

func (br *BookRepositoryDB) UpdateBook(ctx context.Context, title string, availableCopies int) (*models.Book, error) {
	query := "UPDATE books SET available_copies = $1 WHERE title = $2 RETURNING id, title, available_copies"
	row := br.DB.UpdateRecord(ctx, query, availableCopies, title)

	var bookDetail models.Book
	if err := row.Scan(&bookDetail.Id, &bookDetail.Title, &bookDetail.AvailableCopies); err != nil {
		return nil, err
	}
	return &bookDetail, nil
}
