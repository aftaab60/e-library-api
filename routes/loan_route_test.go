package routes

import (
	"context"
	"encoding/json"
	"github.com/aftaab60/e-library-api/models"
	"github.com/aftaab60/e-library-api/repositories"
	"github.com/aftaab60/e-library-api/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestLoanRoute_BorrowBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Register the route
	loanRoute := NewLoanRoute(services.NewLoanService(repositories.NewLoanRepository(), repositories.NewBookRepository()))
	router.POST("/borrow", loanRoute.BorrowBook)

	t.Run("Successfully borrow a book", func(t *testing.T) {
		requestBody := `{"title": "book1", "borrower_name": "borrower1"}`
		req, err := http.NewRequest(http.MethodPost, "/borrow", strings.NewReader(requestBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("borrow non-available book", func(t *testing.T) {
		requestBody := `{"title": "book10", "borrower_name": "borrower1"}`
		req, err := http.NewRequest(http.MethodPost, "/borrow", strings.NewReader(requestBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		requestBody := `{"title": "", "borrower_name": "borrower1"}`
		req, err := http.NewRequest(http.MethodPost, "/borrow", strings.NewReader(requestBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestLoanRoute_ExtendLoan(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Register the route
	loanRepository := repositories.NewLoanRepository()
	loanRoute := NewLoanRoute(services.NewLoanService(loanRepository, repositories.NewBookRepository()))
	router.POST("/extend", loanRoute.ExtendLoan)

	t.Run("Extend a loan where book doesn't exist", func(t *testing.T) {
		requestBody := `{"title": "book100", "borrower_name": "borrower1"}`
		req, err := http.NewRequest(http.MethodPost, "/extend", strings.NewReader(requestBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("successfully extend a loan", func(t *testing.T) {
		currTime := time.Now()
		_, err := loanRepository.CreateLoan(context.Background(), "book1", &models.Loan{
			Id:           1,
			BookId:       1,
			BorrowerName: "user1",
			LoanDate:     currTime,
			ReturnDate:   currTime,
			IsReturn:     false,
		})
		assert.NoError(t, err)

		requestBody := `{"title": "book1", "borrower_name": "user1"}`
		req, err := http.NewRequest(http.MethodPost, "/extend", strings.NewReader(requestBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response models.LoanDetail
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, response.ReturnDate, currTime.AddDate(0, 0, 21))

	})
}

func TestLoanRoute_ReturnBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Register the route
	loanRepository := repositories.NewLoanRepository()
	loanRoute := NewLoanRoute(services.NewLoanService(loanRepository, repositories.NewBookRepository()))
	router.POST("/return", loanRoute.ReturnBook)

	t.Run("Return an invalid loan", func(t *testing.T) {
		requestBody := `{"title": "book100", "borrower_name": "borrower1"}`
		req, err := http.NewRequest(http.MethodPost, "/return", strings.NewReader(requestBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("successfully return a loan", func(t *testing.T) {
		currTime := time.Now()
		_, err := loanRepository.CreateLoan(context.Background(), "book1", &models.Loan{
			Id:           1,
			BookId:       1,
			BorrowerName: "user2",
			LoanDate:     currTime,
			ReturnDate:   currTime,
			IsReturn:     false,
		})
		assert.NoError(t, err)

		requestBody := `{"title": "book1", "borrower_name": "user2"}`
		req, err := http.NewRequest(http.MethodPost, "/return", strings.NewReader(requestBody))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

}
