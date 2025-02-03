package routes

import (
	"errors"
	"fmt"
	"github.com/aftaab60/e-library-api/models"
	"github.com/aftaab60/e-library-api/repositories"
	"github.com/aftaab60/e-library-api/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type LoanRoute struct {
	LoanService services.LoanService
}

func NewLoanRoute(LoanService services.LoanService) *LoanRoute {
	return &LoanRoute{LoanService}
}

func (r *LoanRoute) BorrowBook(c *gin.Context) {
	var request models.LoanRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request body, err: %s", err.Error())})
		return
	}
	request.Title = strings.TrimSpace(request.Title)
	request.BorrowerName = strings.TrimSpace(request.BorrowerName)

	//validate request body for certain parameters
	if err := request.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	LoanDetail, err := r.LoanService.BorrowBook(request.Title, request.BorrowerName)
	if err != nil {
		if errors.Is(err, repositories.ErrBookNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		} else if errors.Is(err, services.ErrNoAvailableCopiesFound) {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, LoanDetail)
}

func (r *LoanRoute) ExtendLoan(c *gin.Context) {
	var request models.LoanRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request body, err: %s", err.Error())})
		return
	}
	request.Title = strings.TrimSpace(request.Title)
	request.BorrowerName = strings.TrimSpace(request.BorrowerName)

	//validate request body for certain parameters
	if err := request.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	LoanDetail, err := r.LoanService.ExtendLoan(request.Title, request.BorrowerName)
	if err != nil {
		if errors.Is(err, repositories.ErrLoanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, LoanDetail)
}

func (r *LoanRoute) ReturnBook(c *gin.Context) {
	var request models.LoanRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request body, err: %s", err.Error())})
		return
	}
	request.Title = strings.TrimSpace(request.Title)
	request.BorrowerName = strings.TrimSpace(request.BorrowerName)

	//validate request body for certain parameters
	if err := request.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.LoanService.ReturnBook(request.Title, request.BorrowerName); err != nil {
		if errors.Is(err, repositories.ErrLoanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "book returned"})
}
