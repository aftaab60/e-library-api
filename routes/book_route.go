package routes

import (
	"errors"
	"github.com/aftaab60/e-library-api/repositories"
	"github.com/aftaab60/e-library-api/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type BookRoute struct {
	BookService services.BookService
}

func NewBookRoute(bookService services.BookService) *BookRoute {
	return &BookRoute{bookService}
}

func (r *BookRoute) GetBookByTitle(c *gin.Context) {
	ctx := c.Request.Context()
	title := strings.TrimSpace(c.Param("title"))
	if err := r.validateTitle(title); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	book, err := r.BookService.GetBookByTitle(ctx, title)
	if err != nil {
		if errors.Is(err, repositories.ErrBookNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, book)
}

var ErrTitleEmpty = errors.New("title is empty")

func (r *BookRoute) validateTitle(title string) error {
	if len(title) == 0 {
		return ErrTitleEmpty
	}
	//additional validations if needed...
	return nil
}
