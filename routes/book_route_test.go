package routes

import (
	"github.com/aftaab60/e-library-api/repositories"
	"github.com/aftaab60/e-library-api/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBookByTitle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	// Register the route
	bookRoute := NewBookRoute(services.NewBookService(repositories.NewBookRepository()))
	router.GET("/book/:title", bookRoute.GetBookByTitle)

	t.Run("success", func(t *testing.T) {
		//create request
		req, err := http.NewRequest(http.MethodGet, "/book/book1", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		expectedBody := `{"title": "book1", "available_copies": 5}`
		assert.JSONEq(t, expectedBody, rec.Body.String())
	})

	t.Run("invalid book title", func(t *testing.T) {
		//create request
		req, err := http.NewRequest(http.MethodGet, "/book/ ", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		expectedBody := `{"message": "title is empty"}`
		assert.JSONEq(t, expectedBody, rec.Body.String())
	})

	t.Run("book not found", func(t *testing.T) {
		//create request
		req, err := http.NewRequest(http.MethodGet, "/book/book100", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		expectedBody := `{"message":"book not found"}`
		assert.JSONEq(t, expectedBody, rec.Body.String())
	})
}
